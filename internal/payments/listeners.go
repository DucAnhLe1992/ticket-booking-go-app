package payments

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
)

// RegisterNATSListeners subscribes to order events to maintain local order state.
func RegisterNATSListeners(ctx context.Context, sub pubsub.Subscriber, repo Repository) error {
	if err := sub.Subscribe(string(events.SubjectOrderCreated), func(msg []byte) {
		var d events.OrderCreatedData
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("order:created unmarshal: %v", err)
			return
		}
		if err := repo.UpsertOrder(ctx, d.ID, int64(d.Ticket.Price), d.Status, d.UserID, d.Version); err != nil {
			log.Printf("order:created upsert: %v", err)
		}
	}); err != nil {
		return err
	}

	if err := sub.Subscribe(string(events.SubjectOrderCancelled), func(msg []byte) {
		var d events.OrderCancelledData
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("order:cancelled unmarshal: %v", err)
			return
		}
		if err := repo.CancelOrder(ctx, d.ID, d.Version); err != nil {
			log.Printf("order:cancelled cancel: %v", err)
		}
	}); err != nil {
		return err
	}

	return nil
}
