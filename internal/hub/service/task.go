package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/models"
)

// TaskService структура для работы с задачами для агента.
type TaskService struct {
	log   *slog.Logger
	store tasksRepository
	specs commandSpecsRepository
	genID generatorID
}

func NewTasksService(logger *slog.Logger, storage tasksRepository, specs commandSpecsRepository) *TaskService {
	return &TaskService{store: storage, log: logger, genID: uuid.NewV7, specs: specs}
}

// Add создает задачу и назначает ее указанному агенту.
// Если нужного агента не существует, возвращается ErrAgentNotFound.
// Если задачу невозможно сохранить, возвращается ошибка.
func (t *TaskService) Add(ctx context.Context, taskCreate modelsHub.CreateTask) (*modelsHub.Task, error) {
	id, err := t.genID()
	if err != nil {
		return nil, fmt.Errorf("generate id: %w", err)
	}

	taskSpecs := models.Task{ID: id, Command: taskCreate.Command, Args: taskCreate.Args}

	specs, err := t.specs.GetCommandSpecs(ctx, taskCreate.AgentID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrAgentNotFound
		}

		return nil, fmt.Errorf("get command specification: %w", err)
	}

	if err = specs.ValidateTask(taskSpecs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	task := modelsHub.Task{
		ID:      id,
		AgentID: taskCreate.AgentID,
		Status:  modelsHub.TaskStatusPending,
		Command: taskCreate.Command,
		Args:    taskCreate.Args,
		ExecMs:  nil,
		Result:  nil,
	}

	taskDB, err := t.store.CreateTask(ctx, task)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrAgentNotFound
		}

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

// GetByID отдает задачу по id.
// Если задача не найдена возращается ErrTaskNotFound.
func (t *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*modelsHub.Task, error) {
	taskDB, err := t.store.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("get task: %w", err)
	}

	return taskDB, nil
}

// GetByAgentID отдает задачи по id агента.
// Если задача не найдена возращается пустой слайс.
func (t *TaskService) GetByAgentID(ctx context.Context, id uuid.UUID) ([]modelsHub.Task, error) {
	taskDB, err := t.store.GetTasksByAgentID(ctx, id)
	if err != nil {
		return []modelsHub.Task{}, fmt.Errorf("get tasks: %w", err)
	}

	return taskDB, nil
}
