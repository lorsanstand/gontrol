package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/models"
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

	allAgents, err := a.agent.GetAll(ctx)
	if err != nil {
		ReturnErrors(w, "error get agents", err)
		return
	}

	utils.WriteJSON(w, allAgents, http.StatusOK)
}

func (a *AgentHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agentID := r.PathValue("id")
	if agentID == "" {
		a.log.Debug("agent ID is empty")
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
	if err != nil {
		ReturnErrors(w, "error get agent", err)
		return
	}

	utils.WriteJSON(w, agent, http.StatusOK)
}

func (a *AgentHandler) EditName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	agentID := r.PathValue("id")
	if agentID == "" {
		a.log.Debug("agent ID is empty")
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
		ReturnErrors(w, "error edit name agent", err)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		ReturnErrors(w, "error edit agent description", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *AgentHandler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createAgent models.CreateAgent

	err := json.NewDecoder(r.Body).Decode(&createAgent)
	if err != nil {
		a.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	agent, err := a.agent.Register(ctx, createAgent)
	if err != nil {
		ReturnErrors(w, "error create agent", err)
		return
	}

	newAgent := models.Agent{ID: agent.ID, Hostname: agent.Hostname, OS: agent.OS}

	utils.WriteJSON(w, newAgent, http.StatusCreated)
}
