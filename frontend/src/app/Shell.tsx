import { useCallback, useEffect, useState } from 'react';
import { LogOut } from 'lucide-react';
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
import type { RecordNavigationTarget } from './navigation';
import { appName, labelFor, roleLabel } from '../i18n/labels';

export function Shell() {
  const { user, loading, logout } = useSession();
  const [view, setView] = useState<AppView>('overview');
  const [focusMode, setFocusMode] = useState(false);
  const [recordTarget, setRecordTarget] = useState<RecordNavigationTarget | null>(null);

  useEffect(() => {
    if (view !== 'overview') {
      setFocusMode(false);
    }
  }, [view]);

  const selectView = (nextView: AppView) => {
    setView(nextView);
    setRecordTarget(null);
    if (nextView !== 'overview') {
      setFocusMode(false);
    }
  };

  const navigateToRecord = useCallback((target: RecordNavigationTarget) => {
    setView(target.view);
    setRecordTarget(target.recordId ? target : null);
    if (target.view !== 'overview') {
      setFocusMode(false);
    }
  }, []);

  if (loading && !user) {
    return <div className="boot">加载中</div>;
  }

  if (!user) {
    return <SignIn />;
  }

  return (
    <div className={`shell ${focusMode && view === 'overview' ? 'focusMode' : ''}`}>
      <aside className="sidebar">
        <div className="brand">
          <span className="brandMark" aria-hidden="true">
            CRM
          </span>
          <span className="brandText">
            <strong>{appName}</strong>
            <span>{user.role === 'Sales' ? '个人销售工作区' : user.role === 'Sales Manager' ? '团队管理工作区' : '全局管理工作区'}</span>
          </span>
        </div>
        <Nav role={user.role} activeView={view} onSelect={selectView} />
      </aside>
      <div className={`workspace ${focusMode && view === 'overview' ? 'focusWorkspace' : ''}`}>
        <header className="topbar">
          <div className="topbarIdentity">
            <span className="avatar" aria-hidden="true">
              {user.displayName.slice(0, 1)}
            </span>
            <div>
              <strong>{user.displayName}</strong>
              <span>{labelFor(roleLabel, user.role)}</span>
            </div>
          </div>
          <div className="topbarSpacer" />
          <button className="secondaryButton" type="button" onClick={logout}>
            <LogOut size={16} aria-hidden="true" />
            退出登录
          </button>
        </header>
        {view === 'leads' ? (
          <LeadList targetRecordId={recordTarget?.view === 'leads' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'accounts' ? (
          <AccountList targetRecordId={recordTarget?.view === 'accounts' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'contacts' ? (
          <ContactList targetRecordId={recordTarget?.view === 'contacts' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'opportunities' ? (
          <OpportunityList targetRecordId={recordTarget?.view === 'opportunities' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'quotes' ? (
          <QuoteList targetRecordId={recordTarget?.view === 'quotes' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'contracts' ? (
          <ContractList targetRecordId={recordTarget?.view === 'contracts' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'payments' ? (
          <PaymentList targetRecordId={recordTarget?.view === 'payments' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'tasks' ? (
          <TaskList targetRecordId={recordTarget?.view === 'tasks' ? recordTarget.recordId : undefined} onTargetHandled={() => setRecordTarget(null)} />
        ) : view === 'reminders' ? (
          <ReminderCenter onNavigate={navigateToRecord} />
        ) : view === 'managerOverview' ? <ManagerOverview /> : view === 'importExport' ? <ImportExportPage /> : view === 'userManagement' ? <UserManagement /> : view === 'operationLogs' ? <OperationLogs /> : <WorkOverview user={user} onFocusChange={setFocusMode} onNavigate={navigateToRecord} />}
      </div>
    </div>
  );
}
