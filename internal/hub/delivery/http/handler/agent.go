package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/utils"
)

type AgentHandler struct {
	log   *slog.Logger
	agent *service.AgentService
}

func NewAgentHandler(logger *slog.Logger, store *service.AgentService) *AgentHandler {
	return &AgentHandler{agent: store, log: logger}
}

func (a *AgentHandler) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// как появятся новые ошибки сделать нормальную обработку
	allAgents, err := a.agent.GetAll(ctx)
	if err != nil {
		a.log.Error("error get agents", slog.Any("error", err))
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, allAgents, http.StatusOK)
}

func (a *AgentHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agentID := r.PathValue("id")
	if agentID == "" {
		a.log.Debug("agent ID not found")
		utils.WriteError(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	agent, err := a.agent.Get(ctx, agentID)
	if err != nil && !errors.Is(err, store.ErrAgentNotExist) {
		a.log.Error("error get agent", slog.Any("error", err), slog.String("agentID", agentID))
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, store.ErrAgentNotExist) {
		a.log.Debug("agent not found", slog.String("agentID", agentID))
		utils.WriteError(w, "Agent not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, agent, http.StatusOK)
}
