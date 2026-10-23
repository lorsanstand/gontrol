package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/service"
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
		ReturnErrors(w, "error create new task", err)
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
		ReturnErrors(w, "error save task result", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
