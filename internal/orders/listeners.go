package orders

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
)

// RegisterNATSListeners subscribes to ticket and payment events.
func RegisterNATSListeners(ctx context.Context, sub pubsub.Subscriber, repo Repository) error {
	// Listen for ticket:created to replicate tickets locally
	if err := sub.Subscribe(string(events.SubjectTicketCreated), func(msg []byte) {
		var d events.TicketCreatedData
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("ticket:created unmarshal: %v", err)
			return
		}
		if err := repo.UpsertTicket(ctx, d.ID, d.Title, d.Price, d.UserID, d.Version); err != nil {
			log.Printf("ticket:created upsert: %v", err)
		}
	}); err != nil {
		return err
	}

	// Listen for ticket:updated to keep tickets in sync
	if err := sub.Subscribe(string(events.SubjectTicketUpdated), func(msg []byte) {
		var d events.TicketUpdatedData
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("ticket:updated unmarshal: %v", err)
			return
		}
		if err := repo.UpsertTicket(ctx, d.ID, d.Title, d.Price, d.UserID, d.Version); err != nil {
			log.Printf("ticket:updated upsert: %v", err)
		}
	}); err != nil {
		return err
	}

	// Listen for expiration:complete to cancel expired orders
	if err := sub.Subscribe(string(events.SubjectExpirationComplete), func(msg []byte) {
		var d events.ExpirationCompleteData
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("expiration:complete unmarshal: %v", err)
			return
		}
		order, err := repo.GetOrder(ctx, d.OrderID)
		if err != nil || order == nil {
			log.Printf("expiration:complete get order: %v", err)
			return
		}
		// Only cancel if still in created state (not already complete)
		if order.Status == "created" {
			if err := repo.CancelOrder(ctx, order.ID, order.Version); err != nil {
				log.Printf("expiration:complete cancel: %v", err)
			} else {
				log.Printf("order %s expired and cancelled", order.ID)
			}
		}
	}); err != nil {
		return err
	}

	// Listen for payment:created to mark order complete
	if err := sub.Subscribe(string(events.SubjectPaymentCreated), func(msg []byte) {
		var d events.PaymentCreatedData
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("payment:created unmarshal: %v", err)
			return
		}
		if err := repo.CompleteOrder(ctx, d.OrderID); err != nil {
			log.Printf("payment:created complete: %v", err)
		}
	}); err != nil {
		return err
	}

	return nil
}
