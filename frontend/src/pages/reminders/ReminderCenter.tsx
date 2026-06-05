import { FormEvent, useEffect, useMemo, useState } from 'react';
import { ApiError } from '../../api/client';
import { ReminderRow, listReminders } from '../../api/reminders';
import { labelFor, reminderTypeLabel, taskStatusLabel, contractStatusLabel, paymentStatusLabel, localizeError } from '../../i18n/labels';

export function ReminderCenter() {
  const [businessDate, setBusinessDate] = useState(today());
  const [rows, setRows] = useState<ReminderRow[]>([]);
  const [timezone, setTimezone] = useState('Asia/Shanghai');
  const [error, setError] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  const grouped = useMemo(() => ({
    tasks: rows.filter((row) => row.type === 'task_due' || row.type === 'task_overdue'),
    contracts: rows.filter((row) => row.type === 'contract_pending_signature'),
    payments: rows.filter((row) => row.type === 'payment_due' || row.type === 'payment_overdue')
  }), [rows]);

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

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>提醒中心</h1>
          <p>{timezone}</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <form className="toolbar" onSubmit={submit}>
        <label>
          业务日期
          <input type="date" value={businessDate} onChange={(event) => setBusinessDate(event.target.value)} />
        </label>
        <button className="primaryButton" type="submit">刷新提醒</button>
      </form>
      <section className="leadLayout">
        <ReminderGroup title="任务提醒" rows={grouped.tasks} />
        <ReminderGroup title="合同提醒" rows={grouped.contracts} />
        <ReminderGroup title="回款提醒" rows={grouped.payments} />
      </section>
    </main>
  );
}

function ReminderGroup({ title, rows }: { title: string; rows: ReminderRow[] }) {
  return (
    <section className="detailPane" aria-label={title}>
      <h2>{title}</h2>
      {rows.length === 0 ? <p className="emptyState">暂无有效提醒。</p> : (
        <div className="recordList">
          {rows.map((row) => (
            <div className="recordRow staticRow" key={`${row.type}-${row.id}`}>
              <strong>{row.relatedRecord.display}</strong>
              <span>{labelForType(row)} · {row.dueDate}</span>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}

function labelForType(row: ReminderRow) {
  const statusLabels = { ...taskStatusLabel, ...contractStatusLabel, ...paymentStatusLabel };
  return `${labelFor(reminderTypeLabel, row.type)} · ${labelFor(statusLabels, row.status)}`;
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
