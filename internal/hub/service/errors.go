package service

import "errors"

var ErrAgentNotFound error = errors.New("agent not found")
var ErrTaskNotFound error = errors.New("task not found")
var ErrValidation error = errors.New("failed validation")
