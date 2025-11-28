package payments

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
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

// StripeClient defines methods used from stripe client.
type StripeClient interface {
	CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (string, error)
	VerifyWebhookSignature(payload []byte, sigHeader string) (bool, error)
	SetWebhookSecret(secret string)
}

// Publisher is expected to publish domain events.
type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
}

type Service struct {
	db     *sql.DB
	pub    Publisher
	stripe StripeClient
	repo   Repository
}

func NewService(db *sql.DB, pub Publisher, stripe StripeClient) *Service {
	return &Service{db: db, pub: pub, stripe: stripe, repo: NewRepository(db)}
}

func (s *Service) CreateCharge(ctx context.Context, orderID string, amount int64, currency string) (*Payment, error) {
	// Validate order exists and is payable
	ord, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if ord.ID == "" {
		return nil, errors.New("order not found")
	}
	if ord.Status == "cancelled" {
		return nil, errors.New("order cancelled")
	}
	if ord.Price != amount {
		return nil, errors.New("amount mismatch")
	}

	// Create a payment intent (stubbed)
	stripeID, err := s.stripe.CreatePaymentIntent(amount, currency, map[string]string{"orderId": orderID})
	if err != nil {
		return nil, err
	}

	// Persist payment
	now := time.Now().UTC()
	pay := &Payment{
		ID:        "pay_" + orderID,
		OrderID:   orderID,
		Amount:    amount,
		Currency:  currency,
		Status:    "created",
		CreatedAt: now,
	}
	if err := s.repo.InsertPayment(ctx, pay.ID, orderID, amount, currency, stripeID); err != nil {
		return nil, err
	}

	// Publish payment:created event
	evt := events.PaymentCreatedData{ID: pay.ID, OrderID: orderID, StripeID: stripeID}
	b, _ := json.Marshal(evt)
	_ = s.pub.Publish(ctx, string(events.SubjectPaymentCreated), b)

	return pay, nil
}

func (s *Service) ProcessWebhook(ctx context.Context, payload []byte, sigHeader string) error {
	// Verify webhook signature
	valid, err := s.stripe.VerifyWebhookSignature(payload, sigHeader)
	if err != nil || !valid {
		return errors.New("invalid webhook signature")
	}

	// Parse the webhook event
	var event struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		// Extract payment intent details and update payment status
		var data struct {
			Object struct {
				ID       string            `json:"id"`
				Metadata map[string]string `json:"metadata"`
			} `json:"object"`
		}
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}

		// Update payment status to succeeded if needed
		// This is where you'd mark payment as complete in your system

	case "payment_intent.payment_failed":
		// Handle failed payment

	default:
		// Unhandled event type
	}

	return nil
}
