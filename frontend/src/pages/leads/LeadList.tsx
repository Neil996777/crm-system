import { FormEvent, useEffect, useMemo, useState } from 'react';
import { ArrowRight, ListChecks, Plus, RotateCcw } from 'lucide-react';
import { ApiError } from '../../api/client';
import { DuplicateWarningResult } from '../../api/duplicates';
import { ConversionResult, Lead, archiveLead, checkLeadDuplicate, createLead, getLead, listLeads, transferLeadOwner } from '../../api/leads';
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
import { labelFor, leadStatusLabel, localizeError } from '../../i18n/labels';
import { LeadDetail } from './LeadDetail';

type Mode = 'list' | 'create' | 'detail';
const statuses = Object.keys(leadStatusLabel);

export function LeadList({ targetRecordId, onTargetHandled }: { targetRecordId?: string; onTargetHandled?: () => void }) {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [leads, setLeads] = useState<Lead[]>([]);
  const [selected, setSelected] = useState<Lead | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [includeArchived, setIncludeArchived] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [duplicateWarning, setDuplicateWarning] = useState<DuplicateWarningResult | null>(null);
  const [transferOwnerId, setTransferOwnerId] = useState('');
  const [form, setForm] = useState({ leadName: '', companyName: '', email: '', phone: '', source: '', ownerId: '', needSummary: '' });

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (!targetRecordId) return;
    void selectLead(targetRecordId).finally(() => onTargetHandled?.());
  }, [targetRecordId, onTargetHandled]);

  async function refresh(nextSearch = search, nextIncludeArchived = includeArchived) {
    const response = await listLeads(nextSearch, nextIncludeArchived);
    setLeads(response.items);
    setPage(1);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    await saveLead();
  }

  async function saveLead(proceedWarningToken?: string) {
    setError('');
    try {
      if (!proceedWarningToken) {
        const warning = await checkLeadDuplicate({ companyName: form.companyName, email: form.email, phone: form.phone });
        if (warning.result === 'PossibleDuplicate' && warning.warningToken) {
          setDuplicateWarning(warning);
          return;
        }
      }
      const created = await createLead({
        leadName: form.leadName,
        companyName: form.companyName,
        email: form.email,
        phone: form.phone,
        source: form.source,
        ownerId: form.ownerId,
        needSummary: form.needSummary,
        proceedWarningToken
      });
      setCreating(false);
      setMode('detail');
      setForm({ leadName: '', companyName: '', email: '', phone: '', source: '', ownerId: '', needSummary: '' });
      setDuplicateWarning(null);
      setSelected(created);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectLead(id: string) {
    setError('');
    setSelected(await getLead(id));
    setMode('detail');
  }

  async function updateSelected(lead: Lead) {
    setSelected(lead);
    await refresh();
  }

  async function converted(result: ConversionResult) {
    const lead = await getLead(result.leadId);
    setSelected(lead);
    setMode('detail');
    await refresh();
  }

  const rows = useMemo(() => leads.filter((lead) => {
    if (status && lead.status !== status) return false;
    if (!includeArchived && lead.archived) return false;
    return true;
  }), [leads, status, includeArchived]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((lead) => selectedIds.includes(lead.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function toggleRow(lead: Lead, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, lead.id])] : value.filter((id) => id !== lead.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((lead) => lead.id) : []);
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
      for (const lead of selectedRows) {
        if (lead.archived) continue;
        await archiveLead(lead, '批量归档线索记录');
        archived += 1;
      }
      setSelectedIds([]);
      setNotice(`已归档 ${archived} 条线索。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function transferSelected() {
    if (selectedRows.length === 0 || user?.role === 'Sales' || transferOwnerId.trim() === '') return;
    setError('');
    try {
      let transferred = 0;
      for (const lead of selectedRows) {
        await transferLeadOwner(lead, transferOwnerId.trim(), '批量转移线索负责人');
        transferred += 1;
      }
      setSelectedIds([]);
      setNotice(`已转移 ${transferred} 条线索负责人。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function archiveRow(lead: Lead) {
    if (user?.role === 'Sales' || lead.archived) return;
    if (!window.confirm(`确认归档线索 ${lead.companyName || lead.leadName || lead.id}？`)) return;
    setError('');
    try {
      await archiveLead(lead, '行操作归档线索记录');
      setNotice(`已归档线索 ${lead.companyName || lead.leadName || lead.id}。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function transferRow(lead: Lead) {
    if (user?.role === 'Sales' || lead.archived) return;
    const toOwnerId = window.prompt('请输入新负责人 ID', lead.ownerId || '');
    if (!toOwnerId?.trim()) return;
    setError('');
    try {
      const updated = await transferLeadOwner(lead, toOwnerId.trim(), '行操作转移线索负责人');
      setSelected(updated);
      setNotice(`已转移线索 ${lead.companyName || lead.leadName || lead.id} 的负责人。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('leads-selected.csv', selectedRows.map((lead) => ({
      线索: lead.companyName || lead.leadName,
      状态: labelFor(leadStatusLabel, lead.status),
      来源: lead.source,
      负责人: lead.ownerId,
      邮箱: lead.email,
      电话: lead.phone
    })));
    setNotice(`已导出 ${selectedRows.length} 条线索。`);
  }

  if (mode === 'create' || creating) {
    return (
      <FormShell
        title="新建线索"
        description="录入线索基础信息，保存前执行重复检查。"
        badge="新建"
        onCancel={() => { setMode('list'); setCreating(false); setDuplicateWarning(null); }}
        actions={<Button variant="primary" form="lead-form" type="submit">保存</Button>}
        side={<LeadFormRules />}
      >
        {error && <div role="alert" className="error">{error}</div>}
        <form id="lead-form" className="actionBand" onSubmit={submit}>
          <FormSection title="线索基本信息">
            <div className="formFields">
              <TextField label="线索名称" value={form.leadName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, leadName: event.target.value }); }} />
              <TextField label="公司名称" value={form.companyName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, companyName: event.target.value }); }} />
              <TextField label="邮箱" value={form.email} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, email: event.target.value }); }} />
              <TextField label="电话" value={form.phone} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, phone: event.target.value }); }} />
              <TextField label="来源" value={form.source} onChange={(event) => setForm({ ...form, source: event.target.value })} />
              <TextField label="负责人 ID" value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              <TextField className="full" label="需求摘要" value={form.needSummary} onChange={(event) => setForm({ ...form, needSummary: event.target.value })} />
            </div>
          </FormSection>
          <div className="saveBar">
            <Button onClick={() => { setMode('list'); setCreating(false); }}>取消</Button>
            <Button variant="primary" type="submit">保存线索</Button>
          </div>
          {duplicateWarning ? (
            <DuplicateWarning
              warning={duplicateWarning}
              onProceed={() => void saveLead(duplicateWarning.warningToken)}
              onCancel={() => setDuplicateWarning(null)}
            />
          ) : null}
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    return (
      <main className="content crudPage">
        <LeadDetail lead={selected} onUpdated={updateSelected} onConverted={converted} onError={setError} onBack={() => setMode('list')} error={error} />
      </main>
    );
  }

  return (
    <CrudListShell
      title="线索"
      description={`${scopeLabel} · 共 ${rows.length} 条 · 默认按更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={
        <>
          <Button onClick={() => void refresh(search, includeArchived)}>
            <RotateCcw size={16} aria-hidden="true" />
            刷新
          </Button>
          <Button variant="primary" onClick={() => { setDuplicateWarning(null); setCreating(true); setMode('create'); }}>
            <Plus size={16} aria-hidden="true" />
            新建线索
          </Button>
        </>
      }
      toolbar={
        <Toolbar
          searchValue={search}
          onSearchChange={setSearch}
          searchPlaceholder="搜索线索名、公司名"
          filters={
            <>
              <label className="inlineCheckbox">
                <span>状态</span>
                <select value={status} onChange={(event) => setStatus(event.target.value)}>
                  <option value="">全部状态</option>
                  {statuses.map((statusOption) => (
                    <option value={statusOption} key={statusOption}>{labelFor(leadStatusLabel, statusOption)}</option>
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
          <span className="chip">状态：{status ? labelFor(leadStatusLabel, status) : '全部状态'}</span>
          <span className="chip">归档：{includeArchived ? '包含已归档' : '仅活动记录'}</span>
        </ActiveFilterSummary>
      }
      bulkBar={
        <BulkActionBar>
          <div className="bulkSummary">
            <span className="bulkCount">已选择 {selectedRows.length} 条</span>
            <span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '线索批量转移/归档通过现有单条接口逐条执行。')}</span>
          </div>
          <div className="bulkActions">
            {user?.role !== 'Sales' ? (
              <>
                <input className="bulkOwnerInput" aria-label="批量转移负责人 ID" placeholder="新负责人 ID" value={transferOwnerId} onChange={(event) => setTransferOwnerId(event.target.value)} />
                <Button className="bulkButton" disabled={selectedRows.length === 0 || transferOwnerId.trim() === ''} onClick={() => void transferSelected()}>批量转移负责人</Button>
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
          caption="线索结果表"
          rows={slice.items}
          rowKey={(lead) => lead.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(lead) => selected?.id === lead.id ? 'selected' : undefined}
          onRowClick={(lead) => void selectLead(lead.id)}
          getRowAriaLabel={(lead) => `打开线索 ${lead.companyName || lead.leadName || lead.id}`}
          empty="没有符合当前筛选条件的线索。"
          columns={[
            {
              key: 'lead',
              header: '线索',
              width: '250px',
              render: (lead) => <RecordIdentity icon={<ListChecks size={17} aria-hidden="true" />} title={lead.companyName || lead.leadName || lead.id} subtitle={`${lead.archived ? '已归档 · ' : ''}${lead.source || '未填写来源'}`} />
            },
            { key: 'contact', header: '联系方式', render: (lead) => lead.email || lead.phone || '未填写' },
            { key: 'status', header: '状态', render: (lead) => <StatusPill tone={leadTone(lead.status)}>{labelFor(leadStatusLabel, lead.status)}</StatusPill> },
            { key: 'owner', header: '负责人', render: (lead) => lead.ownerId || '未分配' },
            { key: 'updated', header: '更新时间', align: 'right', sortable: true, sortDirection: 'desc', render: (lead) => formatDate(lead.updatedAt) }
          ]}
          actions={(lead) => (
            <div className="rowActions">
              <button className="rowAction" type="button" aria-label={`查看 ${lead.companyName || lead.leadName || lead.id}`} onClick={() => void selectLead(lead.id)}>
                <ArrowRight size={16} aria-hidden="true" />
              </button>
              <ActionMenu
                label={`打开 ${lead.companyName || lead.leadName || lead.id} 的行操作菜单`}
                items={[
                  { label: '查看', onSelect: () => void selectLead(lead.id) },
                  {
                    label: '转移负责人',
                    onSelect: () => void transferRow(lead),
                    disabled: user?.role === 'Sales' || Boolean(lead.archived),
                    reason: user?.role === 'Sales' ? '销售角色不能转移负责人。' : lead.archived ? '已归档记录只读。' : undefined
                  },
                  {
                    label: '归档',
                    onSelect: () => void archiveRow(lead),
                    disabled: user?.role === 'Sales' || Boolean(lead.archived),
                    reason: user?.role === 'Sales' ? '销售角色不能归档线索。' : lead.archived ? '已归档。' : undefined,
                    tone: 'danger'
                  },
                  {
                    label: '转为商机',
                    onSelect: () => void selectLead(lead.id),
                    disabled: lead.status !== 'Valid' || lead.ownerId === '' || Boolean(lead.archived),
                    reason: lead.status !== 'Valid' ? '仅有效线索可转换。' : lead.ownerId === '' ? '未分配线索不能转换。' : lead.archived ? '已归档记录只读。' : '在详情页填写预计金额和日期。'
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

function LeadFormRules() {
  return (
    <>
      <div className="sideCard">
        <h3>字段状态</h3>
        <div className="rule"><span>1</span><p>公司名称或线索名称至少填写一项。</p></div>
        <div className="rule"><span>2</span><p>保存前执行重复检查；可能重复时展示确认令牌。</p></div>
        <div className="rule"><span>3</span><p>未分配线索不能确认或转换。</p></div>
      </div>
    </>
  );
}

function leadTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Valid' || status === 'Converted To Opportunity') return 'success';
  if (status === 'Invalid') return 'danger';
  if (status === 'Pending Qualification') return 'warning';
  return 'neutral';
}

function formatDate(value: string) {
  if (!value) return '未更新';
  return value.length > 10 ? value.slice(0, 10) : value;
}
