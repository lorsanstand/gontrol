package in_memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AgentStateStore struct {
	mu    sync.RWMutex
	store map[uuid.UUID]time.Time
	ttl   time.Duration
}

// NewAgentStateStore создает in-memory хранилище статусов агентов.
// ttl время жизни статуса в секундах.
func NewAgentStateStore(ttl int) *AgentStateStore {
	return &AgentStateStore{store: make(map[uuid.UUID]time.Time), ttl: time.Duration(ttl) * time.Second}
}

// GetStatus Возращает статус конкретного агента
func (a *AgentStateStore) GetStatus(id uuid.UUID) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	_, ok := a.store[id]
	return ok
}

// SetOnline меняет статус агента на онлайн
func (a *AgentStateStore) SetOnline(id uuid.UUID) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.setStatus(id, true)
}

// AutoCleaner автоматически чистит агентов у которых истек ttl
// Чистка происходит по указанному ttl. Если ttl<=0 берется, чистка происходит с промежутком в 1 секунду
func (a *AgentStateStore) AutoCleaner(ctx context.Context) error {
	ttl := a.ttl
	if ttl <= 0 {
		ttl = 1*time.Second
	}

	ticker := time.NewTicker(ttl)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		a.CleanUp()
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// CleanUp чистит агентов у которых истек ttl
func (a *AgentStateStore) CleanUp() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for id, times := range a.store {
		if time.Since(times) > a.ttl {
			a.setStatus(id, false)
		}
	}
}

// setStatus меняет статус агента.
// ВНУТРИ НЕ ВЫЗЫВАЕТСЯ MUTEX БУДЬТЕ И МОЖЕТ ПРИВЕСТИ К ПАНИКЕ.
func (a *AgentStateStore) setStatus(id uuid.UUID, status bool) {
	if status {
		a.store[id] = time.Now()
		return
	}

	delete(a.store, id)
}
