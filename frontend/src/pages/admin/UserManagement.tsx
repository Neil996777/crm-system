import { FormEvent, useEffect, useMemo, useState } from 'react';
import { UserRole } from '../../api/auth';
import { ApiError } from '../../api/client';
import { ManagedUser, UserStatus, changeUserRole, changeUserStatus, createUser, listUsers } from '../../api/users';
import { RoleStatusChangeDialog } from '../../components/RoleStatusChangeDialog';
import { labelFor, roleLabel, userStatusLabel } from '../../i18n/labels';

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
      setError(apiError.safeMessage || '请求失败。');
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
      setError(apiError.safeMessage || '请求失败。');
      setConfirming(false);
      await refresh();
    }
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>用户管理</h1>
          <p>仅管理员可维护角色和状态。</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="createPanel" onSubmit={submit}>
            <label>
              邮箱
              <input value={form.email} onChange={(event) => setForm({ ...form, email: event.target.value })} />
            </label>
            <label>
              显示名称
              <input value={form.displayName} onChange={(event) => setForm({ ...form, displayName: event.target.value })} />
            </label>
            <label>
              密码
              <input type="password" value={form.password} onChange={(event) => setForm({ ...form, password: event.target.value })} />
            </label>
            <label>
              角色
              <select value={form.role} onChange={(event) => setForm({ ...form, role: event.target.value as UserRole })}>
                {roles.map((role) => <option key={role} value={role}>{labelFor(roleLabel, role)}</option>)}
              </select>
            </label>
            <button className="primaryButton" type="submit">创建用户</button>
          </form>
          <table className="dataTable" aria-label="用户表">
            <thead>
              <tr>
                <th>用户</th>
                <th>角色</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.displayName}<br />{user.email}</td>
                  <td>{labelFor(roleLabel, user.role)}</td>
                  <td>{labelFor(userStatusLabel, user.status)}</td>
                  <td><button className="secondaryButton" type="button" onClick={() => edit(user)}>编辑 {user.displayName}</button></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <section className="detailPane" aria-label="用户详情">
          {selected ? (
            <>
              <div className="detailHeader">
                <div>
                  <h2>{selected.displayName}</h2>
                  <p>{selected.email}</p>
                </div>
                <span className="statusPill">{labelFor(userStatusLabel, selected.status)}</span>
              </div>
              <label>
                新角色
                <select value={nextRole} onChange={(event) => setNextRole(event.target.value as UserRole)}>
                  {roles.map((role) => <option key={role} value={role}>{labelFor(roleLabel, role)}</option>)}
                </select>
              </label>
              <label>
                新状态
                <select value={nextStatus} onChange={(event) => setNextStatus(event.target.value as UserStatus)}>
                  {statuses.map((status) => <option key={status} value={status}>{labelFor(userStatusLabel, status)}</option>)}
                </select>
              </label>
              {lastAdminBlocked && <div role="alert" className="error">不能变更最后一个启用的管理员。</div>}
              <button className="primaryButton" type="button" disabled={lastAdminBlocked || (nextRole === selected.role && nextStatus === selected.status)} onClick={() => setConfirming(true)}>
                复核角色/状态变更
              </button>
              {confirming && <RoleStatusChangeDialog user={selected} nextRole={nextRole} nextStatus={nextStatus} onCancel={() => setConfirming(false)} onConfirm={() => void confirm()} />}
            </>
          ) : (
            <p className="emptyState">选择用户以复核角色和状态。</p>
          )}
        </section>
      </section>
    </main>
  );
}
