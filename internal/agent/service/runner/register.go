package runner

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lorsanstand/gontrol/internal/models"
)

// RegisterOrUpdate Является блокирующей функцией пока агент не обновится или не зарегистрируется.
func (a *AgentRunner) RegisterOrUpdate(ctx context.Context) (*models.Agent, error) {
	if a.system.GetAgentID() != nil {
		return a.retry(ctx, a.update, "update runner failed")
	}

	return a.retry(ctx, a.register, "register runner failed")
}

// register регистрирует агента и сохраняет AgentID.
func (a *AgentRunner) register(ctx context.Context) (*models.Agent, error) {
	systemSpecs := a.system.GetSpecs()
	commandSpecs := a.taskExec.GetCommands()

	create := models.CreateAgent{
		OS:           systemSpecs.OS,
		Hostname:     systemSpecs.Hostname,
		CommandSpecs: commandSpecs,
	}

	agent, err := a.hub.RegisterAgent(ctx, create)
	if err != nil {
		return nil, fmt.Errorf("register runner in hub: %w", err)
	}

	if err = a.system.SaveAgentID(agent.ID); err != nil {
		return nil, fmt.Errorf("save runner ID: %w", err)
	}

	return agent, nil
}

// update обновляет данные об агенте.
// На будущее когда появится обновление агента через ручку
func (a *AgentRunner) update(ctx context.Context) (*models.Agent, error) {
	specs := a.system.GetSpecs()

	return &models.Agent{ID: *a.system.GetAgentID(), OS: specs.OS, Hostname: specs.Hostname}, nil
}

// Написать обработку ошибки если запрос не <=500, что бы избежать бесконечного update.
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
