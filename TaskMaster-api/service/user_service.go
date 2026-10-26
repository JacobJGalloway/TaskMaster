package service

import (
	"context"
	"TaskMaster-api/model"
	"TaskMaster-api/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, u *model.User) error
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
}

type baseUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &baseUserService{repo: repo}
}

func (s *baseUserService) CreateUser(ctx context.Context, u *model.User) error {
	return s.repo.Create(ctx, u)
}

func (s *baseUserService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *baseUserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
