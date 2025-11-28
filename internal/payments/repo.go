package payments

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	EnsureSchema(ctx context.Context) error
	UpsertOrder(ctx context.Context, id string, price int64, status string, userID string, version int) error
	CancelOrder(ctx context.Context, id string, version int) error
	GetOrder(ctx context.Context, id string) (order struct {
		ID      string
		Price   int64
		Status  string
		UserID  string
		Version int
	}, err error)
	InsertPayment(ctx context.Context, id string, orderID string, amount int64, currency string, stripeID string) error
}

type repo struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repo{db: db} }

func (r *repo) EnsureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    stripe_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS payments_orders (
    id TEXT PRIMARY KEY,
    price BIGINT NOT NULL,
    status TEXT NOT NULL,
    user_id TEXT NOT NULL,
    version INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`)
	return err
}

func (r *repo) UpsertOrder(ctx context.Context, id string, price int64, status string, userID string, version int) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO payments_orders (id, price, status, user_id, version, updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (id) DO UPDATE SET price=EXCLUDED.price, status=EXCLUDED.status, user_id=EXCLUDED.user_id, version=EXCLUDED.version, updated_at=EXCLUDED.updated_at
`, id, price, status, userID, version, time.Now().UTC())
	return err
}

func (r *repo) CancelOrder(ctx context.Context, id string, version int) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE payments_orders SET status='cancelled', version=$2, updated_at=$3 WHERE id=$1
`, id, version, time.Now().UTC())
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *repo) GetOrder(ctx context.Context, id string) (order struct {
	ID      string
	Price   int64
	Status  string
	UserID  string
	Version int
}, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, price, status, user_id, version FROM payments_orders WHERE id=$1`, id)
	err = row.Scan(&order.ID, &order.Price, &order.Status, &order.UserID, &order.Version)
	return
}

func (r *repo) InsertPayment(ctx context.Context, id string, orderID string, amount int64, currency string, stripeID string) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO payments (id, order_id, amount, currency, stripe_id) VALUES ($1,$2,$3,$4,$5)
`, id, orderID, amount, currency, stripeID)
	return err
}
