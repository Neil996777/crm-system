import { useEffect, useState } from 'react';
import { ApiError } from '../api/client';
import { WorkTask, changeTaskStatus, listTasks } from '../api/work';

export function TaskList() {
  const [tasks, setTasks] = useState<WorkTask[]>([]);
  const [selected, setSelected] = useState<WorkTask | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh() {
    try {
      const response = await listTasks({ businessDate: today() });
      setTasks(response.items);
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
    }
  }

  async function complete() {
    if (!selected) return;
    setError('');
    try {
      const updated = await changeTaskStatus(selected.id, 'Completed');
      setSelected(updated);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
    }
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Tasks</h1>
          <p>Open and completed follow-up work.</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <div className="recordList" aria-label="Task records">
            {tasks.length === 0 ? <p className="emptyState">No tasks found.</p> : tasks.map((task) => (
              <button className="recordRow" type="button" key={task.id} onClick={() => setSelected(task)}>
                <strong>{task.title}</strong>
                <span>{task.status} · {task.dueDate}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? (
            <section className="detailPane" aria-label="Task detail">
              <h2>{selected.title}</h2>
              <p>{selected.status}</p>
              <dl className="detailGrid">
                <div>
                  <dt>Related record</dt>
                  <dd>{selected.relatedType} {selected.relatedId}</dd>
                </div>
                <div>
                  <dt>Due date</dt>
                  <dd>{selected.dueDate}</dd>
                </div>
              </dl>
              <button className="primaryButton" type="button" disabled={selected.status === 'Completed' || selected.status === 'Cancelled'} onClick={() => void complete()}>Complete task</button>
            </section>
          ) : <p className="emptyState">Select a task to update status.</p>}
        </div>
      </section>
    </main>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
