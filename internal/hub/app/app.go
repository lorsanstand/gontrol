package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := a.di.GetLogger()

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

	if err := srv.ListenAndServe(); err != nil {
		logger.Info("server stopped", slog.Any("error", err))
	}

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
