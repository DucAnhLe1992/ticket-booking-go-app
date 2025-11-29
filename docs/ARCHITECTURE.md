# Architecture & Reference

This document centralizes deep reference material so the root `README.md` can stay concise.

## System Overview

- Microservices: Auth, Tickets, Orders, Payments, Expiration.
- Communication: NATS pub/sub for domain events; Redis-backed delayed jobs for expiration.
- Persistence: PostgreSQL per service schema; replicated data where needed.
- Frontend: Next.js App Router with BFF API routes.

## Backend Architecture

- Auth: JWT issuance/verification, bcrypt password hashing.
- Tickets: CRUD with optimistic concurrency control (version field) to prevent stale writes.
- Orders: ticket reservation, status lifecycle (created, cancelled, complete), 15-minute expiration.
- Payments: Stripe charge creation, webhook verification, order completion.
- Expiration: schedules delayed jobs, publishes cancellation when timers elapse.

## Event Flow

- `ticket:created` / `ticket:updated`: emitted by Tickets; consumed by Orders to keep local replica.
- `order:created`: emitted by Orders; consumed by Expiration to schedule timeout.
- `order:cancelled`: emitted by Orders; consumed by Tickets to release reservation.
- `payment:created`: emitted by Payments; consumed by Orders to mark complete.

## API Endpoints (BFF)

Frontend talks to Next.js API routes which proxy to Go services.

- Auth: `POST /api/auth/signup`, `POST /api/auth/signin`, `POST /api/auth/signout`, `GET /api/auth/currentuser`
- Tickets: `GET/POST /api/tickets`, `GET/PUT /api/tickets/:id`
- Orders: `GET/POST /api/orders`, `GET/DELETE /api/orders/:id`
- Payments: `POST /api/payments`

## Database Schema Highlights

- Tickets: `id`, `title`, `price`, `version` (OCC)
- Orders: `id`, `user_id`, `status`, `expires_at`, replicated `ticket` data
- Payments: `id`, `order_id`, `stripe_id`, `amount`

## Security Notes

- JWT cookies (HttpOnly) for auth; validate on server.
- Stripe keys handled via environment variables/secrets; never checked into VCS.
- BFF routes mitigate CORS and centralize auth handling.

## Frontend Architecture

- State: Zustand (auth), TanStack Query (server state).
- Forms: React Hook Form + Zod.
- UI: Tailwind CSS + shadcn/ui.
- Pages: `/tickets`, `/orders`, `/orders/[id]/payment`, auth pages.

## Performance & DevX

- Hot reload via Air for Go services.
- Fast Refresh for Next.js.
- Docker Compose for local infra; Kubernetes manifests in `infra/k8s/`.
