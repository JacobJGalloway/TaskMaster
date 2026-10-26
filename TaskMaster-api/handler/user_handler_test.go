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

type fakeUserService struct {
	createErr error
	listErr   error
	getErr    error

	createdUser   *model.User
	usersToReturn []model.User
	userToReturn  *model.User
}

func (f *fakeUserService) CreateUser(ctx context.Context, u *model.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.createdUser = u
	return nil
}

func (f *fakeUserService) GetUsers(ctx context.Context) ([]model.User, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.usersToReturn, nil
}

func (f *fakeUserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.userToReturn, nil
}

func TestUserHandler_Create(t *testing.T) {
	svc := &fakeUserService{}
	h := NewUserHandler(svc)

	body, _ := json.Marshal(model.User{Name: "Ada"})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if svc.createdUser == nil || svc.createdUser.Name != "Ada" {
		t.Fatalf("expected the decoded user to reach the service, got %+v", svc.createdUser)
	}

	var got model.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Name != "Ada" {
		t.Fatalf("expected the response body to echo the created user, got %+v", got)
	}
}

func TestUserHandler_Create_InvalidBody(t *testing.T) {
	svc := &fakeUserService{}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserHandler_Create_ServiceError(t *testing.T) {
	svc := &fakeUserService{createErr: errors.New("db down")}
	h := NewUserHandler(svc)

	body, _ := json.Marshal(model.User{Name: "Ada"})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	svc := &fakeUserService{usersToReturn: []model.User{{ID: 1, Name: "Ada"}}}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	h.GetUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []model.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Ada" {
		t.Fatalf("expected the response body to contain the service's users, got %+v", got)
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	svc := &fakeUserService{userToReturn: &model.User{ID: 1, Name: "Ada"}}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.GetUserByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got model.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Name != "Ada" {
		t.Fatalf("expected the response body to contain the service's user, got %+v", got)
	}
}

func TestUserHandler_GetUserByID_InvalidID(t *testing.T) {
	svc := &fakeUserService{}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.GetUserByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
