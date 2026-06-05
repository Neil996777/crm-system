import { CurrentUser } from '../api/auth';

export function WorkOverview({ user }: { user: CurrentUser }) {
  const scopeLabel = user.role === 'Administrator' ? '全局 CRM 范围' : user.role === 'Sales Manager' ? '团队范围' : '本人负责和分配范围';

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>工作台</h1>
          <p>{scopeLabel}</p>
        </div>
      </section>
      <section className="metrics" aria-label="汇总">
        <div className="metric">
          <span>进行中的工作</span>
          <strong>0</strong>
        </div>
        <div className="metric">
          <span>到期提醒</span>
          <strong>0</strong>
        </div>
        <div className="metric">
          <span>待跟进事项</span>
          <strong>0</strong>
        </div>
      </section>
      <section className="workGrid">
        <div>
          <h2>已分配工作</h2>
          <p className="emptyState">暂无已分配的进行中工作。</p>
        </div>
        <div>
          <h2>提醒</h2>
          <p className="emptyState">暂无到期提醒。</p>
        </div>
      </section>
    </main>
  );
}
