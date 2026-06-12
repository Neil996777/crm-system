import { FormEvent, useEffect, useMemo, useState } from 'react';
import { FileSignature, Plus, RotateCcw } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Contract, archiveContract, changeContractStatus, createContract, getContract, listContracts } from '../../api/contracts';
import {
  ActiveFilterSummary,
  CrudListShell,
  CrudPagination,
  DEFAULT_PAGE_SIZE,
  ExportSelectedButton,
  FormSection,
  FormShell,
  ListAsyncFeedback,
  ListTableLoading,
  RecordIdentity,
  StatusPill,
  exportRows,
  paginate
} from '../../components/CrudScaffold';
import { ActionMenu, BulkActionBar, Button, DataTable, TextAreaField, TextField, Toolbar } from '../../components/ui';
import { useSession } from '../../auth/SessionProvider';
import { contractStatusLabel, labelFor, localizeError } from '../../i18n/labels';
import { ContractDetail } from './ContractDetail';

type Mode = 'list' | 'create' | 'detail';
const statuses = Object.keys(contractStatusLabel);

export function ContractList({ targetRecordId, onTargetHandled }: { targetRecordId?: string; onTargetHandled?: () => void }) {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [selected, setSelected] = useState<Contract | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [includeArchived, setIncludeArchived] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [loading, setLoading] = useState(true);
  const [selectingId, setSelectingId] = useState('');
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [form, setForm] = useState({
    quoteId: '',
    opportunityId: '',
    customerId: '',
    amount: '',
    expectedSignedDate: '',
    contractNote: '',
    amountDifferenceReason: '',
    ownerId: ''
  });

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (!targetRecordId) return;
    void selectContract(targetRecordId).finally(() => onTargetHandled?.());
  }, [targetRecordId, onTargetHandled]);

  async function refresh(nextSearch = search, nextIncludeArchived = includeArchived) {
    setLoading(true);
    setError('');
    try {
      const response = await listContracts(nextSearch, nextIncludeArchived);
      setContracts(response.items);
      setPage(1);
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    } finally {
      setLoading(false);
    }
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createContract(form);
      setSelected(created);
      setMode('detail');
      setForm({ quoteId: '', opportunityId: '', customerId: '', amount: '', expectedSignedDate: '', contractNote: '', amountDifferenceReason: '', ownerId: '' });
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectContract(id: string) {
    setError('');
    setSelectingId(id);
    try {
      setSelected(await getContract(id));
      setMode('detail');
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    } finally {
      setSelectingId('');
    }
  }

  async function updateSelected(contract: Contract) {
    setSelected(contract);
    await refresh();
  }

  const rows = useMemo(() => contracts.filter((contract) => {
    if (status && contract.status !== status) return false;
    if (!includeArchived && contract.archived) return false;
    return true;
  }), [contracts, status, includeArchived]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((contract) => selectedIds.includes(contract.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function toggleRow(contract: Contract, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, contract.id])] : value.filter((id) => id !== contract.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((contract) => contract.id) : []);
  }

  async function archiveSelected() {
    if (selectedRows.length === 0 || user?.role === 'Sales') return;
    setError('');
    try {
      let archived = 0;
      for (const contract of selectedRows) {
        if (contract.archived) continue;
        await archiveContract(contract, '批量归档合同记录');
        archived += 1;
      }
      setSelectedIds([]);
      setNotice(`已归档 ${archived} 条合同。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function archiveRow(contract: Contract) {
    if (user?.role === 'Sales' || contract.archived) return;
    if (!window.confirm(`确认归档合同 ${contract.id}？`)) return;
    setError('');
    try {
      await archiveContract(contract, '行操作归档合同记录');
      setNotice(`已归档合同 ${contract.id}。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function changeRowStatus(contract: Contract, toStatus: string) {
    if (!canChangeContractTo(contract, toStatus)) return;
    let signedEffectiveDate = contract.signedEffectiveDate ?? '';
    if (toStatus === 'Signed') {
      const input = window.prompt('请输入签署/生效日期', signedEffectiveDate || today());
      if (!input?.trim()) return;
      signedEffectiveDate = input.trim();
    }
    if (!window.confirm(`确认将合同 ${contract.id} 标记为${labelFor(contractStatusLabel, toStatus)}？`)) return;
    setError('');
    try {
      const updated = await changeContractStatus(contract.id, contract.version, toStatus, signedEffectiveDate);
      setSelected(updated);
      setNotice(`合同 ${contract.id} 已标记为${labelFor(contractStatusLabel, toStatus)}。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('contracts-selected.csv', selectedRows.map((contract) => ({
      合同: contract.id,
      报价: contract.quoteId,
      商机: contract.opportunityId,
      状态: labelFor(contractStatusLabel, contract.status),
      金额: contract.amount,
      预计签署: contract.expectedSignedDate
    })));
    setNotice(`已导出 ${selectedRows.length} 条合同。`);
  }

  if (mode === 'create') {
    return (
      <FormShell title="新建合同" description="基于已接受报价创建待签署合同。" badge="新建" onCancel={() => setMode('list')} actions={<Button variant="primary" form="contract-form" type="submit">保存</Button>} side={<ContractFormRules />}>
        {error && <div role="alert" className="error">{error}</div>}
        <form id="contract-form" className="actionBand" onSubmit={submit}>
          <FormSection title="合同基本信息">
            <div className="formFields">
              <TextField label="报价 ID" value={form.quoteId} onChange={(event) => setForm({ ...form, quoteId: event.target.value })} />
              <TextField label="商机 ID" value={form.opportunityId} onChange={(event) => setForm({ ...form, opportunityId: event.target.value })} />
              <TextField label="客户 ID" value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              <TextField label="金额" value={form.amount} onChange={(event) => setForm({ ...form, amount: event.target.value })} />
              <TextField label="预计签署日期" type="date" value={form.expectedSignedDate} onChange={(event) => setForm({ ...form, expectedSignedDate: event.target.value })} />
              <TextField label="负责人 ID" value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              <TextAreaField className="full" label="合同备注" value={form.contractNote} onChange={(event) => setForm({ ...form, contractNote: event.target.value })} />
              <TextField className="full" label="金额差异原因" value={form.amountDifferenceReason} onChange={(event) => setForm({ ...form, amountDifferenceReason: event.target.value })} />
            </div>
          </FormSection>
          <div className="saveBar"><Button onClick={() => setMode('list')}>取消</Button><Button variant="primary" type="submit">保存合同</Button></div>
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    return <ContractDetail contract={selected} onUpdated={updateSelected} onError={setError} onBack={() => setMode('list')} error={error} />;
  }

  return (
    <CrudListShell
      title="合同"
      description={`${scopeLabel} · 共 ${rows.length} 条 · 默认按更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={<><Button disabled={loading} aria-busy={loading || undefined} onClick={() => void refresh(search, includeArchived)}><RotateCcw size={16} aria-hidden="true" />{loading ? '刷新中' : '刷新'}</Button><Button variant="primary" onClick={() => setMode('create')}><Plus size={16} aria-hidden="true" />新建合同</Button></>}
      toolbar={<Toolbar searchValue={search} onSearchChange={setSearch} searchPlaceholder="搜索合同、报价或客户" filters={<><label className="inlineCheckbox"><span>状态</span><select value={status} onChange={(event) => setStatus(event.target.value)}><option value="">全部状态</option>{statuses.map((item) => <option value={item} key={item}>{labelFor(contractStatusLabel, item)}</option>)}</select></label><label className="inlineCheckbox"><input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />包含已归档</label></>} actions={<Button disabled={loading} aria-busy={loading || undefined} onClick={() => void refresh(search, includeArchived)}>应用筛选</Button>} />}
      feedback={<ListAsyncFeedback error={error} loading={loading} selecting={Boolean(selectingId)} />}
      activeFilters={<ActiveFilterSummary onClear={() => { setSearch(''); setStatus(''); setIncludeArchived(false); setSelectedIds([]); void refresh('', false); }}><span className="chip">负责人：{scopeLabel}</span><span className="chip">状态：{status ? labelFor(contractStatusLabel, status) : '全部状态'}</span><span className="chip">归档：{includeArchived ? '包含已归档' : '仅活动记录'}</span></ActiveFilterSummary>}
      bulkBar={<BulkActionBar><div className="bulkSummary"><span className="bulkCount">已选择 {selectedRows.length} 条</span><span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '合同批量归档通过现有单条归档接口逐条执行。')}</span></div><div className="bulkActions">{user?.role !== 'Sales' ? <Button className="bulkButton" disabled={selectedRows.length === 0} onClick={() => void archiveSelected()}>批量归档</Button> : null}<ExportSelectedButton disabled={selectedRows.length === 0} onExport={exportSelected} /><Button className="bulkButton" variant="primary" disabled={selectedRows.length === 0} onClick={() => setSelectedIds([])}>清除选择</Button></div></BulkActionBar>}
      table={
        loading ? <ListTableLoading label="正在加载合同列表..." /> : <DataTable
          caption="合同结果表"
          rows={slice.items}
          rowKey={(contract) => contract.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(contract) => selected?.id === contract.id ? 'selected' : undefined}
          onRowClick={(contract) => void selectContract(contract.id)}
          getRowAriaLabel={(contract) => `打开合同 ${contract.id}`}
          empty="没有符合当前筛选条件的合同。"
          columns={[
            {
              key: 'contract',
              header: '合同',
              width: '220px',
              render: (contract) => (
                <RecordIdentity
                  icon={<FileSignature size={17} aria-hidden="true" />}
                  title={contract.opportunityId}
                  titleAriaLabel={`打开合同 ${contract.id}`}
                  titleBusy={selectingId === contract.id}
                  subtitle={`${contract.archived ? '已归档 · ' : ''}${contract.id}`}
                  tone="peach"
                  onTitleClick={() => void selectContract(contract.id)}
                />
              )
            },
            { key: 'quote', header: '报价', render: (contract) => contract.quoteId },
            { key: 'status', header: '状态', render: (contract) => <StatusPill tone={contractTone(contract.status)}>{labelFor(contractStatusLabel, contract.status)}</StatusPill> },
            { key: 'amount', header: '金额', align: 'right', render: (contract) => <span className="amountText">{money(contract.amount)}</span> },
            { key: 'signed', header: '预计签署', align: 'right', render: (contract) => contract.expectedSignedDate },
            { key: 'updated', header: '更新时间', align: 'right', sortable: true, sortDirection: 'desc', render: (contract) => formatDate(contract.updatedAt) }
          ]}
          actions={(contract) => (
            <div className="rowActions">
              <ActionMenu
                label={`打开合同 ${contract.id} 的行操作菜单`}
                items={[
                  { label: '查看', onSelect: () => void selectContract(contract.id) },
                  {
                    label: '签署',
                    onSelect: () => void changeRowStatus(contract, 'Signed'),
                    disabled: !canChangeContractTo(contract, 'Signed'),
                    reason: contract.status !== 'Pending Signature' ? '仅待签署合同可签署。' : undefined
                  },
                  {
                    label: '启用',
                    onSelect: () => void changeRowStatus(contract, 'Active'),
                    disabled: !canChangeContractTo(contract, 'Active'),
                    reason: contract.status !== 'Signed' ? '仅已签署合同可启用。' : undefined
                  },
                  {
                    label: '完成',
                    onSelect: () => void changeRowStatus(contract, 'Completed'),
                    disabled: !canChangeContractTo(contract, 'Completed'),
                    reason: contract.status !== 'Active' ? '仅启用合同可完成。' : undefined
                  },
                  {
                    label: '终止',
                    onSelect: () => void changeRowStatus(contract, 'Terminated'),
                    disabled: !canChangeContractTo(contract, 'Terminated'),
                    reason: contract.status === 'Completed' || contract.status === 'Terminated' ? '已完成或已终止合同不能再终止。' : undefined,
                    tone: 'danger'
                  },
                  {
                    label: '归档',
                    onSelect: () => void archiveRow(contract),
                    disabled: user?.role === 'Sales' || Boolean(contract.archived),
                    reason: user?.role === 'Sales' ? '销售角色不能归档合同。' : contract.archived ? '已归档。' : undefined,
                    tone: 'danger'
                  }
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

function canChangeContractTo(contract: Contract, toStatus: string) {
  if (toStatus === 'Signed') return contract.status === 'Pending Signature';
  if (toStatus === 'Active') return contract.status === 'Signed';
  if (toStatus === 'Completed') return contract.status === 'Active';
  if (toStatus === 'Terminated') return contract.status !== 'Completed' && contract.status !== 'Terminated';
  return false;
}

function ContractFormRules() {
  return <div className="sideCard"><h3>字段状态</h3><div className="rule"><span>1</span><p>合同保存时状态固定为待签署。</p></div><div className="rule"><span>2</span><p>金额与报价不一致时必须填写差异原因。</p></div></div>;
}

function contractTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Signed' || status === 'Active' || status === 'Completed') return 'success';
  if (status === 'Terminated') return 'danger';
  if (status === 'Pending Signature') return 'warning';
  return 'neutral';
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

function today() {
  return new Date().toISOString().slice(0, 10);
}
