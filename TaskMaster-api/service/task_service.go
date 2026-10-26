package service

import (
	"context"
	"TaskMaster-api/model"
	"TaskMaster-api/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, t *model.StudyTask) error
	ListActiveTasksByUser(ctx context.Context, userID int64) ([]model.StudyTask, error)
	CompleteTask(ctx context.Context, taskID int64) error
}

type baseTaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &baseTaskService{repo: repo}
}

func (s *baseTaskService) CreateTask(ctx context.Context, t *model.StudyTask) error {
	return s.repo.Create(ctx, t)
}

func (s *baseTaskService) ListActiveTasksByUser(ctx context.Context, userID int64) ([]model.StudyTask, error) {
	return s.repo.ListActiveByUser(ctx, userID)
}

func (s *baseTaskService) CompleteTask(ctx context.Context, taskID int64) error {
	return s.repo.MarkComplete(ctx, taskID)
}
