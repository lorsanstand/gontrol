package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lorsanstand/gontrol/internal/models"
)

type AgentRunner struct {
	log      *slog.Logger
	taskExec *TaskExecutor
	hub      *HubClient
	system   *System
}

func NewAgentRunner(log *slog.Logger, tasker *TaskExecutor, hub *HubClient, system *System) *AgentRunner {
	return &AgentRunner{log: log, taskExec: tasker, hub: hub, system: system}
}

// Serve блокирующая функция которая отсылает heartbeat и отправляет задачи на выполнение.
func (a *AgentRunner) Serve(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	agent, err := a.RegisterOrUpdate(ctx)
	if err != nil {
		return fmt.Errorf("register or update agent: %w", err)
	}

	hb := models.HeartbeatRequest{AgentID: agent.ID, OS: agent.OS, Hostname: agent.Hostname}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err = a.Heartbeat(ctx, hb)
			slog.Error("failed send heartbeat", slog.Any("error", err))
		}
	}
}

// Heartbeat отправялет запрос и отдает задачи на выполнение.
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

	if err = a.hub.SendTaskResult(ctx, result.ID, result); err != nil {
		return fmt.Errorf("send result task: %w", err)
	}

	return nil
}

// RegisterOrUpdate Является блокирующей функцией пока агент не обновится или не зарегитсрируется.
func (a *AgentRunner) RegisterOrUpdate(ctx context.Context) (*models.Agent, error) {
	if a.system.GetAgentID() != nil {
		return a.retry(ctx, a.update, "update agent failed")
	}

	return a.retry(ctx, a.register, "register agent failed")
}

// register регистрирует агента и сохраняет AgentID.
func (a *AgentRunner) register(ctx context.Context) (*models.Agent, error) {
	specs := a.system.GetSpecs()

	agent, err := a.hub.RegisterAgent(ctx, specs)
	if err != nil {
		return nil, fmt.Errorf("register agent in hub: %w", err)
	}

	if err = a.system.SaveAgentID(agent.ID); err != nil {
		return nil, fmt.Errorf("save agent ID: %w", err)
	}

	return agent, nil
}

// update обновляет данные об агенте.
// На будущее когда появится обновление агента через ручку
func (a *AgentRunner) update(ctx context.Context) (*models.Agent, error) {
	specs := a.system.GetSpecs()

	return &models.Agent{ID: *a.system.GetAgentID(), OS: specs.OS, Hostname: specs.Hostname}, nil
}

// Написать обработку ошибки если запрос не <=500 что бы избежать бесконечного update.
func (a *AgentRunner) retry(
	ctx context.Context,
	f func(ctx context.Context) (*models.Agent, error),
	msg string) (*models.Agent, error) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	agent, err := f(ctx)
	if err == nil {
		return agent, nil
	}

	a.log.Warn(msg, slog.Any("error", err))
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-ticker.C:
			agent, err = f(ctx)
			if err == nil {
				return agent, nil
			}
			a.log.Warn("retry failed", "error", err)
		}
	}
}
