package models

import (
	"fmt"
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

func TaskStatusParse(str string) (TaskStatus, error) {
	switch str {
	case "pending":
		return TaskStatusPending, nil
	case "processing":
		return TaskStatusProcessing, nil
	case "completed":
		return TaskStatusCompleted, nil
	case "canceled":
		return TaskStatusCanceled, nil
	case "error":
		return TaskStatusError, nil
	default:
		return "", fmt.Errorf("failed convert to status")
	}
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

type CreateTask struct {
	AgentID uuid.UUID         `json:"agent_id"`
	Command string            `json:"command"`
	Args    map[string]string `json:"args"`
}
