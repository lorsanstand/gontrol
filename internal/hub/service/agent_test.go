package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"math/rand/v2"
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
	return NewAgentService(&slog.Logger{}, storage, state)
}

type genIDMock struct {
	err error
	id  uuid.UUID
}

func (g genIDMock) gen() (uuid.UUID, error) {
	return g.id, g.err
}

func TestAgentService_Get_Successful_Online(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.New()
	desc := "test_desc"
	state := true
	agentDB := modelsHub.Agent{
		ID:          id,
		Name:        "test_name",
		OS:          "test_os",
		Hostname:    "test_hostname",
		Description: &desc,
		State:       false,
		CreatedAt:   time.Now(),
	}

	expected := modelsHub.Agent{
		ID:          id,
		Name:        agentDB.Name,
		OS:          agentDB.OS,
		Hostname:    agentDB.Hostname,
		Description: agentDB.Description,
		State:       state,
		CreatedAt:   agentDB.CreatedAt,
	}

	r.EXPECT().GetAgentByID(ctx, id).Return(&agentDB, nil).Times(1)
	s.EXPECT().GetStatus(id).Return(state).Times(1)
	agent, err := a.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, *agent)
}

func TestAgentService_Get_Successful_Offline(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.New()
	desc := "test_desc"
	state := false
	agentDB := modelsHub.Agent{
		ID:          id,
		Name:        "test_name",
		OS:          "test_os",
		Hostname:    "test_hostname",
		Description: &desc,
		State:       false,
		CreatedAt:   time.Now(),
	}

	expected := modelsHub.Agent{
		ID:          id,
		Name:        agentDB.Name,
		OS:          agentDB.OS,
		Hostname:    agentDB.Hostname,
		Description: agentDB.Description,
		State:       state,
		CreatedAt:   agentDB.CreatedAt,
	}

	r.EXPECT().GetAgentByID(ctx, id).Return(&agentDB, nil).Times(1)
	s.EXPECT().GetStatus(id).Return(state).Times(1)
	agent, err := a.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, *agent)
}

func TestAgentService_Get_ErrAgentNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.New()

	r.EXPECT().GetAgentByID(ctx, id).Return(nil, store.ErrNotFound).Times(1)
	agent, err := a.Get(ctx, id)
	require.ErrorIs(t, err, ErrAgentNotFound)
	require.Nil(t, agent)
}

func TestAgentService_Get_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.New()

	r.EXPECT().GetAgentByID(ctx, id).Return(nil, sql.ErrConnDone).Times(1)
	agent, err := a.Get(ctx, id)
	require.Error(t, err)
	require.Nil(t, agent)
}

func TestAgentService_GetAll_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	desc := "test_desc"
	agentsDB := []modelsHub.Agent{
		{
			ID:          uuid.New(),
			Name:        "test_name",
			OS:          "test_os",
			Hostname:    "test_hostname",
			Description: &desc,
			State:       false,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "test_name_2",
			OS:          "test_os_2",
			Hostname:    "test_hostname_2",
			Description: &desc,
			State:       false,
			CreatedAt:   time.Now(),
		},
	}

	expecteds := make([]modelsHub.Agent, 0, len(agentsDB))

	for _, agentDB := range agentsDB {
		//nolint:gosec // G404
		state := rand.IntN(2) == 1

		expected := modelsHub.Agent{
			ID:          agentDB.ID,
			Name:        agentDB.Name,
			OS:          agentDB.OS,
			Hostname:    agentDB.Hostname,
			Description: agentDB.Description,
			State:       state,
			CreatedAt:   agentDB.CreatedAt,
		}

		s.EXPECT().GetStatus(expected.ID).Return(state).Times(1)

		expecteds = append(expecteds, expected)
	}

	r.EXPECT().GetAllAgents(ctx).Return(agentsDB, nil).Times(1)
	agent, err := a.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, expecteds, agent)
}

func TestAgentService_GetAll_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()

	r.EXPECT().GetAllAgents(ctx).Return(nil, sql.ErrConnDone).Times(1)
	agent, err := a.GetAll(ctx)
	require.Error(t, err)
	require.Nil(t, agent)
}

func TestAgentService_SetOnlineAgent(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)
	id := uuid.New()

	s.EXPECT().SetOnline(id).Times(1)
	a.SetOnlineAgent(id)
}

func TestAgentService_Register_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	idMock := genIDMock{nil, id}

	a := AgentService{
		genID: idMock.gen,
		state: s,
		store: r,
		log:   &slog.Logger{},
	}

	ctx := context.Background()
	desc := "test_desc"
	agentDB := modelsHub.Agent{
		ID:          id,
		Name:        "test_name",
		OS:          "test_os",
		Hostname:    "test_hostname",
		Description: &desc,
		State:       false,
		CreatedAt:   time.Now(),
	}

	createAgent := models.CreateAgent{Hostname: agentDB.Hostname, OS: agentDB.OS}
	createAgentDB := modelsHub.Agent{ID: id, Hostname: agentDB.Hostname, OS: agentDB.OS}

	expected := modelsHub.Agent{
		ID:          id,
		Name:        agentDB.Name,
		OS:          agentDB.OS,
		Hostname:    agentDB.Hostname,
		Description: agentDB.Description,
		State:       false,
		CreatedAt:   agentDB.CreatedAt,
	}

	r.EXPECT().CreateAgent(ctx, createAgentDB).Return(&agentDB, nil).Times(1)
	agent, err := a.Register(ctx, createAgent)
	require.NoError(t, err)
	require.Equal(t, expected, *agent)
}

func TestAgentService_Register_ErrGenID(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	errTest := errors.New("test")
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	idMock := genIDMock{errTest, id}

	a := AgentService{
		genID: idMock.gen,
		state: s,
		store: r,
		log:   &slog.Logger{},
	}

	ctx := context.Background()
	createAgent := models.CreateAgent{Hostname: "test_hostname", OS: "test_os"}

	agent, err := a.Register(ctx, createAgent)
	require.ErrorIs(t, err, errTest)
	require.Nil(t, agent)
}

func TestAgentService_Register_ErrCreateAgentInRepo(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	idMock := genIDMock{nil, id}

	a := AgentService{
		genID: idMock.gen,
		state: s,
		store: r,
		log:   &slog.Logger{},
	}

	errTest := errors.New("test")

	ctx := context.Background()
	createAgent := models.CreateAgent{Hostname: "test_hostname", OS: "test_os"}
	createAgentDB := modelsHub.Agent{ID: id, Hostname: createAgent.Hostname, OS: createAgent.OS}

	r.EXPECT().CreateAgent(ctx, createAgentDB).Return(nil, errTest).Times(1)
	agent, err := a.Register(ctx, createAgent)
	require.ErrorIs(t, err, errTest)
	require.Nil(t, agent)
}

func TestAgentService_EditNameAgent_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	name := "name_test"

	r.EXPECT().EditNameByID(ctx, id, name).Return(nil).Times(1)
	err := a.EditNameAgent(ctx, id, name)
	require.NoError(t, err)
}

func TestAgentService_EditNameAgent_ErrNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	name := "name_test"

	r.EXPECT().EditNameByID(ctx, id, name).Return(store.ErrChangeData).Times(1)
	err := a.EditNameAgent(ctx, id, name)
	require.ErrorIs(t, err, ErrAgentNotFound)
}

func TestAgentService_EditNameAgent_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	name := "name_test"
	errTest := errors.New("test")

	r.EXPECT().EditNameByID(ctx, id, name).Return(errTest).Times(1)
	err := a.EditNameAgent(ctx, id, name)
	require.ErrorIs(t, err, errTest)
}

func TestAgentService_EditDescAgent_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	desc := "desc_test"

	r.EXPECT().EditDescByID(ctx, id, desc).Return(nil).Times(1)
	err := a.EditDescAgent(ctx, id, desc)
	require.NoError(t, err)
}

func TestAgentService_EditDescAgent_ErrAgentNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	desc := "desc_test"

	r.EXPECT().EditDescByID(ctx, id, desc).Return(store.ErrChangeData).Times(1)
	err := a.EditDescAgent(ctx, id, desc)
	require.ErrorIs(t, err, ErrAgentNotFound)
}

func TestAgentService_EditDescAgent_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	r := mocks.NewMockagentRepository(ctl)
	s := mocks.NewMockagentStateRepository(ctl)
	defer ctl.Finish()

	a := newAgentServiceMock(r, s)

	ctx := context.Background()
	id := uuid.MustParse("018f3a3b-3a0b-7123-8456-426614174000")
	desc := "desc_test"
	errTest := errors.New("test")

	r.EXPECT().EditDescByID(ctx, id, desc).Return(errTest).Times(1)
	err := a.EditDescAgent(ctx, id, desc)
	require.ErrorIs(t, err, errTest)
}
