package repository

import (
	"context"
	"database/sql"
	"testing"

	"TaskMaster-api/model"
)

func TestTaskRepository_Create(t *testing.T) {
	db := newTestDB(t)
	userID := seedUser(t, db, "Ada")
	repo := NewSQLiteTaskRepository(db)

	task := &model.StudyTask{
		UserID:          userID,
		Title:           "Read chapter 1",
		DurationMinutes: 30,
		PeriodStart:     "2026-01-01T09:00",
		PeriodEnd:       "2026-01-01T09:30",
	}

	if err := repo.Create(context.Background(), task); err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM study_tasks WHERE user_id = ?`, userID).Scan(&count); err != nil {
		t.Fatalf("failed to count tasks: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 task to be inserted, got %d", count)
	}
}

func TestTaskRepository_ListActiveByUser(t *testing.T) {
	db := newTestDB(t)
	userID := seedUser(t, db, "Ada")
	otherUserID := seedUser(t, db, "Grace")
	repo := NewSQLiteTaskRepository(db)
	ctx := context.Background()

	tasks := []*model.StudyTask{
		{UserID: userID, Title: "Active task", DurationMinutes: 20, PeriodStart: "s", PeriodEnd: "e"},
		{UserID: userID, Title: "Done task", DurationMinutes: 10, PeriodStart: "s", PeriodEnd: "e", Completed: true},
		{UserID: otherUserID, Title: "Not mine", DurationMinutes: 15, PeriodStart: "s", PeriodEnd: "e"},
	}
	for _, task := range tasks {
		if err := repo.Create(ctx, task); err != nil {
			t.Fatalf("failed to seed task: %v", err)
		}
	}

	got, err := repo.ListActiveByUser(ctx, userID)
	if err != nil {
		t.Fatalf("ListActiveByUser() returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 active task, got %d", len(got))
	}
	if got[0].Title != "Active task" {
		t.Fatalf("expected %q, got %q", "Active task", got[0].Title)
	}
}

func TestTaskRepository_MarkComplete(t *testing.T) {
	db := newTestDB(t)
	userID := seedUser(t, db, "Ada")
	repo := NewSQLiteTaskRepository(db)
	ctx := context.Background()

	task := &model.StudyTask{UserID: userID, Title: "Task", DurationMinutes: 20, PeriodStart: "s", PeriodEnd: "e"}
	if err := repo.Create(ctx, task); err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	var taskID int64
	if err := db.QueryRow(`SELECT id FROM study_tasks WHERE user_id = ?`, userID).Scan(&taskID); err != nil {
		t.Fatalf("failed to read seeded task id: %v", err)
	}

	if err := repo.MarkComplete(ctx, taskID); err != nil {
		t.Fatalf("MarkComplete() returned error: %v", err)
	}

	var completed int
	var completedAt sql.NullString
	if err := db.QueryRow(`SELECT completed, completed_at FROM study_tasks WHERE id = ?`, taskID).Scan(&completed, &completedAt); err != nil {
		t.Fatalf("failed to read task: %v", err)
	}
	if completed != 1 {
		t.Fatalf("expected completed = 1, got %d", completed)
	}
	if !completedAt.Valid || completedAt.String == "" {
		t.Fatal("expected completed_at to be set")
	}

	active, err := repo.ListActiveByUser(ctx, userID)
	if err != nil {
		t.Fatalf("ListActiveByUser() returned error: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("expected the completed task to drop out of the active list, got %d", len(active))
	}
}
