import type { AppView } from './Nav';

export type RecordNavigationTarget = {
  view: AppView;
  recordId?: string;
};

export function targetForRelatedRecord(relatedType: string, relatedId: string): RecordNavigationTarget | null {
  if (!relatedId) return null;
  const normalized = relatedType.toLowerCase();
  if (normalized === 'opportunity') return { view: 'opportunities', recordId: relatedId };
  if (normalized === 'lead') return { view: 'leads', recordId: relatedId };
  if (normalized === 'account') return { view: 'accounts', recordId: relatedId };
  if (normalized === 'contact') return { view: 'contacts', recordId: relatedId };
  if (normalized === 'quote') return { view: 'quotes', recordId: relatedId };
  if (normalized === 'contract') return { view: 'contracts', recordId: relatedId };
  if (normalized === 'payment') return { view: 'payments', recordId: relatedId };
  if (normalized === 'task') return { view: 'tasks', recordId: relatedId };
  return null;
}
