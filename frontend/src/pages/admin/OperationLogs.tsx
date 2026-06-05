import { useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { HistoryEvent } from '../../api/history';
import { getOperationLog } from '../../api/oplog';
import { labelFor, objectTypeLabel, resultLabel, summaryTextZh } from '../../i18n/labels';

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
        setError(apiError.safeMessage || '权限不足。');
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
          <h1>操作日志</h1>
          <p>仅管理员可见的审计事件。</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      {loading && <p className="emptyState">正在加载操作日志...</p>}
      {!loading && !error && (
        <table className="dataTable" aria-label="操作日志表">
          <thead>
            <tr>
              <th>事件</th>
              <th>动作</th>
              <th>操作者</th>
              <th>资源</th>
              <th>发生时间</th>
              <th>结果</th>
              <th>变更前</th>
              <th>变更后</th>
            </tr>
          </thead>
          <tbody>
            {events.map((event) => (
              <tr key={event.eventUid}>
                <td>{event.eventId}</td>
                <td>{event.action}</td>
                <td>{event.actorUserId}</td>
                <td>{labelFor(objectTypeLabel, event.resourceType)} {event.resourceId}</td>
                <td>{formatDate(event.occurredAt)}</td>
                <td>{labelFor(resultLabel, event.result)}</td>
                <td>变更前：{summaryTextZh(event.beforeSummary)}</td>
                <td>变更后：{summaryTextZh(event.afterSummary)}</td>
              </tr>
            ))}
            {events.length === 0 && (
              <tr>
                <td colSpan={8}>暂无操作日志。</td>
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
