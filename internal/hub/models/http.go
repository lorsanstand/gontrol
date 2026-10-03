package models

type CreateTask struct {
	AgentID string            `json:"agent_id"`
	Command string            `json:"command"`
	Args    map[string]string `json:"args"`
}
