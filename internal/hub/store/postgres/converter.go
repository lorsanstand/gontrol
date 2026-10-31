package postgres

import (
	"encoding/json"
	"fmt"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	gen "github.com/lorsanstand/gontrol/internal/hub/store/postgres/gen"
)

func toTask(taskDB gen.Task) (*modelsHub.Task, error) {
	argsDB := new(map[string]string)
	status, err := modelsHub.TaskStatusParse(taskDB.Status)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(taskDB.Args, argsDB)
	if err != nil {
		return nil, fmt.Errorf("unmarshal args: %w", err)
	}

	return &modelsHub.Task{
		ID:        taskDB.ID,
		AgentID:   taskDB.AgentID,
		Status:    status,
		Command:   taskDB.Command,
		Args:      *argsDB,
		ExecMs:    taskDB.ExecMs,
		Result:    taskDB.Result,
		CreatedAt: taskDB.CreatedAt,
	}, nil
}

func toAgent(agentDB gen.Agent) (*modelsHub.Agent, error) {
	return &modelsHub.Agent{
		ID:          agentDB.ID,
		Name:        agentDB.Name,
		OS:          agentDB.Os,
		Hostname:    agentDB.Hostname,
		Description: agentDB.Description,
		CreatedAt:   agentDB.CreatedAt,
	}, nil
}
