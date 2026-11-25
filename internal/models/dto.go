package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RegisterAgentRequest struct {
	OS           string       `json:"OS" validate:"required,max=10"`
	Hostname     string       `json:"hostname" validate:"required,max=30,hostname"`
	CommandSpecs CommandSpecs `json:"command_specs" validate:"required"`
}

type RegisterAgentResponse struct {
	ID           uuid.UUID     `json:"id"`
	HeartbeatTTL time.Duration `json:"heartbeat_ttl"`
}

type Task struct {
	ID      uuid.UUID                  `json:"id"`
	Command string                     `json:"command"`
	Args    map[string]json.RawMessage `json:"args"`
}

type TaskResult struct {
	ID         uuid.UUID `json:"id" validate:"required,uuid7"`
	TimeExecMS int       `json:"time_exec_ms" validate:"required,gt=0"`
	Result     string    `json:"result"`
}

type HeartbeatRequest struct {
	AgentID uuid.UUID `json:"agent_id" validate:"required,uuid7"`
}

type HeartbeatResponse struct {
	AgentID uuid.UUID `json:"agent_id"`
	Task    *Task     `json:"task"`
}
