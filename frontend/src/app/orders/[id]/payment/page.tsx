'use client';

import { use } from 'react';
import { useQuery } from '@tanstack/react-query';
import Link from 'next/link';
import { ordersApi } from '@/lib/api/orders';
import { AuthGuard } from '@/components/layout/AuthGuard';
import { StripeCheckout } from '@/components/payment/StripeCheckout';
import { OrderCountdown } from '@/components/orders/OrderCountdown';
import { Button } from '@/components/ui/button';
import { OrderStatus } from '@/lib/types/order';

export default function PaymentPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);

  const { data: order, isLoading, error } = useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersApi.get(id),
  });

  if (isLoading) {
    return (
      <AuthGuard>
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
        </div>
      </AuthGuard>
    );
  }

  if (error || !order) {
    return (
      <AuthGuard>
        <div className="text-center py-12">
          <p className="text-destructive mb-4">Failed to load order</p>
          <Link href="/orders">
            <Button>Back to Orders</Button>
          </Link>
        </div>
      </AuthGuard>
    );
  }

  if (order.status !== OrderStatus.Created) {
    return (
      <AuthGuard>
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">
            This order is no longer available for payment
          </p>
          <Link href={`/orders/${order.id}`}>
            <Button>View Order</Button>
          </Link>
        </div>
      </AuthGuard>
    );
  }

  return (
    <AuthGuard>
      <div className="space-y-6">
        <div className="text-center">
          <h1 className="text-3xl font-bold mb-2">Complete Your Purchase</h1>
          <p className="text-muted-foreground">Order for: {order.ticket.title}</p>
        </div>

        <div className="max-w-md mx-auto">
          <OrderCountdown expiresAt={order.expiresAt} />
        </div>

        <StripeCheckout order={order} />

        <div className="text-center">
          <Link href={`/orders/${order.id}`}>
            <Button variant="ghost">Back to Order</Button>
          </Link>
        </div>
      </div>
    </AuthGuard>
  );
}
