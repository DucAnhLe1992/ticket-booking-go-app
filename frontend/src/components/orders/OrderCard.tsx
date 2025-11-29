'use client';

import Link from 'next/link';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Order } from '@/lib/types/order';
import { OrderStatus } from '@/lib/types/order';
import { formatPrice, formatDate, formatStatus } from '@/lib/utils/format';

interface OrderCardProps {
  order: Order;
}

export function OrderCard({ order }: OrderCardProps) {
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
    <Card>
      <CardHeader>
        <div className="flex justify-between items-start">
          <CardTitle className="text-lg">{order.ticket.title}</CardTitle>
          <Badge variant={getBadgeVariant(order.status)}>
            {formatStatus(order.status)}
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="flex justify-between">
            <span className="text-muted-foreground">Price</span>
            <span className="font-semibold">{formatPrice(order.ticket.price)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Created</span>
            <span className="text-sm">{formatDate(order.createdAt)}</span>
          </div>
          {order.status === OrderStatus.Created && (
            <div className="flex justify-between">
              <span className="text-muted-foreground">Expires</span>
              <span className="text-sm text-amber-600">{formatDate(order.expiresAt)}</span>
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter>
        <Link href={`/orders/${order.id}`} className="w-full">
          <Button variant="outline" className="w-full">
            View Details
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
