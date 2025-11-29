# Ticket Booking Frontend

A modern, responsive frontend application for the ticket booking platform built with Next.js 14, TypeScript, and Tailwind CSS.

For the full end-to-end local development workflow (infrastructure, backend services, and frontend), see `../docs/DEVELOPMENT.md`.

## 🚀 Features

- **Modern Stack**: Next.js 14 with App Router, TypeScript, Tailwind CSS
- **UI Components**: Beautiful UI with shadcn/ui components
- **State Management**: Zustand for auth, TanStack Query for server state
- **Form Handling**: React Hook Form with Zod validation
- **Payment Integration**: Stripe checkout
- **Real-time Features**: Order countdown timers
- **BFF Pattern**: Backend-for-Frontend API routes
- **Responsive Design**: Mobile-first approach

## 📋 Prerequisites

- Node.js 18+ and npm
- Running Go backend services (Auth, Tickets, Orders, Payments)

## 🛠️ Installation

```bash
npm install
```

## 🔧 Configuration

Update `.env.local`:

```env
# Backend API URL (your Go services)
API_URL=http://localhost:8080

# Stripe public key
NEXT_PUBLIC_STRIPE_PUBLIC_KEY=pk_test_your_key_here
```

## 🏃 Running the Application

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000)

## 📁 Key Pages

- `/` - Home page
- `/tickets` - Browse tickets
- `/tickets/[id]` - Ticket details
- `/orders` - My orders (protected)
- `/orders/[id]/payment` - Payment page (protected)
- `/signin` - Sign in
- `/signup` - Sign up

## 🔌 API Integration

BFF routes proxy to Go microservices:
- `/api/auth/*` → Auth service
- `/api/tickets/*` → Tickets service
- `/api/orders/*` → Orders service
- `/api/payments/*` → Payments service

## 🚀 Deployment

### Vercel
```bash
vercel
```

### Docker
```bash
docker build -t ticket-booking-frontend .
docker run -p 3000:3000 ticket-booking-frontend
```

## 📚 Tech Stack

- Next.js 14, TypeScript, Tailwind CSS v4
- shadcn/ui, Zustand, TanStack Query
- React Hook Form, Zod, Axios, Stripe

