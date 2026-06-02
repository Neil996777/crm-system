import { useEffect, useState } from 'react';
import { ApiError } from '../api/client';
import { HistoryEvent, getRecordHistory } from '../api/history';

type Props = {
  resource: string;
  recordId: string;
  reloadKey?: string | number;
};

export function HistoryTimeline({ resource, recordId, reloadKey }: Props) {
  const [events, setEvents] = useState<HistoryEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let cancelled = false;
    async function load() {
      setLoading(true);
      setError('');
      try {
        const response = await getRecordHistory(resource, recordId);
        if (!cancelled) setEvents(response.events);
      } catch (caught) {
        const apiError = caught as ApiError;
        if (!cancelled) setError(apiError.safeMessage || 'Permission denied.');
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    void load();
    return () => {
      cancelled = true;
    };
  }, [resource, recordId, reloadKey]);

  return (
    <section className="historyTimeline" aria-label="Record history">
      <div className="sectionTitle">
        <h3>History</h3>
        <span>Read-only</span>
      </div>
      {loading && <p className="emptyState">Loading history...</p>}
      {error && <p role="alert" className="error">{error}</p>}
      {!loading && !error && events.length === 0 && <p className="emptyState">No history events found.</p>}
      {!loading && !error && events.length > 0 && (
        <ol className="timelineList">
          {events.map((event) => (
            <li key={event.eventUid} className="timelineItem">
              <div className="timelineHeader">
                <strong>{event.eventId}</strong>
                <span>{event.result}</span>
              </div>
              <p>{event.safeSummary}</p>
              <dl className="timelineMeta">
                <div>
                  <dt>Actor</dt>
                  <dd>Actor: {event.actorUserId}</dd>
                </div>
                <div>
                  <dt>Resource</dt>
                  <dd>Resource: {event.resourceType} {event.resourceId}</dd>
                </div>
                <div>
                  <dt>Occurred</dt>
                  <dd>Occurred: {formatDate(event.occurredAt)}</dd>
                </div>
                <div>
                  <dt>Before</dt>
                  <dd>Before: {summaryText(event.beforeSummary)}</dd>
                </div>
                <div>
                  <dt>After</dt>
                  <dd>After: {summaryText(event.afterSummary)}</dd>
                </div>
              </dl>
            </li>
          ))}
        </ol>
      )}
    </section>
  );
}

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toISOString();
}

function summaryText(summary: Record<string, unknown> | undefined) {
  if (!summary || Object.keys(summary).length === 0) return 'None';
  return Object.entries(summary)
    .filter(([, value]) => value !== '' && value !== null && value !== undefined)
    .map(([key, value]) => `${key}: ${String(value)}`)
    .join(', ');
}
