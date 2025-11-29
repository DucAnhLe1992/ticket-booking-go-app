# Ticket Booking Platform - Full Stack Monorepo

A complete ticket booking system with Go microservices backend and Next.js frontend. This monorepo demonstrates event-driven architecture using NATS for pub/sub messaging, PostgreSQL for persistence, Redis for task queuing, and a modern React-based UI.

## 📚 Documentation

- Development workflow: see `docs/DEVELOPMENT.md`
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

- Go 1.20 or higher
- Docker & Docker Compose
- Kubernetes cluster (minikube, kind, or Docker Desktop)
- kubectl
- Skaffold (optional, for k8s dev workflow)
- Air (optional, for hot-reload: `go install github.com/cosmtrek/air@latest`)

## 🚀 Quick Start

### Full Stack Development

1. **Start infrastructure services**:
```bash
make dev-infra
# Starts Postgres, NATS, and Redis in Docker
```

2. **Run database migrations**:
```bash
make migrate-up
```

3. **Start backend services** (in separate terminals):
```bash
# Terminal 1 - Auth Service
make dev-auth

# Terminal 2 - Tickets Service
make dev-tickets

# Terminal 3 - Orders Service
make dev-orders

# Terminal 4 - Payments Service
make dev-payments

# Terminal 5 - Expiration Worker
make dev-expiration
```

4. **Start frontend** (in a new terminal):
```bash
cd frontend
npm install
npm run dev
```

5. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend services: http://localhost:8080

### Using the Interactive Dev Script

```bash
./dev.sh
# Menu-driven interface for starting services
```

### Option 2: Kubernetes with Skaffold

1. **Start local Kubernetes cluster**:
```bash
# Using minikube
minikube start

# Or using Docker Desktop (enable Kubernetes in settings)
```

2. **Update secrets in `infra/k8s/*-deployment.yaml`** files with your credentials

3. **Run with Skaffold** (builds, deploys, and watches for changes):
```bash
skaffold dev
```

Services will be available at:
- Auth: http://localhost:3001
- Tickets: http://localhost:3002
- Orders: http://localhost:3003
- Payments: http://localhost:3004

## 🧪 Testing

Run all tests:
```bash
make test
```

Generate coverage report:
```bash
make test-coverage
# Opens coverage.html in browser
```

Run tests for specific service:
```bash
go test -v ./internal/auth/...
go test -v ./internal/tickets/...
go test -v ./internal/orders/...
go test -v ./internal/payments/...
```

## 📦 Building

Build all services:
```bash
make build
```

Build Docker images:
```bash
# Build all
make docker-build-all

# Build specific service
make docker-build-auth
make docker-build-tickets
make docker-build-orders
make docker-build-payments
make docker-build-expiration
```

## 📡 API Endpoints

### Auth Service (Port 3001)
- `POST /api/users/signup` - Create account
- `POST /api/users/signin` - Login
- `POST /api/users/signout` - Logout
- `GET /api/users/currentuser` - Get current user

### Tickets Service (Port 3002)
- `POST /api/tickets` - Create ticket (auth required)
- `GET /api/tickets` - List all tickets
- `GET /api/tickets/:id` - Get ticket details
- `PUT /api/tickets/:id` - Update ticket (auth required, owner only)

### Orders Service (Port 3003)
- `POST /api/orders` - Create order (auth required)
- `GET /api/orders` - List user's orders (auth required)
- `GET /api/orders/:id` - Get order details (auth required)
- `DELETE /api/orders/:id` - Cancel order (auth required)

### Payments Service (Port 3004)
- `POST /api/payments` - Create payment (auth required)
- `POST /api/payments/webhook` - Stripe webhook

## 🔄 Event Flow

```
1. User creates ticket
   → Tickets Service publishes "ticket:created"
   → Orders Service replicates ticket data

2. User creates order
   → Orders Service reserves ticket, publishes "order:created"
   → Expiration Service schedules expiration job
   → Payments Service replicates order data

3. Order expires (15 minutes)
   → Expiration Worker publishes "expiration:complete"
   → Orders Service cancels order, releases ticket

4. User pays for order
   → Payments Service creates Stripe PaymentIntent
   → Stripe webhook triggers "payment:created"
   → Orders Service marks order as complete
```

## 🗄️ Database Schema

Each service has its own database/schema:

- **Auth**: `users` table
- **Tickets**: `tickets` table with version field (OCC)
- **Orders**: `orders` and replicated `tickets` tables
- **Payments**: `payments` and replicated `orders` tables

Migrations are in `migrations/*.sql`.

## 🐛 Troubleshooting

**Port conflicts**:
```bash
# Check what's using ports
lsof -i :3000,4222,5432,6379
```

**Database connection issues**:
```bash
# Test Postgres connection
psql postgres://postgres:password@localhost:5432/tickets

# Check running containers
docker-compose ps
```

**NATS connection issues**:
```bash
# Check NATS health
curl http://localhost:8222/healthz
```

**Kubernetes pod issues**:
```bash
# Check pod status
kubectl get pods

# View logs
kubectl logs <pod-name>

# Describe pod
kubectl describe pod <pod-name>
```

## 🏃‍♂️ Development Workflow

1. Make code changes
2. Air automatically rebuilds and restarts the service (if using `make dev-*`)
3. Tests run automatically in watch mode (if using air with test config)
4. Docker Compose keeps infrastructure services running

For Kubernetes development:
1. Make code changes
2. Skaffold detects changes, rebuilds images, and redeploys
3. Port forwarding provides local access to services

## 📝 Environment Variables

### All Services
- `DATABASE_URL` - Postgres connection string
- `JWT_SECRET` - Secret for JWT signing/verification
- `NATS_URL` - NATS server URL (default: nats://localhost:4222)

### Payments Service
- `STRIPE_KEY` - Stripe secret key
- `STRIPE_WEBHOOK_SECRET` - Stripe webhook signing secret

### Expiration Service
- `REDIS_HOST` - Redis connection string (default: localhost:6379)
- `REDIS_PASSWORD` - Redis password (optional)

## 🔐 Security Notes

- **Production**: Use Kubernetes Secrets or external secret management (Vault, AWS Secrets Manager)
- **JWT Secret**: Use a strong random value, rotate regularly
- **Database**: Use strong passwords, enable SSL in production
- **Stripe**: Never commit real API keys, use test keys for development

## 📚 Project Structure

```
.
├── cmd/                    # Service entrypoints
│   ├── auth/
│   ├── tickets/
│   ├── orders/
│   ├── payments/
│   └── expiration/
├── internal/               # Internal packages
│   ├── common/            # Shared code
│   │   ├── errors/
│   │   ├── events/
│   │   ├── middleware/
│   │   └── models/
│   ├── auth/              # Auth service logic
│   ├── tickets/           # Tickets service logic
│   ├── orders/            # Orders service logic
│   ├── payments/          # Payments service logic
│   └── expiration/        # Expiration worker logic
├── migrations/            # SQL migrations
├── infra/
│   └── k8s/              # Kubernetes manifests
├── docker-compose.yaml   # Local dev infrastructure
├── skaffold.yaml         # Skaffold configuration
├── .air.toml             # Air hot-reload config
└── Makefile              # Build automation
```

## 🎯 Advantages of Go over Node.js

✅ **Performance**: 2-3x faster execution, lower memory footprint
✅ **Concurrency**: Built-in goroutines and channels
✅ **Type Safety**: Static typing catches errors at compile time
✅ **Single Binary**: Easy deployment, no runtime dependencies
✅ **Standard Library**: Comprehensive stdlib reduces dependencies

## 🎨 Frontend Application

### Overview

The frontend is a modern Next.js 14 application with:
- **Framework**: Next.js with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **UI Components**: shadcn/ui (Radix UI)
- **State Management**: Zustand + TanStack Query
- **Forms**: React Hook Form + Zod
- **Payments**: Stripe Elements

### Getting Started

```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Configure environment
cp .env.local.example .env.local
# Edit .env.local with your configuration

# Start development server
npm run dev
```

Visit: **http://localhost:3000**

### Frontend Routes

**Public Pages:**
- `/` - Home page with features
- `/tickets` - Browse all tickets
- `/tickets/[id]` - Ticket details & purchase
- `/signin` - User sign in
- `/signup` - User registration

**Protected Pages:**
- `/orders` - View all orders
- `/orders/[id]` - Order details
- `/orders/[id]/payment` - Complete payment

### BFF API Pattern

The frontend includes Backend-for-Frontend (BFF) API routes that proxy to Go services:

```
/api/auth/*          → Auth service
/api/tickets/*       → Tickets service
/api/orders/*        → Orders service
/api/payments/*      → Payments service
```

**Benefits:**
- Simplifies authentication (JWT cookies)
- Reduces CORS complexity
- Aggregates backend calls
- Provides consistent API

### Frontend Documentation

See `frontend/` directory:
- `README.md` - Project overview
- `SETUP_COMPLETE.md` - Complete setup guide
- `INTEGRATION_GUIDE.md` - Backend integration
- `PROJECT_STRUCTURE.md` - Architecture details

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
