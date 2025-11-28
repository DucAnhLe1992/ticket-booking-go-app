-- Orders Service: Create orders and replicated tickets tables
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  ticket_id UUID NOT NULL,
  status TEXT NOT NULL DEFAULT 'created',
  expires_at TIMESTAMPTZ NOT NULL,
  version INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_ticket_id ON orders(ticket_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

-- Replicated tickets table for order service
CREATE TABLE IF NOT EXISTS tickets (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  price BIGINT NOT NULL,
  user_id UUID NOT NULL,
  order_id UUID,
  version INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tickets_order_id ON tickets(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tickets_version ON tickets(version);
