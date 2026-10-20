package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	gen "github.com/lorsanstand/gontrol/internal/hub/store/postgres/gen"
)

type TaskRepository struct {
	db *sql.DB
	q  gen.Querier
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db, q: gen.New(db)}
}

// CreateTask создает новую задачу для определенного агента.
func (t *TaskRepository) CreateTask(ctx context.Context, task modelsHub.Task) (*modelsHub.Task, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate id: %w", err)
	}

	args, err := json.Marshal(task.Args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	taskCreate := gen.CreateTaskParams{
		ID:      id,
		AgentID: task.AgentID,
		Status:  string(task.Status),
		Command: task.Command,
		Args:    args,
		ExecMs:  task.ExecMs,
		Result:  task.Result,
	}

	taskDB, err := t.q.CreateTask(ctx, taskCreate)
	if err != nil {
		return nil, fmt.Errorf("insert task into db: %w", err)
	}

	return toTask(taskDB)
}

// GetTaskByID возращает задачу по id.
func (t *TaskRepository) GetTaskByID(ctx context.Context, id uuid.UUID) (*modelsHub.Task, error) {
	taskDB, err := t.q.GetTaskByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("select agent into DB: %w", err)
	}

	return toTask(taskDB)
}

// GetTasksByAgentID возращает задачи по id агента.
func (t *TaskRepository) GetTasksByAgentID(ctx context.Context, agentID uuid.UUID) ([]modelsHub.Task, error) {
	tasksDB, err := t.q.GetTasksByAgentID(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("select agent into DB: %w", err)
	}

	tasks := make([]modelsHub.Task, 0, len(tasksDB))

	for _, taskDB := range tasksDB {
		task, err := toTask(taskDB)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	return tasks, nil
}

// NextTaskByAgentID возращает следующую задачу в статусе pending по id агента.
func (t *TaskRepository) NextTaskByAgentID(ctx context.Context, agentID uuid.UUID) (*modelsHub.Task, error) {
	taskDB, err := t.q.NextTaskByAgentID(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("select task into db: %w", err)
	}

	return toTask(taskDB)
}

// EditStatusByID изменяет статус по id задачи.
// Возращает ошибку ErrChangeData если ни одна строка не изменилась
func (t *TaskRepository) EditStatusByID(ctx context.Context, id uuid.UUID, status modelsHub.TaskStatus) error {
	taskReq := gen.EditStatusTaskByIDParams{ID: id, Status: string(status)}

	result, err := t.q.EditStatusTaskByID(ctx, taskReq)
	if err != nil {
		return fmt.Errorf("update task into DB: %w", err)
	}

	return ifChangesInDB(result)
}

// EditResultByID изменяет результат статус и время выполнения задачи по id задачи.
// Возращает ошибку ErrChangeData если ни одна строка не изменилась
func (t *TaskRepository) EditResultByID(ctx context.Context,
	id uuid.UUID,
	resultStr string,
	status modelsHub.TaskStatus,
	execTimeMS int) error {

	taskResultDB := gen.EditResultTaskByIDParams {
		ID: id,
		Result: &resultStr,
		Status: string(status),
		ExecMs: &execTimeMS,
	}

	result, err := t.q.EditResultTaskByID(ctx, taskResultDB)
	if err != nil {
		return fmt.Errorf("update task into DB: %w", err)
	}

	return ifChangesInDB(result)
}

func (t *TaskRepository) DeleteTaskByID(ctx context.Context, id uuid.UUID) error {
	result, err := t.q.DeleteTaskByID(ctx, id)
	if err != nil {
		return fmt.Errorf("delete task into DB: %w", err)
	}

	return ifChangesInDB(result)
}