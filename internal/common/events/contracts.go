package events

import "time"

// TicketCreatedEvent
type TicketCreatedData struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Price   int64  `json:"price"`
	UserID  string `json:"userId"`
	Version int    `json:"version"`
}

// TicketUpdatedEvent
type TicketUpdatedData struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Price   int64   `json:"price"`
	UserID  string  `json:"userId"`
	OrderID *string `json:"orderId,omitempty"`
	Version int     `json:"version"`
}

// OrderCreatedEvent
type OrderCreatedData struct {
	ID        string            `json:"id"`
	Version   int               `json:"version"`
	Status    string            `json:"status"`
	UserID    string            `json:"userId"`
	ExpiresAt time.Time         `json:"expiresAt"`
	Ticket    OrderTicketDetail `json:"ticket"`
}

type OrderTicketDetail struct {
	ID    string `json:"id"`
	Price int64  `json:"price"`
}

// OrderCancelledEvent
type OrderCancelledData struct {
	ID      string            `json:"id"`
	Version int               `json:"version"`
	Ticket  OrderTicketDetail `json:"ticket"`
}

// ExpirationCompleteEvent
type ExpirationCompleteData struct {
	OrderID string `json:"orderId"`
}

// PaymentCreatedEvent
type PaymentCreatedData struct {
	ID       string `json:"id"`
	OrderID  string `json:"orderId"`
	StripeID string `json:"stripeId"`
}

// UserCreatedEvent
type UserCreatedData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
