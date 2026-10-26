package repository

import (
	"context"
	"database/sql"
	"TaskMaster-api/model"
)

type TaskRepository interface {
	Create(ctx context.Context, t *model.StudyTask) error
	ListActiveByUser(ctx context.Context, userID int64) ([]model.StudyTask, error)
	MarkComplete(ctx context.Context, taskID int64) error
}

type sqliteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(db *sql.DB) TaskRepository {
	return &sqliteTaskRepository{db: db}
}

func (r *sqliteTaskRepository) Create(ctx context.Context, t *model.StudyTask) error {
	query := `INSERT INTO study_tasks (user_id, title, duration_minutes, period_start, period_end, completed) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, t.UserID, t.Title, t.DurationMinutes, t.PeriodStart, t.PeriodEnd, t.Completed)
	return err
}

func (r *sqliteTaskRepository) ListActiveByUser(ctx context.Context, userID int64) ([]model.StudyTask, error) {
	query := `SELECT id, user_id, title, duration_minutes, period_start, period_end, completed, completed_at FROM study_tasks WHERE user_id = ? AND completed = 0`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.StudyTask
	for rows.Next() {
		var task model.StudyTask
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.DurationMinutes, &task.PeriodStart, &task.PeriodEnd, &task.Completed, &task.CompletedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *sqliteTaskRepository) MarkComplete(ctx context.Context, taskID int64) error {
	query := `UPDATE study_tasks SET completed = 1, completed_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, taskID)
	return err
}