import { apiRequest } from './client';
import { DuplicateWarningResult } from './duplicates';

export type LeadStatus = 'Unassigned' | 'Pending Qualification' | 'Valid' | 'Invalid' | 'Converted To Opportunity';

export type Lead = {
  id: string;
  leadName: string;
  companyName: string;
  email: string;
  phone: string;
  source: string;
  status: LeadStatus;
  ownerId: string;
  needSummary: string;
  invalidReason: string;
  convertedAccountId: string;
  convertedOpportunityId: string;
  version: number;
  updatedAt: string;
};

type GatewayEnvelope<T> = {
  data: T;
};

export type LeadCreateInput = {
  leadName?: string;
  companyName: string;
  email?: string;
  phone?: string;
  source: string;
  ownerId?: string;
  needSummary?: string;
  proceedWarningToken?: string;
};

export type ConversionResult = {
  leadId: string;
  accountId: string;
  contactIds: string[];
  opportunityId: string;
  status: LeadStatus;
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

export async function listLeads(search: string) {
  const query = search.trim() ? `?search=${encodeURIComponent(search.trim())}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: Lead[] }>>(`/api/leads${query}`);
  return unwrap(response);
}

export async function getLead(id: string) {
  const response = await apiRequest<GatewayEnvelope<Lead>>(`/api/leads/${id}`);
  return unwrap(response);
}

export async function createLead(input: LeadCreateInput) {
  const response = await apiRequest<GatewayEnvelope<Lead>>('/api/leads', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function checkLeadDuplicate(input: { companyName?: string; email?: string; phone?: string }) {
  const response = await apiRequest<GatewayEnvelope<DuplicateWarningResult>>('/api/leads/duplicate-checks', {
    method: 'POST',
    body: JSON.stringify({ targetType: 'lead', candidate: input })
  });
  return unwrap(response);
}

export async function qualifyValid(lead: Lead) {
  const response = await apiRequest<GatewayEnvelope<Lead>>(`/api/leads/${lead.id}/qualify-valid`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion: lead.version })
  });
  return unwrap(response);
}

export async function qualifyInvalid(lead: Lead, invalidReason: string) {
  const response = await apiRequest<GatewayEnvelope<Lead>>(`/api/leads/${lead.id}/qualify-invalid`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion: lead.version, invalidReason })
  });
  return unwrap(response);
}

export async function restoreInvalid(lead: Lead) {
  const response = await apiRequest<GatewayEnvelope<Lead>>(`/api/leads/${lead.id}/restore-invalid`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion: lead.version })
  });
  return unwrap(response);
}

export async function convertLead(lead: Lead, expectedAmount: string, expectedCloseDate: string) {
  const response = await apiRequest<GatewayEnvelope<ConversionResult>>(`/api/leads/${lead.id}/convert`, {
    method: 'POST',
    body: JSON.stringify({
      idempotencyKey: `ui-${lead.id}-${Date.now()}`,
      target: {
        accountInput: {
          companyName: lead.companyName || lead.leadName,
          customerStatus: 'Prospect',
          ownerId: lead.ownerId
        },
        opportunityInput: {
          ownerId: lead.ownerId,
          stage: 'New Opportunity',
          expectedAmount,
          expectedCloseDate,
          title: lead.companyName || lead.leadName
        }
      }
    })
  });
  return unwrap(response);
}
