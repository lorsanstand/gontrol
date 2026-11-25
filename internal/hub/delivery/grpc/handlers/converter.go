package handlers

import (
	"fmt"

	protov1 "github.com/lorsanstand/gontrol/contracts/gen/proto/v1"
	"github.com/lorsanstand/gontrol/internal/models"
)

func toRegisterAgentRequestDTO(req *protov1.RegisterAgentRequest) (*models.RegisterAgentRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil pointer")
	}

	specs := make(models.CommandSpecs, len(req.CommandSpecs))
	for _, specProto := range req.CommandSpecs {
		if specProto == nil {
			continue
		}

		args := make([]models.ArgSpec, 0, len(specProto.Args))

		for _, argProto := range specProto.Args {
			if argProto == nil {
				continue
			}

			args = append(args, toArgSpec(argProto))
		}

		specs[specProto.Command] = models.CommandSpec{Command: specProto.Command, Label: specProto.Label, Args: args}
	}

	return &models.RegisterAgentRequest{
		OS:           req.Os,
		Hostname:     req.Hostname,
		CommandSpecs: specs,
	}, nil
}

func ptrInt64ToInt(p *int64) *int {
	if p == nil {
		return nil
	}
	v := int(*p)
	return &v
}

func toArgSpec(argProto *protov1.ArgSpec) models.ArgSpec {
	if argProto == nil {
		return models.ArgSpec{}
	}

	return models.ArgSpec{
		Label:     argProto.Label,
		Type:      models.ArgType(argProto.Type),
		Name:      argProto.Name,
		Required:  argProto.Required,
		Min:       ptrInt64ToInt(argProto.Min),
		Max:       ptrInt64ToInt(argProto.Max),
		MinLength: ptrInt64ToInt(argProto.MinLength),
		MaxLength: ptrInt64ToInt(argProto.MaxLength),
	}
}

func toProtoHearbeatResponse(hb models.HeartbeatResponse) *protov1.HeartbeatResponse {
	var task *protov1.Task

	if hb.Task != nil {
		args := make(map[string][]byte, len(hb.Task.Args))

		for name, raw := range hb.Task.Args {
			args[name] = []byte(raw)
		}

		task = &protov1.Task{Id: hb.Task.ID.String(), Command: hb.Task.Command, Args: args}
	}

	return &protov1.HeartbeatResponse{Id: hb.AgentID.String(), Task: task}
}
