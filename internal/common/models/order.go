package models

// OrderStatus mirrors the TS enum used in the Node project.
type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusAwaiting  OrderStatus = "awaiting:payment"
	OrderStatusComplete  OrderStatus = "complete"
)

type Order struct {
	ID       string      `json:"id"`
	UserID   string      `json:"userId"`
	Status   OrderStatus `json:"status"`
	TicketID string      `json:"ticketId"`
	Version  int         `json:"version"`
}
