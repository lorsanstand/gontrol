package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCanceled   TaskStatus = "canceled"
	TaskStatusError      TaskStatus = "error"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusCanceled, TaskStatusCompleted, TaskStatusError, TaskStatusPending, TaskStatusProcessing:
		return true
	default:
		return false
	}
}

type Agent struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	OS          string    `json:"os"`
	Hostname    string    `json:"hostname"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Task struct {
	ID        uuid.UUID         `json:"id"`
	AgentID   uuid.UUID         `json:"agent_id"`
	Status    TaskStatus        `json:"status"`
	Command   string            `json:"command"`
	Args      map[string]string `json:"args"`
	ExecMs    *int              `json:"exec_ms"`
	Result    *string           `json:"result"`
	CreatedAt time.Time         `json:"created_at"`
}
