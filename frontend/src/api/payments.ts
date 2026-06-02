import { Contract } from './contracts';
import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type PaymentPlan = {
  id: string;
  contractId: string;
  dueAmount: string;
  dueDate: string;
  currency: string;
  status: string;
  version: number;
  updatedAt: string;
};

export type ActualPayment = {
  paymentId: string;
  contractId: string;
  amount: string;
  paymentDate: string;
  paymentStatus: string;
  remainingAmount: string;
  version: number;
  updatedAt: string;
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

export async function listPaymentContracts(search: string) {
  const query = search.trim() ? `?search=${encodeURIComponent(search.trim())}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: Contract[] }>>(`/api/contracts${query}`);
  return unwrap(response);
}

export async function createPaymentPlan(contractId: string, input: { dueAmount: string; dueDate: string; currency: string }) {
  const response = await apiRequest<GatewayEnvelope<PaymentPlan>>(`/api/contracts/${contractId}/payment-plans`, {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function recordPayment(contractId: string, input: { idempotencyKey: string; amount: string; paymentDate: string; note: string; currency: string }) {
  const response = await apiRequest<GatewayEnvelope<ActualPayment>>(`/api/contracts/${contractId}/payments`, {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}
