import { BarChart3, BriefcaseBusiness, Building2, ClipboardList, FileText, Landmark, ListChecks, LockKeyhole, ReceiptText, UsersRound } from 'lucide-react';
import type { ComponentType } from 'react';
import { UserRole } from '../api/auth';

type NavItem = {
  label: string;
  roles: UserRole[];
  icon: ComponentType<{ size?: number }>;
  view: AppView;
};

export type AppView = 'overview' | 'leads' | 'accounts' | 'contacts' | 'opportunities' | 'quotes' | 'contracts' | 'payments' | 'tasks' | 'reminders' | 'userManagement' | 'operationLogs';

const items: NavItem[] = [
  { label: 'Work Overview', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: BriefcaseBusiness, view: 'overview' },
  { label: 'Leads', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ClipboardList, view: 'leads' },
  { label: 'Companies/Customers', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: Building2, view: 'accounts' },
  { label: 'Contacts', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: UsersRound, view: 'contacts' },
  { label: 'Opportunities', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: Landmark, view: 'opportunities' },
  { label: 'Quotes', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: FileText, view: 'quotes' },
  { label: 'Contracts', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ReceiptText, view: 'contracts' },
  { label: 'Payments', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ListChecks, view: 'payments' },
  { label: 'Tasks', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ListChecks, view: 'tasks' },
  { label: 'Reminder Center', roles: ['Administrator', 'Sales Manager', 'Sales'], icon: ListChecks, view: 'reminders' },
  { label: 'Reports', roles: ['Administrator', 'Sales Manager'], icon: BarChart3, view: 'overview' },
  { label: 'Import/Export', roles: ['Administrator', 'Sales Manager'], icon: FileText, view: 'overview' },
  { label: 'Admin: Users/Roles', roles: ['Administrator'], icon: LockKeyhole, view: 'userManagement' },
  { label: 'Operation Logs', roles: ['Administrator'], icon: ListChecks, view: 'operationLogs' }
];

export function Nav({ role, activeView, onSelect }: { role: UserRole; activeView: AppView; onSelect: (view: AppView) => void }) {
  return (
    <nav className="nav" aria-label="Primary">
      {items
        .filter((item) => item.roles.includes(role))
        .map((item) => {
          const Icon = item.icon;
          return (
            <button className={`navItem ${activeView === item.view ? 'active' : ''}`} type="button" key={item.label} onClick={() => onSelect(item.view)}>
              <Icon size={18} />
              <span>{item.label}</span>
            </button>
          );
        })}
    </nav>
  );
}
