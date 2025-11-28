package auth

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	EnsureSchema(ctx context.Context) error
	Create(ctx context.Context, email, passwordHash string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

type userRepo struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) UserRepository { return &userRepo{db: db} }

func (r *userRepo) EnsureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`)
	if err == nil {
		// ensure pgcrypto for gen_random_uuid if available; ignore error if missing
		_, _ = r.db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	}
	return err
}

func (r *userRepo) Create(ctx context.Context, email, passwordHash string) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, password_hash, created_at
`, email, passwordHash)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, email, password_hash, created_at FROM users WHERE email=$1
`, email)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, email, password_hash, created_at FROM users WHERE id=$1
`, id)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
