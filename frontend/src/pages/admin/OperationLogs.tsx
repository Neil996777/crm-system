import { useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { HistoryEvent } from '../../api/history';
import { getOperationLog } from '../../api/oplog';

export function OperationLogs() {
  const [events, setEvents] = useState<HistoryEvent[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError('');
      try {
        const response = await getOperationLog();
        setEvents(response.events);
      } catch (caught) {
        const apiError = caught as ApiError;
        setError(apiError.safeMessage || 'Permission denied.');
      } finally {
        setLoading(false);
      }
    }
    void load();
  }, []);

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Operation Logs</h1>
          <p>Administrator-only audit events.</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      {loading && <p className="emptyState">Loading operation logs...</p>}
      {!loading && !error && (
        <table className="dataTable" aria-label="Operation log table">
          <thead>
            <tr>
              <th>Event</th>
              <th>Action</th>
              <th>Actor</th>
              <th>Resource</th>
              <th>Occurred</th>
              <th>Result</th>
              <th>Before</th>
              <th>After</th>
            </tr>
          </thead>
          <tbody>
            {events.map((event) => (
              <tr key={event.eventUid}>
                <td>{event.eventId}</td>
                <td>{event.action}</td>
                <td>{event.actorUserId}</td>
                <td>{event.resourceType} {event.resourceId}</td>
                <td>{formatDate(event.occurredAt)}</td>
                <td>{event.result}</td>
                <td>Before: {summaryText(event.beforeSummary)}</td>
                <td>After: {summaryText(event.afterSummary)}</td>
              </tr>
            ))}
            {events.length === 0 && (
              <tr>
                <td colSpan={8}>No operation logs found.</td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </main>
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
