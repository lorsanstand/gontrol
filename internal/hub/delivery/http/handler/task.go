package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/utils"
)

type TaskHandler struct {
	log  *slog.Logger
	task *service.TaskService
}

func NewTaskHandler(logger *slog.Logger, task *service.TaskService) *TaskHandler {
	return &TaskHandler{task: task, log: logger}
}

func (t *TaskHandler) PostCreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var taskRequest modelsHub.CreateTask

	err := json.NewDecoder(r.Body).Decode(&taskRequest)
	if err != nil {
		t.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	task, err := t.task.Add(ctx, taskRequest)
	if err != nil {
		if errors.Is(err, store.ErrAgentNotExist) {
			utils.WriteError(w, "Agent not found", http.StatusNotFound)
			return
		}

		utils.WriteError(w, "Failed to save task", http.StatusBadGateway)
		return
	}

	utils.WriteJSON(w, task, http.StatusCreated)
}

// Написать логику сохранения результата и выдача пользователю
func (t *TaskHandler) PostResultTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
