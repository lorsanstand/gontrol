package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lorsanstand/gontrol/internal/agent/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("starting gontrol-agent...")

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	hub := service.NewHubClient("http://127.0.0.1:8089/api/v1/", client)
	task := &service.TaskExecutor{}
	agent := service.NewAgentRunner(logger, task, hub)

	agent.Serve()
}
