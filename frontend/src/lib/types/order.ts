export enum OrderStatus {
  Created = 'created',
  Cancelled = 'cancelled',
  Complete = 'complete',
}

export interface Order {
  id: string;
  userId: string;
  status: OrderStatus;
  expiresAt: string;
  ticket: {
    id: string;
    title: string;
    price: number;
  };
  version: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateOrderInput {
  ticketId: string;
}
