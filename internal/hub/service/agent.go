package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lorsanstand/gontrol/internal/hub/store"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/models"
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
}

// AgentService структура для работы с агентами.
type AgentService struct {
	store agentRepository
	state agentStateRepository
	log   *slog.Logger
}

func NewAgentService(logger *slog.Logger, storage agentRepository, state agentStateRepository) *AgentService {
	return &AgentService{store: storage, log: logger, state: state}
}

// GetAll возращает всех агентов в хранилище.
// Если агентов нет возращает пустой список.
func (a *AgentService) GetAll(ctx context.Context) ([]modelsHub.Agent, error) {
	agentsDB, err := a.store.GetAllAgents(ctx)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("get agent: %w", err)
	}

	// В будущем попытаться оптимизировать
	for _, agentDB := range agentsDB {
		agentDB.State = a.state.GetStatus(agentDB.ID)
	}

	return agentsDB, nil
}

// Get возращает агента по его ID.
// Если агента не существует, вернет ошибку ErrAgentNotFound.
func (a *AgentService) Get(ctx context.Context, id uuid.UUID) (*modelsHub.Agent, error) {
	agentDB, err := a.store.GetAgentByID(ctx, id)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("get agent: %w", err)
	}

	agentDB.State = a.state.GetStatus(agentDB.ID)

	return agentDB, nil
}

// RegisterOrUpdate регистрирует агента если передан пустой AgentID, и обновляет данные если агент уже зарегистрирован.
func (a *AgentService) RegisterOrUpdate(ctx context.Context, heartbeat models.HeartbeatRequest) (*modelsHub.Agent, error) {
	if heartbeat.AgentID != nil {
		a.log.Debug("agent exists", slog.Any("agentID", heartbeat.AgentID))
		return a.Get(ctx, *heartbeat.AgentID)
	}

	newAgent := modelsHub.Agent{Hostname: heartbeat.Hostname, OS: heartbeat.OS}

	agentDB, err := a.store.CreateAgent(ctx, newAgent)
	if err != nil {
		return nil, fmt.Errorf("create agent: %w", err)
	}

	return agentDB, nil
}

// EditNameAgent изменяет имя агента.
// Возращает ErrAgentNotFound если агент не найден.
func (a *AgentService) EditNameAgent(ctx context.Context, id uuid.UUID, newName string) error {
	if err := a.store.EditNameByID(ctx, id, newName); err != nil {
		if errors.Is(err, store.ErrChangeData) {
			return ErrAgentNotFound
		}
		return fmt.Errorf("edit agent name: %w", err)
	}

	return nil
}

// EditDescAgent изменяет описание агента
// Возращает ErrAgentNotFound если агент не найден
func (a *AgentService) EditDescAgent(ctx context.Context, id uuid.UUID, newDesc string) error {
	if err := a.store.EditDescByID(ctx, id, newDesc); err != nil {
		if errors.Is(err, store.ErrChangeData) {
			return ErrAgentNotFound
		}
		return fmt.Errorf("edit agent description: %w", err)
	}

	return nil
}
