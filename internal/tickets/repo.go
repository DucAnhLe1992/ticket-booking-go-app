package tickets

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	EnsureSchema(ctx context.Context) error
	Create(ctx context.Context, title string, price int64, userID string) (*Ticket, error)
	Get(ctx context.Context, id string) (*Ticket, error)
	List(ctx context.Context) ([]*Ticket, error)
	UpdateWithVersion(ctx context.Context, id string, expectedVersion int, title string, price int64, userID string) (*Ticket, error)
}

type repo struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repo{db: db} }

func (r *repo) EnsureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    price BIGINT NOT NULL,
    user_id TEXT NOT NULL,
    order_id TEXT NULL,
    version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE EXTENSION IF NOT EXISTS pgcrypto;
`)
	return err
}

func (r *repo) Create(ctx context.Context, title string, price int64, userID string) (*Ticket, error) {
	row := r.db.QueryRowContext(ctx, `
INSERT INTO tickets (title, price, user_id)
VALUES ($1,$2,$3)
RETURNING id, title, price, user_id, order_id, version, created_at
`, title, price, userID)
	var t Ticket
	if err := row.Scan(&t.ID, &t.Title, &t.Price, &t.UserID, &t.OrderID, &t.Version, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repo) Get(ctx context.Context, id string) (*Ticket, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, title, price, user_id, order_id, version, created_at FROM tickets WHERE id=$1`, id)
	var t Ticket
	if err := row.Scan(&t.ID, &t.Title, &t.Price, &t.UserID, &t.OrderID, &t.Version, &t.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *repo) List(ctx context.Context) ([]*Ticket, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, price, user_id, order_id, version, created_at FROM tickets ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Ticket
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.ID, &t.Title, &t.Price, &t.UserID, &t.OrderID, &t.Version, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, nil
}

func (r *repo) UpdateWithVersion(ctx context.Context, id string, expectedVersion int, title string, price int64, userID string) (*Ticket, error) {
	// OCC: update only if current version matches expected, then bump version
	res, err := r.db.ExecContext(ctx, `
UPDATE tickets SET title=$3, price=$4, version=version+1
WHERE id=$1 AND user_id=$2 AND version=$5
`, id, userID, title, price, expectedVersion)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, errors.New("conflict or not found")
	}

	return r.Get(ctx, id)
}
