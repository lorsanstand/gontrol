package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/models"
)

type HeartbeatHandler struct {
	Store *store.AgentStore
}

func (h *HeartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var agentRequest models.HeartbeatRequest

	err := json.NewDecoder(r.Body).Decode(&agentRequest)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
	}

	h.Store.RegisterOrUpdate(agentRequest)

	w.WriteHeader(http.StatusOK)
}
