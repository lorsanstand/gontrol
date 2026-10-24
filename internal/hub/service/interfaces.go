package service

import (
	"context"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
)

type agentRepository interface {
	CreateAgent(ctx context.Context, agent modelsHub.Agent) (*modelsHub.Agent, error)
	GetAgentByID(ctx context.Context, id uuid.UUID) (*modelsHub.Agent, error)
	GetAllAgents(ctx context.Context) ([]modelsHub.Agent, error)
	EditDescByID(ctx context.Context, id uuid.UUID, desc string) error
	EditNameByID(ctx context.Context, id uuid.UUID, name string) error
}

type agentStateRepository interface {
	GetStatus(id uuid.UUID) bool
	SetOnline(id uuid.UUID)
}
