package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	gen "github.com/lorsanstand/gontrol/internal/hub/store/postgres/gen"
	"github.com/lorsanstand/gontrol/internal/hub/store/postgres/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newAgentRepo(q gen.Querier) *AgentRepository {
	return &AgentRepository{q: q}
}

type sqlResultMock struct {
	rows int64
}

func (s *sqlResultMock) RowsAffected() (int64, error) {
	return s.rows, nil
}

func (s *sqlResultMock) LastInsertId() (int64, error) {
	return 0, nil
}

func TestAgentRepository_GetAgentByID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	desc := "desc_test"
	ctx := context.Background()
	agentID := uuid.New()

	agentDB := gen.Agent{
		ID:          agentID,
		Os:          "os_test",
		Hostname:    "host_test",
		Name:        "name_test",
		Description: &desc,
		CreatedAt:   time.Now(),
	}

	expected, err := toAgent(agentDB)
	require.NoError(t, err)

	q.EXPECT().GetAgentByID(ctx, agentID).Return(agentDB, nil).Times(1)
	repo := newAgentRepo(q)

	agent, err := repo.GetAgentByID(ctx, agentID)
	require.NoError(t, err)
	require.EqualValues(t, expected, agent)
}

func TestAgentRepository_GetAgentByID_ErrNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	ctx := context.Background()
	agentID := uuid.New()

	q.EXPECT().GetAgentByID(ctx, agentID).Return(gen.Agent{}, sql.ErrNoRows).Times(1)
	repo := newAgentRepo(q)

	agent, err := repo.GetAgentByID(ctx, agentID)
	require.ErrorIs(t, err, store.ErrNotFound)
	require.Nil(t, agent)
}

func TestAgentRepository_GetAgentByID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	ctx := context.Background()
	agentID := uuid.New()

	q.EXPECT().GetAgentByID(ctx, agentID).Return(gen.Agent{}, sql.ErrConnDone).Times(1)
	repo := newAgentRepo(q)

	agent, err := repo.GetAgentByID(ctx, agentID)
	require.Error(t, err)
	require.Nil(t, agent)
}

func TestAgentRepository_GetAllAgents_EmptyAgents(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	ctx := context.Background()

	q.EXPECT().GetAllAgents(ctx).Return([]gen.Agent{}, nil).Times(1)
	repo := newAgentRepo(q)
	agents, err := repo.GetAllAgents(ctx)
	require.NoError(t, err)
	require.Equal(t, []modelsHub.Agent{}, agents)
}

func TestAgentRepository_GetAllAgents_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	ctx := context.Background()
	desc := "desc_test"

	agentsDB := []gen.Agent{
		{
			ID:          uuid.New(),
			Name:        "test_name",
			Os:          "test_os",
			Hostname:    "test_host",
			Description: &desc,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "test_name_2",
			Os:          "test_os_2",
			Hostname:    "test_host_2",
			Description: nil,
			CreatedAt:   time.Now(),
		},
	}

	expected := make([]modelsHub.Agent, 0, len(agentsDB))
	for _, agentDB := range agentsDB {
		agent, err := toAgent(agentDB)
		require.NoError(t, err)

		expected = append(expected, *agent)
	}

	q.EXPECT().GetAllAgents(ctx).Return(agentsDB, nil).Times(1)
	repo := newAgentRepo(q)
	agents, err := repo.GetAllAgents(ctx)
	require.NoError(t, err)
	require.Equal(t, expected, agents)
}

func TestAgentRepository_GetAllAgents_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	ctx := context.Background()

	q.EXPECT().GetAllAgents(ctx).Return([]gen.Agent{}, sql.ErrConnDone).Times(1)
	repo := newAgentRepo(q)
	agents, err := repo.GetAllAgents(ctx)
	require.Error(t, err)
	require.Nil(t, agents)
}

func TestAgentRepository_CreateAgent_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	desc := "desc_test"

	ctx := context.Background()
	createAgent := modelsHub.Agent{
		ID:          uuid.New(),
		Name:        "test_name",
		OS:          "test_os",
		Hostname:    "test_host",
		Description: &desc,
	}

	createAgentDB := gen.CreateAgentParams{
		ID:          createAgent.ID,
		Name:        createAgent.Name,
		Os:          createAgent.OS,
		Hostname:    createAgent.Hostname,
		Description: createAgent.Description}

	agentDB := gen.Agent{
		ID:          createAgent.ID,
		Name:        createAgent.Name,
		Os:          createAgent.OS,
		Hostname:    createAgent.Hostname,
		Description: createAgent.Description,
		CreatedAt:   time.Now(),
	}

	expected, err := toAgent(agentDB)
	require.NoError(t, err)

	q.EXPECT().CreateAgent(ctx, createAgentDB).Return(agentDB, nil).Times(1)
	repo := newAgentRepo(q)
	agent, err := repo.CreateAgent(ctx, createAgent)
	require.NoError(t, err)
	require.Equal(t, expected, agent)
}

func TestAgentRepository_CreateAgent_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	desc := "desc_test"

	ctx := context.Background()
	createAgent := modelsHub.Agent{
		ID:          uuid.New(),
		Name:        "test_name",
		OS:          "test_os",
		Hostname:    "test_host",
		Description: &desc,
	}

	createAgentDB := gen.CreateAgentParams{
		ID:          createAgent.ID,
		Name:        createAgent.Name,
		Os:          createAgent.OS,
		Hostname:    createAgent.Hostname,
		Description: createAgent.Description}

	q.EXPECT().CreateAgent(ctx, createAgentDB).Return(gen.Agent{}, sql.ErrConnDone).Times(1)
	repo := newAgentRepo(q)
	agent, err := repo.CreateAgent(ctx, createAgent)
	require.Error(t, err)
	require.Nil(t, agent)
}

func TestAgentRepository_EditDescByID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	desc := "desc_test"
	ctx := context.Background()
	id := uuid.New()

	editDescDB := gen.EditDescAgentByIDParams{ID: id, Description: &desc}

	q.EXPECT().EditDescAgentByID(ctx, editDescDB).Return(&sqlResultMock{rows: 1}, nil).Times(1)
	repo := newAgentRepo(q)
	err := repo.EditDescByID(ctx, id, desc)
	require.Nil(t, err)
}

func TestAgentRepository_EditDescByID_ErrChangeData(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	desc := "desc_test"
	ctx := context.Background()
	id := uuid.New()

	editDescDB := gen.EditDescAgentByIDParams{ID: id, Description: &desc}

	q.EXPECT().EditDescAgentByID(ctx, editDescDB).Return(&sqlResultMock{rows: 0}, nil).Times(1)
	repo := newAgentRepo(q)
	err := repo.EditDescByID(ctx, id, desc)
	require.ErrorIs(t, store.ErrChangeData, err)
}

func TestAgentRepository_EditDescByID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	desc := "desc_test"
	ctx := context.Background()
	id := uuid.New()

	editDescDB := gen.EditDescAgentByIDParams{ID: id, Description: &desc}

	q.EXPECT().EditDescAgentByID(ctx, editDescDB).Return(&sqlResultMock{rows: 1}, sql.ErrConnDone).Times(1)
	repo := newAgentRepo(q)
	err := repo.EditDescByID(ctx, id, desc)
	require.Error(t, err)
}

func TestAgentRepository_EditNameByID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	name := "test_name"
	ctx := context.Background()
	id := uuid.New()

	editNameDB := gen.EditNameAgentByIDParams{ID: id, Name: name}

	q.EXPECT().EditNameAgentByID(ctx, editNameDB).Return(&sqlResultMock{rows: 1}, nil).Times(1)
	repo := newAgentRepo(q)
	err := repo.EditNameByID(ctx, id, name)
	require.Nil(t, err)
}

func TestAgentRepository_EditNameByID_ErrChangeData(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	name := "test_name"
	ctx := context.Background()
	id := uuid.New()

	editNameDB := gen.EditNameAgentByIDParams{ID: id, Name: name}

	q.EXPECT().EditNameAgentByID(ctx, editNameDB).Return(&sqlResultMock{rows: 0}, nil).Times(1)
	repo := newAgentRepo(q)
	err := repo.EditNameByID(ctx, id, name)
	require.ErrorIs(t, store.ErrChangeData, err)
}

func TestAgentRepository_EditNameByID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)

	name := "test_name"
	ctx := context.Background()
	id := uuid.New()

	editNameDB := gen.EditNameAgentByIDParams{ID: id, Name: name}

	q.EXPECT().EditNameAgentByID(ctx, editNameDB).Return(&sqlResultMock{rows: 1}, sql.ErrConnDone).Times(1)
	repo := newAgentRepo(q)
	err := repo.EditNameByID(ctx, id, name)
	require.Error(t, err)
}
