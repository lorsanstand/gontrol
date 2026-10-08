package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lorsanstand/gontrol/internal/models"
)

type HubClient struct {
	addr   string
	client *http.Client
}

func NewHubClient(address string, client *http.Client) *HubClient {
	return &HubClient{client: client, addr: address}
}

func (h *HubClient) SendHeartbeat(ctx context.Context, heartbeat models.HeartbeatRequest) (models.HeartbeatResponse, error) {
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		return models.HeartbeatResponse{}, fmt.Errorf("marshal heartbeat request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.addr+fmt.Sprintf("heartbeat"),
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

func (h *HubClient) SendTaskResult(ctx context.Context, id string, taskResult models.TaskResult) error {
	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		return fmt.Errorf("marshal task result request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.addr+fmt.Sprintf("tasks/%v/result", id),
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
