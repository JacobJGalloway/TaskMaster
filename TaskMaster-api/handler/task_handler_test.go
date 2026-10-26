package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"TaskMaster-api/model"
)

type fakeTaskService struct {
	createErr   error
	listErr     error
	completeErr error

	createdTask   *model.StudyTask
	tasksToReturn []model.StudyTask
	completedID   int64
}

func (f *fakeTaskService) CreateTask(ctx context.Context, t *model.StudyTask) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.createdTask = t
	return nil
}

func (f *fakeTaskService) ListActiveTasksByUser(ctx context.Context, userID int64) ([]model.StudyTask, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.tasksToReturn, nil
}

func (f *fakeTaskService) CompleteTask(ctx context.Context, taskID int64) error {
	if f.completeErr != nil {
		return f.completeErr
	}
	f.completedID = taskID
	return nil
}

func TestTaskHandler_Create(t *testing.T) {
	svc := &fakeTaskService{}
	h := NewTaskHandler(svc)

	body, _ := json.Marshal(model.StudyTask{Title: "Read", DurationMinutes: 30})
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if svc.createdTask == nil || svc.createdTask.Title != "Read" {
		t.Fatalf("expected the decoded task to reach the service, got %+v", svc.createdTask)
	}
}

func TestTaskHandler_Create_InvalidBody(t *testing.T) {
	svc := &fakeTaskService{}
	h := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestTaskHandler_Create_ServiceError(t *testing.T) {
	svc := &fakeTaskService{createErr: errors.New("db down")}
	h := NewTaskHandler(svc)

	body, _ := json.Marshal(model.StudyTask{Title: "Read"})
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestTaskHandler_ListActiveByUser(t *testing.T) {
	svc := &fakeTaskService{tasksToReturn: []model.StudyTask{{ID: 1, Title: "Task 1"}}}
	h := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/1/tasks", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.ListActiveByUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []model.StudyTask
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(got) != 1 || got[0].Title != "Task 1" {
		t.Fatalf("expected the response body to contain the service's tasks, got %+v", got)
	}
}

func TestTaskHandler_ListActiveByUser_InvalidID(t *testing.T) {
	svc := &fakeTaskService{}
	h := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/abc/tasks", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.ListActiveByUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestTaskHandler_MarkComplete(t *testing.T) {
	svc := &fakeTaskService{}
	h := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/tasks/7/complete", nil)
	req.SetPathValue("id", "7")
	rec := httptest.NewRecorder()

	h.MarkComplete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
	if svc.completedID != 7 {
		t.Fatalf("expected task 7 to be completed, got %d", svc.completedID)
	}
}

func TestTaskHandler_MarkComplete_InvalidID(t *testing.T) {
	svc := &fakeTaskService{}
	h := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/tasks/abc/complete", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.MarkComplete(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
