package models

import "github.com/lorsanstand/gontrol/internal/models"

type Agent struct {
	AgentID  string
	OS       string
	Hostname string
	Tasks    []models.Task
}
