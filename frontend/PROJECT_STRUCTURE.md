# Frontend Project Structure

## 📊 Statistics

- **Total Files**: 56 TypeScript/TSX files
- **API Routes**: 9 BFF endpoints
- **Components**: 23 React components
- **Pages**: 8 application pages
- **UI Components**: 10 shadcn/ui components

## 📁 Complete Structure

```
ticket-booking-frontend/
├── src/
│   ├── app/                                # Next.js App Router
│   │   ├── api/                           # BFF API Routes (9 files)
│   │   │   ├── auth/
│   │   │   │   ├── signup/route.ts        # POST /api/auth/signup
│   │   │   │   ├── signin/route.ts        # POST /api/auth/signin
│   │   │   │   ├── signout/route.ts       # POST /api/auth/signout
│   │   │   │   └── currentuser/route.ts   # GET /api/auth/currentuser
│   │   │   ├── tickets/
│   │   │   │   ├── route.ts               # GET/POST /api/tickets
│   │   │   │   └── [id]/route.ts          # GET/PUT /api/tickets/:id
│   │   │   ├── orders/
│   │   │   │   ├── route.ts               # GET/POST /api/orders
│   │   │   │   └── [id]/route.ts          # GET/DELETE /api/orders/:id
│   │   │   └── payments/
│   │   │       └── route.ts               # POST /api/payments
│   │   │
│   │   ├── tickets/                       # Ticket Pages
│   │   │   ├── page.tsx                   # GET /tickets - List all tickets
│   │   │   └── [id]/page.tsx              # GET /tickets/:id - Ticket details
│   │   │
│   │   ├── orders/                        # Order Pages
│   │   │   ├── page.tsx                   # GET /orders - My orders
│   │   │   └── [id]/
│   │   │       ├── page.tsx               # GET /orders/:id - Order details
│   │   │       └── payment/page.tsx       # GET /orders/:id/payment - Payment
│   │   │
│   │   ├── signin/page.tsx                # GET /signin
│   │   ├── signup/page.tsx                # GET /signup
│   │   ├── layout.tsx                     # Root layout with providers
│   │   ├── page.tsx                       # GET / - Home page
│   │   └── globals.css                    # Global styles
│   │
│   ├── components/
│   │   ├── ui/                            # shadcn/ui Components (10)
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── input.tsx
│   │   │   ├── form.tsx
│   │   │   ├── label.tsx
│   │   │   ├── sonner.tsx
│   │   │   ├── dropdown-menu.tsx
│   │   │   ├── avatar.tsx
│   │   │   ├── badge.tsx
│   │   │   ├── alert.tsx
│   │   │   └── separator.tsx
│   │   │
│   │   ├── auth/                          # Auth Components
│   │   │   ├── SignupForm.tsx             # Registration form
│   │   │   └── SigninForm.tsx             # Login form
│   │   │
│   │   ├── tickets/                       # Ticket Components
│   │   │   ├── TicketCard.tsx             # Ticket card in list
│   │   │   ├── TicketList.tsx             # Ticket grid display
│   │   │   └── TicketDetail.tsx           # Full ticket view with purchase
│   │   │
│   │   ├── orders/                        # Order Components
│   │   │   ├── OrderCard.tsx              # Order card in list
│   │   │   ├── OrderList.tsx              # Orders grid display
│   │   │   └── OrderCountdown.tsx         # 15-minute countdown timer
│   │   │
│   │   ├── payment/                       # Payment Components
│   │   │   └── StripeCheckout.tsx         # Stripe payment form
│   │   │
│   │   └── layout/                        # Layout Components
│   │       ├── Header.tsx                 # Navigation header
│   │       ├── Footer.tsx                 # Footer
│   │       └── AuthGuard.tsx              # Protected route wrapper
│   │
│   └── lib/
│       ├── api/                           # API Client Layer
│       │   ├── client.ts                  # Axios config + error handling
│       │   ├── auth.ts                    # Auth API functions
│       │   ├── tickets.ts                 # Tickets API functions
│       │   ├── orders.ts                  # Orders API functions
│       │   └── payments.ts                # Payments API functions
│       │
│       ├── types/                         # TypeScript Definitions
│       │   ├── auth.ts                    # User, SignupInput, etc.
│       │   ├── ticket.ts                  # Ticket, CreateTicketInput, etc.
│       │   ├── order.ts                   # Order, OrderStatus, etc.
│       │   └── payment.ts                 # Payment, CreatePaymentInput, etc.
│       │
│       ├── store/                         # State Management
│       │   └── auth.ts                    # Zustand auth store
│       │
│       ├── hooks/                         # Custom React Hooks
│       │   └── useCountdown.ts            # Countdown timer hook
│       │
│       └── utils/                         # Utility Functions
│           ├── utils.ts                   # General utilities (from shadcn)
│           └── format.ts                  # Price, date, status formatters
│
├── public/                                # Static Assets
│
├── Configuration Files
├── .env.local                             # Environment variables
├── .env.local.example                     # Example env file
├── .gitignore                             # Git ignore rules
├── next.config.ts                         # Next.js configuration
├── tailwind.config.ts                     # Tailwind CSS configuration
├── tsconfig.json                          # TypeScript configuration
├── eslint.config.mjs                      # ESLint configuration
├── postcss.config.mjs                     # PostCSS configuration
├── components.json                        # shadcn/ui configuration
├── package.json                           # Dependencies and scripts
├── package-lock.json                      # Locked dependencies
│
├── Deployment
├── Dockerfile                             # Docker build configuration
│
└── Documentation
    ├── README.md                          # Project overview
    ├── SETUP_COMPLETE.md                  # Complete setup guide
    ├── INTEGRATION_GUIDE.md               # Backend integration guide
    └── PROJECT_STRUCTURE.md               # This file
```

## 🎯 Component Hierarchy

### Page Components

```
RootLayout (layout.tsx)
├── QueryClientProvider
├── Header
│   ├── Navigation Links
│   └── User Menu (if authenticated)
│       └── Avatar + Dropdown
├── Main Content Area
│   ├── HomePage (/)
│   ├── SignupPage (/signup)
│   │   └── SignupForm
│   ├── SigninPage (/signin)
│   │   └── SigninForm
│   ├── TicketsPage (/tickets)
│   │   └── TicketList
│   │       └── TicketCard (multiple)
│   ├── TicketDetailPage (/tickets/:id)
│   │   └── TicketDetail
│   ├── OrdersPage (/orders) [Protected]
│   │   └── AuthGuard
│   │       └── OrderList
│   │           └── OrderCard (multiple)
│   ├── OrderDetailPage (/orders/:id) [Protected]
│   │   └── AuthGuard
│   │       ├── OrderCountdown (if active)
│   │       └── Order Details
│   └── PaymentPage (/orders/:id/payment) [Protected]
│       └── AuthGuard
│           ├── OrderCountdown
│           └── StripeCheckout
├── Footer
└── Toaster (notifications)
```

## 🔌 API Routes Flow

### Authentication Flow
```
Frontend Component
  ↓ (user action)
SignupForm / SigninForm
  ↓ (POST request)
/api/auth/signup or /api/auth/signin (BFF)
  ↓ (proxy)
Go Auth Service (localhost:3000)
  ↓ (returns JWT cookie)
BFF extracts cookie
  ↓ (sets cookie)
Frontend receives user data
  ↓ (updates state)
Zustand auth store updated
```

### Ticket Purchase Flow
```
TicketDetail Component
  ↓ (user clicks "Purchase")
POST /api/orders (BFF)
  ↓ (proxy with JWT cookie)
Go Orders Service
  ↓ (creates order)
  ├── Reserve ticket in DB
  ├── Set 15-min expiration
  └── Publish order:created event
  ↓ (returns order)
Frontend receives order
  ↓ (redirect)
PaymentPage with countdown
  ↓ (user submits payment)
POST /api/payments (BFF)
  ↓ (proxy)
Go Payments Service
  ↓ (process payment)
  ├── Charge via Stripe
  ├── Create payment record
  └── Publish payment:created event
  ↓ (returns success)
Frontend shows success
  ↓ (update state)
Orders marked complete
```

## 📦 Dependencies

### Core
- next: ^16.0.5
- react: ^19.0.0
- typescript: ^5.7.2

### UI & Styling
- tailwindcss: ^4.0.7
- @radix-ui/*: Multiple components
- tailwind-merge: ^2.6.0
- clsx: ^2.1.1

### State & Data
- zustand: ^5.0.2
- @tanstack/react-query: ^5.62.11
- axios: ^1.7.9

### Forms & Validation
- react-hook-form: ^7.54.2
- @hookform/resolvers: ^3.9.1
- zod: ^3.24.1

### Payments
- @stripe/stripe-js: ^4.14.0
- @stripe/react-stripe-js: ^2.10.0

### Utilities
- date-fns: ^4.1.0
- sonner: ^1.7.3 (toast notifications)

## 🎨 Styling System

- **Framework**: Tailwind CSS v4
- **Components**: shadcn/ui (Radix UI primitives)
- **Theme**: Custom CSS variables defined in globals.css
- **Responsive**: Mobile-first breakpoints (sm, md, lg, xl, 2xl)
- **Dark Mode**: Configured (can be enabled)

## 🔒 Security Features

1. **HTTP-only Cookies**: JWT stored securely
2. **CSRF Protection**: Built into Next.js
3. **XSS Prevention**: React's built-in escaping
4. **Environment Variables**: Sensitive data in .env.local
5. **API Proxying**: Backend URLs hidden from client
6. **Input Validation**: Zod schemas on all forms

## 📈 Performance Features

1. **Code Splitting**: Automatic with Next.js
2. **Image Optimization**: Next.js Image component
3. **React Query Caching**: Reduces API calls
4. **Optimistic Updates**: Instant UI feedback
5. **Standalone Output**: Optimized Docker builds

## 🧪 Testing Strategy

Currently set up for:
- ESLint for code quality
- TypeScript for type safety
- Manual testing via UI

Can be extended with:
- Jest + React Testing Library
- Playwright for E2E tests
- Storybook for component development

## 🚀 Deployment Options

1. **Vercel** (Recommended)
   - Zero-config deployment
   - Automatic HTTPS
   - Edge functions
   - Global CDN

2. **Docker**
   - Standalone build
   - Multi-stage Dockerfile
   - Production-ready

3. **Traditional Server**
   - npm build + npm start
   - PM2 for process management
   - Nginx reverse proxy

## 📝 Scripts

- `npm run dev` - Development server (port 3000)
- `npm run build` - Production build
- `npm start` - Start production server
- `npm run lint` - Run ESLint
- `npm run type-check` - TypeScript checking (not in package.json yet)

## 🎓 Learning Resources

Each component is documented with:
- Clear prop interfaces
- TypeScript types
- Usage examples in pages
- Integration with backend APIs

Good starting points for learning:
1. `src/app/layout.tsx` - Providers setup
2. `src/components/auth/SignupForm.tsx` - Form handling
3. `src/components/tickets/TicketDetail.tsx` - API mutations
4. `src/app/api/orders/route.ts` - BFF pattern
5. `src/lib/store/auth.ts` - State management

## 🔄 State Management

### Server State (TanStack Query)
- Tickets data
- Orders data
- API requests
- Caching & invalidation

### Client State (Zustand)
- User authentication
- Current user data
- Auth loading states

### Local State (useState)
- Form inputs
- UI toggles
- Temporary data

## �� Next Features to Add

1. **User Profile Page** - Edit user settings
2. **Ticket Creation** - Sell tickets UI
3. **Search & Filter** - Find specific tickets
4. **Order History** - Detailed transaction log
5. **Notifications** - Real-time updates via WebSocket
6. **Reviews** - Rate events/tickets
7. **Admin Dashboard** - Manage platform
8. **Mobile App** - React Native version

---

**Created**: November 29, 2025  
**Framework**: Next.js 16.0.5  
**Total Development Time**: ~1 hour  
**Files Created**: 56 TypeScript/TSX files
