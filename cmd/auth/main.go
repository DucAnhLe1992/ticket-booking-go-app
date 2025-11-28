package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	authsvc "github.com/DucAnhLe1992/ticket-booking-go-app/internal/auth"
	cmw "github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/middleware"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// NATS publisher (optional)
	natsURL := os.Getenv("NATS_URL")
	var pub pubsub.Publisher
	if natsURL != "" {
		if n, err := pubsub.NewNATS(natsURL); err == nil {
			pub = n
			defer n.Close()
		} else {
			log.Printf("warn: NATS connect failed (%v), continuing without pub", err)
		}
	}

	repo := authsvc.NewUserRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := repo.EnsureSchema(ctx); err != nil {
		log.Printf("warning: EnsureSchema failed: %v", err)
	}
	cancel()

	svc := authsvc.NewService(repo, pub)
	h := authsvc.NewHTTPHandler(svc, repo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cmw.RecoverAndJSON)
	r.Use(cmw.CurrentUser)

	r.Post("/api/users/signup", h.Signup)
	r.Post("/api/users/signin", h.Signin)
	r.Post("/api/users/signout", h.Signout)
	r.Get("/api/users/currentuser", h.CurrentUser)

	srv := &http.Server{Addr: ":3000", Handler: r}

	go func() {
		log.Printf("auth service listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down auth service...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	_ = srv.Shutdown(ctxShutdown)
}
