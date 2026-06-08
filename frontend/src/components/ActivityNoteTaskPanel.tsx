import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../api/client';
import { Activity, Note, WorkTask, createActivity, createNote, createTask, listActivities, listNotes, listTasks } from '../api/work';
import { labelFor, localizeError, taskStatusLabel } from '../i18n/labels';

export function ActivityNoteTaskPanel({
  relatedType,
  relatedId,
  ownerId,
  readOnly = false,
  onError
}: {
  relatedType: string;
  relatedId: string;
  ownerId: string;
  readOnly?: boolean;
  onError: (message: string) => void;
}) {
  const [activities, setActivities] = useState<Activity[]>([]);
  const [notes, setNotes] = useState<Note[]>([]);
  const [tasks, setTasks] = useState<WorkTask[]>([]);
  const [activityForm, setActivityForm] = useState({ activityType: '', content: '' });
  const [noteContent, setNoteContent] = useState('');
  const [taskForm, setTaskForm] = useState({ title: '', dueDate: '' });
  const entries = timelineEntries(activities, notes, tasks);

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
      onError(localizeError(apiError));
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
      onError(localizeError(apiError));
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
      onError(localizeError(apiError));
    }
  }

  return (
    <section className="relatedSection" aria-label="动态、备注、任务">
      <h2>动态、备注、任务</h2>
      {readOnly ? (
        <p className="inlineNotice">终态记录只读，不能新增备注、动态或任务。</p>
      ) : (
        <>
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
              <input className="dateControl" type="date" value={taskForm.dueDate} onChange={(event) => setTaskForm({ ...taskForm, dueDate: event.target.value })} />
            </label>
            <button className="primaryButton" type="submit">保存任务</button>
          </form>
        </>
      )}
      <div className="recordList activityTimeline" aria-label="活动时间线">
        {entries.length === 0 ? <p className="inlineNotice">暂无工作记录。</p> : entries.map((entry) => (
          <article className="timelineItem" key={entry.id}>
            <div className="timelineHeader">
              <span>{entry.type}</span>
              <time dateTime={entry.dateTime}>{entry.timeLabel}</time>
            </div>
            <p>{entry.title}</p>
            <small>{entry.meta}</small>
          </article>
        ))}
      </div>
    </section>
  );
}

function timelineEntries(activities: Activity[], notes: Note[], tasks: WorkTask[]) {
  const entries = [
    ...notes.map((note) => ({
      id: `note-${note.id}`,
      type: '备注',
      title: note.content,
      meta: `负责人 ${note.ownerId}`,
      dateTime: note.occurredAt,
      timeLabel: formatTimelineTime(note.occurredAt)
    })),
    ...activities.map((activity) => ({
      id: `activity-${activity.id}`,
      type: activity.activityType || '动态',
      title: activity.content,
      meta: `负责人 ${activity.ownerId}`,
      dateTime: activity.occurredAt,
      timeLabel: formatTimelineTime(activity.occurredAt)
    })),
    ...tasks.map((task) => ({
      id: `task-${task.id}`,
      type: '任务',
      title: task.title,
      meta: `${labelFor(taskStatusLabel, task.status)} · 负责人 ${task.ownerId}`,
      dateTime: task.dueDate,
      timeLabel: task.dueDate ? `到期 ${task.dueDate}` : '未设置到期日'
    }))
  ];
  return entries.sort((left, right) => timestamp(right.dateTime) - timestamp(left.dateTime));
}

function timestamp(value: string) {
  const time = Date.parse(value);
  return Number.isFinite(time) ? time : 0;
}

function formatTimelineTime(value: string) {
  if (!value) return '未设置时间';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString('zh-CN', { hour12: false });
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
