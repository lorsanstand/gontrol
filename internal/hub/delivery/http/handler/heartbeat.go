package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uuid"

	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/models"
)

type HeartbeatHandler struct {
	store *store.AgentStore
}

func NewHeartbeatHandler(store *store.AgentStore) *HeartbeatHandler {
	return &HeartbeatHandler{store: store}
}

func (h *HeartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var agentRequest models.HeartbeatRequest

	err := json.NewDecoder(r.Body).Decode(&agentRequest)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if agentRequest.AgentID == "" {
		agentRequest.AgentID = uuid.New().String()
	}

	h.store.RegisterOrUpdate(agentRequest)

	task, err := h.store.PopTask(agentRequest.AgentID)
	if err != nil && !errors.Is(err, store.ErrTasksNotExist) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	agentResponse := models.HeartbeatResponse{AgentID: agentRequest.AgentID}
	if err == nil {
		agentResponse.Task = &task
	}

	if err := json.NewEncoder(w).Encode(agentResponse); err != nil {
		log.Printf("error encoding response: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}
