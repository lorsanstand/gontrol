package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/lorsanstand/gontrol/internal/migrator"
	"github.com/lorsanstand/gontrol/migrations"
)

type App struct {
	di *DI
}

func NewApp() *App {
	return &App{di: &DI{}}
}

func (a *App) Run() {
	// Ловим сигнал от системы.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger := a.di.GetLogger()
	cfg := a.di.GetCfg()

	logger.Info("starting gontrol-hub...")

	// Запускаем миграции.
	if err := a.applyMigrations(); err != nil {
		logger.Error("failed to apply migrations", slog.Any("error", err))
	}

	// Запускаем AutoCleaner в AgentStateRepo.
	go func() {
		stateRepo := a.di.GetAgentStateRepo()
		err := stateRepo.AutoCleaner(ctx)
		logger.Warn("auto cleaner state store stopped", slog.Any("error", err))
	}()

	// Закрываем базу данных в конце.
	defer a.di.GetDB().Close()

	srvHttp := a.di.GetHttpServer()
	srvGrpc := a.di.GetGrpcServer()

	// Запускаем HTTP сервер.
	go func() {
		logger.Info("starting http server")

		if err := srvHttp.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server error", slog.Any("error", err))
		}
	}()

	// Запускаем GRPC сервер.
	go func() {
		logger.Info("starting gRPC server")

		if err := srvGrpc.Serve(); err != nil {
			logger.Error("gRPC server error", slog.Any("error", err))
		}
	}()

	logger.Info("gontrol-hub is running",
		slog.String("HTTP-address", fmt.Sprintf("0.0.0.0:%d", cfg.Port)),
		slog.String("gRPC-address", fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort)),
	)

	// Делаем паузу, пока не получим сигнал от системы.
	<-ctx.Done()
	logger.Info("signal received, closing...")

	// Для повторного нажатия CTRL+C моментальное выключение.
	stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Запускаем graceful shutdown в http сервере.
	if err := srvHttp.Shutdown(shutdownCtx); err != nil {
		logger.Error("error http server stopped", slog.Any("error", err))
	}

	if err := srvGrpc.GracefulShutdown(ctx); err != nil {
		logger.Error("error gRPC server stopped", slog.Any("error", err))
	}

	logger.Info("server stop")
}

func (a *App) applyMigrations() error {
	cfg := a.di.GetCfg()

	migrationDB, err := sql.Open("pgx", cfg.DatabaseURL())
	if err != nil {
		return fmt.Errorf("open migration database: %w", err)
	}
	defer migrationDB.Close()

	migrate, err := migrator.New(migrations.MigrationsFs, ".")
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err = migrate.Apply(migrationDB); err != nil {
		return fmt.Errorf("apply to migration: %w", err)
	}

	return nil
}
