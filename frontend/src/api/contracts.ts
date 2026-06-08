import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type Contract = {
  id: string;
  quoteId: string;
  opportunityId: string;
  customerId: string;
  amount: string;
  status: string;
  contractNote: string;
  expectedSignedDate: string;
  signedEffectiveDate?: string;
  amountDifferenceReason?: string;
  ownerId: string;
  archived?: boolean;
  archivedAt?: string;
  archivedBy?: string;
  archiveReason?: string;
  version: number;
  updatedAt: string;
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

export async function listContracts(search: string, includeArchived = false) {
  const params = new URLSearchParams();
  if (search.trim()) params.set('search', search.trim());
  if (includeArchived) params.set('includeArchived', 'true');
  const query = params.toString() ? `?${params.toString()}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: Contract[] }>>(`/api/contracts${query}`);
  return unwrap(response);
}

export async function getContract(id: string) {
  const response = await apiRequest<GatewayEnvelope<Contract>>(`/api/contracts/${id}`);
  return unwrap(response);
}

export async function createContract(input: {
  quoteId: string;
  opportunityId: string;
  customerId: string;
  amount: string;
  contractNote: string;
  expectedSignedDate: string;
  amountDifferenceReason: string;
  ownerId: string;
}) {
  const response = await apiRequest<GatewayEnvelope<Contract>>('/api/contracts', {
    method: 'POST',
    body: JSON.stringify({ ...input, status: 'Pending Signature' })
  });
  return unwrap(response);
}

export async function changeContractStatus(id: string, expectedVersion: number, toStatus: string, signedEffectiveDate: string) {
  const response = await apiRequest<GatewayEnvelope<Contract>>(`/api/contracts/${id}/status`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion, toStatus, signedEffectiveDate })
  });
  return unwrap(response);
}

export async function archiveContract(contract: Contract, reason: string) {
  const response = await apiRequest<GatewayEnvelope<Contract>>(`/api/contracts/${contract.id}/archive`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion: contract.version, reason })
  });
  return unwrap(response);
}
