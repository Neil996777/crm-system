import { FormEvent, useEffect, useMemo, useState } from 'react';
import { ApiError } from '../../api/client';
import { ReminderRow, listReminders } from '../../api/reminders';

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
      setError(apiError.safeMessage || 'Request failed.');
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
          <h1>Reminder Center</h1>
          <p>{timezone}</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <form className="toolbar" onSubmit={submit}>
        <label>
          Business date
          <input type="date" value={businessDate} onChange={(event) => setBusinessDate(event.target.value)} />
        </label>
        <button className="primaryButton" type="submit">Refresh reminders</button>
      </form>
      <section className="leadLayout">
        <ReminderGroup title="Task Reminders" rows={grouped.tasks} />
        <ReminderGroup title="Contract Reminders" rows={grouped.contracts} />
        <ReminderGroup title="Payment Reminders" rows={grouped.payments} />
      </section>
    </main>
  );
}

function ReminderGroup({ title, rows }: { title: string; rows: ReminderRow[] }) {
  return (
    <section className="detailPane" aria-label={title}>
      <h2>{title}</h2>
      {rows.length === 0 ? <p className="emptyState">No active reminders.</p> : (
        <div className="recordList">
          {rows.map((row) => (
            <div className="recordRow staticRow" key={`${row.type}-${row.id}`}>
              <strong>{row.relatedRecord.display}</strong>
              <span>{labelForType(row.type)} · {row.status} · {row.dueDate}</span>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}

function labelForType(type: ReminderRow['type']) {
  switch (type) {
    case 'task_due':
    case 'task_overdue':
      return 'Task';
    case 'contract_pending_signature':
      return 'Pending Signature';
    case 'payment_due':
    case 'payment_overdue':
      return 'Payment';
  }
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
