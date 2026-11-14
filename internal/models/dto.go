package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateAgent struct {
	OS           string       `json:"OS"`
	Hostname     string       `json:"hostname"`
	CommandSpecs CommandSpecs `json:"command_specs"`
}

type Agent struct {
	ID       uuid.UUID `json:"id"`
	OS       string    `json:"OS"`
	Hostname string    `json:"hostname"`
}

type Task struct {
	ID      uuid.UUID                  `json:"id"`
	Command string                     `json:"command"`
	Args    map[string]json.RawMessage `json:"args"`
}

type TaskResult struct {
	ID         uuid.UUID `json:"id"`
	TimeExecMS int       `json:"time_exec_ms"`
	Result     string    `json:"result"`
}

type HeartbeatRequest struct {
	AgentID  uuid.UUID `json:"agent_id"`
	Hostname string    `json:"hostname"`
	OS       string    `json:"os"`
}

type HeartbeatResponse struct {
	AgentID uuid.UUID `json:"agent_id"`
	Task    *Task     `json:"task"`
}
