package main

import (
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/lorsanstand/gontrol/internal/agent/service"
	"github.com/lorsanstand/gontrol/internal/utils/config"
)

type Config struct {
	HubURL   url.URL    `env:"HUB_URL,required"`
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"WARN"`
}

func main() {
	var cfg Config
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	logger.Info("starting gontrol-agent...")

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	hub := service.NewHubClient(cfg.HubURL, client)
	task := &service.TaskExecutor{}
	agent := service.NewAgentRunner(logger, task, hub)

	err := agent.Serve()
	logger.Error("stop agent serve", slog.Any("error", err))
}
