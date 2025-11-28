package payments

import (
	"context"
	"database/sql"
	"time"
)

// Payment is a very small payment record used for responses.
type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// StripeClientStub defines methods used from stripe client.
type StripeClient interface {
	CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (string, error)
	VerifyWebhookSignature(payload []byte, sigHeader string) (bool, error)
}

// Publisher is expected to publish domain events.
type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
}

type Service struct {
	db     *sql.DB
	pub    Publisher
	stripe StripeClient
}

func NewService(db *sql.DB, pub Publisher, stripe StripeClient) *Service {
	return &Service{db: db, pub: pub, stripe: stripe}
}

func (s *Service) CreateCharge(ctx context.Context, orderID string, amount int64, currency string) (*Payment, error) {
	// TODO: implement real Stripe creation, DB transaction, idempotency and publish events.
	now := time.Now().UTC()
	return &Payment{
		ID:       "pay_" + orderID,
		OrderID:  orderID,
		Amount:   amount,
		Currency: currency,
		Status:   "created",
		CreatedAt: now,
	}, nil
}

func (s *Service) ProcessWebhook(ctx context.Context, payload []byte) error {
	// TODO: verify signature and process event idempotently
	return nil
}
