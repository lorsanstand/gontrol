package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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

func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agentID := r.PathValue("id")
	if agentID == "" {
		t.log.Debug("task ID is empty")
		utils.WriteError(w, "task ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(agentID)
	if err != nil {
		t.log.Debug("invalid task ID", slog.String("taskID", agentID), slog.Any("error", err))
		utils.WriteError(w, "invalid task ID format: must be a valid UUID", http.StatusBadRequest)
		return
	}

	task, err := t.task.GetByID(ctx, id)
	if err != nil {
		ReturnErrors(w, "error get task", err)
	}

	utils.WriteJSON(w, task, http.StatusOK)
}

func (t *TaskHandler) GetTaskByAgentID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agentID := r.PathValue("id")
	if agentID == "" {
		t.log.Debug("agent ID is empty")
		utils.WriteError(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(agentID)
	if err != nil {
		t.log.Debug("invalid agent ID", slog.String("AgentID", agentID), slog.Any("error", err))
		utils.WriteError(w, "invalid agent ID format: must be a valid UUID", http.StatusBadRequest)
		return
	}

	tasks, err := t.task.GetByAgentID(ctx, id)
	if err != nil {
		ReturnErrors(w, "error get task", err)
	}

	utils.WriteJSON(w, tasks, http.StatusOK)
}