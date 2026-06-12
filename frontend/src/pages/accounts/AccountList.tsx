import { FormEvent, useEffect, useMemo, useState } from 'react';
import { Building2, Plus, RotateCcw } from 'lucide-react';
import { Account, archiveAccount, checkAccountDuplicate, createAccount, getAccount, listAccounts } from '../../api/accounts';
import { ApiError } from '../../api/client';
import { DuplicateWarningResult } from '../../api/duplicates';
import {
  ActiveFilterSummary,
  CrudListShell,
  CrudPagination,
  DEFAULT_PAGE_SIZE,
  ExportSelectedButton,
  FormSection,
  FormShell,
  RecordIdentity,
  StatusPill,
  exportRows,
  paginate
} from '../../components/CrudScaffold';
import { DuplicateWarning } from '../../components/DuplicateWarning';
import { ActionMenu, BulkActionBar, Button, DataTable, TextField, Toolbar } from '../../components/ui';
import { useSession } from '../../auth/SessionProvider';
import { accountStatusLabel, archiveStatusLabel, labelFor, localizeError } from '../../i18n/labels';
import { AccountDetail } from './AccountDetail';

type Mode = 'list' | 'create' | 'detail';
const statuses = Object.keys(accountStatusLabel);

export function AccountList({ targetRecordId, onTargetHandled }: { targetRecordId?: string; onTargetHandled?: () => void }) {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [selected, setSelected] = useState<Account | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [includeArchived, setIncludeArchived] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [duplicateWarning, setDuplicateWarning] = useState<DuplicateWarningResult | null>(null);
  const [form, setForm] = useState({ companyName: '', customerStatus: '', ownerId: '' });

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (!targetRecordId) return;
    void selectAccount(targetRecordId).finally(() => onTargetHandled?.());
  }, [targetRecordId, onTargetHandled]);

  async function refresh(nextSearch = search, nextIncludeArchived = includeArchived) {
    const response = await listAccounts(nextSearch, nextIncludeArchived);
    setAccounts(response.items);
    setPage(1);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    await saveAccount();
  }

  async function saveAccount(proceedWarningToken?: string) {
    setError('');
    try {
      if (!proceedWarningToken) {
        const warning = await checkAccountDuplicate({ companyName: form.companyName });
        if (warning.result === 'PossibleDuplicate' && warning.warningToken) {
          setDuplicateWarning(warning);
          return;
        }
      }
      const account = await createAccount({ ...form, proceedWarningToken });
      setSelected(account);
      setMode('detail');
      setForm({ companyName: '', customerStatus: '', ownerId: '' });
      setDuplicateWarning(null);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectAccount(id: string) {
    setError('');
    setSelected(await getAccount(id));
    setMode('detail');
  }

  const rows = useMemo(() => accounts.filter((account) => {
    if (status && account.customerStatus !== status) return false;
    if (!includeArchived && account.archived) return false;
    return true;
  }), [accounts, status, includeArchived]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((account) => selectedIds.includes(account.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function toggleRow(account: Account, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, account.id])] : value.filter((id) => id !== account.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((account) => account.id) : []);
  }

  function clearFilters() {
    setSearch('');
    setStatus('');
    setIncludeArchived(false);
    setSelectedIds([]);
    void refresh('', false);
  }

  async function archiveSelected() {
    if (selectedRows.length === 0 || user?.role === 'Sales') return;
    setError('');
    try {
      let archived = 0;
      for (const account of selectedRows) {
        if (account.archived) continue;
        await archiveAccount(account.id, account.version, '批量归档客户记录');
        archived += 1;
      }
      setSelectedIds([]);
      setNotice(`已归档 ${archived} 条客户。`);
      await refresh();
    } catch (caught) {
      const error = caught as ApiError;
      setError(localizeError(error));
    }
  }

  async function archiveRow(account: Account) {
    if (user?.role === 'Sales' || account.archived) return;
    if (!window.confirm(`确认归档客户 ${account.companyName}？`)) return;
    setError('');
    try {
      await archiveAccount(account.id, account.version, '行操作归档客户记录');
      setNotice(`已归档客户 ${account.companyName}。`);
      await refresh();
    } catch (caught) {
      const error = caught as ApiError;
      setError(localizeError(error));
    }
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('accounts-selected.csv', selectedRows.map((account) => ({
      客户: account.companyName,
      状态: account.archived ? labelFor(archiveStatusLabel, 'Archived') : labelFor(accountStatusLabel, account.customerStatus),
      负责人: account.ownerId,
      更新时间: account.updatedAt
    })));
    setNotice(`已导出 ${selectedRows.length} 条客户。`);
  }

  if (mode === 'create') {
    return (
      <FormShell
        title="新建客户"
        description="创建公司客户记录，保存前执行重复检查。"
        badge="新建"
        onCancel={() => { setMode('list'); setDuplicateWarning(null); }}
        actions={<Button variant="primary" form="account-form" type="submit">保存</Button>}
        side={<AccountFormRules />}
      >
        {error && <div role="alert" className="error">{error}</div>}
        <form id="account-form" className="actionBand" onSubmit={submit}>
          <FormSection title="客户基本信息">
            <div className="formFields">
              <TextField label="公司名称" value={form.companyName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, companyName: event.target.value }); }} />
              <TextField label="客户状态" value={form.customerStatus} onChange={(event) => setForm({ ...form, customerStatus: event.target.value })} />
              <TextField label="负责人 ID" value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
            </div>
          </FormSection>
          <div className="saveBar">
            <Button onClick={() => setMode('list')}>取消</Button>
            <Button variant="primary" type="submit">保存客户</Button>
          </div>
          {duplicateWarning ? (
            <DuplicateWarning
              warning={duplicateWarning}
              onProceed={() => void saveAccount(duplicateWarning.warningToken)}
              onCancel={() => setDuplicateWarning(null)}
            />
          ) : null}
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    return <AccountDetail account={selected} onArchived={() => { setSelected(null); setMode('list'); void refresh(); }} onError={setError} onBack={() => setMode('list')} error={error} />;
  }

  return (
    <CrudListShell
      title="公司/客户"
      description={`${scopeLabel} · 共 ${rows.length} 条 · 默认按更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={
        <>
          <Button onClick={() => void refresh(search, includeArchived)}>
            <RotateCcw size={16} aria-hidden="true" />
            刷新
          </Button>
          <Button variant="primary" onClick={() => { setDuplicateWarning(null); setMode('create'); }}>
            <Plus size={16} aria-hidden="true" />
            新建客户
          </Button>
        </>
      }
      toolbar={
        <Toolbar
          searchValue={search}
          onSearchChange={setSearch}
          searchPlaceholder="搜索公司客户"
          filters={
            <>
              <label className="inlineCheckbox">
                <span>状态</span>
                <select value={status} onChange={(event) => setStatus(event.target.value)}>
                  <option value="">全部状态</option>
                  {statuses.map((statusOption) => (
                    <option value={statusOption} key={statusOption}>{labelFor(accountStatusLabel, statusOption)}</option>
                  ))}
                </select>
              </label>
              <label className="inlineCheckbox">
                <input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />
                包含已归档
              </label>
            </>
          }
          actions={<Button onClick={() => void refresh(search, includeArchived)}>应用筛选</Button>}
        />
      }
      activeFilters={
        <ActiveFilterSummary onClear={clearFilters}>
          <span className="chip">负责人：{scopeLabel}</span>
          <span className="chip">状态：{status ? labelFor(accountStatusLabel, status) : '全部状态'}</span>
          <span className="chip">归档：{includeArchived ? '包含已归档' : '仅活动记录'}</span>
        </ActiveFilterSummary>
      }
      bulkBar={
        <BulkActionBar>
          <div className="bulkSummary">
            <span className="bulkCount">已选择 {selectedRows.length} 条</span>
            <span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '客户批量归档通过现有单条归档接口逐条执行。')}</span>
          </div>
          <div className="bulkActions">
            {user?.role !== 'Sales' ? (
              <>
                <Button className="bulkButton" disabled title="客户当前无负责人转移接口；按 A3 禁用。">批量转移负责人</Button>
                <Button className="bulkButton" disabled={selectedRows.length === 0} onClick={() => void archiveSelected()}>批量归档</Button>
              </>
            ) : null}
            <ExportSelectedButton disabled={selectedRows.length === 0} onExport={exportSelected} />
            <Button className="bulkButton" variant="primary" onClick={() => setSelectedIds([])} disabled={selectedRows.length === 0}>清除选择</Button>
          </div>
        </BulkActionBar>
      }
      table={
        <DataTable
          caption="客户结果表"
          rows={slice.items}
          rowKey={(account) => account.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(account) => selected?.id === account.id ? 'selected' : undefined}
          onRowClick={(account) => void selectAccount(account.id)}
          getRowAriaLabel={(account) => `打开客户 ${account.companyName}`}
          empty="没有符合当前筛选条件的客户。"
          columns={[
            {
              key: 'account',
              header: '客户',
              width: '250px',
              render: (account) => (
                <RecordIdentity
                  icon={<Building2 size={17} aria-hidden="true" />}
                  title={account.companyName}
                  titleAriaLabel={`打开客户 ${account.companyName}`}
                  subtitle={`${account.archived ? '已归档 · ' : ''}更新于 ${formatDate(account.updatedAt)}`}
                  tone={account.archived ? 'peach' : 'sky'}
                  onTitleClick={() => void selectAccount(account.id)}
                />
              )
            },
            { key: 'status', header: '状态', render: (account) => <StatusPill tone={account.archived ? 'warning' : accountTone(account.customerStatus)}>{account.archived ? labelFor(archiveStatusLabel, 'Archived') : labelFor(accountStatusLabel, account.customerStatus)}</StatusPill> },
            { key: 'owner', header: '负责人', render: (account) => account.ownerId },
            { key: 'version', header: '版本', align: 'right', render: (account) => `v${account.version}` },
            { key: 'updated', header: '更新时间', align: 'right', sortable: true, sortDirection: 'desc', render: (account) => formatDate(account.updatedAt) }
          ]}
          actions={(account) => (
            <div className="rowActions">
              <ActionMenu
                label={`打开 ${account.companyName} 的行操作菜单`}
                items={[
                  { label: '查看', onSelect: () => void selectAccount(account.id) },
                  {
                    label: '归档',
                    onSelect: () => void archiveRow(account),
                    disabled: user?.role === 'Sales' || account.archived,
                    reason: user?.role === 'Sales' ? '销售角色不能归档客户。' : account.archived ? '已归档。' : undefined,
                    tone: 'danger'
                  }
                ]}
              />
            </div>
          )}
        />
      }
      pagination={
        <CrudPagination
          slice={slice}
          onPageChange={setPage}
          onPageSizeChange={(next) => {
            setPageSize(next);
            setPage(1);
          }}
        />
      }
    />
  );
}

function AccountFormRules() {
  return (
    <div className="sideCard">
      <h3>字段状态</h3>
      <div className="rule"><span>1</span><p>公司名称必填，保存前执行客户重复检查。</p></div>
      <div className="rule"><span>2</span><p>客户状态使用真实枚举值，例如 Prospect / Active / Inactive。</p></div>
      <div className="rule"><span>3</span><p>归档前会按现有接口检查未完成事项。</p></div>
    </div>
  );
}

function accountTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Active') return 'success';
  if (status === 'Inactive') return 'warning';
  return 'neutral';
}

function formatDate(value: string) {
  if (!value) return '未更新';
  return value.length > 10 ? value.slice(0, 10) : value;
}
