package service

import (
	"fmt"
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/lorsanstand/gontrol/internal/agent/config"
	modelsAgent "github.com/lorsanstand/gontrol/internal/agent/models"
)

type SystemService struct {
	cfg *config.Config
}

func NewSystemService(cfg *config.Config) *SystemService {
	return &SystemService{cfg: cfg}
}

// GetSpecs Возращает системные характеристики устройства.
// Если не удалось их получить возращает ошибку.
func (s *SystemService) GetSpecs() (*modelsAgent.Specs, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname: %w", err)
	}
	OS := runtime.GOOS

	return &modelsAgent.Specs{OS: OS, Hostname: hostname}, nil
}

func (s *SystemService) SaveAgentID(AgentID uuid.UUID) error {
	s.cfg.AgentID = &AgentID
	return s.cfg.Save()
}