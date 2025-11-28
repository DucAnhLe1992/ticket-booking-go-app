package orders

import "time"

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expiresAt"`
	TicketID  string    `json:"ticketId"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
}

type Ticket struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Price   int64   `json:"price"`
	UserID  string  `json:"userId"`
	OrderID *string `json:"orderId,omitempty"`
	Version int     `json:"version"`
}
