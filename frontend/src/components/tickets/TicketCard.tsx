'use client';

import Link from 'next/link';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Ticket } from '@/lib/types/ticket';
import { formatPrice } from '@/lib/utils/format';

interface TicketCardProps {
  ticket: Ticket;
}

export function TicketCard({ ticket }: TicketCardProps) {
  const isReserved = !!ticket.orderId;

  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader>
        <div className="flex justify-between items-start">
          <CardTitle className="text-xl">{ticket.title}</CardTitle>
          {isReserved && (
            <Badge variant="secondary">Reserved</Badge>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <p className="text-3xl font-bold text-primary">
          {formatPrice(ticket.price)}
        </p>
      </CardContent>
      <CardFooter>
        <Link href={`/tickets/${ticket.id}`} className="w-full">
          <Button className="w-full" disabled={isReserved}>
            {isReserved ? 'Not Available' : 'View Details'}
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
