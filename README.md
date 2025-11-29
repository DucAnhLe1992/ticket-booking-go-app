# Ticket Booking Platform - Full Stack Monorepo

A complete ticket booking system with Go microservices backend and Next.js frontend. This monorepo demonstrates event-driven architecture using NATS for pub/sub messaging, PostgreSQL for persistence, Redis for task queuing, and a modern React-based UI.

## ✨ Features

- User auth with JWT cookies
- Tickets CRUD with optimistic concurrency control
- Orders with 15-minute expiration and reservation workflow
- Stripe payments with webhook handling
- Event-driven communication via NATS
- Next.js BFF API routes to avoid CORS
- Kubernetes-ready manifests and local Docker Compose

## 🔎 Motivation

Building a real-world ticketing platform is a great way to demonstrate modern microservices patterns end-to-end. This project shows how independent services (Auth, Tickets, Orders, Payments, Expiration) communicate via events, maintain their own data, and scale independently—while providing a unified user experience through a Next.js frontend. It’s designed for learning, rapid iteration, and clear separation of concerns without sacrificing developer ergonomics.

Key goals:
- Event-driven architecture with clear domain boundaries
- Robust data models with optimistic concurrency control
- Production-friendly patterns (Docker, Kubernetes, CI-ready)
- Developer-friendly monorepo with a simple, repeatable workflow

## 📚 Documentation

- Development workflow: see `docs/DEVELOPMENT.md`
- Architecture & deep reference: see `docs/ARCHITECTURE.md`
- Frontend overview: see `frontend/README.md`
- Frontend structure reference: see `frontend/PROJECT_STRUCTURE.md`

## 📁 Monorepo Structure

```
ticket-booking-go-app/
├── cmd/                          # Go service entry points
├── internal/                     # Shared Go packages
├── migrations/                   # Database migrations
├── infra/k8s/                   # Kubernetes manifests
├── frontend/                     # Next.js application
│   ├── src/
│   │   ├── app/                 # Pages & API routes
│   │   ├── components/          # React components
│   │   └── lib/                 # Utilities & types
│   └── package.json
├── docker-compose.yaml           # Local infrastructure
├── skaffold.yaml                 # Kubernetes development
└── Makefile                      # Build & run commands
```

## 🧰 Tech Stack

- Backend: Go (Chi), PostgreSQL, NATS, Redis, Stripe
- Frontend: Next.js 14, TypeScript, Tailwind, shadcn/ui
- Tooling: Docker, Kubernetes, Skaffold, Makefile

Details and commands live in `docs/DEVELOPMENT.md`.

## 🗺️ Screens & Flows

- `/` Home: overview and navigation
- `/signin`, `/signup`: authentication
- `/tickets`: browse tickets
- `/tickets/[id]`: ticket details + purchase
- `/orders`: list my orders
- `/orders/[id]`: order details + countdown
- `/orders/[id]/payment`: Stripe checkout

API calls go through Next.js BFF routes under `/app/api/*`.

## ▶️ Usage

Quick local run:

```bash
# From repository root
make dev-infra         # Start Postgres, NATS, Redis
make migrate-up        # Apply DB migrations

# Start backend services (use multiple terminals)
make dev-auth
make dev-tickets
make dev-orders
make dev-payments
make dev-expiration

# Start the frontend
make frontend-install
make frontend-dev
```

Then open `http://localhost:3000`.

Configure the frontend by creating `frontend/.env.local`:

```env
API_URL=http://localhost:8080
NEXT_PUBLIC_STRIPE_PUBLIC_KEY=pk_test_your_key_here
```

For more details, see `docs/DEVELOPMENT.md`.

<!-- Detailed architecture moved to docs/ARCHITECTURE.md to keep README concise. -->
## 🏗️ Architecture

### Microservices

- **Auth Service** - User authentication with JWT and bcrypt
- **Tickets Service** - CRUD operations for tickets with optimistic concurrency control (OCC)
- **Orders Service** - Order management with ticket reservation and expiration
- **Payments Service** - Stripe payment processing with webhook handling
- **Expiration Service** - Redis-based delayed job worker for order expiration

### Technologies

- **Language**: Go 1.20+
- **Web Framework**: Chi router v5
- **Database**: PostgreSQL with lib/pq driver
- **Message Bus**: NATS (pub/sub)
- **Task Queue**: Asynq (Redis-backed)
- **Payment**: Stripe SDK
- **Auth**: JWT + bcrypt
- **Container Orchestration**: Kubernetes
- **Dev Tools**: Skaffold, Docker Compose, Air (hot-reload)

## 📋 Prerequisites

Go 1.20+, Node.js 18+, Docker & Docker Compose. See `docs/DEVELOPMENT.md` for complete prerequisites and tooling.

## 🚀 Quick Start

For the complete step-by-step quick start (infrastructure, migrations, backend services, frontend, and optional Kubernetes workflow), see `docs/DEVELOPMENT.md`.



For API endpoints, event flows, database schemas, and security notes, see `docs/ARCHITECTURE.md`.



## 🎯 Advantages of Go over Node.js

✅ **Performance**: 2-3x faster execution, lower memory footprint
✅ **Concurrency**: Built-in goroutines and channels
✅ **Type Safety**: Static typing catches errors at compile time
✅ **Single Binary**: Easy deployment, no runtime dependencies
✅ **Standard Library**: Comprehensive stdlib reduces dependencies

## 🎨 Frontend

The Next.js 14 frontend lives in `frontend/`. See `frontend/README.md` for setup and `frontend/PROJECT_STRUCTURE.md` for architecture details.

## ⚠️ Considerations

- Learning curve for developers new to Go
- Verbose error handling compared to try/catch
- Less flexibility in dynamic scenarios
- Smaller ecosystem for some specialized use cases

## 📄 License

MIT

## 👥 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make test` to verify
5. Submit a pull request

## 📞 Support

For issues and questions, please open a GitHub issue
