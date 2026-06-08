import { BarChart3, BriefcaseBusiness, Building2, ClipboardList, FileText, Landmark, ListChecks, LockKeyhole, ReceiptText, UsersRound } from 'lucide-react';
import type { ComponentType } from 'react';
import { UserRole } from '../api/auth';
import { navLabels } from '../i18n/labels';

type NavItem = {
  label: string;
  roles: UserRole[];
  icon: ComponentType<{ size?: number }>;
  view: AppView;
};

export type AppView = 'overview' | 'leads' | 'accounts' | 'contacts' | 'opportunities' | 'quotes' | 'contracts' | 'payments' | 'tasks' | 'reminders' | 'managerOverview' | 'importExport' | 'userManagement' | 'operationLogs';

const items: NavItem[] = [
  { label: navLabels.overview, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: BriefcaseBusiness, view: 'overview' },
  { label: navLabels.leads, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ClipboardList, view: 'leads' },
  { label: navLabels.accounts, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: Building2, view: 'accounts' },
  { label: navLabels.contacts, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: UsersRound, view: 'contacts' },
  { label: navLabels.opportunities, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: Landmark, view: 'opportunities' },
  { label: navLabels.quotes, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: FileText, view: 'quotes' },
  { label: navLabels.contracts, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ReceiptText, view: 'contracts' },
  { label: navLabels.payments, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ListChecks, view: 'payments' },
  { label: navLabels.tasks, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ListChecks, view: 'tasks' },
  { label: navLabels.reminders, roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ListChecks, view: 'reminders' },
  { label: navLabels.managerOverview, roles: ['Administrator', 'Sales Manager'], icon: BarChart3, view: 'managerOverview' },
  { label: navLabels.importExport, roles: ['Administrator', 'Sales Manager'], icon: FileText, view: 'importExport' },
  { label: navLabels.userManagement, roles: ['Administrator'], icon: LockKeyhole, view: 'userManagement' },
  { label: navLabels.operationLogs, roles: ['Administrator'], icon: ListChecks, view: 'operationLogs' }
];

export function Nav({ role, activeView, onSelect }: { role: UserRole; activeView: AppView; onSelect: (view: AppView) => void }) {
  return (
    <nav className="nav" aria-label="主导航">
      {items
        .filter((item) => item.roles.includes(role))
        .map((item) => {
          const Icon = item.icon;
          return (
            <button
              aria-current={activeView === item.view ? 'page' : undefined}
              aria-label={item.label}
              className={`navItem ${activeView === item.view ? 'active' : ''}`}
              type="button"
              key={item.label}
              onClick={() => onSelect(item.view)}
            >
              <Icon size={18} />
              <span>{item.label}</span>
            </button>
          );
        })}
    </nav>
  );
}
