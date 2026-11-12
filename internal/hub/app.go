package hub

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lorsanstand/gontrol/internal/hub/config"
	deliveryHttp "github.com/lorsanstand/gontrol/internal/hub/delivery/http"
	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/middlewares"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	in_memory "github.com/lorsanstand/gontrol/internal/hub/store/in-memory"
	"github.com/lorsanstand/gontrol/internal/hub/store/postgres"
	"github.com/lorsanstand/gontrol/internal/migrator"
	"github.com/lorsanstand/gontrol/internal/utils/logger"
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

	log := logger.Init(cfg.LogLevel, cfg.LogFormat)

	return &App{logger: log, cfg: cfg}, nil
}

func (a *App) Run() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	specsRepo := postgres.NewCommandSpecsRepository(db)

	agentStateRepo := in_memory.NewAgentStateStore(TimeTTLOnlineAgent)
	go func() {
		err = agentStateRepo.AutoCleaner(ctx)
		a.logger.Warn("auto cleaner state store stopped", slog.Any("error", err))
	}()

	agentService := service.NewAgentService(a.logger, agentRepo, agentStateRepo, specsRepo)
	taskService := service.NewTasksService(a.logger, taskRepo, specsRepo)
	heartbeatService := service.NewHeartbeatService(agentService, taskService)

	r := chi.NewRouter()

	middlewares.ConnectMiddlewares(r, a.logger)
	apiRouter := deliveryHttp.RegisterRoutes(a.logger, taskService, agentService, heartbeatService)

	r.Mount("/api/v1", apiRouter)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.Port),
		Handler:      r,
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
