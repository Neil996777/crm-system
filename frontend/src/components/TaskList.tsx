import { useEffect, useState } from 'react';
import { ApiError } from '../api/client';
import { WorkTask, changeTaskStatus, listTasks } from '../api/work';
import { labelFor, objectTypeLabel, taskStatusLabel } from '../i18n/labels';

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
      setError(apiError.safeMessage || '请求失败。');
    }
  }

  async function complete() {
    if (!selected) return;
    setError('');
    try {
      const updated = await changeTaskStatus(selected.id, 'Completed', selected.version);
      setSelected(updated);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || '请求失败。');
    }
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>任务</h1>
          <p>查看待处理和已完成的跟进工作。</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <div className="recordList" aria-label="任务记录">
            {tasks.length === 0 ? <p className="emptyState">暂无任务。</p> : tasks.map((task) => (
              <button className="recordRow" type="button" key={task.id} onClick={() => setSelected(task)}>
                <strong>{task.title}</strong>
                <span>{labelFor(taskStatusLabel, task.status)} · {task.dueDate}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? (
            <section className="detailPane" aria-label="任务详情">
              <h2>{selected.title}</h2>
              <p>{labelFor(taskStatusLabel, selected.status)}</p>
              <dl className="detailGrid">
                <div>
                  <dt>关联记录</dt>
                  <dd>{labelFor(objectTypeLabel, selected.relatedType)} {selected.relatedId}</dd>
                </div>
                <div>
                  <dt>到期日</dt>
                  <dd>{selected.dueDate}</dd>
                </div>
              </dl>
              <button className="primaryButton" type="button" disabled={selected.status === 'Completed' || selected.status === 'Cancelled'} onClick={() => void complete()}>完成任务</button>
            </section>
          ) : <p className="emptyState">选择任务以更新状态。</p>}
        </div>
      </section>
    </main>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
