package expiration

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
)

// RegisterNATSListeners subscribes to order:created events to schedule expirations.
func RegisterNATSListeners(ctx context.Context, sub pubsub.Subscriber, queue *ExpirationQueue) error {
	// Listen for order:created to schedule expiration jobs
	if err := sub.Subscribe(string(events.SubjectOrderCreated), func(msg []byte) {
		var d events.OrderCreatedData
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("order:created unmarshal: %v", err)
			return
		}

		log.Printf("Scheduling expiration for order %s at %v", d.ID, d.ExpiresAt)
		if err := queue.ScheduleOrderExpiration(d.ID, d.ExpiresAt); err != nil {
			log.Printf("Failed to schedule expiration: %v", err)
		}
	}); err != nil {
		return err
	}

	return nil
}
