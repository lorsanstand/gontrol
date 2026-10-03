package http

import (
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handler"
	"github.com/lorsanstand/gontrol/internal/hub/store"
)

func RegisterRoutes(mux *http.ServeMux, store *store.AgentStore) {
	heartbeatHandler := handler.NewHeartbeatHandler(store)
	mux.Handle("POST /api/v1/heartbeat", heartbeatHandler)

	taskHandler := handler.NewTaskHandler(store)
	mux.HandleFunc("POST /api/v1/tasks", taskHandler.PostCreateTask)
	mux.HandleFunc("POST /api/v1/tasks/{id}/result", taskHandler.PostResultTask)

	agentHandler := handler.NewAgentHandler(store)
	mux.HandleFunc("GET /api/v1/agents/{id}", agentHandler.GetAgent)
	mux.HandleFunc("GET /api/v1/agents/", agentHandler.GetAllAgents)
}
