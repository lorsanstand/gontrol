package hub_client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	protov1 "github.com/lorsanstand/gontrol/contracts/gen/proto/v1"
	"github.com/lorsanstand/gontrol/internal/models"
)

func toHeartbeatResponseDTO(hb *protov1.HeartbeatResponse) (*models.HeartbeatResponse, error) {
	if hb == nil {
		return nil, fmt.Errorf("nil pointer deference")
	}

	id, err := uuid.Parse(hb.Id)
	if err != nil {
		return nil, fmt.Errorf("parse agent id: %w", err)
	}

	var task *models.Task

	if hb.Task != nil {
		args := make(map[string]json.RawMessage, len(hb.Task.Args))
		for name, raw := range hb.Task.Args {
			args[name] = json.RawMessage(raw)
		}

		id, err = uuid.Parse(hb.Task.Id)
		if err != nil {
			return nil, fmt.Errorf("parse id: %w", err)
		}

		task = &models.Task{ID: id, Command: hb.Task.Command, Args: args}
	}

	return &models.HeartbeatResponse{AgentID: id, Task: task}, nil
}

func toRegisterAgentRequestProto(req models.RegisterAgentRequest) *protov1.RegisterAgentRequest {
	specsProto := make([]*protov1.CommandSpec, 0, len(req.CommandSpecs))
	for _, spec := range req.CommandSpecs {
		args := make([]*protov1.ArgSpec, 0, len(spec.Args))

		for _, arg := range spec.Args {
			args = append(args, toArgSpecProto(arg))
		}

		specProto := &protov1.CommandSpec{Command: spec.Command, Label: spec.Label, Args: args}

		specsProto = append(specsProto, specProto)
	}

	return &protov1.RegisterAgentRequest{Os: req.OS, Hostname: req.Hostname, CommandSpecs: specsProto}
}

func toArgSpecProto(arg models.ArgSpec) *protov1.ArgSpec {
	return &protov1.ArgSpec{
		Label: arg.Label,
		Name: arg.Name,
		Type: string(arg.Type),
		Required: arg.Required,
		Placeholder: arg.Placeholder,
		Min: ptrIntToInt64(arg.Min),
		Max: ptrIntToInt64(arg.Max),
		MinLength: ptrIntToInt64(arg.MinLength),
		MaxLength: ptrIntToInt64(arg.MaxLength),
		Enum: arg.Enum,
	}
}

func toRegisterAgentResponseDTO(resp *protov1.RegisterAgentResponse) (*models.RegisterAgentResponse, error) {
	if resp == nil {
		return nil, fmt.Errorf("nil pointer deference")
	}

	id, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, fmt.Errorf("parse id: %w", err)
	}

	return &models.RegisterAgentResponse{ID: id, HeartbeatTTL: time.Duration(resp.HeartbeatTtl)}, nil
}

func ptrIntToInt64(p *int) *int64 {
	if p == nil {
		return nil
	}
	v := int64(*p)
	return &v
}