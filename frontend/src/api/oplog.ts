import { apiRequest } from './client';
import { HistoryEvent } from './history';

type GatewayEnvelope<T> = {
  data: T;
};

export async function getOperationLog() {
  const response = await apiRequest<GatewayEnvelope<{ events: HistoryEvent[] }>>('/api/operation-log');
  return response.data;
}
