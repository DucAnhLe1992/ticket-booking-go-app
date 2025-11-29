export interface Payment {
  id: string;
  orderId: string;
  stripeId: string;
  createdAt: string;
}

export interface CreatePaymentInput {
  orderId: string;
  token: string;
}

export interface CreatePaymentResponse {
  payment: Payment;
  success: boolean;
}
