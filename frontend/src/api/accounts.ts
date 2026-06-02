import { apiRequest } from './client';
import { DuplicateWarningResult } from './duplicates';

type GatewayEnvelope<T> = { data: T };

export type Account = {
  id: string;
  companyName: string;
  customerStatus: string;
  ownerId: string;
  archived: boolean;
  archivedAt: string;
  archivedBy: string;
  archiveReason: string;
  version: number;
  updatedAt: string;
};

export type Contact = {
  id: string;
  accountId: string;
  accountName: string;
  contactName: string;
  email: string;
  phone: string;
  roleNote: string;
  version: number;
  updatedAt: string;
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

export type ArchiveObligation = {
  type: string;
  id: string;
  service: string;
  status: string;
  dueDate: string;
  ownerDisplay: string;
  blocking: boolean;
  safeMessage: string;
};

export type ArchiveEligibility = {
  resourceType: string;
  resourceId: string;
  canArchive: boolean;
  recordVersion: number;
  obligations: ArchiveObligation[];
};

export async function listAccounts(search: string, includeArchived = false) {
  const params = new URLSearchParams();
  if (search.trim()) params.set('search', search.trim());
  if (includeArchived) params.set('includeArchived', 'true');
  const query = params.toString() ? `?${params.toString()}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: Account[] }>>(`/api/accounts${query}`);
  return unwrap(response);
}

export async function getAccount(id: string) {
  const response = await apiRequest<GatewayEnvelope<Account>>(`/api/accounts/${id}`);
  return unwrap(response);
}

export type AccountCreateInput = { companyName: string; customerStatus: string; ownerId: string; proceedWarningToken?: string };

export async function createAccount(input: AccountCreateInput) {
  const response = await apiRequest<GatewayEnvelope<Account>>('/api/accounts', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function getAccountArchiveEligibility(accountId: string) {
  const response = await apiRequest<GatewayEnvelope<ArchiveEligibility>>(`/api/accounts/${accountId}/archive-eligibility`);
  return unwrap(response);
}

export async function archiveAccount(accountId: string, expectedVersion: number, reason: string) {
  const response = await apiRequest<GatewayEnvelope<Account>>(`/api/accounts/${accountId}/archive`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion, reason })
  });
  return unwrap(response);
}

export async function checkAccountDuplicate(input: { companyName: string }) {
  const response = await apiRequest<GatewayEnvelope<DuplicateWarningResult>>('/api/accounts/duplicate-checks', {
    method: 'POST',
    body: JSON.stringify({ targetType: 'account', candidate: input })
  });
  return unwrap(response);
}

export async function checkContactDuplicate(input: { email?: string; phone?: string }) {
  const response = await apiRequest<GatewayEnvelope<DuplicateWarningResult>>('/api/accounts/duplicate-checks', {
    method: 'POST',
    body: JSON.stringify({ targetType: 'contact', candidate: input })
  });
  return unwrap(response);
}

export async function listContacts(accountId: string) {
  const response = await apiRequest<GatewayEnvelope<{ items: Contact[] }>>(`/api/accounts/${accountId}/contacts`);
  return unwrap(response);
}

export async function listAllContacts(search: string) {
  const query = search.trim() ? `?search=${encodeURIComponent(search.trim())}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: Contact[] }>>(`/api/contacts${query}`);
  return unwrap(response);
}

export async function getContact(id: string) {
  const response = await apiRequest<GatewayEnvelope<Contact>>(`/api/contacts/${id}`);
  return unwrap(response);
}

export async function createContact(accountId: string, input: { contactName: string; email: string; phone: string; roleNote: string; proceedWarningToken?: string }) {
  const response = await apiRequest<GatewayEnvelope<Contact>>(`/api/accounts/${accountId}/contacts`, {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}
