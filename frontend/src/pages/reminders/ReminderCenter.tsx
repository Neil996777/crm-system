import { FormEvent, useEffect, useMemo, useState } from 'react';
import { Bell, CalendarDays, Clock, CreditCard, FileSignature, RefreshCcw, UserRound } from 'lucide-react';
import { ApiError } from '../../api/client';
import { ReminderRow, listReminders } from '../../api/reminders';
import type { RecordNavigationTarget } from '../../app/navigation';
import { targetForRelatedRecord } from '../../app/navigation';
import { Badge, Button, MetricCard, PageHeader, Panel, PanelHeader, ReminderRowCard } from '../../components/ui';
import { contractStatusLabel, labelFor, localizeError, objectTypeLabel, paymentStatusLabel, priorityLabel, reminderTypeLabel, taskStatusLabel } from '../../i18n/labels';

type ReminderFilter = 'all' | 'due' | 'overdue' | 'contract' | 'payment';

export function ReminderCenter({ onNavigate }: { onNavigate?: (target: RecordNavigationTarget) => void }) {
  const [businessDate, setBusinessDate] = useState(today());
  const [rows, setRows] = useState<ReminderRow[]>([]);
  const [timezone, setTimezone] = useState('Asia/Shanghai');
  const [filter, setFilter] = useState<ReminderFilter>('all');
  const [sortDueFirst, setSortDueFirst] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  const stats = useMemo(() => ({
    taskDue: rows.filter((row) => row.type === 'task_due').length,
    taskOverdue: rows.filter((row) => row.type === 'task_overdue').length,
    contractPending: rows.filter((row) => row.type === 'contract_pending_signature').length,
    paymentDue: rows.filter((row) => row.type === 'payment_due').length,
    paymentOverdue: rows.filter((row) => row.type === 'payment_overdue').length
  }), [rows]);

  const visibleRows = useMemo(() => rows.filter((row) => {
    if (filter === 'due') return row.type === 'task_due' || row.type === 'payment_due';
    if (filter === 'overdue') return row.type === 'task_overdue' || row.type === 'payment_overdue';
    if (filter === 'contract') return row.type === 'contract_pending_signature';
    if (filter === 'payment') return row.type === 'payment_due' || row.type === 'payment_overdue';
    return true;
  }).sort((left, right) => sortDueFirst ? left.dueDate.localeCompare(right.dueDate) : 0), [rows, filter, sortDueFirst]);

  async function refresh() {
    setError('');
    try {
      const response = await listReminders(businessDate);
      setRows(response.rows);
      setTimezone(response.timezone);
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  function submit(event: FormEvent) {
    event.preventDefault();
    void refresh();
  }

  function openReminder(row: ReminderRow) {
    const target = targetForRelatedRecord(row.relatedRecord.type, row.relatedRecord.id);
    if (!target) {
      setError('当前提醒没有可跳转的关联记录。');
      return;
    }
    onNavigate?.(target);
  }

  return (
    <main className="content remindersPage" data-uiux="reminders-center">
      <PageHeader
        title="提醒中心"
        description={`本人提醒 · 业务日期 ${businessDate} · 时区 ${timezone}`}
        actions={(
          <form className="reminderDateForm" onSubmit={submit}>
            <label>
              <span className="srOnly">业务日期</span>
              <input type="date" value={businessDate} onChange={(event) => setBusinessDate(event.target.value)} />
            </label>
            <Button type="button" aria-pressed={sortDueFirst} onClick={() => setSortDueFirst((value) => !value)}>按到期排序</Button>
            <Button variant="primary" type="submit" disabled={false}>
              <RefreshCcw size={16} aria-hidden="true" />
              刷新提醒
            </Button>
          </form>
        )}
      />
      {error && <div role="alert" className="errorBanner">{error}</div>}

      <section className="reminderStatGrid" aria-label="提醒统计">
        <MetricCard label="任务到期" value={stats.taskDue} icon={<CalendarDays size={18} aria-hidden="true" />} />
        <MetricCard label="任务逾期" value={stats.taskOverdue} tone="peach" icon={<Clock size={18} aria-hidden="true" />} />
        <MetricCard label="合同待签署" value={stats.contractPending} tone="purple" icon={<FileSignature size={18} aria-hidden="true" />} />
        <MetricCard label="回款到期" value={stats.paymentDue} tone="mint" icon={<CreditCard size={18} aria-hidden="true" />} />
        <MetricCard label="回款逾期" value={stats.paymentOverdue} tone="peach" icon={<Bell size={18} aria-hidden="true" />} />
      </section>

      <section className="reminderLayout">
        <Panel className="reminderMainPanel" aria-label="待处理提醒">
          <PanelHeader
            title="待处理提醒"
            description="字段：类型 / 关联记录 / 负责人 / 到期日 / 状态 / 优先级 / 版本"
            meta="本人范围"
            actions={<span className="panelIcon peach"><Bell size={18} aria-hidden="true" /></span>}
          />
          <div className="filterChipRow" role="tablist" aria-label="提醒筛选">
            {[
              ['all', '全部'],
              ['due', '到期'],
              ['overdue', '逾期'],
              ['contract', '合同'],
              ['payment', '回款']
            ].map(([key, label]) => (
              <button className={filter === key ? 'chip selected' : 'chip'} key={key} type="button" onClick={() => setFilter(key as ReminderFilter)}>
                {label}
              </button>
            ))}
          </div>
          {visibleRows.length === 0 ? (
            <p className="emptyState">今日无待处理提醒。</p>
          ) : (
            <div className="reminderCardList">
              {visibleRows.map((row) => (
                <ReminderRowCard
                  key={`${row.type}-${row.id}`}
                  title={reminderTitle(row)}
                  description={`关联记录：${labelFor(objectTypeLabel, row.relatedRecord.type)} ${row.relatedRecord.id}`}
                  icon={reminderIcon(row)}
                  meta={`负责人 ${row.ownerDisplay || '未分配'} · 到期 ${row.dueDate} · 版本 v${row.version}`}
                  badges={(
                    <>
                      <Badge tone={reminderTone(row)}>{labelFor(reminderTypeLabel, row.type)}</Badge>
                      <Badge tone={statusTone(row.status)}>{statusLabel(row)}</Badge>
                      <Badge tone={row.priority === 'P0' || row.priority === 'P1' ? 'warning' : 'neutral'}>{labelFor(priorityLabel, row.priority)}</Badge>
                    </>
                  )}
                  tone={reminderAccent(row)}
                  overdue={row.type.includes('overdue')}
                  actions={<Button variant="ghost" onClick={() => openReminder(row)}>查看</Button>}
                />
              ))}
            </div>
          )}
        </Panel>

        <aside className="rightRail" aria-label="提醒数据范围">
          <Panel>
            <PanelHeader title="数据范围" description="销售视图" meta="本人范围" />
            <div className="railCard">
              <span className="panelIcon"><UserRound size={18} aria-hidden="true" /></span>
              <div>
                <strong>仅本人提醒</strong>
                <p>ownerDisplay 按授权返回；不显示团队聚合。</p>
              </div>
            </div>
            <div className="railCard">
              <span className="panelIcon mint"><RefreshCcw size={18} aria-hidden="true" /></span>
              <div>
                <strong>跳转关联记录</strong>
                <p>每条提醒保留“查看”入口，进入对应任务、合同或回款。</p>
              </div>
            </div>
          </Panel>
          <Panel>
            <PanelHeader title="空态" description="本帐下单独出图" />
            <p className="emptyState">今日无待处理提醒</p>
          </Panel>
        </aside>
      </section>
    </main>
  );
}

function reminderTitle(row: ReminderRow) {
  return `${labelFor(reminderTypeLabel, row.type)} · ${labelFor(objectTypeLabel, row.relatedRecord.type)} ${row.relatedRecord.id}`;
}

function statusLabel(row: ReminderRow) {
  const statusLabels = { ...taskStatusLabel, ...contractStatusLabel, ...paymentStatusLabel };
  return labelFor(statusLabels, row.status);
}

function reminderTone(row: ReminderRow): 'neutral' | 'primary' | 'success' | 'warning' | 'danger' {
  if (row.type.includes('overdue')) return 'danger';
  if (row.type.includes('payment')) return 'warning';
  if (row.type.includes('contract')) return 'primary';
  return 'primary';
}

function statusTone(status: string): 'neutral' | 'primary' | 'success' | 'warning' | 'danger' {
  if (status === 'Overdue') return 'danger';
  if (status === 'Open' || status === 'Pending Signature' || status === 'Unpaid') return 'warning';
  if (status === 'Completed' || status === 'Paid' || status === 'Signed') return 'success';
  return 'neutral';
}

function reminderAccent(row: ReminderRow): 'primary' | 'sky' | 'mint' | 'peach' | 'purple' | 'success' | 'warning' | 'danger' {
  if (row.type.includes('overdue')) return 'peach';
  if (row.type.includes('payment')) return 'mint';
  if (row.type.includes('contract')) return 'purple';
  return 'primary';
}

function reminderIcon(row: ReminderRow) {
  if (row.type.includes('payment')) return <CreditCard size={18} aria-hidden="true" />;
  if (row.type.includes('contract')) return <FileSignature size={18} aria-hidden="true" />;
  if (row.type.includes('overdue')) return <Clock size={18} aria-hidden="true" />;
  return <CalendarDays size={18} aria-hidden="true" />;
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
