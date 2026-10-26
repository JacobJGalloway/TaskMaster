package service

import (
	"context"
	"errors"
	"testing"

	"TaskMaster-api/model"
)

type fakeUserRepository struct {
	createErr error
	listErr   error
	getErr    error

	createdUsers []*model.User
	usersToReturn []model.User
	userToReturn  *model.User
}

func (f *fakeUserRepository) Create(ctx context.Context, u *model.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.createdUsers = append(f.createdUsers, u)
	return nil
}

func (f *fakeUserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.usersToReturn, nil
}

func (f *fakeUserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.userToReturn, nil
}

func TestUserService_CreateUser(t *testing.T) {
	repo := &fakeUserRepository{}
	svc := NewUserService(repo)

	user := &model.User{Name: "Ada"}
	if err := svc.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("CreateUser() returned error: %v", err)
	}
	if len(repo.createdUsers) != 1 || repo.createdUsers[0] != user {
		t.Fatalf("expected the user to be forwarded to the repository, got %+v", repo.createdUsers)
	}
}

func TestUserService_CreateUser_RepositoryError(t *testing.T) {
	wantErr := errors.New("db down")
	repo := &fakeUserRepository{createErr: wantErr}
	svc := NewUserService(repo)

	if err := svc.CreateUser(context.Background(), &model.User{}); !errors.Is(err, wantErr) {
		t.Fatalf("expected the repository error to propagate, got %v", err)
	}
}

func TestUserService_GetUsers(t *testing.T) {
	want := []model.User{{ID: 1, Name: "Ada"}}
	repo := &fakeUserRepository{usersToReturn: want}
	svc := NewUserService(repo)

	got, err := svc.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("GetUsers() returned error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Ada" {
		t.Fatalf("expected users to be passed through unchanged, got %+v", got)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	want := &model.User{ID: 1, Name: "Ada"}
	repo := &fakeUserRepository{userToReturn: want}
	svc := NewUserService(repo)

	got, err := svc.GetUserByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUserByID() returned error: %v", err)
	}
	if got != want {
		t.Fatalf("expected the repository's user to be returned unchanged, got %+v", got)
	}
}

func TestUserService_GetUserByID_RepositoryError(t *testing.T) {
	wantErr := errors.New("not found")
	repo := &fakeUserRepository{getErr: wantErr}
	svc := NewUserService(repo)

	if _, err := svc.GetUserByID(context.Background(), 1); !errors.Is(err, wantErr) {
		t.Fatalf("expected the repository error to propagate, got %v", err)
	}
}
