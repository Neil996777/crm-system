import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type Activity = {
  id: string;
  relatedType: string;
  relatedId: string;
  activityType: string;
  content: string;
  ownerId: string;
  occurredAt: string;
};

export type Note = {
  id: string;
  relatedType: string;
  relatedId: string;
  content: string;
  ownerId: string;
  occurredAt: string;
};

export type WorkTask = {
  id: string;
  relatedType: string;
  relatedId: string;
  title: string;
  dueDate: string;
  status: string;
  ownerId: string;
  version: number;
};

function unwrap<T>(response: GatewayEnvelope<T>) {
  return response.data;
}

function relatedQuery(relatedType: string, relatedId: string) {
  const params = new URLSearchParams();
  if (relatedType) params.set('relatedType', relatedType);
  if (relatedId) params.set('relatedId', relatedId);
  return params.toString() ? `?${params.toString()}` : '';
}

export async function listActivities(relatedType: string, relatedId: string) {
  const response = await apiRequest<GatewayEnvelope<{ items: Activity[] }>>(`/api/activities${relatedQuery(relatedType, relatedId)}`);
  return unwrap(response);
}

export async function createActivity(input: { relatedType: string; relatedId: string; activityType: string; content: string; ownerId: string }) {
  const response = await apiRequest<GatewayEnvelope<Activity>>('/api/activities', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function listNotes(relatedType: string, relatedId: string) {
  const response = await apiRequest<GatewayEnvelope<{ items: Note[] }>>(`/api/notes${relatedQuery(relatedType, relatedId)}`);
  return unwrap(response);
}

export async function createNote(input: { relatedType: string; relatedId: string; content: string; ownerId: string }) {
  const response = await apiRequest<GatewayEnvelope<Note>>('/api/notes', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function listTasks(input: { relatedType?: string; relatedId?: string; activeOnly?: boolean; businessDate?: string } = {}) {
  const params = new URLSearchParams();
  if (input.relatedType) params.set('relatedType', input.relatedType);
  if (input.relatedId) params.set('relatedId', input.relatedId);
  if (input.activeOnly) params.set('activeOnly', 'true');
  if (input.businessDate) params.set('businessDate', input.businessDate);
  const query = params.toString() ? `?${params.toString()}` : '';
  const response = await apiRequest<GatewayEnvelope<{ items: WorkTask[] }>>(`/api/tasks${query}`);
  return unwrap(response);
}

export async function createTask(input: { relatedType: string; relatedId: string; title: string; dueDate: string; ownerId: string }) {
  const response = await apiRequest<GatewayEnvelope<WorkTask>>('/api/tasks', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return unwrap(response);
}

export async function changeTaskStatus(id: string, toStatus: string) {
  const response = await apiRequest<GatewayEnvelope<WorkTask>>(`/api/tasks/${id}/status`, {
    method: 'POST',
    body: JSON.stringify({ toStatus })
  });
  return unwrap(response);
}
