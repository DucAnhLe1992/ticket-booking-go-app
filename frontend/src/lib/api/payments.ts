import { apiClient } from './client';
import type { CreatePaymentInput, CreatePaymentResponse } from '../types/payment';

export const paymentsApi = {
  create: async (input: CreatePaymentInput): Promise<CreatePaymentResponse> => {
    const { data } = await apiClient.post<CreatePaymentResponse>('/payments', input);
    return data;
  },
};
