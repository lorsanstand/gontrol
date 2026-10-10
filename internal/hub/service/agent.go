package service

import (
	"context"
	"log/slog"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/models"
)

type agentRepository interface {
	RegisterOrUpdate(agentRequest models.HeartbeatRequest)
	GetAgent(agentID string) (modelsHub.Agent, error)
	GetAllAgents() []modelsHub.Agent
}

type AgentService struct {
	store agentRepository
	log   *slog.Logger
}

func NewAgentService(logger *slog.Logger, storage agentRepository) *AgentService {
	return &AgentService{store: storage, log: logger}
}

func (a *AgentService) GetAll(ctx context.Context) ([]modelsHub.Agent, error) {
	return a.store.GetAllAgents(), nil
}

func (a *AgentService) Get(ctx context.Context, agentID string) (modelsHub.Agent, error) {
	return a.store.GetAgent(agentID)
}

func (a *AgentService) RegisterOrUpdate(ctx context.Context, heartbeat models.HeartbeatRequest) error {
	a.store.RegisterOrUpdate(heartbeat)
	return nil
}
