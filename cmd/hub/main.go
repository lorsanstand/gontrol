package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lorsanstand/gontrol/internal/hub/config"
	deliveryHttp "github.com/lorsanstand/gontrol/internal/hub/delivery/http"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/migrator"
	"github.com/lorsanstand/gontrol/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &cfg.LogLevel}))
	logger.Info("starting gontrol-hub...")

	db, err := sql.Open("pgx", cfg.DatabaseURL())
	if err != nil {
		logger.Error("unable connect to database", slog.Any("error", err))
		return
	}
	defer db.Close()

	migrate, err := migrator.New(migrations.MigrationsFs, ".")
	if err != nil {
		logger.Error("failed to create migrator", slog.Any("error", err))
		return
	}

	if err = migrate.Apply(db); err != nil {
		logger.Error("unable apply to migration", slog.Any("error", err))
	}

	agentStore := store.NewAgentStore()
	agentService := service.NewAgentService(logger, agentStore)
	taskService := service.NewTasksService(logger, agentStore)
	heartbeatService := service.NewHeartbeatService(agentService, taskService)

	mux := http.NewServeMux()

	deliveryHttp.RegisterRoutes(mux, logger, taskService, agentService, heartbeatService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Info("server stopped", slog.Any("error", err))
		return
	}
}
