'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { paymentsApi } from '@/lib/api/payments';
import type { Order } from '@/lib/types/order';
import { formatPrice } from '@/lib/utils/format';
import { toast } from 'sonner';

interface StripeCheckoutProps {
  order: Order;
}

export function StripeCheckout({ order }: StripeCheckoutProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [cardNumber, setCardNumber] = useState('4242424242424242');
  const [isProcessing, setIsProcessing] = useState(false);

  const createPayment = useMutation({
    mutationFn: () => paymentsApi.create({ orderId: order.id, token: 'tok_visa' }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.invalidateQueries({ queryKey: ['order', order.id] });
      toast.success('Payment successful!');
      router.push(`/orders/${order.id}?success=true`);
    },
    onError: (error) => {
      const message = error instanceof Error ? error.message : 'Payment failed';
      toast.error(message);
      setIsProcessing(false);
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsProcessing(true);
    createPayment.mutate();
  };

  return (
    <Card className="max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Complete Payment</CardTitle>
        <CardDescription>
          Enter your payment details to complete the purchase
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="rounded-lg border p-4 bg-muted/50">
            <p className="text-sm text-muted-foreground mb-1">Total Amount</p>
            <p className="text-3xl font-bold">{formatPrice(order.ticket.price)}</p>
          </div>

          <div className="space-y-4">
            <div>
              <Label htmlFor="cardNumber">Card Number</Label>
              <Input
                id="cardNumber"
                value={cardNumber}
                onChange={(e) => setCardNumber(e.target.value)}
                placeholder="4242 4242 4242 4242"
                maxLength={16}
              />
              <p className="text-xs text-muted-foreground mt-1">
                Use 4242424242424242 for testing
              </p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="expiry">Expiry Date</Label>
                <Input
                  id="expiry"
                  placeholder="MM/YY"
                  defaultValue="12/25"
                />
              </div>
              <div>
                <Label htmlFor="cvc">CVC</Label>
                <Input
                  id="cvc"
                  placeholder="123"
                  maxLength={3}
                  defaultValue="123"
                />
              </div>
            </div>
          </div>

          <Button
            type="submit"
            className="w-full"
            size="lg"
            disabled={isProcessing}
          >
            {isProcessing ? 'Processing...' : `Pay ${formatPrice(order.ticket.price)}`}
          </Button>

          <p className="text-xs text-center text-muted-foreground">
            Your payment is secure and encrypted
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
