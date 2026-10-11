package service

import (
	"context"
	"fmt"

	"github.com/lorsanstand/gontrol/internal/models"
)

type HeartbeatService struct {
	agent *AgentService
	task  *TaskService
}

func NewHeartbeatService(agent *AgentService, task *TaskService) *HeartbeatService {
	return &HeartbeatService{agent: agent, task: task}
}

func (h *HeartbeatService) Handle(ctx context.Context, heartbeatReq models.HeartbeatRequest) (models.HeartbeatResponse, error) {
	agent, err := h.agent.RegisterOrUpdate(ctx, heartbeatReq)
	if err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("register or update agent: %w", err)
	}

	task, err := h.task.Pop(ctx, heartbeatReq.AgentID)
	if err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("pop task: %w", err)
	}

	return models.HeartbeatResponse{AgentID: agent.AgentID, Task: task}, nil
}
