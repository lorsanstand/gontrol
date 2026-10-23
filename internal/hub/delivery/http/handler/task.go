package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/hub/store/postgres"
	"github.com/lorsanstand/gontrol/internal/models"
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

func (t *TaskHandler) PostResultTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var taskResult models.TaskResult

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&taskResult)
	if err != nil {
		t.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	err = t.task.SaveResult(ctx, taskResult.ID, taskResult.Result, taskResult.TimeExecMS)
	if err != nil {
		if errors.Is(err, postgres.ErrChangeData) {
			t.log.Debug("error change data", slog.Any("error", err))
			utils.WriteError(w, "Task not found", http.StatusNotFound)
			return
		}
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
