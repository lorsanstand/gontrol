package utils

import (
	"fmt"
	"os"
	"runtime"

	modelsAgent "github.com/lorsanstand/gontrol/internal/agent/models"
)

type Specification struct {

}

func NewSpecification() *Specification {
	return &Specification{}
}

// GetSpecs Возращает системные характеристики устройства.
// Если не удалось их получить возращает ошибку.
func (s *Specification) GetSpecs() (*modelsAgent.Specs, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname: %w", err)
	}
	OS := runtime.GOOS

	return &modelsAgent.Specs{OS: OS, Hostname: hostname}, nil
}
