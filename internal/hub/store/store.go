package store

import (
	"errors"
	"slices"
	"sync"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/models"
)

var ErrAgentNotExist = errors.New("agent not found")
var ErrTasksNotExist = errors.New("task not exist")

type AgentStore struct {
	m     sync.RWMutex
	store map[string]modelsHub.Agent
}

func NewAgentStore() *AgentStore {
	return &AgentStore{store: make(map[string]modelsHub.Agent)}
}

func (a *AgentStore) RegisterOrUpdate(agentRequest models.HeartbeatRequest) {
	a.m.Lock()
	defer a.m.Unlock()

	agentDB, ok := a.store[agentRequest.AgentID]
	if ok {
		a.store[agentRequest.AgentID] = modelsHub.Agent{
			AgentID:  agentRequest.AgentID,
			OS:       agentRequest.OS,
			Hostname: agentRequest.Hostname,
			Tasks:    make([]models.Task, 0),
		}
		return
	}

	a.store[agentRequest.AgentID] = modelsHub.Agent{
		AgentID:  agentRequest.AgentID,
		OS:       agentRequest.OS,
		Hostname: agentRequest.Hostname,
		Tasks:    agentDB.Tasks,
	}
}

func (a *AgentStore) AddTask(agentID string, task models.Task) error {
	a.m.Lock()
	defer a.m.Unlock()

	agentDB, ok := a.store[agentID]
	if !ok {
		return ErrAgentNotExist
	}

	a.store[agentID] = modelsHub.Agent{
		AgentID:  agentDB.AgentID,
		OS:       agentDB.OS,
		Hostname: agentDB.Hostname,
		Tasks:    append(agentDB.Tasks, task),
	}
	return nil
}

func (a *AgentStore) PopTask(agentID string) (models.Task, error) {
	a.m.Lock()
	defer a.m.Unlock()

	agentDB, ok := a.store[agentID]
	if !ok {
		return models.Task{}, ErrAgentNotExist
	}

	if len(agentDB.Tasks) == 0 {
		return models.Task{}, ErrTasksNotExist
	}

	a.store[agentID] = modelsHub.Agent{
		AgentID:  agentDB.AgentID,
		OS:       agentDB.OS,
		Hostname: agentDB.Hostname,
		Tasks:    agentDB.Tasks[1:],
	}
	return agentDB.Tasks[0], nil
}

func (a *AgentStore) GetAllAgents() []modelsHub.Agent {
	a.m.RLock()
	defer a.m.RUnlock()

	allAgents := make([]modelsHub.Agent, 0, len(a.store))

	for _, agent := range a.store {
		allAgents = append(allAgents, agent)
	}

	slices.SortFunc(allAgents, func(a, b modelsHub.Agent) int {
		if a.AgentID < b.AgentID {
			return -1
		}
		if a.AgentID > b.AgentID {
			return 1
		}
		return 0
	})

	return allAgents
}

func (a *AgentStore) GetAgent(agentID string) (modelsHub.Agent, error) {
	a.m.RLock()
	defer a.m.RUnlock()

	agentDB, ok := a.store[agentID]
	if !ok {
		return modelsHub.Agent{}, ErrAgentNotExist
	}

	return agentDB, nil
}
