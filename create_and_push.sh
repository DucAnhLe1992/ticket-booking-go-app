#!/usr/bin/env bash
set -euo pipefail

# Usage:
# 1) git clone git@github.com:DucAnhLe1992/ticket-booking-go-app.git
# 2) cd ticket-booking-go-app
# 3) chmod +x create_and_push.sh
# 4) ./create_and_push.sh
#
# The script will create files, commit them, and push to origin/main.
# Make sure your local git is configured and you have push permission.

GIT_BRANCH="main"
COMMIT_MSG="infra-first Go monorepo scaffold: initial commit"

echo "Creating scaffold files..."

# Create directories
mkdir -p cmd/payments cmd/orders cmd/expiration cmd/auth cmd/tickets
mkdir -p internal/pubsub internal/store internal/payments
mkdir -p infra/k8s migrations .github/workflows

# README.md
cat > README.md <<'EOF'
# Ticket Booking — Go Monorepo (starter scaffold)

This repository is a starter scaffold to port and implement the ticket-booking microservices in Go.
It is infra-first: Dockerfiles, Kubernetes manifests, CI workflow, migrations and minimal service skeletons
are included so you can run, test, and iterate quickly.

Goals
- Provide infra (Docker + k8s + CI) so services can be deployed consistently.
- Provide minimal service skeletons for payments, orders, expiration (worker), auth and tickets.
- Let you swap stubs for real clients (Stripe, NATS, Postgres) incrementally.

How to use
1. Install Go 1.20+.
2. Run `make build` to build all services.
3. Configure environment (DATABASE_URL, NATS_URL, STRIPE_KEY) for local run.
4. Use Docker and kubectl (or skaffold) to deploy infra/k8s manifests.

Project layout (selected)
- cmd/
  - payments/, orders/, expiration/, auth/, tickets/ — service entrypoints
- internal/
  - store/ — DB connection / migration helpers
  - pubsub/ — pub/sub interfaces and stubs (NATS)
  - payments/ — payments service, handler and stripe client stub
- infra/
  - k8s/ — Kubernetes manifests
- migrations/ — example SQL migrations
- .github/workflows/ci.yml — CI workflow

Next steps
- Replace stubs with real clients: nats.go, stripe-go.
- Implement DB models and run migrations.
- Implement HTTP handlers and tests for each service.
- Add secrets management (Vault / cloud secrets / GitHub Secrets), observability (Prometheus / OpenTelemetry).

License: MIT
EOF

# .gitignore
cat > .gitignore <<'EOF'
# Binaries
/bin/
/build/

# Go
*.test
vendor/
*.exe
*.dll
*.so
*.dylib

# Node / TS artifacts (if present)
node_modules/
dist/
coverage/

# Editor
.vscode/
.idea/
.DS_Store
EOF

# go.mod
cat > go.mod <<'EOF'
module github.com/DucAnhLe1992/ticket-booking-go-app

go 1.20
EOF

# Makefile
cat > Makefile <<'EOF'
.PHONY: all build docker

all: build

build:
	go build ./cmd/...

docker-build-payments:
	docker build -t ghcr.io/ducanhle1992/ticket-payments:dev -f ./cmd/payments/Dockerfile .

docker-build-orders:
	docker build -t ghcr.io/ducanhle1992/ticket-orders:dev -f ./cmd/orders/Dockerfile .

docker-build-all: docker-build-payments docker-build-orders
EOF

# GitHub Actions CI
cat > .github/workflows/ci.yml <<'EOF'
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
      - name: Build
        run: |
          go version
          go test ./... -v
          go build ./...
      - name: Lint (basic)
        run: |
          echo "Add linters as needed"
EOF

# cmd/payments/main.go
cat > cmd/payments/main.go <<'EOF'
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
EOF

# cmd/orders/main.go
cat > cmd/orders/main.go <<'EOF'
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := ":3001"
	if p := os.Getenv("ORDERS_PORT"); p != "" {
		port = ":" + p
	}
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("orders ok"))
	})
	log.Printf("orders service listening on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
EOF

# cmd/expiration/main.go
cat > cmd/expiration/main.go <<'EOF'
package main

import (
	"log"
	"os"
	"time"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/expiration"
)

func main() {
	interval := time.Minute
	log.Println("expiration worker starting")
	for {
		if err := expiration.RunOnce(); err != nil {
			log.Printf("expiration.RunOnce error: %v", err)
		}
		time.Sleep(interval)
		if os.Getenv("ONE_SHOT") == "true" {
			break
		}
	}
}
EOF

# cmd/auth/main.go
cat > cmd/auth/main.go <<'EOF'
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := ":3002"
	if p := os.Getenv("AUTH_PORT"); p != "" {
		port = ":" + p
	}
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("auth ok"))
	})
	log.Printf("auth service listening on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
EOF

# cmd/tickets/main.go
cat > cmd/tickets/main.go <<'EOF'
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := ":3003"
	if p := os.Getenv("TICKETS_PORT"); p != "" {
		port = ":" + p
	}
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("tickets ok"))
	})
	log.Printf("tickets service listening on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
EOF

# internal/pubsub/pubsub.go
cat > internal/pubsub/pubsub.go <<'EOF'
package pubsub

import "context"

// Publisher is a minimal pub/sub publisher interface used by services.
type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Close() error
}

// Subscriber is a minimal subscriber interface.
type Subscriber interface {
	Subscribe(subject string, handler func(msg []byte)) error
	Close() error
}

// NewNoopPublisher returns a publisher that does nothing (useful for dev).
func NewNoopPublisher() Publisher { return noop{} }

type noop struct{}

func (n noop) Publish(ctx context.Context, subject string, data []byte) error { return nil }
func (n noop) Close() error                                                  { return nil }
func (n noop) Subscribe(subject string, handler func(msg []byte)) error     { return nil }
EOF

# internal/store/store.go
cat > internal/store/store.go <<'EOF'
package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // optional: uncomment when integrating Postgres
)

// NewPostgres connects to Postgres and returns *sql.DB.
// Example DSN: postgres://user:pw@localhost:5432/dbname?sslmode=disable
func NewPostgres(dsn string) (*sql.DB, error) {
	// NOTE: the pq/pgx driver must be added in go.mod if you call this function
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// quick ping
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}
EOF

# internal/payments/handler.go
cat > internal/payments/handler.go <<'EOF'
package payments

import (
	"encoding/json"
	"io"
	"net/http"
)

// HTTPHandler exposes HTTP endpoints for payments.
type HTTPHandler struct {
	svc *Service
}

func NewHTTPHandler(svc *Service) *HTTPHandler { return &HTTPHandler{svc: svc} }

type CreateChargeRequest struct {
	OrderID      string `json:"order_id"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
}

func (h *HTTPHandler) CreateCharge(w http.ResponseWriter, r *http.Request) {
	var req CreateChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	payment, err := h.svc.CreateCharge(r.Context(), req.OrderID, req.Amount, req.Currency)
	if err != nil {
		http.Error(w, "create charge error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payment)
}

// Webhook reads raw body and forwards to service for verification/processing.
func (h *HTTPHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	// In real impl: verify signature header and call ProcessWebhook
	_ = h.svc.ProcessWebhook(r.Context(), body)
	w.WriteHeader(http.StatusOK)
}
EOF

# internal/payments/service.go
cat > internal/payments/service.go <<'EOF'
package payments

import (
	"context"
	"database/sql"
	"time"
)

// Payment is a very small payment record used for responses.
type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// StripeClientStub defines methods used from stripe client.
type StripeClient interface {
	CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (string, error)
	VerifyWebhookSignature(payload []byte, sigHeader string) (bool, error)
}

// Publisher is expected to publish domain events.
type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
}

type Service struct {
	db     *sql.DB
	pub    Publisher
	stripe StripeClient
}

func NewService(db *sql.DB, pub Publisher, stripe StripeClient) *Service {
	return &Service{db: db, pub: pub, stripe: stripe}
}

func (s *Service) CreateCharge(ctx context.Context, orderID string, amount int64, currency string) (*Payment, error) {
	// TODO: implement real Stripe creation, DB transaction, idempotency and publish events.
	now := time.Now().UTC()
	return &Payment{
		ID:       "pay_" + orderID,
		OrderID:  orderID,
		Amount:   amount,
		Currency: currency,
		Status:   "created",
		CreatedAt: now,
	}, nil
}

func (s *Service) ProcessWebhook(ctx context.Context, payload []byte) error {
	// TODO: verify signature and process event idempotently
	return nil
}
EOF

# internal/payments/stripe_client.go
cat > internal/payments/stripe_client.go <<'EOF'
package payments

// Stripe client stub. Replace with a real stripe-go wrapper.

type stripeClientStub struct {
	key string
}

func NewStripeClient(key string) StripeClient {
	return &stripeClientStub{key: key}
}

func (s *stripeClientStub) CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (string, error) {
	// In real implementation: use stripe-go to create a PaymentIntent and return its ID.
	return "pi_stub_123", nil
}

func (s *stripeClientStub) VerifyWebhookSignature(payload []byte, sigHeader string) (bool, error) {
	// In real implementation: use stripe.Webhook.ConstructEvent or similar to verify.
	return true, nil
}
EOF

# migrations
cat > migrations/001_create_payments.sql <<'EOF'
-- Create payments table
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL,
  customer_id UUID,
  stripe_id TEXT,
  amount BIGINT NOT NULL,
  currency TEXT NOT NULL DEFAULT 'usd',
  status TEXT NOT NULL,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS webhook_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  provider TEXT NOT NULL,
  event_id TEXT NOT NULL,
  payload JSONB NOT NULL,
  processed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS idempotency_keys (
  key TEXT PRIMARY KEY,
  response JSONB,
  created_at TIMESTAMPTZ DEFAULT now()
);
EOF

# infra/k8s/payments-deployment.yaml
cat > infra/k8s/payments-deployment.yaml <<'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payments
spec:
  replicas: 1
  selector:
    matchLabels:
      app: payments
  template:
    metadata:
      labels:
        app: payments
    spec:
      containers:
      - name: payments
        image: ghcr.io/ducanhle1992/ticket-payments:dev
        imagePullPolicy: IfNotPresent
        env:
        - name: DATABASE_URL
          value: "postgres://postgres:password@postgres:5432/tickets?sslmode=disable"
        - name: NATS_URL
          value: "nats://nats:4222"
        - name: STRIPE_KEY
          value: "sk_test_..."
        ports:
        - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: payments
spec:
  selector:
    app: payments
  ports:
  - protocol: TCP
    port: 3000
    targetPort: 3000
EOF

# cmd/payments/Dockerfile
cat > cmd/payments/Dockerfile <<'EOF'
# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/payments ./cmd/payments

# Final stage
FROM alpine:3.18
RUN addgroup -S app && adduser -S -G app app
COPY --from=builder /bin/payments /bin/payments
USER app
EXPOSE 3000
ENTRYPOINT ["/bin/payments"]
EOF

# Optional: ensure go.sum exists (empty to begin)
touch go.sum

echo "All files created."

# Initialize git if needed
if [ ! -d .git ]; then
  echo ".git directory not found. Initializing git repository..."
  git init
  git branch -M "${GIT_BRANCH}"
  git remote add origin "$(git config --get remote.origin.url || echo git@github.com:DucAnhLe1992/ticket-booking-go-app.git)"
fi

# Ensure on correct branch
current_branch=$(git rev-parse --abbrev-ref HEAD || true)
if [ "$current_branch" != "${GIT_BRANCH}" ]; then
  echo "Switching to branch ${GIT_BRANCH} (creating if necessary)..."
  git checkout -B "${GIT_BRANCH}"
fi

git add -A
git commit -m "${COMMIT_MSG}" || true

echo "Pushing to origin/${GIT_BRANCH}..."
git push -u origin "${GIT_BRANCH}"

echo "Done. Please verify the repository on GitHub."