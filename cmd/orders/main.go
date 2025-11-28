package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	cmw "github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/middleware"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/orders"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/store"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/tickets?sslmode=disable"
	}

	db, err := store.NewPostgres(dbURL)
	if err != nil {
		log.Fatalf("store.NewPostgres: %v", err)
	}

	// NATS (required for both pub and sub)
	var pub pubsub.Publisher
	var sub pubsub.Subscriber
	natsURL := os.Getenv("NATS_URL")
	if natsURL != "" {
		if n, err := pubsub.NewNATS(natsURL); err == nil {
			pub = n
			sub = n
			defer n.Close()
		} else {
			log.Printf("warn: NATS connect failed (%v), continuing without pub/sub", err)
		}
	}

	repo := orders.NewRepository(db)
	if err := repo.EnsureSchema(context.Background()); err != nil {
		log.Printf("orders.EnsureSchema: %v", err)
	}

	svc := orders.NewService(repo, pub)
	h := orders.NewHTTPHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cmw.RecoverAndJSON)
	r.Use(cmw.CurrentUser)

	r.Group(func(r chi.Router) {
		r.Use(cmw.RequireAuth)
		r.Get("/api/orders", h.Index)
		r.Post("/api/orders", h.Create)
		r.Get("/api/orders/{orderId}", h.Show)
		r.Delete("/api/orders/{orderId}", h.Delete)
	})

	srv := &http.Server{Addr: ":3000", Handler: r}

	go func() {
		log.Printf("orders service listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Register NATS listeners
	if sub != nil {
		if err := orders.RegisterNATSListeners(context.Background(), sub, repo); err != nil {
			log.Printf("register listeners: %v", err)
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down orders service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
