import { apiRequest } from './client';

type GatewayEnvelope<T> = { data: T };

export type ImportRowError = {
  rowNumber: number;
  field: string;
  code: string;
  safeMessage: string;
};

export type ImportRun = {
  runId: string;
  objectType: string;
  filename: string;
  status: string;
  totalRows: number;
  successCount: number;
  failureCount: number;
  rowErrors: ImportRowError[];
  operationLogStatus: string;
  cleanupStatus: string;
  retainedUntil: string;
};

export async function startImport(input: { objectType: string; filename: string; content: string }) {
  const response = await apiRequest<GatewayEnvelope<ImportRun>>('/api/imports', {
    method: 'POST',
    body: JSON.stringify(input)
  });
  return response.data;
}
