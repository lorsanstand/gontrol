package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lorsanstand/gontrol/internal/hub/store"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	gen "github.com/lorsanstand/gontrol/internal/hub/store/postgres/gen"
)

//go:generate mockgen -destination=mocks/mock_querier.go -package=mocks  ./gen Querier

type AgentRepository struct {
	db *sql.DB
	q  gen.Querier
}

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db: db, q: gen.New(db)}
}

// CreateAgent создает агента в БД. передается агент без даты создания возвращается готовый агент из БД.
func (a *AgentRepository) CreateAgent(ctx context.Context, agent modelsHub.Agent) (*modelsHub.Agent, error) {

	agentCreate := gen.CreateAgentParams{
		ID:          agent.ID,
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
// Если агента нет возращает store.ErrNotFound
func (a *AgentRepository) GetAgentByID(ctx context.Context, id uuid.UUID) (*modelsHub.Agent, error) {
	agentDB, err := a.q.GetAgentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, fmt.Errorf("select agent into db: %w", err)
	}

	return &modelsHub.Agent{
		ID:          agentDB.ID,
		Name:        agentDB.Name,
		OS:          agentDB.Os,
		Hostname:    agentDB.Hostname,
		Description: agentDB.Description,
		CreatedAt:   agentDB.CreatedAt,
	}, nil
}

// GetAllAgents возращает всех агентов.
// Если агентов нет возращает пустой слайс.
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

// EditDescByID редактирует описание агента.
// Возращает ошибку store.ErrChangeData если ни одна строка не изменилась
func (a *AgentRepository) EditDescByID(ctx context.Context, id uuid.UUID, desc string) error {
	query := gen.EditDescAgentByIDParams{ID: id, Description: &desc}

	result, err := a.q.EditDescAgentByID(ctx, query)
	if err != nil {
		return fmt.Errorf("update agents into db: %w", err)
	}

	return ifChangesInDB(result)
}

// EditNameByID редактирует описание агента
// Возращает ошибку store.ErrChangeData если ни одна строка не изменилась
func (a *AgentRepository) EditNameByID(ctx context.Context, id uuid.UUID, name string) error {
	query := gen.EditNameAgentByIDParams{ID: id, Name: name}

	result, err := a.q.EditNameAgentByID(ctx, query)
	if err != nil {
		return fmt.Errorf("update agent into db: %w", err)
	}

	return ifChangesInDB(result)
}
