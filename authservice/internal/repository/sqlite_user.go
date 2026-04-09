package repository

import (
	"authservice/internal/model"
	"authservice/internal/telemetry"
	"context"
	"database/sql"
	"errors"
	"sync"
)

type UserRepo struct {
	DB *sql.DB
	mu sync.Mutex
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx, span := telemetry.Tracer().Start(ctx, "repo.create")
	defer span.End()

	var existing string

	err := r.DB.QueryRow(
		"SELECT email FROM users WHERE email = ?",
		user.Email,
	).Scan(&existing)
	if err == nil {
		return "", errors.New("user exists")
	}
	if _, err := r.DB.Exec(
		"INSERT INTO users (email, passwordhash, address) VALUES (?, ?, ?)",

		user.Email, user.PasswordHash, user.Address); err != nil {
		return "", err
	}

	return "", nil
}

func (r *UserRepo) Get(ctx context.Context, email string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx, span := telemetry.Tracer().Start(ctx, "repo.get")
	defer span.End()

	var user model.User
	row := r.DB.QueryRow("SELECT email, passwordhash, address FROM users WHERE email = ?", email)
	err := row.Scan(&user.Email, &user.PasswordHash, &user.Address)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
