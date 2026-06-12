import { FormEvent, useEffect, useMemo, useState } from 'react';
import { ArrowRight, CreditCard, Plus, RotateCcw } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Contract, getContract } from '../../api/contracts';
import { createPaymentPlan, listPaymentContracts } from '../../api/payments';
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
import { ActionMenu, BulkActionBar, Button, DataTable, TextField, Toolbar } from '../../components/ui';
import { useSession } from '../../auth/SessionProvider';
import { contractStatusLabel, labelFor, localizeError, paymentStatusLabel } from '../../i18n/labels';
import { PaymentDetail } from './PaymentDetail';

type Mode = 'list' | 'create' | 'detail';

export function PaymentList({ targetRecordId, onTargetHandled }: { targetRecordId?: string; onTargetHandled?: () => void }) {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [selected, setSelected] = useState<Contract | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [includeArchived, setIncludeArchived] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [form, setForm] = useState({ contractId: '', dueAmount: '', dueDate: '', currency: 'CNY' });

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (!targetRecordId) return;
    void selectPayment(targetRecordId).finally(() => onTargetHandled?.());
  }, [targetRecordId, onTargetHandled]);

  async function refresh(nextSearch = search, nextIncludeArchived = includeArchived) {
    try {
      const response = await listPaymentContracts(nextSearch, nextIncludeArchived);
      setContracts(response.items);
      setPage(1);
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function submitPlan(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      await createPaymentPlan(form.contractId, { dueAmount: form.dueAmount, dueDate: form.dueDate, currency: form.currency });
      setNotice(`已为合同 ${form.contractId} 创建回款计划。`);
      setForm({ contractId: '', dueAmount: '', dueDate: '', currency: 'CNY' });
      setMode('list');
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  const rows = useMemo(() => contracts.filter((contract) => includeArchived || !contract.archived), [contracts, includeArchived]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((contract) => selectedIds.includes(contract.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function toggleRow(contract: Contract, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, contract.id])] : value.filter((id) => id !== contract.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((contract) => contract.id) : []);
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('payments-contracts-selected.csv', selectedRows.map((contract) => ({
      合同: contract.id,
      商机: contract.opportunityId,
      合同状态: labelFor(contractStatusLabel, contract.status),
      合同金额: contract.amount
    })));
    setNotice(`已导出 ${selectedRows.length} 条回款合同。`);
  }

  async function selectPayment(contractId: string) {
    setError('');
    try {
      setSelected(await getContract(contractId));
      setMode('detail');
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  function startPlan(contract: Contract) {
    setForm({ contractId: contract.id, dueAmount: '', dueDate: '', currency: 'CNY' });
    setMode('create');
  }

  if (mode === 'create') {
    return (
      <FormShell title="新建回款计划" description="为已授权合同创建回款计划。" badge="新建" onCancel={() => setMode('list')} actions={<Button variant="primary" form="payment-plan-form" type="submit">保存</Button>} side={<PaymentFormRules />}>
        {error && <div role="alert" className="error">{error}</div>}
        <form id="payment-plan-form" className="actionBand" onSubmit={submitPlan}>
          <FormSection title="回款计划信息">
            <div className="formFields">
              <TextField label="合同 ID" value={form.contractId} onChange={(event) => setForm({ ...form, contractId: event.target.value })} />
              <TextField label="计划金额" value={form.dueAmount} onChange={(event) => setForm({ ...form, dueAmount: event.target.value })} />
              <TextField label="计划到期日" type="date" value={form.dueDate} onChange={(event) => setForm({ ...form, dueDate: event.target.value })} />
              <TextField label="计划币种" value={form.currency} onChange={(event) => setForm({ ...form, currency: event.target.value })} />
            </div>
          </FormSection>
          <div className="saveBar"><Button onClick={() => setMode('list')}>取消</Button><Button variant="primary" type="submit">保存回款计划</Button></div>
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    return <PaymentDetail key={selected.id} contract={selected} onError={setError} onBack={() => setMode('list')} error={error} />;
  }

  return (
    <CrudListShell
      title="回款"
      description={`${scopeLabel} · 共 ${rows.length} 条合同 · 默认按更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={<><Button onClick={() => void refresh(search, includeArchived)}><RotateCcw size={16} aria-hidden="true" />刷新</Button><Button variant="primary" onClick={() => setMode('create')}><Plus size={16} aria-hidden="true" />新建回款计划</Button></>}
      toolbar={<Toolbar searchValue={search} onSearchChange={setSearch} searchPlaceholder="搜索合同、商机或客户" filters={<label className="inlineCheckbox"><input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />包含已归档</label>} actions={<Button onClick={() => void refresh(search, includeArchived)}>应用筛选</Button>} />}
      activeFilters={<ActiveFilterSummary onClear={() => { setSearch(''); setIncludeArchived(false); setSelectedIds([]); void refresh('', false); }}><span className="chip">负责人：{scopeLabel}</span><span className="chip">归档：{includeArchived ? '包含已归档' : '仅活动记录'}</span></ActiveFilterSummary>}
      bulkBar={<BulkActionBar><div className="bulkSummary"><span className="bulkCount">已选择 {selectedRows.length} 条</span><span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '回款页当前按合同授权范围展示；转移/归档走合同页面。')}</span></div><div className="bulkActions">{user?.role !== 'Sales' ? <><Button className="bulkButton" disabled title="回款页不执行负责人转移；按 A3 禁用。">批量转移负责人</Button><Button className="bulkButton" disabled title="回款页归档请在合同页面执行；按 A3 禁用。">批量归档</Button></> : null}<ExportSelectedButton disabled={selectedRows.length === 0} onExport={exportSelected} /><Button className="bulkButton" variant="primary" disabled={selectedRows.length === 0} onClick={() => setSelectedIds([])}>清除选择</Button></div></BulkActionBar>}
      table={
        <DataTable
          caption="回款合同结果表"
          rows={slice.items}
          rowKey={(contract) => contract.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(contract) => selected?.id === contract.id ? 'selected' : undefined}
          onRowClick={(contract) => void selectPayment(contract.id)}
          getRowAriaLabel={(contract) => `打开回款合同 ${contract.id}`}
          empty="没有符合当前筛选条件的回款合同。"
          columns={[
            { key: 'contract', header: '合同', width: '230px', render: (contract) => <RecordIdentity icon={<CreditCard size={17} aria-hidden="true" />} title={contract.opportunityId} subtitle={contract.id} tone="mint" /> },
            { key: 'status', header: '合同状态', render: (contract) => <StatusPill tone={contract.status === 'Signed' || contract.status === 'Active' || contract.status === 'Completed' ? 'success' : 'warning'}>{labelFor(contractStatusLabel, contract.status)}</StatusPill> },
            { key: 'payment', header: '回款状态', render: () => <StatusPill>{labelFor(paymentStatusLabel, 'No plan')}</StatusPill> },
            { key: 'amount', header: '合同金额', align: 'right', render: (contract) => <span className="amountText">{money(contract.amount)}</span> },
            { key: 'updated', header: '更新时间', align: 'right', sortable: true, sortDirection: 'desc', render: (contract) => formatDate(contract.updatedAt) }
          ]}
          actions={(contract) => (
            <div className="rowActions">
              <button className="rowAction" type="button" aria-label={`查看回款 ${contract.id}`} onClick={() => { setError(''); setSelected(contract); setMode('detail'); }}><ArrowRight size={16} aria-hidden="true" /></button>
              <ActionMenu
                label={`打开回款 ${contract.id} 的行操作菜单`}
                items={[
                  { label: '查看', onSelect: () => { setError(''); setSelected(contract); setMode('detail'); } },
                  { label: '新建计划', onSelect: () => startPlan(contract), reason: '使用现有回款计划表单。' },
                  { label: '登记回款', onSelect: () => { setError(''); setSelected(contract); setMode('detail'); }, reason: '在详情页填写回款金额和幂等键。' }
                ]}
              />
            </div>
          )}
        />
      }
      pagination={<CrudPagination slice={slice} onPageChange={setPage} onPageSizeChange={(next) => { setPageSize(next); setPage(1); }} />}
    />
  );
}

function PaymentFormRules() {
  return <div className="sideCard"><h3>字段状态</h3><div className="rule"><span>1</span><p>回款计划必须绑定合同 ID。</p></div><div className="rule"><span>2</span><p>回款金额不能超过合同剩余应收。</p></div></div>;
}

function money(value: string) {
  const number = Number(value);
  if (!Number.isFinite(number)) return value || '未填写';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY', maximumFractionDigits: 0 }).format(number);
}

function formatDate(value: string) {
  if (!value) return '未更新';
  return value.length > 10 ? value.slice(0, 10) : value;
}
