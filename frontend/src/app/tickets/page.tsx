'use client';

import { useQuery } from '@tanstack/react-query';
import { ticketsApi } from '@/lib/api/tickets';
import { TicketList } from '@/components/tickets/TicketList';

export default function TicketsPage() {
  const { data: tickets, isLoading, error } = useQuery({
    queryKey: ['tickets'],
    queryFn: ticketsApi.list,
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-destructive">Failed to load tickets</p>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-4xl font-bold">Available Tickets</h1>
      </div>
      <TicketList tickets={tickets || []} />
    </div>
  );
}
