'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { ordersApi } from '@/lib/api/orders';
import type { Ticket } from '@/lib/types/ticket';
import { formatPrice, formatDate } from '@/lib/utils/format';
import { toast } from 'sonner';

interface TicketDetailProps {
  ticket: Ticket;
}

export function TicketDetail({ ticket }: TicketDetailProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [isProcessing, setIsProcessing] = useState(false);

  const createOrder = useMutation({
    mutationFn: () => ordersApi.create({ ticketId: ticket.id }),
    onSuccess: (order) => {
      queryClient.invalidateQueries({ queryKey: ['tickets'] });
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      toast.success('Order created! Complete your purchase.');
      router.push(`/orders/${order.id}/payment`);
    },
    onError: (error) => {
      const message = error instanceof Error ? error.message : 'Failed to create order';
      toast.error(message);
      setIsProcessing(false);
    },
  });

  const handlePurchase = () => {
    setIsProcessing(true);
    createOrder.mutate();
  };

  const isReserved = !!ticket.orderId;

  return (
    <Card className="max-w-2xl mx-auto">
      <CardHeader>
        <div className="flex justify-between items-start">
          <div>
            <CardTitle className="text-3xl mb-2">{ticket.title}</CardTitle>
            <CardDescription>Ticket ID: {ticket.id}</CardDescription>
          </div>
          {isReserved && (
            <Badge variant="secondary" className="text-sm">
              Reserved
            </Badge>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        <div>
          <p className="text-4xl font-bold text-primary">
            {formatPrice(ticket.price)}
          </p>
        </div>

        <Separator />

        <div className="space-y-2">
          <div className="flex justify-between">
            <span className="text-muted-foreground">Created</span>
            <span className="font-medium">{formatDate(ticket.createdAt)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Last Updated</span>
            <span className="font-medium">{formatDate(ticket.updatedAt)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Version</span>
            <span className="font-medium">#{ticket.version}</span>
          </div>
        </div>

        <Separator />

        <Button
          onClick={handlePurchase}
          disabled={isReserved || isProcessing}
          className="w-full"
          size="lg"
        >
          {isProcessing
            ? 'Processing...'
            : isReserved
            ? 'Not Available'
            : 'Purchase Now'}
        </Button>

        {isReserved && (
          <p className="text-sm text-center text-muted-foreground">
            This ticket is currently reserved by another user.
          </p>
        )}
      </CardContent>
    </Card>
  );
}
