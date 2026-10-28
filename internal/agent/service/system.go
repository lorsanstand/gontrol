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
	specs modelsAgent.Specs
	cfg *config.Config
}

// Init Иницилизрует SystemService.
// При создании берет данные о системе.
func Init(cfg *config.Config) (*SystemService, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname: %w", err)
	}
	OS := runtime.GOOS

	specs := modelsAgent.Specs{OS: OS, Hostname: hostname}

	return &SystemService{cfg: cfg, specs: specs}, nil
}

// GetSpecs Возращает системные характеристики устройства.
// Если не удалось их получить возращает ошибку.
func (s *SystemService) GetSpecs() modelsAgent.Specs {
	return s.specs
}

// SaveAgentID сохраняет ID агента в файле конфига.
func (s *SystemService) SaveAgentID(AgentID uuid.UUID) error {
	s.cfg.AgentID = &AgentID
	return s.cfg.Save()
}

// GetAgentID возращает ID агента.
func (s *SystemService) GetAgentID() *uuid.UUID {
	return s.cfg.AgentID
}