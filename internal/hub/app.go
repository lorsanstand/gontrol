package hub

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lorsanstand/gontrol/internal/hub/config"
	deliveryHttp "github.com/lorsanstand/gontrol/internal/hub/delivery/http"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/hub/store/postgres"
	"github.com/lorsanstand/gontrol/internal/migrator"
	"github.com/lorsanstand/gontrol/migrations"
)

type App struct {
	logger *slog.Logger
	cfg    *config.Config
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &cfg.LogLevel}))

	return &App{logger: logger, cfg: cfg}, nil
}

func (a *App) Run() int {
	a.logger.Info("starting gontrol-hub...")

	db, err := sql.Open("pgx", a.cfg.DatabaseURL())
	if err != nil {
		a.logger.Error("unable connect to database", slog.Any("error", err))
		return 1
	}
	defer db.Close()

	if err = a.applyMigrations(); err != nil {
		a.logger.Error("failed to apply migrations", slog.Any("error", err))
		return 1
	}

	agentRepo := postgres.NewAgentRepository(db)
	taskRepo := postgres.NewTaskRepository(db)

	agentService := service.NewAgentService(a.logger, agentRepo)
	taskService := service.NewTasksService(a.logger, taskRepo)
	heartbeatService := service.NewHeartbeatService(agentService, taskService)

	mux := http.NewServeMux()

	deliveryHttp.RegisterRoutes(mux, a.logger, taskService, agentService, heartbeatService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		a.logger.Info("server stopped", slog.Any("error", err))
		return 1
	}

	return 0
}

func (a *App) applyMigrations() error {
	migrationDB, err := sql.Open("pgx", a.cfg.DatabaseURL())
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
