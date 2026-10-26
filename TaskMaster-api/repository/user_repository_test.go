package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"TaskMaster-api/model"
)

const testSchema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS study_tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL REFERENCES users(id),
	title TEXT NOT NULL,
	duration_minutes INTEGER NOT NULL,
	period_start TEXT NOT NULL,
	period_end TEXT NOT NULL,
	completed INTEGER NOT NULL DEFAULT 0,
	completed_at TEXT
);`

// newTestDB opens a fresh in-memory SQLite database with the application
// schema applied, so each test runs against its own isolated database.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(testSchema); err != nil {
		t.Fatalf("failed to apply schema: %v", err)
	}
	return db
}

func seedUser(t *testing.T, db *sql.DB, name string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO users (name) VALUES (?)`, name)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("failed to read seeded user id: %v", err)
	}
	return id
}

func TestUserRepository_CreateAndGetUsers(t *testing.T) {
	db := newTestDB(t)
	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	if err := repo.Create(ctx, &model.User{Name: "Ada"}); err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	users, err := repo.GetUsers(ctx)
	if err != nil {
		t.Fatalf("GetUsers() returned error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Name != "Ada" {
		t.Fatalf("expected name %q, got %q", "Ada", users[0].Name)
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	if err := repo.Create(ctx, &model.User{Name: "Grace"}); err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	var id int64
	if err := db.QueryRow(`SELECT id FROM users WHERE name = ?`, "Grace").Scan(&id); err != nil {
		t.Fatalf("failed to read seeded user id: %v", err)
	}

	got, err := repo.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("GetUserByID() returned error: %v", err)
	}
	if got.Name != "Grace" {
		t.Fatalf("expected name %q, got %q", "Grace", got.Name)
	}
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewSQLiteUserRepository(db)

	_, err := repo.GetUserByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected an error for a missing user, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}
