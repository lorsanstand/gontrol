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

	argsDB := new(map[string]string)

	err = json.Unmarshal(taskDB.Args, argsDB)
	if err != nil {
		return nil, fmt.Errorf("unmarshal args: %w", err)
	}

	return &modelsHub.Task{
		ID:        taskDB.ID,
		AgentID:   taskDB.AgentID,
		Status:    modelsHub.TaskStatus(taskDB.Status),
		Command:   taskDB.Command,
		Args:      *argsDB,
		ExecMs:    taskDB.ExecMs,
		Result:    taskDB.Result,
		CreatedAt: taskDB.CreatedAt,
	}, nil
}
