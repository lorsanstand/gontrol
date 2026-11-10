package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/service/mocks"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/lorsanstand/gontrol/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newAgentServiceMock(storage agentRepository, state agentStateRepository) *AgentService {
	return &AgentService{log: &slog.Logger{}, store: storage, state: state}
}

type genIDMock struct {
	err error
	id  uuid.UUID
}

func (g genIDMock) gen() (uuid.UUID, error) {
	return g.id, g.err
}

func TestAgentService_Get(t *testing.T) {
	id := uuid.New()
	desc := "test_desc"
	createdAt := time.Now()

	agentDB := modelsHub.Agent{
		ID:          id,
		Name:        "test_name",
		OS:          "test_os",
		Hostname:    "test_hostname",
		Description: &desc,
		State:       false,
		CreatedAt:   createdAt,
	}

	expectedOnline := agentDB
	expectedOnline.State = true

	expectedOffline := agentDB
	expectedOffline.State = false

	tests := []struct {
		name       string
		expectErr  error
		expected   *modelsHub.Agent
		setupMocks func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository)
	}{
		{
			name:      "successful online",
			expectErr: nil,
			expected:  &expectedOnline,
			setupMocks: func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository) {
				r.EXPECT().GetAgentByID(gomock.Any(), id).Return(&agentDB, nil).Times(1)
				s.EXPECT().GetStatus(id).Return(true).Times(1)
			},
		},
		{
			name:      "successful offline",
			expectErr: nil,
			expected:  &expectedOffline,
			setupMocks: func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository) {
				r.EXPECT().GetAgentByID(gomock.Any(), id).Return(&agentDB, nil).Times(1)
				s.EXPECT().GetStatus(id).Return(false).Times(1)
			},
		},
		{
			name:      "agent not found",
			expectErr: ErrAgentNotFound,
			expected:  nil,
			setupMocks: func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository) {
				r.EXPECT().GetAgentByID(gomock.Any(), id).Return(nil, store.ErrNotFound).Times(1)
			},
		},
		{
			name:      "unhandled error",
			expectErr: sql.ErrConnDone,
			expected:  nil,
			setupMocks: func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository) {
				r.EXPECT().GetAgentByID(gomock.Any(), id).Return(nil, sql.ErrConnDone).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMockagentRepository(ctl)
			s := mocks.NewMockagentStateRepository(ctl)

			a := newAgentServiceMock(r, s)
			ctx := context.Background()

			tt.setupMocks(r, s)

			agent, err := a.Get(ctx, id)
			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expected, agent)
		})
	}
}

func TestAgentService_GetAll(t *testing.T) {
	desc := "test_desc"
	createdAt := time.Now()

	agentsDB := []modelsHub.Agent{
		{
			ID:          uuid.New(),
			Name:        "test_name_1",
			OS:          "test_os_1",
			Hostname:    "test_hostname_1",
			Description: &desc,
			State:       false,
			CreatedAt:   createdAt,
		},
		{
			ID:          uuid.New(),
			Name:        "test_name_2",
			OS:          "test_os_2",
			Hostname:    "test_hostname_2",
			Description: &desc,
			State:       false,
			CreatedAt:   createdAt,
		},
	}

	expectedAgents := []modelsHub.Agent{
		agentsDB[0],
		agentsDB[1],
	}
	expectedAgents[0].State = true
	expectedAgents[1].State = false

	tests := []struct {
		name       string
		expectErr  error
		expected   []modelsHub.Agent
		setupMocks func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			expected:  expectedAgents,
			setupMocks: func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository) {
				r.EXPECT().GetAllAgents(gomock.Any()).Return(agentsDB, nil).Times(1)
				s.EXPECT().GetStatus(expectedAgents[0].ID).Return(true).Times(1)
				s.EXPECT().GetStatus(expectedAgents[1].ID).Return(false).Times(1)
			},
		},
		{
			name:      "unhandled error",
			expectErr: sql.ErrConnDone,
			expected:  nil,
			setupMocks: func(r *mocks.MockagentRepository, s *mocks.MockagentStateRepository) {
				r.EXPECT().GetAllAgents(gomock.Any()).Return(nil, sql.ErrConnDone).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMockagentRepository(ctl)
			s := mocks.NewMockagentStateRepository(ctl)

			a := newAgentServiceMock(r, s)
			ctx := context.Background()

			tt.setupMocks(r, s)

			agents, err := a.GetAll(ctx)
			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expected, agents)
		})
	}
}

func TestAgentService_SetOnlineAgent(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name       string
		setupMocks func(s *mocks.MockagentStateRepository)
	}{
		{
			name: "successful",
			setupMocks: func(s *mocks.MockagentStateRepository) {
				s.EXPECT().SetOnline(id).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMockagentRepository(ctl)
			s := mocks.NewMockagentStateRepository(ctl)

			a := newAgentServiceMock(r, s)

			tt.setupMocks(s)

			a.SetOnlineAgent(id)
		})
	}
}

func TestAgentService_Register(t *testing.T) {
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	desc := "test_desc"
	errTest := errors.New("test")
	commands := models.CommandSpecs{
		"test_command": models.CommandSpec{
			Command: "test_command",
			Args:    []models.ArgSpec{},
			Label:   "test_label",
		},
	}

	createAgent := models.CreateAgent{
		Hostname:     "test_hostname",
		OS:           "test_os",
		CommandSpecs: commands,
	}

	createAgentDB := modelsHub.Agent{
		ID:       id,
		Hostname: createAgent.Hostname,
		OS:       createAgent.OS,
	}

	agentDB := modelsHub.Agent{
		ID:          id,
		Name:        "test_name",
		OS:          "test_os",
		Hostname:    "test_hostname",
		Description: &desc,
		State:       false,
		CreatedAt:   time.Now(),
	}

	expected := agentDB

	tests := []struct {
		name       string
		expectErr  error
		expected   *modelsHub.Agent
		genID      genIDMock
		setupMocks func(r *mocks.MockagentRepository, c *mocks.MockcommandSpecsRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			expected:  &expected,
			genID:     genIDMock{err: nil, id: id},
			setupMocks: func(r *mocks.MockagentRepository, c *mocks.MockcommandSpecsRepository) {
				r.EXPECT().CreateAgent(gomock.Any(), createAgentDB).Return(&agentDB, nil).Times(1)
				c.EXPECT().SetCommandSpecs(gomock.Any(), createAgentDB.ID, commands).Return(nil)
			},
		},
		{
			name:      "error gen id",
			expectErr: errTest,
			expected:  nil,
			genID:     genIDMock{err: errTest, id: id},
			setupMocks: func(r *mocks.MockagentRepository, c *mocks.MockcommandSpecsRepository) {
			},
		},
		{
			name:      "error create agent in repo",
			expectErr: errTest,
			expected:  nil,
			genID:     genIDMock{err: nil, id: id},
			setupMocks: func(r *mocks.MockagentRepository, c *mocks.MockcommandSpecsRepository) {
				r.EXPECT().CreateAgent(gomock.Any(), createAgentDB).Return(nil, errTest).Times(1)
			},
		},
		{
			name:      "error save command specification in repo",
			expectErr: errTest,
			expected:  nil,
			genID:     genIDMock{err: nil, id: id},
			setupMocks: func(r *mocks.MockagentRepository, c *mocks.MockcommandSpecsRepository) {
				r.EXPECT().CreateAgent(gomock.Any(), createAgentDB).Return(&agentDB, nil).Times(1)
				c.EXPECT().SetCommandSpecs(gomock.Any(), createAgentDB.ID, commands).Return(errTest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMockagentRepository(ctl)
			s := mocks.NewMockagentStateRepository(ctl)
			c := mocks.NewMockcommandSpecsRepository(ctl)

			a := AgentService{
				genID: tt.genID.gen,
				state: s,
				store: r,
				log:   &slog.Logger{},
				specs: c,
			}

			ctx := context.Background()

			tt.setupMocks(r, c)

			agent, err := a.Register(ctx, createAgent)
			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expected, agent)
		})
	}
}

func TestAgentService_EditNameAgent(t *testing.T) {
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	name := "name_test"
	errTest := errors.New("test")

	tests := []struct {
		name       string
		expectErr  error
		setupMocks func(r *mocks.MockagentRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			setupMocks: func(r *mocks.MockagentRepository) {
				r.EXPECT().EditNameByID(gomock.Any(), id, name).Return(nil).Times(1)
			},
		},
		{
			name:      "error agent not found",
			expectErr: ErrAgentNotFound,
			setupMocks: func(r *mocks.MockagentRepository) {
				r.EXPECT().EditNameByID(gomock.Any(), id, name).Return(store.ErrChangeData).Times(1)
			},
		},
		{
			name:      "unhandled error",
			expectErr: errTest,
			setupMocks: func(r *mocks.MockagentRepository) {
				r.EXPECT().EditNameByID(gomock.Any(), id, name).Return(errTest).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMockagentRepository(ctl)
			s := mocks.NewMockagentStateRepository(ctl)

			a := newAgentServiceMock(r, s)
			ctx := context.Background()

			tt.setupMocks(r)

			err := a.EditNameAgent(ctx, id, name)
			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}

func TestAgentService_EditDescAgent(t *testing.T) {
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	desc := "desc_test"
	errTest := errors.New("test")

	tests := []struct {
		name       string
		expectErr  error
		setupMocks func(r *mocks.MockagentRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			setupMocks: func(r *mocks.MockagentRepository) {
				r.EXPECT().EditDescByID(gomock.Any(), id, desc).Return(nil).Times(1)
			},
		},
		{
			name:      "error agent not found",
			expectErr: ErrAgentNotFound,
			setupMocks: func(r *mocks.MockagentRepository) {
				r.EXPECT().EditDescByID(gomock.Any(), id, desc).Return(store.ErrChangeData).Times(1)
			},
		},
		{
			name:      "unhandled error",
			expectErr: errTest,
			setupMocks: func(r *mocks.MockagentRepository) {
				r.EXPECT().EditDescByID(gomock.Any(), id, desc).Return(errTest).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMockagentRepository(ctl)
			s := mocks.NewMockagentStateRepository(ctl)

			a := newAgentServiceMock(r, s)
			ctx := context.Background()

			tt.setupMocks(r)

			err := a.EditDescAgent(ctx, id, desc)
			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}
