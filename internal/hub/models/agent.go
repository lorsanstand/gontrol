package models

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	OS          string    `json:"os"`
	Hostname    string    `json:"hostname"`
	Description *string   `json:"description"`
	State       bool      `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
}
