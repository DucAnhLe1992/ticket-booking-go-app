package expiration

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/events"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
)

const (
	TypeOrderExpiration = "order:expiration"
)

// ExpirationQueue manages delayed job processing for order expirations.
type ExpirationQueue struct {
	client *asynq.Client
	pub    pubsub.Publisher
}

// NewExpirationQueue creates a new expiration queue client.
func NewExpirationQueue(redisAddr string, pub pubsub.Publisher) *ExpirationQueue {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &ExpirationQueue{
		client: client,
		pub:    pub,
	}
}

// Close closes the queue client.
func (q *ExpirationQueue) Close() error {
	return q.client.Close()
}

// ScheduleOrderExpiration schedules an expiration job for an order.
func (q *ExpirationQueue) ScheduleOrderExpiration(orderID string, expiresAt time.Time) error {
	payload, err := json.Marshal(map[string]string{"orderId": orderID})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeOrderExpiration, payload)
	delay := time.Until(expiresAt)
	if delay < 0 {
		delay = 0 // Process immediately if already expired
	}

	_, err = q.client.Enqueue(task, asynq.ProcessIn(delay))
	return err
}

// ExpirationWorker processes expiration jobs.
type ExpirationWorker struct {
	server *asynq.Server
	pub    pubsub.Publisher
}

// NewExpirationWorker creates a worker that processes expiration jobs.
func NewExpirationWorker(redisAddr string, pub pubsub.Publisher, concurrency int) *ExpirationWorker {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	return &ExpirationWorker{
		server: srv,
		pub:    pub,
	}
}

// Start starts processing jobs.
func (w *ExpirationWorker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeOrderExpiration, w.handleOrderExpiration)
	return w.server.Start(mux)
}

// Stop gracefully stops the worker.
func (w *ExpirationWorker) Stop() {
	w.server.Stop()
	w.server.Shutdown()
}

// handleOrderExpiration processes an order expiration job.
func (w *ExpirationWorker) handleOrderExpiration(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		OrderID string `json:"orderId"`
	}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Processing expiration for order: %s", payload.OrderID)

	// Publish expiration:complete event
	evt := events.ExpirationCompleteData{
		OrderID: payload.OrderID,
	}
	b, _ := json.Marshal(evt)
	if err := w.pub.Publish(ctx, string(events.SubjectExpirationComplete), b); err != nil {
		log.Printf("Failed to publish expiration event: %v", err)
		return err
	}

	log.Printf("Published expiration:complete for order: %s", payload.OrderID)
	return nil
}
