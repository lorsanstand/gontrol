package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	modelsAgent "github.com/lorsanstand/gontrol/internal/agent/models"
	"github.com/lorsanstand/gontrol/internal/models"
)

type HubClient struct {
	addr   url.URL
	client *http.Client
}

func NewHubClient(address url.URL, client *http.Client) *HubClient {
	return &HubClient{client: client, addr: address}
}

// SendHeartbeat отправялет хабу heartbeat.
// Если статус код не равен 201 возращает ошибку.
func (h *HubClient) SendHeartbeat(ctx context.Context, heartbeat models.HeartbeatRequest) (models.HeartbeatResponse, error) {
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("marshal heartbeat request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.addr.JoinPath("heartbeat").String(),
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("create heartbeat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("send heartbeat request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return models.HeartbeatResponse{}, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var result models.HeartbeatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("unmarshal heartbeat response: %w", err)
	}

	return result, nil
}

// SendTaskResult отправляет хабу результат задачи.
// Если статус код не равен 201 возращает ошибку.
func (h *HubClient) SendTaskResult(ctx context.Context, id uuid.UUID, taskResult models.TaskResult) error {
	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		return fmt.Errorf("marshal task result request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.addr.JoinPath("tasks", id.String(), "result").String(),
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return fmt.Errorf("create task result request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("send task result request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return nil
}

// RegisterAgent регистрирует агента.
// Если статус код не равен 201 возращает ошибку.
func (h *HubClient) RegisterAgent(ctx context.Context, specs modelsAgent.Specs) (*models.Agent, error) {
	jsonData, err := json.Marshal(specs)
	if err != nil {
		return nil, fmt.Errorf("marshal specification: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.addr.JoinPath("agents", "register").String(),
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("create register agent request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send create agent request: %w", err)
	}
	defer resp.Body.Close()

	var agent models.Agent

	if err = json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return nil, fmt.Errorf("unmarshal agent: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return &agent, nil
}
