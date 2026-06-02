import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../api/client';
import { Activity, Note, WorkTask, createActivity, createNote, createTask, listActivities, listNotes, listTasks } from '../api/work';

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
      onError(apiError.safeMessage || 'Request failed.');
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
      onError(apiError.safeMessage || 'Request failed.');
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
      onError(apiError.safeMessage || 'Request failed.');
    }
  }

  return (
    <section className="detailPane" aria-label="Activities, Notes, Tasks">
      <h2>Activities, Notes, Tasks</h2>
      <form className="createPanel" onSubmit={submitNote}>
        <label>
          Note content
          <input value={noteContent} onChange={(event) => setNoteContent(event.target.value)} />
        </label>
        <button className="primaryButton" type="submit">Save note</button>
      </form>
      <form className="createPanel" onSubmit={submitActivity}>
        <label>
          Activity type
          <input value={activityForm.activityType} onChange={(event) => setActivityForm({ ...activityForm, activityType: event.target.value })} />
        </label>
        <label>
          Activity content
          <input value={activityForm.content} onChange={(event) => setActivityForm({ ...activityForm, content: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">Save activity</button>
      </form>
      <form className="createPanel" onSubmit={submitTask}>
        <label>
          Task title
          <input value={taskForm.title} onChange={(event) => setTaskForm({ ...taskForm, title: event.target.value })} />
        </label>
        <label>
          Task due date
          <input type="date" value={taskForm.dueDate} onChange={(event) => setTaskForm({ ...taskForm, dueDate: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">Save task</button>
      </form>
      <div className="recordList" aria-label="Work records">
        {notes.map((note) => <p className="inlineNotice" key={note.id}>{note.content}</p>)}
        {activities.map((activity) => <p className="inlineNotice" key={activity.id}>{activity.activityType}: {activity.content}</p>)}
        {tasks.map((task) => <p className="inlineNotice" key={task.id}>{task.title} · {task.status}</p>)}
      </div>
    </section>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
