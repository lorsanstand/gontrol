package service

import "errors"

var ErrAgentNotFound error = errors.New("agent not found")
var ErrTaskNotFound error = errors.New("task not found")
