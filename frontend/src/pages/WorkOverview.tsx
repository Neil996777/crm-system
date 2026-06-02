import { CurrentUser } from '../api/auth';

export function WorkOverview({ user }: { user: CurrentUser }) {
  const scopeLabel = user.role === 'Administrator' ? 'Governed CRM scope' : user.role === 'Sales Manager' ? 'Team scope' : 'Owned and assigned scope';

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Work Overview</h1>
          <p>{scopeLabel}</p>
        </div>
      </section>
      <section className="metrics" aria-label="Summary">
        <div className="metric">
          <span>Active work</span>
          <strong>0</strong>
        </div>
        <div className="metric">
          <span>Due reminders</span>
          <strong>0</strong>
        </div>
        <div className="metric">
          <span>Open follow-ups</span>
          <strong>0</strong>
        </div>
      </section>
      <section className="workGrid">
        <div>
          <h2>Assigned Work</h2>
          <p className="emptyState">No assigned active work.</p>
        </div>
        <div>
          <h2>Reminders</h2>
          <p className="emptyState">No due reminders.</p>
        </div>
      </section>
    </main>
  );
}
