package service

import (
	"context"
	"log/slog"
	"uuid"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/models"
)

type agentRepository interface {
	RegisterOrUpdate(agentRequest models.HeartbeatRequest)
	GetAgent(agentID string) (modelsHub.Agent, error)
	GetAllAgents() []modelsHub.Agent
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
	return a.store.GetAllAgents(), nil
}

// Get возращает агента по его ID.
// Если агента не существует, вернет ошибку ErrAgentNotExist из store.
func (a *AgentService) Get(ctx context.Context, agentID string) (modelsHub.Agent, error) {
	return a.store.GetAgent(agentID)
}

// RegisterOrUpdate регистрирует агента если передан пустой AgentID, и обновляет данные если агент уже зарегистрирован.
func (a *AgentService) RegisterOrUpdate(ctx context.Context, heartbeat models.HeartbeatRequest) (modelsHub.Agent, error) {
	if heartbeat.AgentID == "" {
		heartbeat.AgentID = uuid.New().String()
	}

	a.store.RegisterOrUpdate(heartbeat)

	return modelsHub.Agent{
		AgentID:  heartbeat.AgentID,
		Hostname: heartbeat.Hostname,
		OS:       heartbeat.OS,
	}, nil
}
