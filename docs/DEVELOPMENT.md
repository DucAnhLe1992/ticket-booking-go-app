# Development Guide

This is the single source of truth for running the full stack locally. It consolidates the previous Quick Start and Frontend Integration docs.

## Prerequisites

- Go 1.20+
- Node.js 18+
- Docker + Docker Compose
- Stripe test public key (for frontend)

## 1) Start Local Infrastructure

```bash
make dev-infra
```

Services started:
- PostgreSQL on `5432`
- NATS on `4222`
- Redis on `6379`

## 2) Run Database Migrations

```bash
make migrate-up
```

## 3) Start Backend Services

Open multiple terminals (or use your preferred process manager):

```bash
make dev-auth
make dev-tickets
make dev-orders
make dev-payments
make dev-expiration
```

Notes:
- All services use JWT cookies for auth.
- Orders expire after 15 minutes; expiration worker releases tickets via Redis-backed jobs.

## 4) Configure and Run Frontend

Option A — via Makefile:

```bash
make frontend-install
make frontend-dev
```

Option B — via npm in `frontend/`:

```bash
cd frontend
npm install
echo "API_URL=http://localhost:8080" > .env.local
echo "NEXT_PUBLIC_STRIPE_PUBLIC_KEY=pk_test_your_key_here" >> .env.local
npm run dev
```

Visit `http://localhost:3000`.

## Environment Variables

Frontend (`frontend/.env.local`):
- `API_URL` — BFF proxy target for Go services (default `http://localhost:8080` if you have an API gateway; otherwise configure per service).
- `NEXT_PUBLIC_STRIPE_PUBLIC_KEY` — Your Stripe test public key.

Backend: configure DB/NATS/Redis/Stripe secrets via your existing env or Kubernetes manifests (see `infra/k8s/`).

## Common Make Targets

```bash
# Infra and DB
make dev-infra        # Start Postgres, NATS, Redis
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations

# Backend services (hot reload with Air when configured)
make dev-auth
make dev-tickets
make dev-orders
make dev-payments
make dev-expiration

# Frontend
make frontend-install
make frontend-dev
make frontend-build
make frontend-start

# All-in dev convenience (if available)
make dev-all
```

## Test Data and Stripe

- Use Stripe test card `4242 4242 4242 4242` with any future expiry, any CVC, any ZIP.
- Create a ticket, place an order, and complete payment via the payment page flow.

## Troubleshooting

- Auth cookie not set: ensure the frontend and backend share the same effective domain/port expectations and that API routes proxy correctly through the Next.js BFF.
- 401s after signin: clear cookies in the browser and retry; check `API_URL` in `frontend/.env.local`.
- CORS issues: the BFF pattern (Next.js `/app/api/*`) avoids CORS by server-side proxying; confirm requests go through the BFF routes rather than browser-direct to Go services.
- Payments failing: verify `NEXT_PUBLIC_STRIPE_PUBLIC_KEY` is set and backend Stripe secret keys are configured for the payments service.

## Repository Map

- Backend services: `cmd/`, shared packages in `internal/`, SQL in `migrations/`
- Frontend app: `frontend/`
- Kubernetes: `infra/k8s/`
- This guide: `docs/DEVELOPMENT.md`
