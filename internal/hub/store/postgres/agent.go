package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	gen "github.com/lorsanstand/gontrol/internal/hub/store/postgres/gen"
)

type AgentRepository struct {
	db *sql.DB
	q  gen.Querier
}

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db: db, q: gen.New(db)}
}

// CreateAgent создает агента в БД. передается агент без ID и даты создания возращается готовый агент из БД.
func (a *AgentRepository) CreateAgent(ctx context.Context, agent modelsHub.Agent) (*modelsHub.Agent, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate id: %w", err)
	}

	agentCreate := gen.CreateAgentParams{
		ID:          id,
		Name:        agent.Name,
		Os:          agent.OS,
		Hostname:    agent.Hostname,
		Description: agent.Description,
	}

	agentDB, err := a.q.CreateAgent(ctx, agentCreate)
	if err != nil {
		return nil, fmt.Errorf("insert agent into db: %w", err)
	}

	return &modelsHub.Agent{
		ID:        agentDB.ID,
		Name:      agentDB.Name,
		OS:        agentDB.Os,
		Hostname:  agentDB.Hostname,
		CreatedAt: agentDB.CreatedAt,
	}, nil
}

// GetAgentByID возращает агента по ID.
func (a *AgentRepository) GetAgentByID(ctx context.Context, id uuid.UUID) (*modelsHub.Agent, error) {
	agentDB, err := a.q.GetAgentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("select agent into db: %w", err)
	}

	return &modelsHub.Agent{
		ID:        agentDB.ID,
		Name:      agentDB.Name,
		OS:        agentDB.Os,
		Hostname:  agentDB.Hostname,
		CreatedAt: agentDB.CreatedAt,
	}, nil
}

// GetAllAgents возращает всех агентов
func (a *AgentRepository) GetAllAgents(ctx context.Context) ([]modelsHub.Agent, error) {
	agentsDB, err := a.q.GetAllAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("select agents into db: %w", err)
	}

	agents := make([]modelsHub.Agent, 0, len(agentsDB))
	for _, agentDB := range agentsDB {
		agent := modelsHub.Agent{
			ID:        agentDB.ID,
			Name:      agentDB.Name,
			OS:        agentDB.Os,
			Hostname:  agentDB.Hostname,
			CreatedAt: agentDB.CreatedAt,
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

// EditDescByID редактирует описание агента
func (a *AgentRepository) EditDescByID(ctx context.Context, id uuid.UUID, desc string) error {
	query := gen.EditDescAgentByIDParams{ID: id, Description: &desc}

	result, err := a.q.EditDescAgentByID(ctx, query)
	if err != nil {
		return fmt.Errorf("update agents into db: %w", err)
	}

	return IFChangesInDB(result)
}

// EditNameByID редактирует описание агента
func (a *AgentRepository) EditNameByID(ctx context.Context, id uuid.UUID, name string) error {
	query := gen.EditNameAgentByIDParams{ID: id, Description: &name}

	result, err := a.q.EditNameAgentByID(ctx, query)
	if err != nil {
		return fmt.Errorf("update agents into db: %w", err)
	}

	return IFChangesInDB(result)
}
