package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/hub/store/postgres"
	"github.com/lorsanstand/gontrol/internal/utils"
)

type AgentHandler struct {
	log   *slog.Logger
	agent *service.AgentService
}

func NewAgentHandler(logger *slog.Logger, store *service.AgentService) *AgentHandler {
	return &AgentHandler{agent: store, log: logger}
}

func (a *AgentHandler) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// как появятся новые ошибки сделать нормальную обработку
	allAgents, err := a.agent.GetAll(ctx)
	if err != nil {
		a.log.Error("error get agents", slog.Any("error", err))
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, allAgents, http.StatusOK)
}

func (a *AgentHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agentID := r.PathValue("id")
	if agentID == "" {
		a.log.Debug("agent ID not found")
		utils.WriteError(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(agentID)
	if err != nil {
		a.log.Debug("invalid agent ID", slog.String("AgentID", agentID), slog.Any("error", err))
		utils.WriteError(w, "invalid agent ID format: must be a valid UUID", http.StatusBadRequest)
		return
	}

	agent, err := a.agent.Get(ctx, id)
	if err != nil && !errors.Is(err, store.ErrAgentNotExist) {
		a.log.Error("error get agent", slog.Any("error", err), slog.String("agentID", agentID), slog.Any("error", err))
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, store.ErrAgentNotExist) {
		a.log.Debug("agent not found", slog.String("agentID", agentID))
		utils.WriteError(w, "Agent not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, agent, http.StatusOK)
}

func (a *AgentHandler) EditName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	agentID := r.PathValue("id")
	if agentID == "" {
		a.log.Debug("agent ID not found")
		utils.WriteError(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(agentID)
	if err != nil {
		a.log.Debug("invalid agent ID", slog.String("AgentID", agentID), slog.Any("error", err))
		utils.WriteError(w, "invalid agent ID format: must be a valid UUID", http.StatusBadRequest)
		return
	}

	body := struct {
		Name string `json:"Name"`
	}{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	err = a.agent.EditNameAgent(ctx, id, body.Name)
	if err != nil {
		if err != nil {
			if errors.Is(err, postgres.ErrChangeData) {
				a.log.Debug("error change data", slog.Any("error", err))
				utils.WriteError(w, "Agent not found", http.StatusNotFound)
				return
			}
			utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *AgentHandler) EditDesc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	agentID := r.PathValue("id")
	if agentID == "" {
		a.log.Debug("agent ID not found")
		utils.WriteError(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(agentID)
	if err != nil {
		a.log.Debug("invalid agent ID", slog.String("AgentID", agentID), slog.Any("error", err))
		utils.WriteError(w, "invalid agent ID format: must be a valid UUID", http.StatusBadRequest)
		return
	}

	body := struct {
		Description string `json:"description"`
	}{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	err = a.agent.EditDescAgent(ctx, id, body.Description)
	if err != nil {
		if errors.Is(err, postgres.ErrChangeData) {
			a.log.Debug("error change data", slog.Any("error", err))
			utils.WriteError(w, "Agent not found", http.StatusNotFound)
			return
		}
		utils.WriteError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
