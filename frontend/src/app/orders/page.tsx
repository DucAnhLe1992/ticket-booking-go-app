'use client';

import { useQuery } from '@tanstack/react-query';
import { ordersApi } from '@/lib/api/orders';
import { OrderList } from '@/components/orders/OrderList';
import { AuthGuard } from '@/components/layout/AuthGuard';

export default function OrdersPage() {
  const { data: orders, isLoading, error } = useQuery({
    queryKey: ['orders'],
    queryFn: ordersApi.list,
  });

  return (
    <AuthGuard>
      <div>
        <h1 className="text-4xl font-bold mb-8">My Orders</h1>
        
        {isLoading && (
          <div className="flex items-center justify-center min-h-[400px]">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
          </div>
        )}

        {error && (
          <div className="text-center py-12">
            <p className="text-destructive">Failed to load orders</p>
          </div>
        )}

        {!isLoading && !error && <OrderList orders={orders || []} />}
      </div>
    </AuthGuard>
  );
}
