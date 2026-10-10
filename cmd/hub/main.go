package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	deliveryHttp "github.com/lorsanstand/gontrol/internal/hub/delivery/http"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/utils/config"
)

type Config struct {
	Port     int        `env:"PORT" envDefault:"5467"`
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"WARN"`
}

func main() {
	var cfg Config
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &cfg.LogLevel}))
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
