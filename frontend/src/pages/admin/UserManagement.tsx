import { FormEvent, useEffect, useMemo, useState } from 'react';
import { Plus, ShieldCheck } from 'lucide-react';
import { UserRole } from '../../api/auth';
import { ApiError } from '../../api/client';
import { ManagedUser, UserStatus, changeUserRole, changeUserStatus, createUser, listUsers } from '../../api/users';
import { RoleStatusChangeDialog } from '../../components/RoleStatusChangeDialog';
import { Badge, Button, DataTable, PageHeader, Pagination, Toolbar } from '../../components/ui';
import { exportRows } from '../../components/CrudScaffold';
import { labelFor, localizeError, roleLabel, userStatusLabel } from '../../i18n/labels';

const roles: UserRole[] = ['Administrator', 'Sales Manager', 'Sales'];
const statuses: UserStatus[] = ['Active', 'Disabled'];
const pageSizeOptions = [5, 10, 20];

export function UserManagement() {
  const [users, setUsers] = useState<ManagedUser[]>([]);
  const [activeAdministratorCount, setActiveAdministratorCount] = useState(0);
  const [selected, setSelected] = useState<ManagedUser | null>(null);
  const [form, setForm] = useState({ email: '', displayName: '', password: '', role: 'Sales' as UserRole });
  const [nextRole, setNextRole] = useState<UserRole>('Sales');
  const [nextStatus, setNextStatus] = useState<UserStatus>('Active');
  const [confirming, setConfirming] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const [search, setSearch] = useState('');
  const [roleFilter, setRoleFilter] = useState<'all' | UserRole>('all');
  const [statusFilter, setStatusFilter] = useState<'all' | UserStatus>('all');
  const [adminOnly, setAdminOnly] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
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
      setShowCreate(false);
      setSelected(created);
      setNextRole(created.role);
      setNextStatus(created.status);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  function edit(user: ManagedUser) {
    setSelected(user);
    setNextRole(user.role);
    setNextStatus(user.status);
    setConfirming(false);
    setShowCreate(false);
    setError('');
  }

  function prepareStatusChange(user: ManagedUser) {
    edit(user);
    setNextStatus(user.status === 'Active' ? 'Disabled' : 'Active');
  }

  function prepareRoleChange(user: ManagedUser) {
    edit(user);
  }

  function isLastActiveAdministrator(user: ManagedUser) {
    return user.role === 'Administrator' && user.status === 'Active' && activeAdministratorCount <= 1;
  }

  const filteredUsers = useMemo(() => users.filter((user) => {
    const query = search.trim().toLowerCase();
    const matchesSearch = !query || user.displayName.toLowerCase().includes(query) || user.email.toLowerCase().includes(query) || labelFor(roleLabel, user.role).includes(search.trim());
    const matchesRole = roleFilter === 'all' || user.role === roleFilter;
    const matchesStatus = statusFilter === 'all' || user.status === statusFilter;
    const matchesAdminOnly = !adminOnly || user.role === 'Administrator';
    return matchesSearch && matchesRole && matchesStatus && matchesAdminOnly;
  }), [users, search, roleFilter, statusFilter, adminOnly]);

  const totalPages = Math.max(1, Math.ceil(filteredUsers.length / pageSize));
  const safePage = Math.min(page, totalPages);
  const pagedUsers = filteredUsers.slice((safePage - 1) * pageSize, safePage * pageSize);

  const lastAdminBlocked = useMemo(() => {
    if (!selected) return false;
    return selected.role === 'Administrator' && selected.status === 'Active' && activeAdministratorCount <= 1 && (nextRole !== 'Administrator' || nextStatus !== 'Active');
  }, [selected, nextRole, nextStatus, activeAdministratorCount]);
  const selectedIsLastActiveAdmin = Boolean(selected && selected.role === 'Administrator' && selected.status === 'Active' && activeAdministratorCount <= 1);

  function exportFilteredUsers() {
    exportRows('users-filtered.csv', filteredUsers.map((user) => ({
      显示名: user.displayName,
      邮箱: user.email,
      角色: labelFor(roleLabel, user.role),
      状态: labelFor(userStatusLabel, user.status)
    })));
  }

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
      setError(localizeError(apiError));
      setConfirming(false);
      await refresh();
    }
  }

  return (
    <main className="content adminUsersPage" data-uiux="admin-users">
      <PageHeader
        title="用户与角色"
        description="仅管理员可访问 · 字段：显示名 / 邮箱 / 角色 / 状态 · 仅显示公开身份字段"
        actions={(
          <>
            <Button onClick={exportFilteredUsers}>导出</Button>
            <Button variant="primary" onClick={() => { setShowCreate(true); setSelected(null); }}>
              <Plus size={16} aria-hidden="true" />
              新建用户
            </Button>
          </>
        )}
      />
      {error && <div role="alert" className="errorBanner">{error}</div>}

      {showCreate ? (
        <form className="createPanel adminCreatePanel" aria-label="新建用户表单" onSubmit={submit}>
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
          <div className="formActions">
            <button className="primaryButton" type="submit">创建用户</button>
            <button className="secondaryButton" type="button" onClick={() => setShowCreate(false)}>取消</button>
          </div>
        </form>
      ) : null}

      <section className="listPanel adminUserPanel">
        <Toolbar
          searchValue={search}
          onSearchChange={(value) => { setSearch(value); setPage(1); }}
          searchPlaceholder="搜索显示名或邮箱"
          filters={(
            <>
              <label className="compactFilter">
                <span className="srOnly">角色筛选</span>
                <select value={roleFilter} onChange={(event) => { setRoleFilter(event.target.value as 'all' | UserRole); setPage(1); }}>
                  <option value="all">角色 全部</option>
                  {roles.map((role) => <option key={role} value={role}>{labelFor(roleLabel, role)}</option>)}
                </select>
              </label>
              <label className="compactFilter">
                <span className="srOnly">状态筛选</span>
                <select value={statusFilter} onChange={(event) => { setStatusFilter(event.target.value as 'all' | UserStatus); setPage(1); }}>
                  <option value="all">状态 全部</option>
                  {statuses.map((status) => <option key={status} value={status}>{labelFor(userStatusLabel, status)}</option>)}
                </select>
              </label>
              <label className="inlineCheckbox adminOnlyFilter">
                <input checked={adminOnly} type="checkbox" onChange={(event) => { setAdminOnly(event.target.checked); setPage(1); }} />
                仅管理员
              </label>
            </>
          )}
          activeFilters={[
            ...(roleFilter === 'all' ? [] : [{ label: '角色', value: labelFor(roleLabel, roleFilter), tone: 'primary' as const }]),
            ...(statusFilter === 'all' ? [] : [{ label: '状态', value: labelFor(userStatusLabel, statusFilter), tone: 'primary' as const }]),
            ...(adminOnly ? [{ label: '仅管理员', tone: 'primary' as const }] : [])
          ]}
          onClearFilters={() => { setSearch(''); setRoleFilter('all'); setStatusFilter('all'); setAdminOnly(false); setPage(1); }}
          summary={`显示 ${filteredUsers.length} / ${users.length} 个用户`}
        />

        {activeAdministratorCount <= 1 ? (
          <section className="protectBanner" aria-label="末位管理员保护">
            <span className="panelIcon">
              <ShieldCheck size={18} aria-hidden="true" />
            </span>
            <div>
              <strong>末位管理员保护</strong>
              <p>唯一启用管理员的停用或降级操作已置灰。</p>
            </div>
          </section>
        ) : null}

        <DataTable
          caption="用户与角色表"
          rows={pagedUsers}
          rowKey={(user) => user.id}
          empty="暂无符合筛选的用户。"
          getRowClassName={(user) => isLastActiveAdministrator(user) ? 'protectedRow' : undefined}
          columns={[
            { key: 'name', header: '显示名', render: (user) => <UserIdentity user={user} lastActiveAdministrator={isLastActiveAdministrator(user)} /> },
            { key: 'email', header: '邮箱', render: (user) => user.email },
            { key: 'role', header: '角色', render: (user) => <Badge tone={user.role === 'Administrator' ? 'primary' : user.role === 'Sales Manager' ? 'warning' : 'neutral'}>{labelFor(roleLabel, user.role)}</Badge> },
            { key: 'status', header: '状态', render: (user) => <Badge tone={user.status === 'Active' ? 'success' : 'neutral'}>{labelFor(userStatusLabel, user.status)}</Badge> }
          ]}
          actions={(user) => (
            <div className="adminRowActions">
              <button className="secondaryButton" type="button" aria-label={`编辑 ${user.displayName}`} onClick={() => edit(user)}>
                编辑
              </button>
              <button
                className="secondaryButton"
                type="button"
                aria-label={`${user.status === 'Active' ? '停用' : '启用'} ${user.displayName}`}
                disabled={isLastActiveAdministrator(user) && user.status === 'Active'}
                onClick={() => prepareStatusChange(user)}
              >
                {user.status === 'Active' ? '停用' : '启用'}
              </button>
              <button className="secondaryButton" type="button" aria-label={`改角色 ${user.displayName}`} onClick={() => prepareRoleChange(user)}>
                改角色
              </button>
            </div>
          )}
        />

        <div className="tableFooterSummary">
          <span>共 {filteredUsers.length} 个用户 · 角色仅：管理员 / 销售经理 / 销售</span>
          <span>状态仅：启用 / 停用</span>
        </div>
        <Pagination
          page={safePage}
          totalPages={totalPages}
          totalItems={filteredUsers.length}
          pageSize={pageSize}
          onPageChange={setPage}
          onPageSizeChange={(nextSize) => { setPageSize(nextSize); setPage(1); }}
          pageSizeOptions={pageSizeOptions}
        />
      </section>

      {selected ? (
        <section className="detailPane adminEditPanel" aria-label="用户详情">
          <div className="detailHeader">
            <div>
              <h2>{selected.displayName}</h2>
              <p>{selected.email}</p>
            </div>
            <Badge tone={selected.status === 'Active' ? 'success' : 'neutral'}>{labelFor(userStatusLabel, selected.status)}</Badge>
          </div>
          <label>
            新角色
            <select value={nextRole} onChange={(event) => setNextRole(event.target.value as UserRole)}>
              {roles.map((role) => <option key={role} value={role} disabled={selectedIsLastActiveAdmin && role !== 'Administrator'}>{labelFor(roleLabel, role)}</option>)}
            </select>
          </label>
          <label>
            新状态
            <select value={nextStatus} onChange={(event) => setNextStatus(event.target.value as UserStatus)}>
              {statuses.map((status) => <option key={status} value={status} disabled={selectedIsLastActiveAdmin && status !== 'Active'}>{labelFor(userStatusLabel, status)}</option>)}
            </select>
          </label>
          {selectedIsLastActiveAdmin ? <p className="inlineNotice">末位启用管理员受保护，不能降级或停用。</p> : null}
          {lastAdminBlocked && <div role="alert" className="error">不能变更最后一个启用的管理员。</div>}
          <button className="primaryButton" type="button" disabled={lastAdminBlocked || (nextRole === selected.role && nextStatus === selected.status)} onClick={() => setConfirming(true)}>
            复核角色/状态变更
          </button>
          {confirming && <RoleStatusChangeDialog user={selected} nextRole={nextRole} nextStatus={nextStatus} onCancel={() => setConfirming(false)} onConfirm={() => void confirm()} />}
        </section>
      ) : null}
    </main>
  );
}

function UserIdentity({ user, lastActiveAdministrator }: { user: ManagedUser; lastActiveAdministrator: boolean }) {
  return (
    <div className="userIdentity">
      <span className="tableAvatar" aria-hidden="true">{user.displayName.slice(0, 1)}</span>
      <span>
        <strong>{user.displayName}</strong>
        <small>{lastActiveAdministrator ? '唯一启用管理员' : labelFor(roleLabel, user.role)}</small>
      </span>
    </div>
  );
}
