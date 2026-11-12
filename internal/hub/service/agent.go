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

// AgentService структура для работы с агентами.
type AgentService struct {
	store agentRepository
	state agentStateRepository
	specs commandSpecsRepository
	genID generatorID
	log   *slog.Logger
}

func NewAgentService(
	logger *slog.Logger,
	storage agentRepository,
	state agentStateRepository,
	specs commandSpecsRepository,
) *AgentService {
	return &AgentService{store: storage, log: logger, state: state, specs: specs, genID: uuid.NewV7}
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
	for idx, agentDB := range agentsDB {
		agentDB.State = a.state.GetStatus(agentDB.ID)
		agentsDB[idx] = agentDB
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

// Register сохраняет агента и спецификацию ее команд.
func (a *AgentService) Register(ctx context.Context, agentReq models.CreateAgent) (*modelsHub.Agent, error) {
	id, err := a.genID()
	if err != nil {
		return nil, fmt.Errorf("generate id: %w", err)
	}
	newAgent := modelsHub.Agent{ID: id, Hostname: agentReq.Hostname, OS: agentReq.OS}

	agentDB, err := a.store.CreateAgent(ctx, newAgent)
	if err != nil {
		return nil, fmt.Errorf("create agent: %w", err)
	}

	if err = a.specs.SetCommandSpecs(ctx, agentDB.ID, agentReq.CommandSpecs); err != nil {
		return nil, fmt.Errorf("save command specification: %w", err)
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

// EditDescAgent изменяет описание агента.
// Возращает ErrAgentNotFound если агент не найден.
func (a *AgentService) EditDescAgent(ctx context.Context, id uuid.UUID, newDesc string) error {
	if err := a.store.EditDescByID(ctx, id, newDesc); err != nil {
		if errors.Is(err, store.ErrChangeData) {
			return ErrAgentNotFound
		}
		return fmt.Errorf("edit agent description: %w", err)
	}

	return nil
}

// SetOnlineAgent меняет статус агента с offline на online.
// Не проверяет, существует ли агент.
func (a *AgentService) SetOnlineAgent(id uuid.UUID) {
	a.state.SetOnline(id)
}

// GetCommandSpecs возвращает спецификацию команд агента.
// Если спецификации не найдено возращает ErrAgentNotFound
func (a *AgentService) GetCommandSpecs(ctx context.Context, agentID uuid.UUID) (models.CommandSpecs, error) {
	specs, err := a.specs.GetCommandSpecs(ctx, agentID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrAgentNotFound
		}

		return nil, fmt.Errorf("get command specification: %w", err)
	}

	return specs, nil
}
