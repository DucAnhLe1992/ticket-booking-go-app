export interface Ticket {
  id: string;
  title: string;
  price: number;
  userId: string;
  version: number;
  orderId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateTicketInput {
  title: string;
  price: number;
}

export interface UpdateTicketInput {
  title?: string;
  price?: number;
}
