package service

import (
	"context"

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
}

func NewAgentService(storage agentRepository) *AgentService {
	return &AgentService{store: storage}
}

func (a *AgentService) GetAllAgents(ctx context.Context) ([]modelsHub.Agent, error) {
	return a.store.GetAllAgents(), nil
}

func (a *AgentService) GetAgent(ctx context.Context, agentID string) (modelsHub.Agent, error) {
	return a.store.GetAgent(agentID)
}

func (a *AgentService) RegisterOrUpdate(ctx context.Context, heartbeat models.HeartbeatRequest) error {
	a.store.RegisterOrUpdate(heartbeat)
	return nil
}
