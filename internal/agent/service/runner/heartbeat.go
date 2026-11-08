package runner

import (
	"context"
	"fmt"

	"github.com/lorsanstand/gontrol/internal/models"
)

// Heartbeat отправляет запрос и отдает задачи на выполнение.
// В будущем нужно сделать асинхронное выполнение задач через go execute
func (a *AgentRunner) Heartbeat(ctx context.Context, hb models.HeartbeatRequest) error {
	resp, err := a.hub.SendHeartbeat(ctx, hb)
	if err != nil {
		return fmt.Errorf("send heartbeat: %w", err)
	}

	if resp.Task == nil {
		a.log.Info("no tasks, continue...")
		return nil
	}

	result, err := a.taskExec.Execute(ctx, *resp.Task)
	if err != nil {
		return fmt.Errorf("execute task: %w", err)
	}

	if err = a.hub.SendTaskResult(ctx, result.ID, *result); err != nil {
		return fmt.Errorf("send result task: %w", err)
	}

	return nil
}
