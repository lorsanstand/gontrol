package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/models"
	"github.com/lorsanstand/gontrol/internal/utils"
)

type HeartbeatHandler struct {
	log *slog.Logger
	hb  *service.HeartbeatService
}

func NewHeartbeatHandler(logger *slog.Logger, heartbeat *service.HeartbeatService) *HeartbeatHandler {
	return &HeartbeatHandler{hb: heartbeat, log: logger}
}

func (h *HeartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var heartbeatReq models.HeartbeatRequest

	err := json.NewDecoder(r.Body).Decode(&heartbeatReq)
	if err != nil {
		h.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	heartbeatResp, err := h.hb.Handle(ctx, heartbeatReq)
	if err != nil {
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, heartbeatResp, http.StatusCreated)
}
