import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../api/client';
import { Activity, Note, WorkTask, createActivity, createNote, createTask, listActivities, listNotes, listTasks } from '../api/work';
import { labelFor, taskStatusLabel } from '../i18n/labels';

export function ActivityNoteTaskPanel({ relatedType, relatedId, ownerId, onError }: { relatedType: string; relatedId: string; ownerId: string; onError: (message: string) => void }) {
  const [activities, setActivities] = useState<Activity[]>([]);
  const [notes, setNotes] = useState<Note[]>([]);
  const [tasks, setTasks] = useState<WorkTask[]>([]);
  const [activityForm, setActivityForm] = useState({ activityType: '', content: '' });
  const [noteContent, setNoteContent] = useState('');
  const [taskForm, setTaskForm] = useState({ title: '', dueDate: '' });

  useEffect(() => {
    void refresh();
  }, [relatedType, relatedId]);

  async function refresh() {
    const [activityResponse, noteResponse, taskResponse] = await Promise.all([
      listActivities(relatedType, relatedId),
      listNotes(relatedType, relatedId),
      listTasks({ relatedType, relatedId, businessDate: today() })
    ]);
    setActivities(activityResponse.items);
    setNotes(noteResponse.items);
    setTasks(taskResponse.items);
  }

  async function submitActivity(event: FormEvent) {
    event.preventDefault();
    onError('');
    try {
      const created = await createActivity({ relatedType, relatedId, ownerId, ...activityForm });
      setActivities([created, ...activities]);
      setActivityForm({ activityType: '', content: '' });
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(apiError.safeMessage || '请求失败。');
    }
  }

  async function submitNote(event: FormEvent) {
    event.preventDefault();
    onError('');
    try {
      const created = await createNote({ relatedType, relatedId, ownerId, content: noteContent });
      setNotes([created, ...notes]);
      setNoteContent('');
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(apiError.safeMessage || '请求失败。');
    }
  }

  async function submitTask(event: FormEvent) {
    event.preventDefault();
    onError('');
    try {
      const created = await createTask({ relatedType, relatedId, ownerId, title: taskForm.title, dueDate: taskForm.dueDate });
      setTasks([created, ...tasks]);
      setTaskForm({ title: '', dueDate: '' });
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(apiError.safeMessage || '请求失败。');
    }
  }

  return (
    <section className="detailPane" aria-label="动态、备注、任务">
      <h2>动态、备注、任务</h2>
      <form className="createPanel" onSubmit={submitNote}>
        <label>
          备注内容
          <input value={noteContent} onChange={(event) => setNoteContent(event.target.value)} />
        </label>
        <button className="primaryButton" type="submit">保存备注</button>
      </form>
      <form className="createPanel" onSubmit={submitActivity}>
        <label>
          动态类型
          <input value={activityForm.activityType} onChange={(event) => setActivityForm({ ...activityForm, activityType: event.target.value })} />
        </label>
        <label>
          动态内容
          <input value={activityForm.content} onChange={(event) => setActivityForm({ ...activityForm, content: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">保存动态</button>
      </form>
      <form className="createPanel" onSubmit={submitTask}>
        <label>
          任务标题
          <input value={taskForm.title} onChange={(event) => setTaskForm({ ...taskForm, title: event.target.value })} />
        </label>
        <label>
          任务到期日
          <input type="date" value={taskForm.dueDate} onChange={(event) => setTaskForm({ ...taskForm, dueDate: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">保存任务</button>
      </form>
      <div className="recordList" aria-label="工作记录">
        {notes.map((note) => <p className="inlineNotice" key={note.id}>{note.content}</p>)}
        {activities.map((activity) => <p className="inlineNotice" key={activity.id}>{activity.activityType}: {activity.content}</p>)}
        {tasks.map((task) => <p className="inlineNotice" key={task.id}>{task.title} · {labelFor(taskStatusLabel, task.status)}</p>)}
      </div>
    </section>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
