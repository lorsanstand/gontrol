package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/store"
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
		http.Error(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	agent, err := a.store.GetAgent(agentID)
	if err != nil && !errors.Is(err, store.ErrAgentNotExist) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(agent); err != nil {
		log.Printf("error encoding response: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}
