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
