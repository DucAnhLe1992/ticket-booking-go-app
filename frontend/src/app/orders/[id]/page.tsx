'use client';

import { use } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import Link from 'next/link';
import { ordersApi } from '@/lib/api/orders';
import { AuthGuard } from '@/components/layout/AuthGuard';
import { OrderCountdown } from '@/components/orders/OrderCountdown';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { OrderStatus } from '@/lib/types/order';
import { formatPrice, formatDate, formatStatus } from '@/lib/utils/format';
import { toast } from 'sonner';

export default function OrderDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const queryClient = useQueryClient();

  const { data: order, isLoading, error } = useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersApi.get(id),
  });

  const cancelOrder = useMutation({
    mutationFn: () => ordersApi.cancel(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['order', id] });
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      toast.success('Order cancelled successfully');
    },
    onError: (error) => {
      const message = error instanceof Error ? error.message : 'Failed to cancel order';
      toast.error(message);
    },
  });

  const handleExpire = () => {
    queryClient.invalidateQueries({ queryKey: ['order', id] });
  };

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

  const getBadgeVariant = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.Complete:
        return 'default';
      case OrderStatus.Cancelled:
        return 'destructive';
      case OrderStatus.Created:
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  return (
    <AuthGuard>
      <div>
        <Link href="/orders">
          <Button variant="ghost" className="mb-6">← Back to Orders</Button>
        </Link>

        <Card className="max-w-2xl mx-auto">
          <CardHeader>
            <div className="flex justify-between items-start">
              <div>
                <CardTitle className="text-3xl mb-2">{order.ticket.title}</CardTitle>
                <CardDescription>Order ID: {order.id}</CardDescription>
              </div>
              <Badge variant={getBadgeVariant(order.status)}>
                {formatStatus(order.status)}
              </Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            {order.status === OrderStatus.Created && (
              <OrderCountdown expiresAt={order.expiresAt} onExpire={handleExpire} />
            )}

            <div>
              <p className="text-4xl font-bold text-primary">
                {formatPrice(order.ticket.price)}
              </p>
            </div>

            <Separator />

            <div className="space-y-2">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Ticket ID</span>
                <span className="font-medium">{order.ticket.id}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Created</span>
                <span className="font-medium">{formatDate(order.createdAt)}</span>
              </div>
              {order.status === OrderStatus.Created && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Expires At</span>
                  <span className="font-medium text-amber-600">
                    {formatDate(order.expiresAt)}
                  </span>
                </div>
              )}
            </div>

            <Separator />

            <div className="space-y-3">
              {order.status === OrderStatus.Created && (
                <>
                  <Link href={`/orders/${order.id}/payment`} className="block">
                    <Button className="w-full" size="lg">
                      Complete Payment
                    </Button>
                  </Link>
                  <Button
                    variant="outline"
                    className="w-full"
                    onClick={() => cancelOrder.mutate()}
                    disabled={cancelOrder.isPending}
                  >
                    {cancelOrder.isPending ? 'Cancelling...' : 'Cancel Order'}
                  </Button>
                </>
              )}

              {order.status === OrderStatus.Complete && (
                <div className="text-center text-green-600 font-semibold">
                  ✓ Payment completed successfully!
                </div>
              )}

              {order.status === OrderStatus.Cancelled && (
                <div className="text-center text-muted-foreground">
                  This order has been cancelled.
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </AuthGuard>
  );
}
