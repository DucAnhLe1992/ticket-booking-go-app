import { apiClient } from './client';
import type { Order, CreateOrderInput } from '../types/order';

export const ordersApi = {
  list: async (): Promise<Order[]> => {
    const { data } = await apiClient.get<Order[]>('/orders');
    return data;
  },

  get: async (id: string): Promise<Order> => {
    const { data } = await apiClient.get<Order>(`/orders/${id}`);
    return data;
  },

  create: async (input: CreateOrderInput): Promise<Order> => {
    const { data } = await apiClient.post<Order>('/orders', input);
    return data;
  },

  cancel: async (id: string): Promise<Order> => {
    const { data } = await apiClient.delete<Order>(`/orders/${id}`);
    return data;
  },
};
