package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	deliveryHttp "github.com/lorsanstand/gontrol/internal/hub/delivery/http"
	"github.com/lorsanstand/gontrol/internal/hub/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting gontrol-hub...")

	agentStore := store.NewAgentStore()

	mux := http.NewServeMux()

	deliveryHttp.RegisterRoutes(mux, agentStore)

	srv := &http.Server{
		Addr:         ":8089",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(":8089", mux); err != nil {
		logger.Info("server stopped", slog.Any("error", err))
		return
	}
}
