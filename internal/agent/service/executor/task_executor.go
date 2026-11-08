package executor

import (
	"context"
	"log/slog"
	"time"

	"github.com/lorsanstand/gontrol/internal/models"
)

type TaskExecutor struct {
	log *slog.Logger
}

// New создает TaskExecutor для выполнения задач с хаба.
func New(log *slog.Logger) *TaskExecutor {
	return &TaskExecutor{log: log}
}

// Execute валидирует, выполняет задачу, и выдает результат
func (t *TaskExecutor) Execute(ctx context.Context, task models.Task) (*models.TaskResult, error) {
	start := time.Now()

	command, err := commands.GetCommandAndValidate(task)
	if err != nil {
		return nil, err
	}

	// Исправить после изменения TaskResult что бы ошибка отдавалась в хаб
	data, err := command.Exec(ctx, task.Args)
	if err != nil {
		return nil, err
	}

	d := time.Since(start)

	return &models.TaskResult{
		ID:         task.ID,
		Result:     string(data),
		TimeExecMS: int(d.Milliseconds()),
	}, nil
}
