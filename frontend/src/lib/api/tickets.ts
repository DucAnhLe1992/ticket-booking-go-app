import { apiClient } from './client';
import type { Ticket, CreateTicketInput, UpdateTicketInput } from '../types/ticket';

export const ticketsApi = {
  list: async (): Promise<Ticket[]> => {
    const { data } = await apiClient.get<Ticket[]>('/tickets');
    return data;
  },

  get: async (id: string): Promise<Ticket> => {
    const { data } = await apiClient.get<Ticket>(`/tickets/${id}`);
    return data;
  },

  create: async (input: CreateTicketInput): Promise<Ticket> => {
    const { data } = await apiClient.post<Ticket>('/tickets', input);
    return data;
  },

  update: async (id: string, input: UpdateTicketInput): Promise<Ticket> => {
    const { data } = await apiClient.put<Ticket>(`/tickets/${id}`, input);
    return data;
  },
};
