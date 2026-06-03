import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type OverviewMetrics = {
  leadCount: number;
  opportunityCount: number;
  taskCount: number;
  wonCount: number;
  lostCount: number;
  quoteAmount: string;
  contractAmount: string;
  paidAmount: string;
  receivableAmount: string;
};

export type GroupRow = {
  key: string;
  label: string;
  count: number;
  amount: string;
};

export type ManagerOverview = {
  scope: string;
  filters: { teamId: string; archived: string };
  currency: string;
  metrics: OverviewMetrics;
  pipeline: GroupRow[];
  emptyState: boolean;
};

export async function getManagerOverview() {
  const response = await apiRequest<GatewayEnvelope<ManagerOverview>>('/api/reports/team-overview');
  return response.data;
}
