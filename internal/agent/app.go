package agent

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lorsanstand/gontrol/internal/agent/config"
	"github.com/lorsanstand/gontrol/internal/agent/service"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrCreateFile) {
			return nil, fmt.Errorf("file config.yaml created, please fill file")
		}
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	return &App{cfg: cfg, logger: logger}, nil
}

func (a *App) Run() int {
	a.logger.Info("starting gontrol-agent...")

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	hub := service.NewHubClient(*a.cfg.HubURL, client)
	task := &service.TaskExecutor{}
	agent := service.NewAgentRunner(a.logger, task, hub)

	err := agent.Serve()
	a.logger.Error("stop agent serve", slog.Any("error", err))

	return 0
}
