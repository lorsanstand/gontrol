package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/store"
)

type tasksRepository interface {
	CreateTask(ctx context.Context, task modelsHub.Task) (*modelsHub.Task, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (*modelsHub.Task, error)
	GetTasksByAgentID(ctx context.Context, agentID uuid.UUID) ([]modelsHub.Task, error)
	NextTaskByAgentID(ctx context.Context, agentID uuid.UUID) (*modelsHub.Task, error)
	EditStatusByID(ctx context.Context, id uuid.UUID, status modelsHub.TaskStatus) error
	EditResultByID(ctx context.Context, id uuid.UUID, resultStr string, status modelsHub.TaskStatus, execTimeMS int) error
}

// TaskService структура для работы с задачами для агента.
type TaskService struct {
	log   *slog.Logger
	store tasksRepository
}

func NewTasksService(logger *slog.Logger, storage tasksRepository) *TaskService {
	return &TaskService{store: storage, log: logger}
}

// Метод Add создает задачу и назначает ее указанному агенту.
// Если задачу невозможно сохранить или агент не существует, возвращается ошибка.
func (t *TaskService) Add(ctx context.Context, taskCreate modelsHub.CreateTask) (*modelsHub.Task, error) {

	task := modelsHub.Task{
		AgentID: taskCreate.AgentID,
		Status:  modelsHub.TaskStatusPending,
		Command: taskCreate.Command,
		Args:    taskCreate.Args,
		ExecMs:  nil,
		Result:  nil,
	}

	taskDB, err := t.store.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return taskDB, nil
}

// Pop меняет статус на processing и возвращает следующую ожидающую выполнения задачу для указанного агента.
// Возвращает nil, если у агента нет ожидающих выполнения задач.
// Возращает ErrTaskNotFound если задач нет.
func (t *TaskService) Pop(ctx context.Context, agentID uuid.UUID) (*modelsHub.Task, error) {
	taskDB, err := t.store.NextTaskByAgentID(ctx, agentID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("get next task: %w", err)
	}

	err = t.store.EditStatusByID(ctx, taskDB.ID, modelsHub.TaskStatusProcessing)
	if err != nil {
		return nil, fmt.Errorf("edit status task: %w", err)
	}

	return taskDB, nil
}

// SaveResult сохраняет результат задачи, время выполнения и меняет статус на completed.
// Возращает ErrTaskNotFound если задача не существует.
func (t *TaskService) SaveResult(ctx context.Context, id uuid.UUID, result string, execMS int) error {
	err := t.store.EditResultByID(ctx, id, result, modelsHub.TaskStatusCompleted, execMS)
	if err != nil {
		if errors.Is(err, store.ErrChangeData) {
			return ErrTaskNotFound
		}

		return fmt.Errorf("save result task: %w", err)
	}

	return nil
}
