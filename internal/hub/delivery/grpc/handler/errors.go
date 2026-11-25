package handler

import (
	"errors"

	"github.com/lorsanstand/gontrol/internal/hub/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func returnCode(err error) error {
	switch {
	case errors.Is(err, service.ErrAgentNotFound):
		return status.Error(codes.NotFound, "agent not found")
	case errors.Is(err, service.ErrTaskNotFound):
		return status.Error(codes.NotFound, "task not found")
	case errors.Is(err, service.ErrValidation):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "unhandled error")
	}
}
