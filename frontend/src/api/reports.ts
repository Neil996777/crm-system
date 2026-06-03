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

export type PaymentGroupRow = GroupRow & {
  dueAmount: string;
  paidAmount: string;
};

export type ManagerOverview = {
  scope: string;
  filters: { teamId: string; archived: string };
  currency: string;
  metrics: OverviewMetrics;
  pipeline: GroupRow[];
  emptyState: boolean;
};

export type BasicReport = {
  scope: string;
  filters: { teamId: string; archived: string; from: string; to: string; groupBy: string };
  currency: string;
  metrics: OverviewMetrics;
  breakdowns: {
    leadsByStatus: GroupRow[];
    opportunitiesByStage: GroupRow[];
    quotesByStatus: GroupRow[];
    contractsByStatus: GroupRow[];
    paymentsByStatus: PaymentGroupRow[];
  };
  groups: GroupRow[];
  emptyState: boolean;
};

export async function getManagerOverview() {
  const response = await apiRequest<GatewayEnvelope<ManagerOverview>>('/api/reports/team-overview');
  return response.data;
}

export async function getBasicReport() {
  const response = await apiRequest<GatewayEnvelope<BasicReport>>('/api/reports/sales-overview');
  return response.data;
}
