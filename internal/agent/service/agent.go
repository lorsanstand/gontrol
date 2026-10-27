package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lorsanstand/gontrol/internal/models"
)

type AgentRunner struct {
	log      *slog.Logger
	taskExec *TaskExecutor
	hub      *HubClient
	system   *SystemService
}

func NewAgentRunner(log *slog.Logger, tasker *TaskExecutor, hub *HubClient, system *SystemService) *AgentRunner {
	return &AgentRunner{log: log, taskExec: tasker, hub: hub, system: system}
}

func (a *AgentRunner) Serve(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	agentID := a.system.GetAgentID()
	if agentID == nil {
		agent, err := a.RegisterAgent(ctx)
		if err != nil {
			return err
		}

		agentID = &agent.ID
	}

	hb := models.HeartbeatRequest{AgentID: *agentID}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ticker.C:
		a.hub.SendHeartbeat(ctx)
	}
}

// RegisterAgent регистрирует агента и сохраняет AgentID.
func (a *AgentRunner) RegisterAgent(ctx context.Context) (*models.Agent, error) {
	specs, err := a.system.GetSpecs()
	if err != nil {
		return nil, fmt.Errorf("get specification: %w", err)
	}

	agent, err := a.hub.RegisterAgent(ctx, *specs)
	if err != nil {
		return nil, fmt.Errorf("register agent in hub: %w", err)
	}

	if err = a.system.SaveAgentID(agent.ID); err != nil {
		return nil, fmt.Errorf("save agent ID: %w", err)
	}

	return agent, nil
}
