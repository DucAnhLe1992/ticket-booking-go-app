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
	// Initialize pubsub (currently a stub)
	pub := pubsub.NewNoopPublisher()

	// Initialize stripe client stub
	sc := payments.NewStripeClient(stripeKey)

	svc := payments.NewService(db, pub, sc)
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
