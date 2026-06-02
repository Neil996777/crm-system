import { apiRequest } from './client';

type GatewayEnvelope<T> = {
  data: T;
};

export type HistoryEvent = {
  eventUid: string;
  eventId: string;
  eventVersion: number;
  actorUserId: string;
  actorRole: string;
  actorDisplay: string;
  action: string;
  resourceType: string;
  resourceId: string;
  result: string;
  beforeSummary: Record<string, unknown>;
  afterSummary: Record<string, unknown>;
  safeSummary: string;
  occurredAt: string;
  eventHash: string;
};

export async function getRecordHistory(resource: string, id: string) {
  const response = await apiRequest<GatewayEnvelope<{ events: HistoryEvent[] }>>(`/api/${resource}/${id}/history`);
  return response.data;
}
