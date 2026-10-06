package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lorsanstand/gontrol/internal/agent/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting gontrol-agent...")

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	service.NewHubService("127.0.0.1:8089", client)
}
