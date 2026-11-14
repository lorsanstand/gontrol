package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	modelsHub "github.com/lorsanstand/gontrol/internal/hub/models"
	"github.com/lorsanstand/gontrol/internal/hub/service/mocks"
	"github.com/lorsanstand/gontrol/internal/hub/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTaskServiceMock(r tasksRepository) *TaskService {
	return &TaskService{log: &slog.Logger{}, store: r}
}

func newDummyTask(id uuid.UUID, agentID uuid.UUID) modelsHub.Task {
	return modelsHub.Task{
		ID:      id,
		AgentID: agentID,
		Status:  modelsHub.TaskStatusPending,
		Command: "test_command",
	}
}

func TestTaskService_GetByID(t *testing.T) {
	id := uuid.New()
	agentID := uuid.New()
	task := newDummyTask(id, agentID)
	errTest := errors.New("test")

	tests := []struct {
		name       string
		expectErr  error
		expected   *modelsHub.Task
		setupMocks func(r *mocks.MocktasksRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			expected:  &task,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().GetTaskByID(gomock.Any(), id).Return(&task, nil).Times(1)
			},
		},
		{
			name:      "error not found",
			expectErr: ErrTaskNotFound,
			expected:  nil,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().GetTaskByID(gomock.Any(), id).Return(nil, store.ErrNotFound).Times(1)
			},
		},
		{
			name:      "unhandled error",
			expectErr: errTest,
			expected:  nil,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().GetTaskByID(gomock.Any(), id).Return(nil, errTest).Times(1)
			},
		},
	}

	for _, tt := range tests {
		runner := func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMocktasksRepository(ctl)
			s := newTaskServiceMock(r)

			ctx := context.Background()

			tt.setupMocks(r)

			taskTest, err := s.GetByID(ctx, id)
			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expected, taskTest)
		}

		t.Run(tt.name, runner)
	}
}

func TestTaskService_GetByAgentID(t *testing.T) {
	agentID := uuid.New()
	task1 := newDummyTask(uuid.New(), agentID)
	task2 := newDummyTask(uuid.New(), agentID)
	errTest := errors.New("test")

	tests := []struct {
		name       string
		expectErr  error
		expected   []modelsHub.Task
		setupMocks func(r *mocks.MocktasksRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			expected:  []modelsHub.Task{task1, task2},
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().GetTasksByAgentID(gomock.Any(), agentID).Return([]modelsHub.Task{task1, task2}, nil).Times(1)
			},
		},
		{
			name:      "unhandled error",
			expectErr: errTest,
			expected:  []modelsHub.Task{},
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().GetTasksByAgentID(gomock.Any(), agentID).Return([]modelsHub.Task{}, errTest).Times(1)
			},
		},
	}

	for _, tt := range tests {
		runner := func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMocktasksRepository(ctl)
			s := newTaskServiceMock(r)

			ctx := context.Background()

			tt.setupMocks(r)

			taskTest, err := s.GetByAgentID(ctx, agentID)
			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expected, taskTest)
		}

		t.Run(tt.name, runner)
	}
}

func TestTaskService_Pop(t *testing.T) {
	id := uuid.New()
	agentID := uuid.New()
	task := newDummyTask(id, agentID)
	errTest := errors.New("test")
	status := modelsHub.TaskStatusProcessing

	tests := []struct {
		name       string
		expectErr  error
		expected   *modelsHub.Task
		setupMocks func(r *mocks.MocktasksRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			expected:  &task,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().NextTaskByAgentID(gomock.Any(), id).Return(&task, nil).Times(1)
				r.EXPECT().EditStatusByID(gomock.Any(), id, status).Return(nil).Times(1)
			},
		},
		{
			name:      "error not found",
			expectErr: ErrTaskNotFound,
			expected:  nil,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().NextTaskByAgentID(gomock.Any(), id).Return(nil, store.ErrNotFound).Times(1)
			},
		},
		{
			name:      "unhandled error get task",
			expectErr: errTest,
			expected:  nil,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().NextTaskByAgentID(gomock.Any(), id).Return(nil, errTest).Times(1)
			},
		},
		{
			name:      "unhandled error set status",
			expectErr: errTest,
			expected:  nil,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().NextTaskByAgentID(gomock.Any(), id).Return(&task, nil).Times(1)
				r.EXPECT().EditStatusByID(gomock.Any(), id, status).Return(errTest).Times(1)
			},
		},
	}

	for _, tt := range tests {
		runner := func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMocktasksRepository(ctl)
			s := newTaskServiceMock(r)

			ctx := context.Background()

			tt.setupMocks(r)

			taskTest, err := s.Pop(ctx, id)
			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expected, taskTest)
		}

		t.Run(tt.name, runner)
	}
}

func TestTaskService_SaveResult(t *testing.T) {
	result := "test_result"
	execMs := 50
	id := uuid.New()
	errTest := errors.New("test")

	tests := []struct {
		name       string
		expectErr  error
		setupMocks func(r *mocks.MocktasksRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().EditResultByID(gomock.Any(), id, result, modelsHub.TaskStatusCompleted, execMs).Return(nil).Times(1)
			},
		},
		{
			name:      "error not found",
			expectErr: ErrTaskNotFound,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().EditResultByID(gomock.Any(), id, result, modelsHub.TaskStatusCompleted, execMs).Return(store.ErrChangeData).Times(1)
			},
		},
		{
			name:      "unhandled error",
			expectErr: errTest,
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().EditResultByID(gomock.Any(), id, result, modelsHub.TaskStatusCompleted, execMs).Return(errTest).Times(1)
			},
		},
	}

	for _, tt := range tests {
		runner := func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMocktasksRepository(ctl)
			s := newTaskServiceMock(r)

			ctx := context.Background()

			tt.setupMocks(r)

			err := s.SaveResult(ctx, id, result, execMs)
			require.ErrorIs(t, err, tt.expectErr)
		}

		t.Run(tt.name, runner)
	}
}

func TestTaskService_Add(t *testing.T) {
	id := uuid.New()
	agentID := uuid.New()
	taskDB := newDummyTask(id, agentID)
	taskDB.ExecMs = nil
	taskDB.Result = nil
	taskCreate := modelsHub.CreateTask{AgentID: taskDB.AgentID, Command: taskDB.Command, Args: taskDB.Args}
	errTest := errors.New("test")

	tests := []struct {
		name       string
		expectErr  error
		expected   *modelsHub.Task
		genID      genIDMock
		setupMocks func(r *mocks.MocktasksRepository)
	}{
		{
			name:      "successful",
			expectErr: nil,
			expected:  &taskDB,
			genID:     genIDMock{err: nil, id: id},
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().CreateTask(gomock.Any(), taskDB).Return(&taskDB, nil).Times(1)
			},
		},
		{
			name:      "generate id",
			expectErr: errTest,
			expected:  nil,
			genID:     genIDMock{err: errTest, id: id},
			setupMocks: func(r *mocks.MocktasksRepository) {
			},
		},
		{
			name:      "unhandled error",
			expectErr: errTest,
			expected:  nil,
			genID:     genIDMock{err: nil, id: id},
			setupMocks: func(r *mocks.MocktasksRepository) {
				r.EXPECT().CreateTask(gomock.Any(), taskDB).Return(nil, errTest).Times(1)
			},
		},
	}

	for _, tt := range tests {
		runner := func(t *testing.T) {
			ctl := gomock.NewController(t)
			r := mocks.NewMocktasksRepository(ctl)
			s := TaskService{log: &slog.Logger{}, store: r, genID: tt.genID.gen}

			ctx := context.Background()

			tt.setupMocks(r)

			taskTest, err := s.Add(ctx, taskCreate)
			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expected, taskTest)
		}

		t.Run(tt.name, runner)
	}
}
