import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type ReminderRow = {
  id: string;
  sourceService: string;
  type: 'task_due' | 'task_overdue' | 'contract_pending_signature' | 'payment_due' | 'payment_overdue';
  relatedRecord: {
    type: string;
    id: string;
    display: string;
  };
  ownerDisplay: string;
  dueDate: string;
  status: string;
  priority: string;
  version: number;
};

export type ReminderResponse = {
  timezone: string;
  businessDate: string;
  rows: ReminderRow[];
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

export async function listReminders(businessDate: string) {
  const query = businessDate.trim() ? `?businessDate=${encodeURIComponent(businessDate.trim())}` : '';
  const response = await apiRequest<GatewayEnvelope<ReminderResponse>>(`/api/reminders${query}`);
  return unwrap(response);
}
