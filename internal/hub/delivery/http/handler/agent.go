package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lorsanstand/gontrol/internal/hub"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/models"
	"github.com/lorsanstand/gontrol/internal/utils"
)

type AgentHandler struct {
	log       *slog.Logger
	agent     *service.AgentService
	validator *validator.Validate
}

func NewAgentHandler(logger *slog.Logger, store *service.AgentService) *AgentHandler {
	return &AgentHandler{agent: store, log: logger, validator: validator.New()}
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

func (a *AgentHandler) PatchEditName(w http.ResponseWriter, r *http.Request) {
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

	var body modelsHub.EditNameAgent
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if err = a.validator.StructCtx(ctx, body); err != nil {
		utils.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.agent.EditNameAgent(ctx, id, body.Name)
	if err != nil {
		ReturnErrors(w, "error edit name agent", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *AgentHandler) PatchEditDesc(w http.ResponseWriter, r *http.Request) {
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

	var body modelsHub.EditDescAgent
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if err = a.validator.StructCtx(ctx, body); err != nil {
		utils.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.agent.EditDescAgent(ctx, id, body.Description)
	if err != nil {
		ReturnErrors(w, "error edit agent description", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *AgentHandler) PostRegisterAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createAgent models.RegisterAgentRequest

	err := json.NewDecoder(r.Body).Decode(&createAgent)
	if err != nil {
		a.log.Debug("invalid json", slog.Any("error", err))
		utils.WriteError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if err = a.validator.StructCtx(ctx, createAgent); err != nil {
		utils.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	agent, err := a.agent.Register(ctx, createAgent)
	if err != nil {
		ReturnErrors(w, "error create agent", err)
		return
	}

	newAgent := models.RegisterAgentResponse{ID: agent.ID, HeartbeatTTL: hub.TimeSendHeartbeat * time.Second}

	utils.WriteJSON(w, newAgent, http.StatusCreated)
}

func (a *AgentHandler) GetCommandSpecs(w http.ResponseWriter, r *http.Request) {
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

	specs, err := a.agent.GetCommandSpecs(ctx, id)
	if err != nil {
		ReturnErrors(w, "error get agent", err)
		return
	}

	utils.WriteJSON(w, specs, http.StatusOK)
}
