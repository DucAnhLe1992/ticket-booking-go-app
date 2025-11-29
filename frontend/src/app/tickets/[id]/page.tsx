'use client';

import { use } from 'react';
import { useQuery } from '@tanstack/react-query';
import { ticketsApi } from '@/lib/api/tickets';
import { TicketDetail } from '@/components/tickets/TicketDetail';
import { Button } from '@/components/ui/button';
import Link from 'next/link';

export default function TicketDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  
  const { data: ticket, isLoading, error } = useQuery({
    queryKey: ['ticket', id],
    queryFn: () => ticketsApi.get(id),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
      </div>
    );
  }

  if (error || !ticket) {
    return (
      <div className="text-center py-12">
        <p className="text-destructive mb-4">Failed to load ticket</p>
        <Link href="/tickets">
          <Button>Back to Tickets</Button>
        </Link>
      </div>
    );
  }

  return (
    <div>
      <Link href="/tickets">
        <Button variant="ghost" className="mb-6">← Back to Tickets</Button>
      </Link>
      <TicketDetail ticket={ticket} />
    </div>
  );
}
