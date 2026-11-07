package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ArgType string

const (
	ArgString ArgType = "string"
	ArgInt    ArgType = "integer"
	ArgBool   ArgType = "boolean"
)

type ArgSpec struct {
	Name     string  `json:"name"`
	Label    string  `json:"label"`
	Type     ArgType `json:"type"`
	Required bool    `json:"required"`

	Min         *int   `json:"min,omitempty"`
	Max         *int   `json:"max,omitempty"`
	MinLength   *int   `json:"min_length,omitempty"`
	MaxLength   *int   `json:"max_length,omitempty"`
	Placeholder string `json:"placeholder"`
}

func validateArg(spec ArgSpec, raw json.RawMessage) error {
	switch spec.Type {
	case ArgBool:
		var value bool
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("%s must be a boolean", spec.Name)
		}
	case ArgInt:
		var value int
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("%s must be a integer", spec.Name)
		}

		if spec.Min != nil && value < *spec.Min {
			return fmt.Errorf("%s is too small", spec.Name)
		}

		if spec.Max != nil && value > *spec.Max {
			return fmt.Errorf("%s is too large", spec.Name)
		}
	case ArgString:
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("%s must be a string", spec.Name)
		}

		if spec.MinLength != nil && len(value) < *spec.MinLength {
			return fmt.Errorf("%s is too short", spec.Name)
		}

		if spec.MaxLength != nil && len(value) > *spec.MaxLength {
			return fmt.Errorf("%s is too long", spec.Name)
		}
	default:
		return fmt.Errorf("unsupported argument type: %s", spec.Type)
	}

	return nil
}

type CommandSpec struct {
	Command string    `json:"command"`
	Label   string    `json:"label"`
	Args    []ArgSpec `json:"args"`
}

type CommandSpecs struct {
	Commands map[string]CommandSpec
}

func (c *CommandSpecs) ValidateTask(task Task) error {
	command, ok := c.Commands[task.Command]
	if !ok {
		return fmt.Errorf("command not found")
	}

	var errs error = nil

	for _, argSpec := range command.Args {
		arg, ok := task.Args[argSpec.Name]
		if argSpec.Required && !ok {
			errs = errors.Join(errs, fmt.Errorf("%s is required", argSpec.Name))
			continue
		}

		errs = errors.Join(errs, validateArg(argSpec, arg))
	}

	return errs
}
