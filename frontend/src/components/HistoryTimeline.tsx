import { useEffect, useState } from 'react';
import { ApiError } from '../api/client';
import { HistoryEvent, getRecordHistory } from '../api/history';
import { labelFor, localizeError, localizeMessage, objectTypeLabel, resultLabel, summaryTextZh } from '../i18n/labels';

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
        if (!cancelled) setError(localizeError(apiError, '权限不足。'));
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
    <section className="historyTimeline" aria-label="记录历史">
      <div className="sectionTitle">
        <h3>历史</h3>
        <span>只读</span>
      </div>
      {loading && <p className="emptyState">正在加载历史...</p>}
      {error && <p role="alert" className="error">{error}</p>}
      {!loading && !error && events.length === 0 && <p className="emptyState">暂无历史事件。</p>}
      {!loading && !error && events.length > 0 && (
        <ol className="timelineList">
          {events.map((event) => (
            <li key={event.eventUid} className="timelineItem">
              <div className="timelineHeader">
                <strong>{event.eventId}</strong>
                <span>{labelFor(resultLabel, event.result)}</span>
              </div>
              <p>{localizeMessage(event.safeSummary, event.safeSummary)}</p>
              <dl className="timelineMeta">
                <div>
                  <dt>操作者</dt>
                  <dd>操作者：{event.actorUserId}</dd>
                </div>
                <div>
                  <dt>资源</dt>
                  <dd>资源：{labelFor(objectTypeLabel, event.resourceType)} {event.resourceId}</dd>
                </div>
                <div>
                  <dt>发生时间</dt>
                  <dd>发生时间：{formatDate(event.occurredAt)}</dd>
                </div>
                <div>
                  <dt>变更前</dt>
                  <dd>变更前：{summaryTextZh(event.beforeSummary)}</dd>
                </div>
                <div>
                  <dt>变更后</dt>
                  <dd>变更后：{summaryTextZh(event.afterSummary)}</dd>
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
