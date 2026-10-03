package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"uuid"

	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/models"
	"github.com/lorsanstand/gontrol/internal/utils"
)

type TaskHandler struct {
	store *store.AgentStore
}

func NewTaskHandler(store *store.AgentStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (t *TaskHandler) PostCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskRequest modelsHub.CreateTask

	err := json.NewDecoder(r.Body).Decode(&taskRequest)
	if err != nil {
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	taskID := uuid.New().String()

	task := models.Task{ID: taskID, Command: taskRequest.Command, Args: taskRequest.Args}

	err = t.store.AddTask(taskRequest.AgentID, task)
	if err != nil {
		if errors.Is(err, store.ErrAgentNotExist) {
			utils.WriteError(w, "Agent not found", http.StatusNotFound)
			return
		}

		utils.WriteError(w, "Failed to save task", http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
