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

// Handle обновляет или регистрирует агента, и отдает ему задачу.
// Если агент не удалось обновить или выдать задачу, возвращается ошибка
func (h *HeartbeatService) Handle(ctx context.Context, heartbeatReq models.HeartbeatRequest) (models.HeartbeatResponse, error) {
	agent, err := h.agent.RegisterOrUpdate(ctx, heartbeatReq)
	if err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("register or update agent: %w", err)
	}

	h.agent.SetOnlineAgent(agent.ID)

	task, err := h.task.Pop(ctx, agent.ID)
	if err != nil && !errors.Is(err, ErrTaskNotFound) {
		return models.HeartbeatResponse{}, fmt.Errorf("pop task: %w", err)
	}
	var taskAgent *models.Task = nil

	if task != nil {
		taskAgent = &models.Task{ID: task.ID, Command: task.Command, Args: task.Args}
	}

	return models.HeartbeatResponse{AgentID: agent.ID, Task: taskAgent}, nil
}
