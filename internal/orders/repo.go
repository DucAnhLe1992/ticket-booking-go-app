package orders

import (
	"context"
	"database/sql"
	"time"
)

// Repository defines the data layer interface for orders.
type Repository interface {
	EnsureSchema(ctx context.Context) error
	CreateOrder(ctx context.Context, userID string, ticketID string, expiresAt time.Time) (*Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	ListOrdersByUser(ctx context.Context, userID string) ([]*Order, error)
	CancelOrder(ctx context.Context, id string, version int) error
	CompleteOrder(ctx context.Context, id string) error

	// Ticket replica management
	UpsertTicket(ctx context.Context, id string, title string, price int64, userID string, version int) error
	UpdateTicketReservation(ctx context.Context, ticketID string, orderID *string, expectedVersion int) error
	GetTicket(ctx context.Context, id string) (*Ticket, error)
	IsTicketReserved(ctx context.Context, ticketID string) (bool, error)
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) EnsureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;
		
		CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id TEXT NOT NULL,
			ticket_id TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'created',
			expires_at TIMESTAMPTZ NOT NULL,
			version INT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		
		CREATE TABLE IF NOT EXISTS orders_tickets (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			price BIGINT NOT NULL,
			user_id TEXT NOT NULL,
			order_id TEXT NULL,
			version INT NOT NULL DEFAULT 0,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	return err
}

func (r *repo) CreateOrder(ctx context.Context, userID string, ticketID string, expiresAt time.Time) (*Order, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, ticket_id, expires_at, status)
		VALUES ($1,$2,$3,'created')
		RETURNING id, user_id, status, expires_at, ticket_id, version, created_at
	`, userID, ticketID, expiresAt)

	var o Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Status, &o.ExpiresAt, &o.TicketID, &o.Version, &o.CreatedAt); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repo) GetOrder(ctx context.Context, id string) (*Order, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, status, expires_at, ticket_id, version, created_at FROM orders WHERE id=$1
	`, id)

	var o Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Status, &o.ExpiresAt, &o.TicketID, &o.Version, &o.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}

func (r *repo) ListOrdersByUser(ctx context.Context, userID string) ([]*Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, status, expires_at, ticket_id, version, created_at FROM orders WHERE user_id=$1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.ExpiresAt, &o.TicketID, &o.Version, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &o)
	}
	return out, nil
}

func (r *repo) CancelOrder(ctx context.Context, id string, version int) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE orders SET status='cancelled', version=version+1 WHERE id=$1 AND version=$2
	`, id, version)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repo) CompleteOrder(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE orders SET status='complete' WHERE id=$1
	`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repo) UpsertTicket(ctx context.Context, id string, title string, price int64, userID string, version int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO orders_tickets (id, title, price, user_id, version, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (id) DO UPDATE SET title=EXCLUDED.title, price=EXCLUDED.price, user_id=EXCLUDED.user_id, version=EXCLUDED.version, updated_at=EXCLUDED.updated_at
	`, id, title, price, userID, version, time.Now().UTC())
	return err
}

func (r *repo) UpdateTicketReservation(ctx context.Context, ticketID string, orderID *string, expectedVersion int) error {
	var oid interface{}
	if orderID != nil {
		oid = *orderID
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE orders_tickets SET order_id=$2, version=version+1, updated_at=$3 WHERE id=$1 AND version=$4
	`, ticketID, oid, time.Now().UTC(), expectedVersion)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repo) GetTicket(ctx context.Context, id string) (*Ticket, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, price, user_id, order_id, version FROM orders_tickets WHERE id=$1
	`, id)

	var t Ticket
	if err := row.Scan(&t.ID, &t.Title, &t.Price, &t.UserID, &t.OrderID, &t.Version); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *repo) IsTicketReserved(ctx context.Context, ticketID string) (bool, error) {
	var orderID *string
	err := r.db.QueryRowContext(ctx, `SELECT order_id FROM orders_tickets WHERE id=$1`, ticketID).Scan(&orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return orderID != nil, nil
}
