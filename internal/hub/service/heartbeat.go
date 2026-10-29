package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/lorsanstand/gontrol/internal/models"
)

// HeartbeatService структура для обработки heartbeat
type HeartbeatService struct {
	agent *AgentService
	task  *TaskService
}

func NewHeartbeatService(agent *AgentService, task *TaskService) *HeartbeatService {
	return &HeartbeatService{agent: agent, task: task}
}

// Handle обновляет статус агента и возращает ID и задачу для выполнения.
// Если не удалось проверить существует ли агент возращается ErrAgentNotFound.
// Если задача нет то возращается nil в поле Tasks
func (h *HeartbeatService) Handle(ctx context.Context, heartbeatReq models.HeartbeatRequest) (*models.HeartbeatResponse, error) {
	agent, err := h.agent.Get(ctx, heartbeatReq.AgentID)
	if err != nil {
		return nil, err
	}

	h.agent.SetOnlineAgent(agent.ID)

	task, err := h.task.Pop(ctx, agent.ID)
	if err != nil && !errors.Is(err, ErrTaskNotFound) {
		return nil, fmt.Errorf("pop task: %w", err)
	}

	var taskAgent *models.Task = nil
	if task != nil {
		taskAgent = &models.Task{ID: task.ID, Command: task.Command, Args: task.Args}
	}

	return &models.HeartbeatResponse{AgentID: agent.ID, Task: taskAgent}, nil
}
