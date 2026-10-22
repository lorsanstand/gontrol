package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

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

// AgentService структура для работы с агентами.
type AgentService struct {
	store agentRepository
	log   *slog.Logger
}

func NewAgentService(logger *slog.Logger, storage agentRepository) *AgentService {
	return &AgentService{store: storage, log: logger}
}

// GetAll возращает всех агентов в хранилище.
// Если агентов нет возращает пустой список.
func (a *AgentService) GetAll(ctx context.Context) ([]modelsHub.Agent, error) {
	return a.store.GetAllAgents(ctx)
}

// Get возращает агента по его ID.
// Если агента не существует, вернет ошибку ErrAgentNotExist из store.
func (a *AgentService) Get(ctx context.Context, agentID uuid.UUID) (*modelsHub.Agent, error) {
	return a.store.GetAgentByID(ctx, agentID)
}

// RegisterOrUpdate регистрирует агента если передан пустой AgentID, и обновляет данные если агент уже зарегистрирован.
func (a *AgentService) RegisterOrUpdate(ctx context.Context, heartbeat models.HeartbeatRequest) (*modelsHub.Agent, error) {
	if heartbeat.AgentID != "" {
		a.log.Debug("agent err")
		agentID, err := uuid.Parse(heartbeat.AgentID)
		if err != nil {
			return nil, fmt.Errorf("parse agent ID to uuid type: %w", err)
		}

		agentDB, err := a.store.GetAgentByID(ctx, agentID)
		if err != nil {
			return nil, fmt.Errorf("get agent: %w", err)
		}
		return agentDB, nil
	}

	newAgent := modelsHub.Agent{Hostname: heartbeat.Hostname, OS: heartbeat.OS}

	agentDB, err := a.store.CreateAgent(ctx, newAgent)
	if err != nil {
		return nil, fmt.Errorf("create agent: %w", err)
	}

	return agentDB, nil
}

// EditNameAgent изменяет имя агента.
// Возращает ErrChangeData если агент не найден.
func (a *AgentService) EditNameAgent(ctx context.Context, id uuid.UUID, newName string) error {
	return a.store.EditNameByID(ctx, id, newName)
}

// EditDescAgent изменяет описание агента
// Возращает ErrChangeData если агент не найден
func (a *AgentService) EditDescAgent(ctx context.Context, id uuid.UUID, newDesc string) error {
	return a.store.EditNameByID(ctx, id, newDesc)
}
