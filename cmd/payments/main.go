package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/payments"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/store"
)

func main() {
	// Load config from env
	dbURL := os.Getenv("DATABASE_URL")
	natsURL := os.Getenv("NATS_URL")
	_ = natsURL
	stripeKey := os.Getenv("STRIPE_KEY")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/tickets?sslmode=disable"
	}
	// Initialize DB (stubbed or real)
	db, err := store.NewPostgres(dbURL)
	if err != nil {
		log.Fatalf("store.NewPostgres: %v", err)
	}
	// Initialize pubsub (use NATS if configured; otherwise noop for local dev)
	var pub pubsub.Publisher
	var sub pubsub.Subscriber
	if natsURL != "" {
		if n, err := pubsub.NewNATS(natsURL); err == nil {
			pub = n
			sub = n
			defer n.Close()
		} else {
			log.Printf("warn: NATS connect failed (%v), using noop pubsub", err)
			pub = pubsub.NewNoopPublisher()
		}
	} else {
		pub = pubsub.NewNoopPublisher()
	}

	// Initialize stripe client with webhook secret
	sc := payments.NewStripeClient(stripeKey)
	if webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET"); webhookSecret != "" {
		sc.SetWebhookSecret(webhookSecret)
	}

	svc := payments.NewService(db, pub, sc)
	// Ensure local repo schema
	if err := payments.NewRepository(db).EnsureSchema(context.Background()); err != nil {
		log.Printf("payments.EnsureSchema: %v", err)
	}
	handler := payments.NewHTTPHandler(svc)

	r := chi.NewRouter()
	r.Post("/api/payments/charge", handler.CreateCharge)
	r.Post("/api/payments/webhook", handler.Webhook)

	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		log.Printf("payments service listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Register NATS listeners
	if sub != nil {
		if err := payments.RegisterNATSListeners(context.Background(), sub, payments.NewRepository(db)); err != nil {
			log.Printf("register listeners: %v", err)
		}
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down payments service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

	// optional: drain pubsub connections
	_ = db // placeholder usage
	_ = uuid.New()
}
