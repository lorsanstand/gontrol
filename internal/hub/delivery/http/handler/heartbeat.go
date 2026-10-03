package handler

import (
	"encoding/json"
	"net/http"

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
	}

	h.store.RegisterOrUpdate(agentRequest)

	w.WriteHeader(http.StatusOK)
}
