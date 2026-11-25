package grpc

import (
	"context"
	"fmt"
	"net"

	protov1 "github.com/lorsanstand/gontrol/contracts/gen/proto/v1"
	"github.com/lorsanstand/gontrol/internal/hub/delivery/grpc/handler"
	"google.golang.org/grpc"
)

type HubServer struct {
	gRPCServer *grpc.Server
	port       int
}

func New(api *handler.ServerAPI, port int) *HubServer {
	gRPC := grpc.NewServer()

	protov1.RegisterHubAgentServiceServer(gRPC, api)

	return &HubServer{gRPCServer: gRPC, port: port}
}

func (h *HubServer) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", h.port))
	if err != nil {
		return fmt.Errorf("init tcp server: %w", err)
	}

	if err = h.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (h *HubServer) GracefulShutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		defer close(done)
		h.gRPCServer.GracefulStop()
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		h.gRPCServer.Stop()
		return ctx.Err()
	}
}
