package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lorsanstand/gontrol/internal/hub/config"
	deliveryHttp "github.com/lorsanstand/gontrol/internal/hub/delivery/http"
	"github.com/lorsanstand/gontrol/internal/hub/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &cfg}))
	logger.Info("starting gontrol-hub...")

	agentStore := store.NewAgentStore()

	mux := http.NewServeMux()

	deliveryHttp.RegisterRoutes(mux, agentStore)

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
