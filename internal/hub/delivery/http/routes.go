package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handler"
)

func RegisterRoutes(
	taskHandler *handler.TaskHandler,
	agentHandler *handler.AgentHandler,
) *chi.Mux {
	mux := chi.NewRouter()

	taskRouter := chi.NewRouter()
	taskRouter.Post("/", taskHandler.PostCreateTask)
	taskRouter.Get("/{id}", taskHandler.GetTask)
	taskRouter.Get("/agent/{id}", taskHandler.GetTaskByAgentID)

	agentRouter := chi.NewRouter()
	agentRouter.Get("/{id}", agentHandler.GetAgent)
	agentRouter.Get("/", agentHandler.GetAllAgents)
	agentRouter.Get("/{id}/commands", agentHandler.GetCommandSpecs)
	agentRouter.Patch("/{id}/name", agentHandler.PatchEditName)
	agentRouter.Patch("/{id}/description", agentHandler.PatchEditDesc)

	mux.Mount("/tasks", taskRouter)
	mux.Mount("/agents", agentRouter)

	return mux
}
