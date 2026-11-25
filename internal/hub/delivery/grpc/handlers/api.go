package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	protov1 "github.com/lorsanstand/gontrol/contracts/gen/proto/v1"
	"github.com/lorsanstand/gontrol/internal/hub"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	"github.com/lorsanstand/gontrol/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	log       *slog.Logger
	task      *service.TaskService
	hb        *service.HeartbeatService
	agent     *service.AgentService
	validator *validator.Validate
}

func NewServerAPI(
	logger *slog.Logger,
	task *service.TaskService,
	heartbeat *service.HeartbeatService,
	agent *service.AgentService,
) *ServerAPI {
	return &ServerAPI{agent: agent, task: task, hb: heartbeat, log: logger, validator: validator.New()}
}

func (s *ServerAPI) RegisterAgent(ctx context.Context, req *protov1.RegisterAgentRequest) (*protov1.RegisterAgentResponse, error) {
	regAgent, err := toRegisterAgentRequestDTO(req)
	if err != nil {
		s.log.Debug("failed convert to dto", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = s.validator.StructCtx(ctx, regAgent); err != nil {
		s.log.Debug("error validate", slog.Any("error", err))
		return nil, status.Errorf(codes.InvalidArgument, "validate: %w", err)
	}

	agent, err := s.agent.Register(ctx, *regAgent)
	if err != nil {
		s.log.Debug("error register agent", slog.Any("error", err))
		return nil, returnCode(err)
	}

	return &protov1.RegisterAgentResponse{Id: agent.ID.String(), HeartbeatTtl: int64(hub.TimeSendHeartbeat * time.Second)}, nil
}

func (s *ServerAPI) SendHeartbeat(ctx context.Context, req *protov1.HeartbeatRequest) (*protov1.HeartbeatResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request not found")
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		slog.Debug("invalid agent id", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, "invalid agent id")
	}

	hb := models.HeartbeatRequest{AgentID: id}

	heartbeatResp, err := s.hb.Handle(ctx, hb)
	if err != nil {
		slog.Debug("error save hb")
		return nil, returnCode(err)
	}

	return toProtoHearbeatResponse(*heartbeatResp), nil
}

func (s *ServerAPI) SendTaskResult(ctx context.Context, req *protov1.TaskResultRequest) (*protov1.TaskResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request not found")
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		slog.Debug("invalid task id", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, "invalid task id")
	}

	taskResult := models.TaskResult{ID: id, TimeExecMS: int(req.TimeExecMs), Result: string(req.Result)}

	if err = s.validator.StructCtx(ctx, taskResult); err != nil {
		s.log.Debug("error validate", slog.Any("error", err))
		return nil, status.Errorf(codes.InvalidArgument, "validate: %w", err)
	}

	err = s.task.SaveResult(ctx, taskResult.ID, taskResult.Result, taskResult.TimeExecMS)
	if err != nil {
		s.log.Debug("error save task result", slog.Any("error", err))
		return nil, returnCode(err)
	}

	return &protov1.TaskResultResponse{}, nil
}
