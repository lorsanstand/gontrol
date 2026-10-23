package http

import (
	"log/slog"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handler"
	"github.com/lorsanstand/gontrol/internal/hub/service"
)

func RegisterRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	task *service.TaskService,
	agent *service.AgentService,
	heartbeat *service.HeartbeatService,
) {
	heartbeatHandler := handler.NewHeartbeatHandler(logger, heartbeat)
	mux.Handle("POST /api/v1/heartbeat", heartbeatHandler)

	taskHandler := handler.NewTaskHandler(logger, task)
	mux.HandleFunc("POST /api/v1/tasks", taskHandler.PostCreateTask)
	mux.HandleFunc("POST /api/v1/tasks/{id}/result", taskHandler.PostResultTask)

	agentHandler := handler.NewAgentHandler(logger, agent)
	mux.HandleFunc("GET /api/v1/agents/{id}", agentHandler.GetAgent)
	mux.HandleFunc("GET /api/v1/agents/", agentHandler.GetAllAgents)
	mux.HandleFunc("PATCH /api/v1/agents/{id}/name", agentHandler.EditName)
	mux.HandleFunc("PATCH /api/v1/agents/{id}/description", agentHandler.EditName)
}
