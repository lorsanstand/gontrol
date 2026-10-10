package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"uuid"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/models"
)

type tasksRepository interface {
	PopTask(agentID string) (models.Task, error)
	AddTask(agentID string, task models.Task) error
}

type TaskService struct {
	log   *slog.Logger
	store tasksRepository
}

func NewTasksService(logger *slog.Logger, storage tasksRepository) *TaskService {
	return &TaskService{store: storage, log: logger}
}

func (t *TaskService) Add(ctx context.Context, agentID string, taskRequest modelsHub.CreateTask) (models.Task, error) {
	taskID := uuid.New().String()
	task := models.Task{ID: taskID, Command: taskRequest.Command, Args: taskRequest.Args}

	err := t.store.AddTask(agentID, task)
	if err != nil {
		return models.Task{}, fmt.Errorf("add task in storage: %w", err)
	}

	return task, nil
}

func (t *TaskService) Pop(ctx context.Context, agentID string) (*models.Task, error) {
	task, err := t.store.PopTask(agentID)
	if err != nil && !errors.Is(err, store.ErrTasksNotExist) {
		return nil, fmt.Errorf("pop task in storage: %w", err)
	}

	if errors.Is(err, store.ErrTasksNotExist) {
		return nil, nil
	}

	return &task, nil
}
