package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/utils"
)

type AgentHandler struct {
	store *store.AgentStore
}

func NewAgentHandler(store *store.AgentStore) *AgentHandler {
	return &AgentHandler{store: store}
}

func (a *AgentHandler) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	allAgents := a.store.GetAllAgents()

	if err := json.NewEncoder(w).Encode(allAgents); err != nil {
		log.Printf("error encoding response: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *AgentHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("id")
	if agentID == "" {
		utils.WriteError(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	agent, err := a.store.GetAgent(agentID)
	if err != nil && !errors.Is(err, store.ErrAgentNotExist) {
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, agent, http.StatusOK)
}
