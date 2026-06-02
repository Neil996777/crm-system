import { FormEvent, useEffect, useMemo, useState } from 'react';
import { UserRole } from '../../api/auth';
import { ApiError } from '../../api/client';
import { ManagedUser, UserStatus, changeUserRole, changeUserStatus, createUser, listUsers } from '../../api/users';
import { RoleStatusChangeDialog } from '../../components/RoleStatusChangeDialog';

const roles: UserRole[] = ['Administrator', 'Sales Manager', 'Sales'];
const statuses: UserStatus[] = ['Active', 'Disabled'];

export function UserManagement() {
  const [users, setUsers] = useState<ManagedUser[]>([]);
  const [activeAdministratorCount, setActiveAdministratorCount] = useState(0);
  const [selected, setSelected] = useState<ManagedUser | null>(null);
  const [form, setForm] = useState({ email: '', displayName: '', password: '', role: 'Sales' as UserRole });
  const [nextRole, setNextRole] = useState<UserRole>('Sales');
  const [nextStatus, setNextStatus] = useState<UserStatus>('Active');
  const [confirming, setConfirming] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh() {
    const response = await listUsers();
    setUsers(response.items);
    setActiveAdministratorCount(response.activeAdministratorCount);
    setSelected((current) => current ? response.items.find((user) => user.id === current.id) ?? current : current);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createUser(form);
      setForm({ email: '', displayName: '', password: '', role: 'Sales' });
      setSelected(created);
      setNextRole(created.role);
      setNextStatus(created.status);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
    }
  }

  function edit(user: ManagedUser) {
    setSelected(user);
    setNextRole(user.role);
    setNextStatus(user.status);
    setConfirming(false);
    setError('');
  }

  const lastAdminBlocked = useMemo(() => {
    if (!selected) return false;
    return selected.role === 'Administrator' && selected.status === 'Active' && activeAdministratorCount <= 1 && (nextRole !== 'Administrator' || nextStatus !== 'Active');
  }, [selected, nextRole, nextStatus, activeAdministratorCount]);

  async function confirm() {
    if (!selected) return;
    setError('');
    try {
      let updated = selected;
      if (nextRole !== selected.role) {
        updated = await changeUserRole(selected.id, nextRole);
      }
      if (nextStatus !== updated.status) {
        updated = await changeUserStatus(selected.id, nextStatus);
      }
      setSelected(updated);
      setConfirming(false);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
      setConfirming(false);
      await refresh();
    }
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>User Management</h1>
          <p>Administrator-only role and status governance.</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="createPanel" onSubmit={submit}>
            <label>
              Email
              <input value={form.email} onChange={(event) => setForm({ ...form, email: event.target.value })} />
            </label>
            <label>
              Display name
              <input value={form.displayName} onChange={(event) => setForm({ ...form, displayName: event.target.value })} />
            </label>
            <label>
              Password
              <input type="password" value={form.password} onChange={(event) => setForm({ ...form, password: event.target.value })} />
            </label>
            <label>
              Role
              <select value={form.role} onChange={(event) => setForm({ ...form, role: event.target.value as UserRole })}>
                {roles.map((role) => <option key={role} value={role}>{role}</option>)}
              </select>
            </label>
            <button className="primaryButton" type="submit">Create user</button>
          </form>
          <table className="dataTable" aria-label="User table">
            <thead>
              <tr>
                <th>User</th>
                <th>Role</th>
                <th>Status</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.displayName}<br />{user.email}</td>
                  <td>{user.role}</td>
                  <td>{user.status}</td>
                  <td><button className="secondaryButton" type="button" onClick={() => edit(user)}>Edit {user.displayName}</button></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <section className="detailPane" aria-label="User detail">
          {selected ? (
            <>
              <div className="detailHeader">
                <div>
                  <h2>{selected.displayName}</h2>
                  <p>{selected.email}</p>
                </div>
                <span className="statusPill">{selected.status}</span>
              </div>
              <label>
                New role
                <select value={nextRole} onChange={(event) => setNextRole(event.target.value as UserRole)}>
                  {roles.map((role) => <option key={role} value={role}>{role}</option>)}
                </select>
              </label>
              <label>
                New status
                <select value={nextStatus} onChange={(event) => setNextStatus(event.target.value as UserStatus)}>
                  {statuses.map((status) => <option key={status} value={status}>{status}</option>)}
                </select>
              </label>
              {lastAdminBlocked && <div role="alert" className="error">Last active Administrator change is blocked.</div>}
              <button className="primaryButton" type="button" disabled={lastAdminBlocked || (nextRole === selected.role && nextStatus === selected.status)} onClick={() => setConfirming(true)}>
                Review role/status change
              </button>
              {confirming && <RoleStatusChangeDialog user={selected} nextRole={nextRole} nextStatus={nextStatus} onCancel={() => setConfirming(false)} onConfirm={() => void confirm()} />}
            </>
          ) : (
            <p className="emptyState">Select a user to review role and status.</p>
          )}
        </section>
      </section>
    </main>
  );
}
