package service

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/lorsanstand/gontrol/internal/models"
)

type TaskExecutor struct {
}

func (t *TaskExecutor) Execute(ctx context.Context, task models.Task) (models.TaskResult, error) {
	start := time.Now()

	//nolint:gosec // G204: запуск команды из задачи, перепишу позже
	cmd := exec.Command(task.Command, task.Args["args"])
	output, err := cmd.CombinedOutput()
	if err != nil {
		return models.TaskResult{}, fmt.Errorf("execute command: %w", err)
	}

	d := time.Since(start)

	return models.TaskResult{
		Result:     string(output),
		TimeExecMS: int(d.Milliseconds()),
	}, nil
}
