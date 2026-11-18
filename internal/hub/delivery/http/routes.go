package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handler"
)

func RegisterRoutes(
	taskHandler *handler.TaskHandler,
	agentHandler *handler.AgentHandler,
	heartbeatHandler *handler.HeartbeatHandler,
) *chi.Mux {
	mux := chi.NewRouter()

	mux.Post("/heartbeat", heartbeatHandler.ServeHTTP)

	taskRouter := chi.NewRouter()
	taskRouter.Post("/", taskHandler.PostCreateTask)
	taskRouter.Post("/result", taskHandler.PostResultTask)
	taskRouter.Get("/{id}", taskHandler.GetTask)
	taskRouter.Get("/agent/{id}", taskHandler.GetTaskByAgentID)

	agentRouter := chi.NewRouter()
	agentRouter.Get("/{id}", agentHandler.GetAgent)
	agentRouter.Get("/", agentHandler.GetAllAgents)
	agentRouter.Get("/{id}/commands", agentHandler.GetCommandSpecs)
	agentRouter.Patch("/{id}/name", agentHandler.PatchEditName)
	agentRouter.Patch("/{id}/description", agentHandler.PatchEditDesc)
	agentRouter.Post("/register", agentHandler.PostRegisterAgent)

	mux.Mount("/tasks", taskRouter)
	mux.Mount("/agents", agentRouter)

	return mux
}
