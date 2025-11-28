package tickets

import "time"

type Ticket struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Price     int64     `json:"price"`
	UserID    string    `json:"userId"`
	OrderID   *string   `json:"orderId,omitempty"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
}
