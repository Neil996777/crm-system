import { useEffect, useMemo, useState } from 'react';
import { CheckCircle2, FileText, ShieldCheck } from 'lucide-react';
import { ApiError } from '../../api/client';
import { HistoryEvent } from '../../api/history';
import { getOperationLog } from '../../api/oplog';
import { AuditEventCard, Badge, Button, PageHeader, Pagination, Panel, PanelHeader, Toolbar } from '../../components/ui';
import { actionLabel, labelFor, localizeError, localizeMessage, objectTypeLabel, resultLabel, roleLabel } from '../../i18n/labels';

const pageSizeOptions = [5, 10, 20];

export function OperationLogs() {
  const [events, setEvents] = useState<HistoryEvent[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [timeFilter, setTimeFilter] = useState<'all' | 'today'>('all');
  const [actionFilter, setActionFilter] = useState('all');
  const [resourceFilter, setResourceFilter] = useState('all');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);

  useEffect(() => {
    void load();
  }, []);

  async function load() {
    setLoading(true);
    setError('');
    try {
      const response = await getOperationLog();
      setEvents(response.events);
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError, '权限不足。'));
    } finally {
      setLoading(false);
    }
  }

  const actionOptions = useMemo(() => Array.from(new Set(events.map((event) => event.action))).sort(), [events]);
  const resourceOptions = useMemo(() => Array.from(new Set(events.map((event) => event.resourceType))).sort(), [events]);

  const filteredEvents = useMemo(() => events.filter((event) => {
    const query = search.trim().toLowerCase();
    const matchesSearch = !query || [
      event.actorDisplay,
      event.actorUserId,
      labelFor(actionLabel, event.action),
      labelFor(objectTypeLabel, event.resourceType),
      event.resourceId,
      safeSummaryText(event)
    ].join(' ').toLowerCase().includes(query);
    const matchesTime = timeFilter === 'all' || event.occurredAt.slice(0, 10) === today();
    const matchesAction = actionFilter === 'all' || event.action === actionFilter;
    const matchesResource = resourceFilter === 'all' || event.resourceType === resourceFilter;
    return matchesSearch && matchesTime && matchesAction && matchesResource;
  }), [events, search, timeFilter, actionFilter, resourceFilter]);

  const totalPages = Math.max(1, Math.ceil(filteredEvents.length / pageSize));
  const safePage = Math.min(page, totalPages);
  const pagedEvents = filteredEvents.slice((safePage - 1) * pageSize, safePage * pageSize);

  return (
    <main className="content operationLogPage" data-uiux="operation-log">
      <PageHeader
        title="操作日志"
        description="仅管理员可访问 · 只读审计 · 使用安全摘要"
        actions={(
          <>
            <Button>今天</Button>
            <Button variant="primary">导出</Button>
          </>
        )}
      />
      {error && <div role="alert" className="errorBanner">{error}</div>}
      {loading && <p className="emptyState">正在加载操作日志...</p>}
      {!loading && !error && (
        <>
          <Toolbar
            searchValue={search}
            onSearchChange={(value) => { setSearch(value); setPage(1); }}
            searchPlaceholder="搜索操作人、对象或摘要"
            filters={(
              <>
                <label className="compactFilter">
                  <span className="srOnly">时间筛选</span>
                  <select value={timeFilter} onChange={(event) => { setTimeFilter(event.target.value as 'all' | 'today'); setPage(1); }}>
                    <option value="all">时间：全部</option>
                    <option value="today">时间：今天</option>
                  </select>
                </label>
                <label className="compactFilter">
                  <span className="srOnly">动作筛选</span>
                  <select value={actionFilter} onChange={(event) => { setActionFilter(event.target.value); setPage(1); }}>
                    <option value="all">操作人：全部</option>
                    {actionOptions.map((action) => <option key={action} value={action}>{labelFor(actionLabel, action)}</option>)}
                  </select>
                </label>
                <label className="compactFilter">
                  <span className="srOnly">类型筛选</span>
                  <select value={resourceFilter} onChange={(event) => { setResourceFilter(event.target.value); setPage(1); }}>
                    <option value="all">类型：全部</option>
                    {resourceOptions.map((type) => <option key={type} value={type}>{labelFor(objectTypeLabel, type)}</option>)}
                  </select>
                </label>
              </>
            )}
            summary={`只读列表 · ${filteredEvents.length} / ${events.length} 条`}
            onClearFilters={() => { setSearch(''); setTimeFilter('all'); setActionFilter('all'); setResourceFilter('all'); setPage(1); }}
          />
          <section className="operationLogLayout">
            <Panel className="operationLogList" aria-label="只读审计列表">
              <PanelHeader
                title="只读审计列表"
                description="操作人 / 动作 / 对象 / 结果 / 安全摘要 / 时间 / 事件哈希"
                meta={`第 ${safePage} 页`}
                actions={<span className="panelIcon"><FileText size={18} aria-hidden="true" /></span>}
              />
              {pagedEvents.length === 0 ? (
                <p className="emptyState">暂无操作日志。</p>
              ) : (
                <div className="auditCardList">
                  {pagedEvents.map((event) => (
                    <AuditEventCard
                      key={event.eventUid}
                      actor={`${event.actorDisplay || event.actorUserId} · ${labelFor(roleLabel, event.actorRole)}`}
                      action={labelFor(actionLabel, event.action)}
                      resource={`${labelFor(objectTypeLabel, event.resourceType)} ${event.resourceId}`}
                      result={labelFor(resultLabel, event.result)}
                      occurredAt={formatDate(event.occurredAt)}
                      safeSummary={`安全摘要：${safeSummaryText(event)}`}
                      eventId={event.eventId}
                      hash={event.eventHash ? event.eventHash.slice(0, 28) : '未返回'}
                      badges={<Badge tone={event.result.toLowerCase() === 'success' ? 'success' : 'danger'}>{labelFor(resultLabel, event.result)}</Badge>}
                    />
                  ))}
                </div>
              )}
              <Pagination
                page={safePage}
                totalPages={totalPages}
                totalItems={filteredEvents.length}
                pageSize={pageSize}
                pageSizeOptions={pageSizeOptions}
                onPageChange={setPage}
                onPageSizeChange={(nextSize) => { setPageSize(nextSize); setPage(1); }}
              />
            </Panel>
            <aside className="rightRail" aria-label="操作日志门控说明">
              <Panel>
                <PanelHeader title="权限门控" description="整页仅管理员" meta="只读" />
                <div className="railCard">
                  <span className="panelIcon"><ShieldCheck size={18} aria-hidden="true" /></span>
                  <div>
                    <strong>非管理员不进入</strong>
                    <p>操作日志读取权限仅管理员可见。</p>
                  </div>
                </div>
                <div className="railCard">
                  <span className="panelIcon mint"><CheckCircle2 size={18} aria-hidden="true" /></span>
                  <div>
                    <strong>不可编辑</strong>
                    <p>页面无编辑、删除、重写日志入口。</p>
                  </div>
                </div>
              </Panel>
              <Panel>
                <PanelHeader title="防篡改提示" description="事件哈希" />
                <div className="railCard">
                  <span className="panelIcon"><FileText size={18} aria-hidden="true" /></span>
                  <div>
                    <strong>事件哈希弱化展示</strong>
                    <p>每条记录保留事件哈希，用于审计链校验。</p>
                  </div>
                </div>
              </Panel>
            </aside>
          </section>
        </>
      )}
    </main>
  );
}

function safeSummaryText(event: HistoryEvent) {
  const text = event.safeSummary || event.action;
  const localized = localizeMessage(text, text);
  if (localized !== text) return localized;

  const importMatch = text.match(/^CSV import completed for ([a-zA-Z]+) with (\d+) successful and (\d+) failed rows\.$/);
  if (importMatch) {
    return `CSV 导入已完成：${labelFor(objectTypeLabel, importMatch[1])}，成功 ${importMatch[2]} 行，失败 ${importMatch[3]} 行。`;
  }
  const exportMatch = text.match(/^CSV export completed for ([a-zA-Z]+) with (\d+) rows\.$/);
  if (exportMatch) {
    return `CSV 导出已完成：${labelFor(objectTypeLabel, exportMatch[1])} ${exportMatch[2]} 行。`;
  }
  const actorActionMatch = text.match(/^(Administrator|Sales Manager|Sales) ([a-z_]+) on ([a-zA-Z]+)$/);
  if (actorActionMatch) {
    return `${labelFor(roleLabel, actorActionMatch[1])}执行${labelFor(actionLabel, actorActionMatch[2])}，对象为${labelFor(objectTypeLabel, actorActionMatch[3])}。`;
  }

  const action = labelFor(actionLabel, text);
  return action || '安全摘要已记录。';
}

function today() {
  return new Date().toISOString().slice(0, 10);
}

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString('zh-CN', { hour12: false });
}
