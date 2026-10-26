package repository

import (
	"context"
	"database/sql"
	"TaskMaster-api/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
}

type sqliteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) Create(ctx context.Context, u *model.User) error {
	query := `INSERT INTO users (name) VALUES (?)`
	_, err := r.db.ExecContext(ctx, query, u.Name)
	return err
}

func (r *sqliteUserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, name FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *sqliteUserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, name FROM users WHERE id = ?`
	var u model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name)
	if err != nil {
		return nil, err
	}
	return &u, nil
}