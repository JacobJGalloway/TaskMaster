package service

import (
	"context"
	"errors"
	"testing"

	"TaskMaster-api/model"
)

type fakeTaskRepository struct {
	createErr   error
	listErr     error
	completeErr error

	createdTasks     []*model.StudyTask
	tasksToReturn    []model.StudyTask
	completedTaskIDs []int64
}

func (f *fakeTaskRepository) Create(ctx context.Context, t *model.StudyTask) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.createdTasks = append(f.createdTasks, t)
	return nil
}

func (f *fakeTaskRepository) ListActiveByUser(ctx context.Context, userID int64) ([]model.StudyTask, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.tasksToReturn, nil
}

func (f *fakeTaskRepository) MarkComplete(ctx context.Context, taskID int64) error {
	if f.completeErr != nil {
		return f.completeErr
	}
	f.completedTaskIDs = append(f.completedTaskIDs, taskID)
	return nil
}

func TestTaskService_CreateTask(t *testing.T) {
	repo := &fakeTaskRepository{}
	svc := NewTaskService(repo)

	task := &model.StudyTask{Title: "Read", DurationMinutes: 15}
	if err := svc.CreateTask(context.Background(), task); err != nil {
		t.Fatalf("CreateTask() returned error: %v", err)
	}
	if len(repo.createdTasks) != 1 || repo.createdTasks[0] != task {
		t.Fatalf("expected the task to be forwarded to the repository, got %+v", repo.createdTasks)
	}
}

func TestTaskService_CreateTask_RepositoryError(t *testing.T) {
	wantErr := errors.New("db down")
	repo := &fakeTaskRepository{createErr: wantErr}
	svc := NewTaskService(repo)

	if err := svc.CreateTask(context.Background(), &model.StudyTask{}); !errors.Is(err, wantErr) {
		t.Fatalf("expected the repository error to propagate, got %v", err)
	}
}

func TestTaskService_ListActiveTasksByUser(t *testing.T) {
	want := []model.StudyTask{{ID: 1, Title: "Task 1"}}
	repo := &fakeTaskRepository{tasksToReturn: want}
	svc := NewTaskService(repo)

	got, err := svc.ListActiveTasksByUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListActiveTasksByUser() returned error: %v", err)
	}
	if len(got) != 1 || got[0].Title != "Task 1" {
		t.Fatalf("expected tasks to be passed through unchanged, got %+v", got)
	}
}

func TestTaskService_ListActiveTasksByUser_RepositoryError(t *testing.T) {
	wantErr := errors.New("db down")
	repo := &fakeTaskRepository{listErr: wantErr}
	svc := NewTaskService(repo)

	if _, err := svc.ListActiveTasksByUser(context.Background(), 1); !errors.Is(err, wantErr) {
		t.Fatalf("expected the repository error to propagate, got %v", err)
	}
}

func TestTaskService_CompleteTask(t *testing.T) {
	repo := &fakeTaskRepository{}
	svc := NewTaskService(repo)

	if err := svc.CompleteTask(context.Background(), 42); err != nil {
		t.Fatalf("CompleteTask() returned error: %v", err)
	}
	if len(repo.completedTaskIDs) != 1 || repo.completedTaskIDs[0] != 42 {
		t.Fatalf("expected task 42 to be marked complete, got %v", repo.completedTaskIDs)
	}
}
