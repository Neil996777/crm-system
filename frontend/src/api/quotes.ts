import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type Quote = {
  id: string;
  opportunityId: string;
  customerId: string;
  amount: string;
  status: string;
  validityEnd: string;
  ownerId: string;
  version: number;
  updatedAt: string;
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

export async function listQuotes(search: string) {
  const query = search.trim() ? `?search=${encodeURIComponent(search.trim())}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: Quote[] }>>(`/api/quotes${query}`);
  return unwrap(response);
}

export async function getQuote(id: string) {
  const response = await apiRequest<GatewayEnvelope<Quote>>(`/api/quotes/${id}`);
  return unwrap(response);
}

export async function createQuote(input: { opportunityId: string; customerId: string; amount: string; validityEnd: string; ownerId: string }) {
  const response = await apiRequest<GatewayEnvelope<Quote>>('/api/quotes', {
    method: 'POST',
    body: JSON.stringify({ ...input, status: 'Draft' })
  });
  return unwrap(response);
}

export async function changeQuoteStatus(id: string, expectedVersion: number, toStatus: string) {
  const response = await apiRequest<GatewayEnvelope<Quote>>(`/api/quotes/${id}/status`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion, toStatus })
  });
  return unwrap(response);
}
