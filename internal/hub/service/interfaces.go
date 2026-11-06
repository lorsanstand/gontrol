package service

import (
	"context"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
)

//go:generate mockgen -destination=mocks/mock_agent_repo.go -package=mocks . agentRepository
type agentRepository interface {
	CreateAgent(ctx context.Context, agent modelsHub.Agent) (*modelsHub.Agent, error)
	GetAgentByID(ctx context.Context, id uuid.UUID) (*modelsHub.Agent, error)
	GetAllAgents(ctx context.Context) ([]modelsHub.Agent, error)
	EditDescByID(ctx context.Context, id uuid.UUID, desc string) error
	EditNameByID(ctx context.Context, id uuid.UUID, name string) error
}

//go:generate mockgen -destination=mocks/mock_agent_state.go -package=mocks . agentStateRepository
type agentStateRepository interface {
	GetStatus(id uuid.UUID) bool
	SetOnline(id uuid.UUID)
}

//go:generate mockgen -destination=mocks/mock_task_repo.go -package=mocks . tasksRepository
type tasksRepository interface {
	CreateTask(ctx context.Context, task modelsHub.Task) (*modelsHub.Task, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (*modelsHub.Task, error)
	GetTasksByAgentID(ctx context.Context, agentID uuid.UUID) ([]modelsHub.Task, error)
	NextTaskByAgentID(ctx context.Context, agentID uuid.UUID) (*modelsHub.Task, error)
	EditStatusByID(ctx context.Context, id uuid.UUID, status modelsHub.TaskStatus) error
	EditResultByID(ctx context.Context, id uuid.UUID, resultStr string, status modelsHub.TaskStatus, execTimeMS int) error
}

type taskService interface {
	Pop(ctx context.Context, agentID uuid.UUID) (*modelsHub.Task, error)
}

type agentService interface {
	Get(ctx context.Context, id uuid.UUID) (*modelsHub.Agent, error)
	SetOnlineAgent(id uuid.UUID)
}

type generatorID func() (uuid.UUID, error)
