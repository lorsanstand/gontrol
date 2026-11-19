package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger := a.di.GetLogger()
	cfg := a.di.GetCfg()

	logger.Info("starting gontrol-hub...")

	if err := a.applyMigrations(); err != nil {
		logger.Error("failed to apply migrations", slog.Any("error", err))
	}
	go func() {
		stateRepo := a.di.GetAgentStateRepo()
		err := stateRepo.AutoCleaner(ctx)
		logger.Warn("auto cleaner state store stopped", slog.Any("error", err))
	}()

	srv := a.di.GetHttpServer()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("server error", slog.Any("error", err))
		}
	}()

	logger.Info("gontrol-hub is running", slog.String("address", fmt.Sprintf("0.0.0.0:%s", cfg.Port)))

	<-ctx.Done()
	logger.Info("signal received, closing...")

	// Для повторного нажатия CTRL+C
	stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("error server stopped", slog.Any("error", err))
	}

	logger.Info("server stop")

	defer a.di.GetDB().Close()
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
