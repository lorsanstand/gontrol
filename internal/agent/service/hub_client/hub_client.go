package hub_client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	protov1 "github.com/lorsanstand/gontrol/contracts/gen/proto/v1"
	"github.com/lorsanstand/gontrol/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type HubClient struct {
	client protov1.HubAgentServiceClient
}

func inteceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields)
	})
}

func NewHubClient(logger *slog.Logger, addr string, timeout time.Duration, retries int) (*HubClient, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retries)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// TODO: исправить защиту в будущем для mtls
	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(inteceptorLogger(logger), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}

	return &HubClient{client: protov1.NewHubAgentServiceClient(cc)}, nil
}

// SendHeartbeat отправялет хабу heartbeat.
// Если статус код не равен 201 возращает ошибку.
func (h *HubClient) SendHeartbeat(ctx context.Context, heartbeat models.HeartbeatRequest) (*models.HeartbeatResponse, error) {
	protoReq := &protov1.HeartbeatRequest{Id: heartbeat.AgentID.String()}

	protoResp, err := h.client.SendHeartbeat(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("send hb to hub: %w", err)
	}

	return toHeartbeatResponseDTO(protoResp)
}

// SendTaskResult отправляет хабу результат задачи.
// Если статус код не равен 201 возращает ошибку.
func (h *HubClient) SendTaskResult(ctx context.Context, taskResult models.TaskResult) error {
	protoReq := &protov1.TaskResultRequest{
		Id:         taskResult.ID.String(),
		TimeExecMs: int64(taskResult.TimeExecMS),
		Result:     []byte(taskResult.Result),
	}

	_, err := h.client.SendTaskResult(ctx, protoReq)
	if err != nil {
		return fmt.Errorf("send task result to hub: %w", err)
	}

	return nil
}

// RegisterAgent регистрирует агента.
// Если статус код не равен 201 возращает ошибку.
func (h *HubClient) RegisterAgent(ctx context.Context, createAgent models.RegisterAgentRequest) (*models.RegisterAgentResponse, error) {
	protoReq := toRegisterAgentRequestProto(createAgent)

	protoResp, err := h.client.RegisterAgent(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("send register agent to hub: %w", err)
	}

	return toRegisterAgentResponseDTO(protoResp)
}
