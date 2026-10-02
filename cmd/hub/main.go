package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handlers"
	"github.com/lorsanstand/gontrol/internal/hub/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting gontrol-hub...")

	agentStore := store.NewAgentStore()

	mux := http.NewServeMux()

	hb := &handlers.HeartbeatHandler{Store: agentStore}

	mux.Handle("POST /api/v1/heartbeat", hb)

	http.ListenAndServe(":8089", mux)
}
