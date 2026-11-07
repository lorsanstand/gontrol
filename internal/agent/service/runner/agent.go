package runner

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lorsanstand/gontrol/internal/agent/service"
	"github.com/lorsanstand/gontrol/internal/models"
)

type AgentRunner struct {
	log      *slog.Logger
	taskExec *service.TaskExecutor
	hub      *service.HubClient
	system   *service.System
}

func NewAgentRunner(log *slog.Logger, tasker *service.TaskExecutor, hub *service.HubClient, system *service.System) *AgentRunner {
	return &AgentRunner{log: log, taskExec: tasker, hub: hub, system: system}
}

// Serve блокирующая функция, которая отсылает heartbeat и отправляет задачи на выполнение.
func (a *AgentRunner) Serve(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	agent, err := a.RegisterOrUpdate(ctx)
	if err != nil {
		return fmt.Errorf("register or update runner: %w", err)
	}

	hb := models.HeartbeatRequest{AgentID: agent.ID, OS: agent.OS, Hostname: agent.Hostname}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err = a.Heartbeat(ctx, hb)
			if err != nil {
				a.log.Error("failed send heartbeat", slog.Any("error", err))
			}
		}
	}
}
