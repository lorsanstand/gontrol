package http

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handler"
	"github.com/lorsanstand/gontrol/internal/hub/service"
)

func RegisterRoutes(
	logger *slog.Logger,
	task *service.TaskService,
	agent *service.AgentService,
	heartbeat *service.HeartbeatService,
) *chi.Mux {
	mux := chi.NewRouter()

	heartbeatHandler := handler.NewHeartbeatHandler(logger, heartbeat)
	mux.Post("/heartbeat", heartbeatHandler.ServeHTTP)

	taskRouter := chi.NewRouter()
	taskHandler := handler.NewTaskHandler(logger, task)
	taskRouter.Post("/", taskHandler.PostCreateTask)
	taskRouter.Post("/{id}/result", taskHandler.PostResultTask)
	taskRouter.Get("/{id}", taskHandler.GetTask)
	taskRouter.Get("/agent/{id}", taskHandler.GetTaskByAgentID)

	agentRouter := chi.NewRouter()
	agentHandler := handler.NewAgentHandler(logger, agent)
	agentRouter.Get("/{id}", agentHandler.GetAgent)
	agentRouter.Get("/", agentHandler.GetAllAgents)
	agentRouter.Patch("/{id}/name", agentHandler.EditName)
	agentRouter.Patch("/{id}/description", agentHandler.EditDesc)
	agentRouter.Post("/register", agentHandler.RegisterAgent)

	mux.Mount("/tasks", taskRouter)
	mux.Mount("/agents", agentRouter)

	return mux
}
