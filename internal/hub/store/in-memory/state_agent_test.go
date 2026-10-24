package in_memory

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAgentStateStore_Set_Get(t *testing.T) {
	s := NewAgentStateStore(1)
	id := uuid.New()
	id2 := uuid.New()

	s.SetOnline(id)

	require.Equal(t, true, s.GetStatus(id))
	require.Equal(t, false, s.GetStatus(id2))
}

func TestAgentStateStore_AutoCleanerCancel(t *testing.T) {
	s := NewAgentStateStore(5)
	ctx, cancel := context.WithCancel(context.Background())
	autoChan := make(chan error)

	go func() {
		err := s.AutoCleaner(ctx)
		autoChan <- err
	}()

	cancel()

	select {
	case err := <-autoChan:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.Tick(1 * time.Second):
		t.Error("auto cleaner dont stopped")
	}

}

func TestAgentStateStore_CleanUp(t *testing.T) {
	s := NewAgentStateStore(0)
	id := uuid.New()

	s.SetOnline(id)
	s.CleanUp()
	require.Equal(t, false, s.GetStatus(id))
}
