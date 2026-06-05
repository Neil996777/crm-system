import { useState } from 'react';
import { useSession } from '../auth/SessionProvider';
import { SignIn } from '../pages/SignIn';
import { WorkOverview } from '../pages/WorkOverview';
import { AccountList } from '../pages/accounts/AccountList';
import { ContactList } from '../pages/accounts/ContactList';
import { LeadList } from '../pages/leads/LeadList';
import { OpportunityList } from '../pages/opportunities/OpportunityList';
import { QuoteList } from '../pages/quotes/QuoteList';
import { ContractList } from '../pages/contracts/ContractList';
import { PaymentList } from '../pages/payments/PaymentList';
import { TaskList } from '../components/TaskList';
import { ReminderCenter } from '../pages/reminders/ReminderCenter';
import { ManagerOverview } from '../pages/reports/ManagerOverview';
import { ImportExportPage } from '../pages/importexport/Import';
import { OperationLogs } from '../pages/admin/OperationLogs';
import { UserManagement } from '../pages/admin/UserManagement';
import { AppView, Nav } from './Nav';
import { appName, labelFor, roleLabel } from '../i18n/labels';

export function Shell() {
  const { user, loading, logout } = useSession();
  const [view, setView] = useState<AppView>('overview');

  if (loading && !user) {
    return <div className="boot">加载中</div>;
  }

  if (!user) {
    return <SignIn />;
  }

  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="brand">{appName}</div>
        <Nav role={user.role} activeView={view} onSelect={setView} />
      </aside>
      <div className="workspace">
        <header className="topbar">
          <div>
            <strong>{user.displayName}</strong>
            <span>{labelFor(roleLabel, user.role)}</span>
          </div>
          <button className="secondaryButton" type="button" onClick={logout}>
            退出登录
          </button>
        </header>
        {view === 'leads' ? <LeadList /> : view === 'accounts' ? <AccountList /> : view === 'contacts' ? <ContactList /> : view === 'opportunities' ? <OpportunityList /> : view === 'quotes' ? <QuoteList /> : view === 'contracts' ? <ContractList /> : view === 'payments' ? <PaymentList /> : view === 'tasks' ? <TaskList /> : view === 'reminders' ? <ReminderCenter /> : view === 'managerOverview' ? <ManagerOverview /> : view === 'importExport' ? <ImportExportPage /> : view === 'userManagement' ? <UserManagement /> : view === 'operationLogs' ? <OperationLogs /> : <WorkOverview user={user} />}
      </div>
    </div>
  );
}
