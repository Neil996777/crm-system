import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type Opportunity = {
  id: string;
  customerId: string;
  ownerId: string;
  stage: string;
  expectedAmount: string;
  expectedCloseDate: string;
  title: string;
  closeDate?: string;
  wonContractId?: string;
  lostReasonCode?: string;
  lostReasonDetail?: string;
  closedAt?: string;
  version: number;
  updatedAt: string;
};

export type CloseResult = {
  opportunityId: string;
  status: string;
  closedAt: string;
  version: number;
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

export async function listOpportunities(search: string, stage = '') {
  const params = new URLSearchParams();
  if (search.trim()) params.set('search', search.trim());
  if (stage.trim()) params.set('stage', stage.trim());
  const query = params.toString() ? `?${params.toString()}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: Opportunity[] }>>(`/api/opportunities${query}`);
  return unwrap(response);
}

export async function getOpportunity(id: string) {
  const response = await apiRequest<GatewayEnvelope<Opportunity>>(`/api/opportunities/${id}`);
  return unwrap(response);
}

export async function createOpportunity(input: {
  customerId: string;
  ownerId: string;
  stage: string;
  expectedAmount: string;
  expectedCloseDate: string;
  title: string;
}) {
  const response = await apiRequest<GatewayEnvelope<Opportunity>>('/api/opportunities', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function changeOpportunityStage(id: string, expectedVersion: number, toStage: string) {
  const response = await apiRequest<GatewayEnvelope<Opportunity>>(`/api/opportunities/${id}/stage`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion, toStage })
  });
  return unwrap(response);
}

export async function closeOpportunityWon(id: string, input: { expectedVersion: number; contractId: string; closeDate: string }) {
  const response = await apiRequest<GatewayEnvelope<CloseResult>>(`/api/opportunities/${id}/close-won`, {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function closeOpportunityLost(id: string, input: { expectedVersion: number; closeDate: string; lostReason: { code: string; detail: string } }) {
  const response = await apiRequest<GatewayEnvelope<CloseResult>>(`/api/opportunities/${id}/close-lost`, {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}
