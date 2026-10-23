package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/utils"
)

func ReturnErrors(w http.ResponseWriter, msg string, err error) {
	switch {
	case errors.Is(err, service.ErrAgentNotFound):
		slog.Debug(msg, slog.Any("error", err))
		utils.WriteError(w, "Agent not found", http.StatusNotFound)
	case errors.Is(err, service.ErrTaskNotFound):
		slog.Debug(msg, slog.Any("error", err))
		utils.WriteError(w, "Task not found", http.StatusNotFound)

	default:
		slog.Error(msg, slog.Any("error", err))
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)

	}
}
