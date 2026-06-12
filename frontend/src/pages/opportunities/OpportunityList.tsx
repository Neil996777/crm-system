import { FormEvent, useEffect, useMemo, useState } from 'react';
import { BriefcaseBusiness, Plus, RotateCcw } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Opportunity, archiveOpportunity, changeOpportunityStage, createOpportunity, getOpportunity, listOpportunities, updateOpportunity } from '../../api/opportunities';
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
import { labelFor, localizeError, opportunityStageLabel } from '../../i18n/labels';
import { OpportunityDetail } from './OpportunityDetail';

const stages = Object.keys(opportunityStageLabel);
const editableStages = ['New Opportunity', 'Needs Confirmed', 'Quote', 'Contract Negotiation'];
type Mode = 'list' | 'create' | 'edit' | 'detail';

export function OpportunityList({ targetRecordId, onTargetHandled }: { targetRecordId?: string; onTargetHandled?: () => void }) {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [opportunities, setOpportunities] = useState<Opportunity[]>([]);
  const [selected, setSelected] = useState<Opportunity | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [stage, setStage] = useState('');
  const [statusFilter, setStatusFilter] = useState('active');
  const [amountBand, setAmountBand] = useState('');
  const [includeArchived, setIncludeArchived] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [form, setForm] = useState({
    title: '',
    customerId: '',
    ownerId: '',
    stage: 'New Opportunity',
    expectedAmount: '',
    expectedCloseDate: ''
  });
  const ownerLocked = user?.role === 'Sales';

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (!targetRecordId) return;
    void selectOpportunity(targetRecordId).finally(() => onTargetHandled?.());
  }, [targetRecordId, onTargetHandled]);

  useEffect(() => {
    if (ownerLocked && user) {
      setForm((value) => ({ ...value, ownerId: user.id }));
    }
  }, [ownerLocked, user]);

  function blankForm() {
    return { title: '', customerId: '', ownerId: ownerLocked && user ? user.id : '', stage: 'New Opportunity', expectedAmount: '', expectedCloseDate: '' };
  }

  async function refresh(nextSearch = search, nextStage = stage, nextIncludeArchived = includeArchived) {
    const response = await listOpportunities(nextSearch, nextStage, nextIncludeArchived);
    setOpportunities(response.items);
    setPage(1);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const ownerId = ownerLocked && user ? user.id : form.ownerId;
      const saved = mode === 'edit' && selected
        ? await updateOpportunity(selected.id, { ...form, ownerId, expectedVersion: selected.version })
        : await createOpportunity({ ...form, ownerId });
      setSelected(saved);
      setMode('detail');
      setForm(blankForm());
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectOpportunity(id: string) {
    setError('');
    setSelected(await getOpportunity(id));
    setMode('detail');
  }

  async function updateSelected(opportunity: Opportunity) {
    setSelected(opportunity);
    await refresh();
  }

  function startCreate() {
    setError('');
    setSelected(null);
    setForm(blankForm());
    setMode('create');
  }

  function startEdit(opportunity: Opportunity) {
    setError('');
    setSelected(opportunity);
    setForm({
      title: opportunity.title,
      customerId: opportunity.customerId,
      ownerId: opportunity.ownerId,
      stage: editableStages.includes(opportunity.stage) ? opportunity.stage : 'New Opportunity',
      expectedAmount: opportunity.expectedAmount,
      expectedCloseDate: opportunity.expectedCloseDate
    });
    setMode('edit');
  }

  function filteredRows() {
    return opportunities.filter((opportunity) => {
      if (statusFilter === 'active' && (opportunity.stage === 'Won' || opportunity.stage === 'Lost')) return false;
      if (statusFilter === 'terminal' && opportunity.stage !== 'Won' && opportunity.stage !== 'Lost') return false;
      if (!includeArchived && opportunity.archived) return false;
      if (amountBand) {
        const amount = Number(opportunity.expectedAmount);
        if (amountBand === 'lt100k' && !(amount < 100000)) return false;
        if (amountBand === '100k-2m' && !(amount >= 100000 && amount <= 2000000)) return false;
        if (amountBand === 'gt2m' && !(amount > 2000000)) return false;
      }
      return true;
    });
  }

  const rows = useMemo(filteredRows, [opportunities, statusFilter, amountBand]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((opportunity) => selectedIds.includes(opportunity.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function clearFilters() {
    setSearch('');
    setStage('');
    setStatusFilter('active');
    setAmountBand('');
    setIncludeArchived(false);
    setSelectedIds([]);
    void refresh('', '', false);
  }

  function toggleRow(opportunity: Opportunity, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, opportunity.id])] : value.filter((id) => id !== opportunity.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((opportunity) => opportunity.id) : []);
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('opportunities-selected.csv', selectedRows.map((opportunity) => ({
      商机: opportunity.title || opportunity.id,
      客户: opportunity.customerId,
      阶段: labelFor(opportunityStageLabel, opportunity.stage),
      负责人: opportunity.ownerId,
      金额: opportunity.expectedAmount,
      预计签约: opportunity.expectedCloseDate,
      更新时间: opportunity.updatedAt
    })));
    setNotice(`已导出 ${selectedRows.length} 条商机。`);
  }

  async function archiveSelected() {
    if (selectedRows.length === 0 || user?.role === 'Sales') return;
    setError('');
    setNotice('');
    try {
      let archived = 0;
      for (const opportunity of selectedRows) {
        if (opportunity.archived) continue;
        await archiveOpportunity(opportunity.id, opportunity.version, '批量归档商机记录');
        archived += 1;
      }
      setSelectedIds([]);
      setNotice(`已归档 ${archived} 条商机。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function archiveRow(opportunity: Opportunity) {
    if (user?.role === 'Sales' || opportunity.archived || isTerminal(opportunity)) return;
    if (!window.confirm(`确认归档商机 ${opportunity.title || opportunity.id}？`)) return;
    setError('');
    setNotice('');
    try {
      await archiveOpportunity(opportunity.id, opportunity.version, '行操作归档商机记录');
      setNotice(`已归档商机 ${opportunity.title || opportunity.id}。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function advanceRow(opportunity: Opportunity) {
    const nextStage = nextEditableStage(opportunity.stage);
    if (!nextStage) return;
    if (!window.confirm(`确认推进到${labelFor(opportunityStageLabel, nextStage)}？`)) return;
    setError('');
    try {
      const updated = await changeOpportunityStage(opportunity.id, opportunity.version, nextStage);
      setSelected(updated);
      setNotice(`已推进 ${opportunity.title || opportunity.id} 到${labelFor(opportunityStageLabel, nextStage)}。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  if (mode === 'create' || mode === 'edit') {
    const editing = mode === 'edit';
    return (
      <FormShell
        title={editing ? '编辑商机' : '新建商机'}
        description={ownerLocked ? '销售视图 · 负责人锁定为本人 · 阶段排除赢单/丢单终态' : '团队视图 · 阶段排除赢单/丢单终态'}
        badge={editing ? `v${selected?.version ?? ''}` : '新建'}
        onCancel={() => editing && selected ? setMode('detail') : setMode('list')}
        actions={<Button variant="primary" form="opportunity-form" type="submit">保存</Button>}
        side={<OpportunityFormRules ownerLocked={ownerLocked} />}
      >
        {error && <div role="alert" className="error">{error}</div>}
        <form id="opportunity-form" className="actionBand" onSubmit={submit}>
          <FormSection title="商机基本信息">
            <div className="formFields">
              <TextField label="商机名 *" aria-label="标题" value={form.title} onChange={(event) => setForm({ ...form, title: event.target.value })} />
              <TextField label="客户 *" aria-label="客户 ID" value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              <TextField
                label="负责人 *"
                aria-label="负责人 ID"
                value={ownerLocked && user ? user.id : form.ownerId}
                onChange={(event) => setForm({ ...form, ownerId: event.target.value })}
                disabled={ownerLocked}
                hint={ownerLocked ? '销售账号新建商机时负责人固定为本人。' : '填写负责人 ID。'}
              />
              <TextField label="预计金额 *" aria-label="预计金额" value={form.expectedAmount} onChange={(event) => setForm({ ...form, expectedAmount: event.target.value })} />
              <TextField label="预计签约日期 *" aria-label="预计关闭日期" type="date" value={form.expectedCloseDate} onChange={(event) => setForm({ ...form, expectedCloseDate: event.target.value })} />
              <fieldset className="choiceField" aria-label="阶段">
                <legend>阶段 *</legend>
                <div className="choiceGroup" role="presentation">
                  {editableStages.map((stageOption) => (
                    <button
                      className={`chip choiceButton ${form.stage === stageOption ? 'selected' : ''}`}
                      type="button"
                      aria-pressed={form.stage === stageOption}
                      key={stageOption}
                      onClick={() => setForm({ ...form, stage: stageOption })}
                    >
                      {labelFor(opportunityStageLabel, stageOption)}
                    </button>
                  ))}
                </div>
                <span className="panelMeta">新建默认新商机；赢单/丢单只能通过关闭动作生成。</span>
              </fieldset>
            </div>
          </FormSection>
          <div className="saveBar">
            <Button onClick={() => editing && selected ? setMode('detail') : setMode('list')}>取消</Button>
            <Button variant="primary" type="submit">{editing ? '保存编辑' : '保存商机'}</Button>
          </div>
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    return <OpportunityDetail opportunity={selected} onUpdated={updateSelected} onError={setError} onBack={() => setMode('list')} onEdit={startEdit} error={error} />;
  }

  return (
    <CrudListShell
      title="商机"
      description={`${scopeLabel} · 共 ${rows.length} 条 · 默认按更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={
        <>
          <Button onClick={() => void refresh(search, stage)}>
            <RotateCcw size={16} aria-hidden="true" />
            刷新
          </Button>
          <Button variant="primary" onClick={startCreate}>
            <Plus size={16} aria-hidden="true" />
            新建商机
          </Button>
        </>
      }
      toolbar={
        <Toolbar
          searchValue={search}
          onSearchChange={setSearch}
          searchPlaceholder="搜索商机名、客户名"
          filters={
            <>
              <label className="inlineCheckbox">
                <span>阶段</span>
                <select value={stage} onChange={(event) => setStage(event.target.value)}>
                  <option value="">全部阶段</option>
                  {stages.map((stageOption) => (
                    <option value={stageOption} key={stageOption}>
                      {labelFor(opportunityStageLabel, stageOption)}
                    </option>
                  ))}
                </select>
              </label>
              <label className="inlineCheckbox">
                <span>金额</span>
                <select value={amountBand} onChange={(event) => setAmountBand(event.target.value)}>
                  <option value="">全部金额</option>
                  <option value="lt100k">小于 100K</option>
                  <option value="100k-2m">100K-2M</option>
                  <option value="gt2m">大于 2M</option>
                </select>
              </label>
              <label className="inlineCheckbox">
                <span>状态</span>
                <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value)}>
                  <option value="">全部状态</option>
                  <option value="active">进行中</option>
                  <option value="terminal">已关闭</option>
                </select>
              </label>
              <label className="inlineCheckbox">
                <input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />
                包含已归档
              </label>
            </>
          }
          actions={<Button onClick={() => void refresh(search, stage, includeArchived)}>应用筛选</Button>}
        />
      }
      activeFilters={
        <ActiveFilterSummary onClear={clearFilters}>
          <span className="chip">负责人：{scopeLabel}</span>
          <span className="chip">阶段：{stage ? labelFor(opportunityStageLabel, stage) : '全部阶段'}</span>
          <span className="chip">金额：{amountBand ? amountLabel(amountBand) : '全部金额'}</span>
          <span className="chip">状态：{statusFilter ? statusLabel(statusFilter) : '全部状态'}</span>
          <span className="chip">归档：{includeArchived ? '包含已归档' : '仅活动记录'}</span>
        </ActiveFilterSummary>
      }
      bulkBar={
        <BulkActionBar>
          <div className="bulkSummary">
            <span className="bulkCount">已选择 {selectedRows.length} 条</span>
            <span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '批量转移/归档需要后端能力；当前按 A3 显示不可提交说明。')}</span>
          </div>
          <div className="bulkActions">
            {user?.role !== 'Sales' ? (
              <>
                <Button className="bulkButton" disabled title="商机当前无负责人转移接口；按 A3 禁用。">批量转移负责人</Button>
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
          caption="商机结果表"
          rows={slice.items}
          rowKey={(opportunity) => opportunity.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(opportunity) => selected?.id === opportunity.id ? 'selected' : undefined}
          onRowClick={(opportunity) => void selectOpportunity(opportunity.id)}
          getRowAriaLabel={(opportunity) => `打开商机 ${opportunity.title || opportunity.id}`}
          empty="没有符合当前筛选条件的商机。"
          columns={[
            {
              key: 'title',
              header: '商机名',
              width: '250px',
              render: (opportunity) => (
                <RecordIdentity
                  icon={<BriefcaseBusiness size={17} aria-hidden="true" />}
                  title={opportunity.title || opportunity.id}
                  titleAriaLabel={`打开商机 ${opportunity.title || opportunity.id}`}
                  subtitle={`${opportunity.archived ? '已归档 · ' : ''}更新于 ${formatDate(opportunity.updatedAt)}`}
                  tone={identityTone(opportunity.stage)}
                  onTitleClick={() => void selectOpportunity(opportunity.id)}
                />
              )
            },
            { key: 'customer', header: '客户', render: (opportunity) => opportunity.customerId },
            {
              key: 'stage',
              header: '阶段',
              render: (opportunity) => <StatusPill tone={stageTone(opportunity.stage)}>{labelFor(opportunityStageLabel, opportunity.stage)}</StatusPill>
            },
            { key: 'owner', header: '负责人', render: (opportunity) => opportunity.ownerId },
            { key: 'amount', header: '金额', align: 'right', render: (opportunity) => <span className="amountText">{money(opportunity.expectedAmount)}</span> },
            { key: 'close', header: '预计签约', align: 'right', render: (opportunity) => opportunity.expectedCloseDate || '未填写' },
            { key: 'updated', header: '更新时间', align: 'right', sortable: true, sortDirection: 'desc', render: (opportunity) => formatDate(opportunity.updatedAt) }
          ]}
          actions={(opportunity) => (
            <div className="rowActions">
              <ActionMenu
                label={`打开 ${opportunity.title || opportunity.id} 的行操作菜单`}
                items={[
                  { label: '查看', onSelect: () => void selectOpportunity(opportunity.id) },
                  {
                    label: '编辑',
                    onSelect: () => startEdit(opportunity),
                    disabled: user?.role === 'Sales' || isTerminal(opportunity) || Boolean(opportunity.archived),
                    reason: user?.role === 'Sales' ? '销售角色不能编辑负责人范围外字段。' : isTerminal(opportunity) ? '终态商机只读。' : opportunity.archived ? '已归档记录只读。' : undefined
                  },
                  {
                    label: nextEditableStage(opportunity.stage) ? `推进到${labelFor(opportunityStageLabel, nextEditableStage(opportunity.stage) ?? '')}` : '推进阶段',
                    onSelect: () => void advanceRow(opportunity),
                    disabled: !nextEditableStage(opportunity.stage) || Boolean(opportunity.archived),
                    reason: !nextEditableStage(opportunity.stage) ? '终态商机无法推进阶段。' : opportunity.archived ? '已归档记录只读。' : undefined
                  },
                  {
                    label: '转移负责人',
                    onSelect: () => void selectOpportunity(opportunity.id),
                    disabled: user?.role === 'Sales' || isTerminal(opportunity) || Boolean(opportunity.archived),
                    reason: user?.role === 'Sales' ? '销售角色不能转移负责人。' : isTerminal(opportunity) ? '终态商机只读。' : opportunity.archived ? '已归档记录只读。' : '在详情页填写新负责人 ID。'
                  },
                  {
                    label: '归档',
                    onSelect: () => void archiveRow(opportunity),
                    disabled: user?.role === 'Sales' || isTerminal(opportunity) || Boolean(opportunity.archived),
                    reason: user?.role === 'Sales' ? '销售角色不能归档商机。' : isTerminal(opportunity) ? '终态商机只读。' : opportunity.archived ? '已归档。' : undefined,
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

function isTerminal(opportunity: Opportunity) {
  return opportunity.stage === 'Won' || opportunity.stage === 'Lost';
}

function nextEditableStage(stage: string) {
  if (stage === 'New Opportunity') return 'Needs Confirmed';
  if (stage === 'Needs Confirmed') return 'Quote';
  if (stage === 'Quote') return 'Contract Negotiation';
  return '';
}

function OpportunityFormRules({ ownerLocked }: { ownerLocked: boolean }) {
  return (
    <>
      <div className="sideCard">
        <h3>字段状态</h3>
        <div className="rule"><span>1</span><p>{ownerLocked ? '负责人为只读预填：当前销售本人。' : '负责人由有权限的角色填写。'}</p></div>
        <div className="rule"><span>2</span><p>阶段下拉只列非终态四项；赢单/丢单不在表单直选。</p></div>
        <div className="rule"><span>3</span><p>保存失败以内联错误保留字段输入，不清空负责人锁定态。</p></div>
      </div>
      <div className="sideCard">
        <h3>并发说明</h3>
        <p className="helper">编辑提交携带当前版本；冲突时提示“记录在你打开后已被他人修改，请刷新重试”。</p>
      </div>
    </>
  );
}

function amountLabel(value: string) {
  if (value === 'lt100k') return '小于 100K';
  if (value === '100k-2m') return '100K-2M';
  if (value === 'gt2m') return '大于 2M';
  return '全部金额';
}

function statusLabel(value: string) {
  if (value === 'active') return '进行中';
  if (value === 'terminal') return '已关闭';
  return '全部状态';
}

function stageTone(stage: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (stage === 'Won') return 'success';
  if (stage === 'Lost') return 'danger';
  if (stage === 'Contract Negotiation') return 'warning';
  if (stage === 'Quote') return 'primary';
  return 'neutral';
}

function identityTone(stage: string): 'sky' | 'mint' | 'peach' | 'purple' | 'primary' {
  if (stage === 'Won') return 'mint';
  if (stage === 'Lost') return 'purple';
  if (stage === 'Contract Negotiation') return 'peach';
  if (stage === 'Quote') return 'primary';
  return 'sky';
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
