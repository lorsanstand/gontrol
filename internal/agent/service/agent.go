package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/lorsanstand/gontrol/internal/models"
)

type AgentRunner struct {
	log      *slog.Logger
	taskExec *TaskExecutor
	hub      *HubClient
}

func NewAgentRunner(log *slog.Logger, tasker *TaskExecutor, hub *HubClient) *AgentRunner {
	return &AgentRunner{log: log, taskExec: tasker, hub: hub}
}

func (a *AgentRunner) Serve() error {
	agentID := ""
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("get hostname: %w", err)
	}
	OS := runtime.GOOS
	hb := models.HeartbeatRequest{AgentID: agentID, Hostname: hostname, OS: OS}

	for {
		time.Sleep(5 * time.Second)

		a.log.Debug("send heartbeat")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		resp, err := a.hub.SendHeartbeat(ctx, hb)
		if err != nil {
			a.log.Warn("send heartbeat to hub", slog.Any("error", err))
			time.Sleep(5 * time.Second)
			continue
		}

		if hb.AgentID == "" {
			slog.Info("new agent id", slog.String("agent id", resp.AgentID))
			hb.AgentID = resp.AgentID
		}

		if resp.Task == nil {
			a.log.Debug("no tasks, continue")
			continue
		}

		result, err := a.taskExec.Execute(ctx, *resp.Task)
		if err != nil {
			a.log.Warn("execute task", slog.Any("error", err))
			continue
		}

		err = a.hub.SendTaskResult(ctx, resp.Task.ID, result)
		if err != nil {
			a.log.Warn("send result task to hub", slog.Any("error", err))
			continue
		}

		cancel()
	}
}
