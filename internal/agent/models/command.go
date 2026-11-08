package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lorsanstand/gontrol/internal/models"
)

type Command struct {
	Exec func(ctx context.Context, args map[string]json.RawMessage) ([]byte, error)
	Spec models.CommandSpec
}

type Commands struct {
	Commands map[string]Command
}

func (c *Commands) GetCommandAndValidate(task models.Task) (*Command, error) {
	command, ok := c.Commands[task.Command]
	if !ok {
		return nil, fmt.Errorf("command not found")
	}

	return &command, command.Spec.ValidateArgs(task.Args)
}
