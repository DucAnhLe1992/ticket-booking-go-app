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
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/store"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/tickets"
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

	// NATS (optional publishing)
	var pub pubsub.Publisher
	if natsURL := os.Getenv("NATS_URL"); natsURL != "" {
		if n, err := pubsub.NewNATS(natsURL); err == nil {
			pub = n
			defer n.Close()
		} else {
			log.Printf("warn: NATS connect failed (%v), continuing without pub", err)
		}
	}

	repo := tickets.NewRepository(db)
	if err := repo.EnsureSchema(context.Background()); err != nil {
		log.Printf("tickets.EnsureSchema: %v", err)
	}
	svc := tickets.NewService(repo, pub)
	h := tickets.NewHTTPHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cmw.RecoverAndJSON)
	r.Use(cmw.CurrentUser)

	r.Get("/api/tickets", h.Index)
	r.Get("/api/tickets/show", h.Show)
	r.Group(func(r chi.Router) {
		r.Use(cmw.RequireAuth)
		r.Post("/api/tickets", h.Create)
		r.Put("/api/tickets", h.Update)
	})

	srv := &http.Server{Addr: ":3000", Handler: r}
	go func() {
		log.Printf("tickets service listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down tickets service...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
