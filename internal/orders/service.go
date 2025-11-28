package orders

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
)

const (
	EXPIRATION_WINDOW_SECONDS = 900 // 15 minutes for orders to expire
)

// Service handles order business logic.
type Service struct {
	repo Repository
	pub  pubsub.Publisher
}

func NewService(repo Repository, pub pubsub.Publisher) *Service {
	return &Service{repo: repo, pub: pub}
}

// CreateOrder reserves a ticket and creates an order with expiration.
func (s *Service) CreateOrder(ctx context.Context, userID string, ticketID string) (*Order, error) {
	// Check if ticket exists and is not reserved
	ticket, err := s.repo.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, errors.New("ticket not found")
	}
	if ticket.OrderID != nil {
		return nil, errors.New("ticket already reserved")
	}

	// Create order with expiration
	expiresAt := time.Now().UTC().Add(time.Duration(EXPIRATION_WINDOW_SECONDS) * time.Second)
	order, err := s.repo.CreateOrder(ctx, userID, ticketID, expiresAt)
	if err != nil {
		return nil, err
	}

	// Reserve the ticket by setting its orderID
	if err := s.repo.UpdateTicketReservation(ctx, ticketID, &order.ID, ticket.Version); err != nil {
		// Rollback: cancel the order if reservation fails
		_ = s.repo.CancelOrder(ctx, order.ID, order.Version)
		return nil, errors.New("failed to reserve ticket")
	}

	// Publish order:created event
	if s.pub != nil {
		evt := events.OrderCreatedData{
			ID:        order.ID,
			Version:   order.Version,
			Status:    order.Status,
			UserID:    order.UserID,
			ExpiresAt: order.ExpiresAt,
			Ticket: events.OrderTicketDetail{
				ID:    ticket.ID,
				Price: ticket.Price,
			},
		}
		b, _ := json.Marshal(evt)
		_ = s.pub.Publish(ctx, string(events.SubjectOrderCreated), b)
	}

	return order, nil
}

// CancelOrder marks an order as cancelled and releases the ticket reservation.
func (s *Service) CancelOrder(ctx context.Context, orderID string, userID string) error {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}
	if order.UserID != userID {
		return errors.New("not authorized")
	}
	if order.Status == "cancelled" || order.Status == "complete" {
		return errors.New("cannot cancel order in status: " + order.Status)
	}

	// Cancel the order
	if err := s.repo.CancelOrder(ctx, order.ID, order.Version); err != nil {
		return err
	}

	// Release ticket reservation
	ticket, _ := s.repo.GetTicket(ctx, order.TicketID)
	if ticket != nil && ticket.OrderID != nil && *ticket.OrderID == order.ID {
		_ = s.repo.UpdateTicketReservation(ctx, ticket.ID, nil, ticket.Version)
	}

	// Publish order:cancelled event
	if s.pub != nil {
		evt := events.OrderCancelledData{
			ID:      order.ID,
			Version: order.Version + 1,
			Ticket: events.OrderTicketDetail{
				ID:    order.TicketID,
				Price: 0, // Price not needed for cancellation
			},
		}
		b, _ := json.Marshal(evt)
		_ = s.pub.Publish(ctx, string(events.SubjectOrderCancelled), b)
	}

	return nil
}

// GetOrder retrieves a single order.
func (s *Service) GetOrder(ctx context.Context, orderID string, userID string) (*Order, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != userID {
		return nil, errors.New("not authorized")
	}
	return order, nil
}

// ListOrders retrieves all orders for a user.
func (s *Service) ListOrders(ctx context.Context, userID string) ([]*Order, error) {
	return s.repo.ListOrdersByUser(ctx, userID)
}
