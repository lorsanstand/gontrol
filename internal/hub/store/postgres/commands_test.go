package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	gen "github.com/lorsanstand/gontrol/internal/hub/store/postgres/gen"
	"github.com/lorsanstand/gontrol/internal/hub/store/postgres/mocks"
	"github.com/lorsanstand/gontrol/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newCommandRepo(q *mocks.MockQuerier) *CommandSpecsRepository {
	return &CommandSpecsRepository{q: q}
}

func TestCommandSpecsRepository_GetCommandSpecs(t *testing.T) {
	commands := models.CommandSpecs{
		"test_command": models.CommandSpec{
			Command: "test_command",
			Args:    []models.ArgSpec{},
			Label:   "test_label",
		},
	}

	data, err := json.Marshal(commands)
	require.NoError(t, err)

	agentID := uuid.New()

	tests := []struct {
		name        string
		expectedErr error
		expected    models.CommandSpecs
		setupMocks  func(q *mocks.MockQuerier)
	}{
		{
			name:        "successful",
			expectedErr: nil,
			expected:    commands,
			setupMocks: func(q *mocks.MockQuerier) {
				q.EXPECT().GetCommandsByAgentID(gomock.Any(), agentID).Return(data, nil)
			},
		},
		{
			name:        "error not found",
			expectedErr: store.ErrNotFound,
			expected:    models.CommandSpecs{},
			setupMocks: func(q *mocks.MockQuerier) {
				q.EXPECT().GetCommandsByAgentID(gomock.Any(), agentID).Return(nil, sql.ErrNoRows)
			},
		},
		{
			name:        "unhandled error",
			expectedErr: sql.ErrConnDone,
			expected:    models.CommandSpecs{},
			setupMocks: func(q *mocks.MockQuerier) {
				q.EXPECT().GetCommandsByAgentID(gomock.Any(), agentID).Return(nil, sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		runner := func(t *testing.T) {
			ctl := gomock.NewController(t)
			q := mocks.NewMockQuerier(ctl)
			r := newCommandRepo(q)

			ctx := context.Background()
			tt.setupMocks(q)

			command, err := r.GetCommandSpecs(ctx, agentID)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.expected, command)
		}

		t.Run(tt.name, runner)
	}
}

func TestCommandSpecsRepository_SetCommandSpecs(t *testing.T) {
	commands := models.CommandSpecs{
		"test_command": models.CommandSpec{
			Command: "test_command",
			Args:    []models.ArgSpec{},
			Label:   "test_label",
		},
	}

	data, err := json.Marshal(commands)
	require.NoError(t, err)

	agentID := uuid.New()

	errTest := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation}

	tests := []struct {
		name        string
		expectedErr error
		setupMocks  func(q *mocks.MockQuerier)
	}{
		{
			name:        "successful",
			expectedErr: nil,
			setupMocks: func(q *mocks.MockQuerier) {
				q.EXPECT().SetCommandsByAgentID(gomock.Any(), gen.SetCommandsByAgentIDParams{AgentID: agentID, Commands: data}).Return(nil)
			},
		},
		{
			name:        "unhandled error",
			expectedErr: sql.ErrConnDone,
			setupMocks: func(q *mocks.MockQuerier) {
				q.EXPECT().SetCommandsByAgentID(gomock.Any(), gen.SetCommandsByAgentIDParams{AgentID: agentID, Commands: data}).Return(sql.ErrConnDone)
			},
		},
		{
			name:        "error not found",
			expectedErr: store.ErrNotFound,
			setupMocks: func(q *mocks.MockQuerier) {
				q.EXPECT().SetCommandsByAgentID(gomock.Any(), gen.SetCommandsByAgentIDParams{AgentID: agentID, Commands: data}).Return(errTest)
			},
		},
	}

	for _, tt := range tests {
		runner := func(t *testing.T) {
			ctl := gomock.NewController(t)
			q := mocks.NewMockQuerier(ctl)
			r := newCommandRepo(q)

			ctx := context.Background()
			tt.setupMocks(q)

			err := r.SetCommandSpecs(ctx, agentID, commands)
			require.ErrorIs(t, err, tt.expectedErr)
		}

		t.Run(tt.name, runner)
	}
}
