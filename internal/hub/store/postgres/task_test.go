package postgres

import (
	"context"
	"database/sql"
	"encoding/json/v2"
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

func newTaskRepo(q gen.Querier) *TaskRepository {
	return &TaskRepository{q: q}
}

func TestTaskRepository_GetTaskByID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	execMs := 50
	result := "task_result"
	m := map[string]string{"test": "test"}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	taskDB := gen.Task{
		ID:        id,
		AgentID:   uuid.New(),
		Status:    string(modelsHub.TaskStatusCompleted),
		Command:   "test_command",
		Args:      data,
		ExecMs:    &execMs,
		Result:    &result,
		CreatedAt: time.Now(),
	}

	expected, err := toTask(taskDB)
	require.NoError(t, err)

	q.EXPECT().GetTaskByID(ctx, id).Return(taskDB, nil).Times(1)
	task, err := repo.GetTaskByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, task)
}

func TestTaskRepository_GetTaskByID_ErrNotfound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()

	q.EXPECT().GetTaskByID(ctx, id).Return(gen.Task{}, sql.ErrNoRows).Times(1)
	task, err := repo.GetTaskByID(ctx, id)
	require.ErrorIs(t, err, store.ErrNotFound)
	require.Nil(t, task)
}

func TestTaskRepository_GetTaskByID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()

	q.EXPECT().GetTaskByID(ctx, id).Return(gen.Task{}, sql.ErrConnDone).Times(1)
	task, err := repo.GetTaskByID(ctx, id)
	require.Error(t, err)
	require.Nil(t, task)
}

func TestTaskRepository_NextTaskByAgentID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	agentID := uuid.New()
	execMs := 50
	result := "task_result"
	m := map[string]string{"test": "test"}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	taskDB := gen.Task{
		ID:        uuid.New(),
		AgentID:   agentID,
		Status:    string(modelsHub.TaskStatusPending),
		Command:   "test_command",
		Args:      data,
		ExecMs:    &execMs,
		Result:    &result,
		CreatedAt: time.Now(),
	}

	expected, err := toTask(taskDB)
	require.NoError(t, err)

	q.EXPECT().NextTaskByAgentID(ctx, agentID).Return(taskDB, nil).Times(1)
	task, err := repo.NextTaskByAgentID(ctx, agentID)
	require.NoError(t, err)
	require.Equal(t, expected, task)
}

func TestTaskRepository_NextTaskByAgentID_ErrNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	agentID := uuid.New()

	q.EXPECT().NextTaskByAgentID(ctx, agentID).Return(gen.Task{}, sql.ErrNoRows).Times(1)
	task, err := repo.NextTaskByAgentID(ctx, agentID)
	require.ErrorIs(t, err, store.ErrNotFound)
	require.Nil(t, task)
}

func TestTaskRepository_NextTaskByAgentID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	agentID := uuid.New()

	q.EXPECT().NextTaskByAgentID(ctx, agentID).Return(gen.Task{}, sql.ErrConnDone).Times(1)
	task, err := repo.NextTaskByAgentID(ctx, agentID)
	require.Error(t, err)
	require.Nil(t, task)
}

func TestTaskRepository_GetTasksByAgentID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	agentID := uuid.New()
	execMs := 50
	result := "task_result"
	m := map[string]string{"test": "test"}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	tasksDB := []gen.Task{
		{
			ID:        uuid.New(),
			AgentID:   agentID,
			Status:    string(modelsHub.TaskStatusPending),
			Command:   "test_command",
			Args:      data,
			ExecMs:    &execMs,
			Result:    &result,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			AgentID:   agentID,
			Status:    string(modelsHub.TaskStatusCompleted),
			Command:   "test_command_2",
			Args:      data,
			ExecMs:    &execMs,
			Result:    &result,
			CreatedAt: time.Now(),
		},
	}

	expected := make([]modelsHub.Task, 0, len(tasksDB))

	for _, taskDB := range tasksDB {
		task, err := toTask(taskDB)
		require.NoError(t, err)

		expected = append(expected, *task)
	}

	q.EXPECT().GetTasksByAgentID(ctx, agentID).Return(tasksDB, nil).Times(1)
	task, err := repo.GetTasksByAgentID(ctx, agentID)
	require.NoError(t, err)
	require.Equal(t, expected, task)
}

func TestTaskRepository_GetTasksByAgentID_EmptySlice(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	agentID := uuid.New()

	q.EXPECT().GetTasksByAgentID(ctx, agentID).Return([]gen.Task{}, nil).Times(1)
	task, err := repo.GetTasksByAgentID(ctx, agentID)
	require.NoError(t, err)
	require.Equal(t, []modelsHub.Task{}, task)
}

func TestTaskRepository_GetTasksByAgentID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	agentID := uuid.New()

	q.EXPECT().GetTasksByAgentID(ctx, agentID).Return([]gen.Task{}, sql.ErrConnDone).Times(1)
	task, err := repo.GetTasksByAgentID(ctx, agentID)
	require.Error(t, err)
	require.Nil(t, task)
}

func TestTaskRepository_CreateTask_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	execMs := 50
	result := "task_result"
	m := map[string]string{"test": "test"}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	taskDB := gen.Task{
		ID:        id,
		AgentID:   uuid.New(),
		Status:    string(modelsHub.TaskStatusCompleted),
		Command:   "test_command",
		Args:      data,
		ExecMs:    &execMs,
		Result:    &result,
		CreatedAt: time.Now(),
	}

	taskCreate := gen.CreateTaskParams{
		ID:      taskDB.ID,
		AgentID: taskDB.AgentID,
		Status:  taskDB.Status,
		Command: taskDB.Command,
		Args:    taskDB.Args,
		ExecMs:  taskDB.ExecMs,
		Result:  taskDB.Result,
	}

	expected, err := toTask(taskDB)
	require.NoError(t, err)

	q.EXPECT().CreateTask(ctx, taskCreate).Return(taskDB, nil).Times(1)
	task, err := repo.CreateTask(ctx, *expected)
	require.NoError(t, err)
	require.Equal(t, expected, task)
}

func TestTaskRepository_CreateTask_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	execMs := 50
	result := "task_result"
	m := map[string]string{"test": "test"}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	taskDB := gen.Task{
		ID:        id,
		AgentID:   uuid.New(),
		Status:    string(modelsHub.TaskStatusCompleted),
		Command:   "test_command",
		Args:      data,
		ExecMs:    &execMs,
		Result:    &result,
		CreatedAt: time.Now(),
	}

	taskCreate := gen.CreateTaskParams{
		ID:      taskDB.ID,
		AgentID: taskDB.AgentID,
		Status:  taskDB.Status,
		Command: taskDB.Command,
		Args:    taskDB.Args,
		ExecMs:  taskDB.ExecMs,
		Result:  taskDB.Result,
	}

	expected, err := toTask(taskDB)
	require.NoError(t, err)

	q.EXPECT().CreateTask(ctx, taskCreate).Return(gen.Task{}, sql.ErrConnDone).Times(1)
	task, err := repo.CreateTask(ctx, *expected)
	require.Error(t, err)
	require.Nil(t, task)
}

func TestTaskRepository_EditStatusByID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	status := modelsHub.TaskStatusPending

	editTask := gen.EditStatusTaskByIDParams{ID: id, Status: string(status)}

	q.EXPECT().EditStatusTaskByID(ctx, editTask).Return(&sqlResultMock{rows: 1}, nil)
	err := repo.EditStatusByID(ctx, id, status)
	require.NoError(t, err)
}

func TestTaskRepository_EditStatusByID_ErrNoChangeData(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	status := modelsHub.TaskStatusPending

	editTask := gen.EditStatusTaskByIDParams{ID: id, Status: string(status)}

	q.EXPECT().EditStatusTaskByID(ctx, editTask).Return(&sqlResultMock{rows: 0}, nil)
	err := repo.EditStatusByID(ctx, id, status)
	require.ErrorIs(t, err, store.ErrChangeData)
}

func TestTaskRepository_EditStatusByID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	status := modelsHub.TaskStatusPending

	editTask := gen.EditStatusTaskByIDParams{ID: id, Status: string(status)}

	q.EXPECT().EditStatusTaskByID(ctx, editTask).Return(&sqlResultMock{rows: 1}, sql.ErrConnDone)
	err := repo.EditStatusByID(ctx, id, status)
	require.Error(t, err)
}

func TestTaskRepository_EditResultByID_Successful(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	status := modelsHub.TaskStatusPending
	result := "test_result"
	execMs := 50

	editTask := gen.EditResultTaskByIDParams{ID: id, Status: string(status), Result: &result, ExecMs: &execMs}

	q.EXPECT().EditResultTaskByID(ctx, editTask).Return(&sqlResultMock{rows: 1}, nil)
	err := repo.EditResultByID(ctx, id, result, status, execMs)
	require.NoError(t, err)
}

func TestTaskRepository_EditResultByID_ErrNoChangeData(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	status := modelsHub.TaskStatusPending
	result := "test_result"
	execMs := 50

	editTask := gen.EditResultTaskByIDParams{ID: id, Status: string(status), Result: &result, ExecMs: &execMs}

	q.EXPECT().EditResultTaskByID(ctx, editTask).Return(&sqlResultMock{rows: 0}, nil)
	err := repo.EditResultByID(ctx, id, result, status, execMs)
	require.ErrorIs(t, err, store.ErrChangeData)
}

func TestTaskRepository_EditResultByID_UnhandledError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	q := mocks.NewMockQuerier(ctl)
	repo := newTaskRepo(q)

	ctx := context.Background()
	id := uuid.New()
	status := modelsHub.TaskStatusPending
	result := "test_result"
	execMs := 50

	editTask := gen.EditResultTaskByIDParams{ID: id, Status: string(status), Result: &result, ExecMs: &execMs}

	q.EXPECT().EditResultTaskByID(ctx, editTask).Return(&sqlResultMock{rows: 1}, sql.ErrConnDone)
	err := repo.EditResultByID(ctx, id, result, status, execMs)
	require.Error(t, err)
}
