package todo

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	appErrors "github.com/sborsh1kmusora/todo/internal/errors"
	"github.com/sborsh1kmusora/todo/internal/model"
	"github.com/sborsh1kmusora/todo/internal/service/todo/mocks"
)

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	task := model.Task{
		ID:    1,
		Title: "test",
	}

	tests := []struct {
		name      string
		task      model.Task
		setupMock func(repo *mocks.MockRepository)
		wantErr   error
	}{
		{
			name: "success",
			task: task,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Exists(ctx, task.ID).
					Return(false)

				repo.EXPECT().
					Create(ctx, task).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "task already exists",
			task: task,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Exists(ctx, task.ID).
					Return(true)
			},
			wantErr: appErrors.ErrTaskAlreadyExist,
		},
		{
			name: "empty task title",
			task: model.Task{
				ID:    2,
				Title: "",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Exists(ctx, 2).
					Return(false)
			},
			wantErr: appErrors.ErrTaskTitleIsEmpty,
		},
		{
			name: "invalid task",
			task: model.Task{
				ID:    -1,
				Title: "test",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Exists(ctx, -1).
					Return(false)
			},
			wantErr: appErrors.ErrInvalidTask,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			svc := New(repo)
			err := svc.Create(ctx, tt.task)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()

	task := model.Task{
		ID:    1,
		Title: "updated",
	}

	tests := []struct {
		name      string
		id        int
		task      model.Task
		setupMock func(repo *mocks.MockRepository)
		wantErr   error
	}{
		{
			name: "success",
			id:   1,
			task: task,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Update(ctx, 1, task).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "invalid task",
			id:      1,
			task:    model.Task{},
			wantErr: appErrors.ErrInvalidTask,
		},
		{
			name: "task not found",
			id:   1,
			task: task,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Update(ctx, 1, task).
					Return(appErrors.ErrTaskNotFound)
			},
			wantErr: appErrors.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			svc := New(repo)
			err := svc.Update(ctx, tt.id, tt.task)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	ctx := context.Background()

	task := model.Task{
		ID:    1,
		Title: "title",
	}

	tests := []struct {
		name      string
		id        int
		setupMock func(repo *mocks.MockRepository)
		wantTask  model.Task
		wantErr   error
	}{
		{
			name: "success",
			id:   1,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Get(ctx, 1).
					Return(task, nil)
			},
			wantTask: task,
			wantErr:  nil,
		},
		{
			name: "task not found",
			id:   2,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Get(ctx, 2).
					Return(model.Task{}, appErrors.ErrTaskNotFound)
			},
			wantErr: appErrors.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			svc := New(repo)
			got, err := svc.Get(ctx, tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil && got != tt.wantTask {
				t.Fatalf("expected task %+v, got %+v", tt.wantTask, got)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		id        int
		setupMock func(repo *mocks.MockRepository)
		wantErr   error
	}{
		{
			name: "success",
			id:   1,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Delete(ctx, 1).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "task not found",
			id:   2,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Delete(ctx, 2).
					Return(appErrors.ErrTaskNotFound)
			},
			wantErr: appErrors.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			svc := New(repo)
			err := svc.Delete(ctx, tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
