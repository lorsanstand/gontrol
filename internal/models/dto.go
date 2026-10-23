package models

import "github.com/google/uuid"

type Task struct {
	ID      uuid.UUID         `json:"id"`
	Command string            `json:"command"`
	Args    map[string]string `json:"args"`
}

type TaskResult struct {
	ID         uuid.UUID `json:"id"`
	TimeExecMS int       `json:"time_exec_ms"`
	Result     string    `json:"result"`
}

type HeartbeatRequest struct {
	AgentID  *uuid.UUID `json:"agent_id"`
	Hostname string     `json:"hostname"`
	OS       string     `json:"os"`
}

type HeartbeatResponse struct {
	AgentID uuid.UUID `json:"agent_id"`
	Task    *Task     `json:"task"`
}
