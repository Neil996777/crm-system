import { FormEvent, useEffect, useMemo, useState } from 'react';
import { FileText, Plus, RotateCcw } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Quote, changeQuoteStatus, createQuote, getQuote, listQuotes } from '../../api/quotes';
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
import { ActionMenu, BulkActionBar, Button, DataTable, TextField, Toolbar } from '../../components/ui';
import { useSession } from '../../auth/SessionProvider';
import { labelFor, localizeError, quoteStatusLabel } from '../../i18n/labels';
import { QuoteDetail } from './QuoteDetail';

type Mode = 'list' | 'create' | 'detail';
const statuses = Object.keys(quoteStatusLabel);

export function QuoteList({ targetRecordId, onTargetHandled }: { targetRecordId?: string; onTargetHandled?: () => void }) {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [quotes, setQuotes] = useState<Quote[]>([]);
  const [selected, setSelected] = useState<Quote | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [loading, setLoading] = useState(true);
  const [selectingId, setSelectingId] = useState('');
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [form, setForm] = useState({ opportunityId: '', customerId: '', amount: '', validityEnd: '', ownerId: '' });

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (!targetRecordId) return;
    void selectQuote(targetRecordId).finally(() => onTargetHandled?.());
  }, [targetRecordId, onTargetHandled]);

  async function refresh(nextSearch = search) {
    setLoading(true);
    setError('');
    try {
      const response = await listQuotes(nextSearch);
      setQuotes(response.items);
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
      const created = await createQuote(form);
      setSelected(created);
      setMode('detail');
      setForm({ opportunityId: '', customerId: '', amount: '', validityEnd: '', ownerId: '' });
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectQuote(id: string) {
    setError('');
    setSelectingId(id);
    try {
      setSelected(await getQuote(id));
      setMode('detail');
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    } finally {
      setSelectingId('');
    }
  }

  async function updateSelected(quote: Quote) {
    setSelected(quote);
    await refresh();
  }

  const rows = useMemo(() => quotes.filter((quote) => !status || quote.status === status), [quotes, status]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((quote) => selectedIds.includes(quote.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function toggleRow(quote: Quote, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, quote.id])] : value.filter((id) => id !== quote.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((quote) => quote.id) : []);
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('quotes-selected.csv', selectedRows.map((quote) => ({
      报价: quote.id,
      商机: quote.opportunityId,
      客户: quote.customerId,
      状态: labelFor(quoteStatusLabel, quote.status),
      金额: quote.amount,
      有效期: quote.validityEnd
    })));
    setNotice(`已导出 ${selectedRows.length} 条报价。`);
  }

  async function changeRowStatus(quote: Quote, toStatus: string) {
    if (!canChangeQuoteTo(quote, toStatus)) return;
    if (!window.confirm(`确认将报价 ${quote.id} 标记为${labelFor(quoteStatusLabel, toStatus)}？`)) return;
    setError('');
    try {
      const updated = await changeQuoteStatus(quote.id, quote.version, toStatus);
      setSelected(updated);
      setNotice(`报价 ${quote.id} 已标记为${labelFor(quoteStatusLabel, toStatus)}。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  if (mode === 'create') {
    return (
      <FormShell
        title="新建报价"
        description="每个商机最多一个报价，保存后初始状态为草稿。"
        badge="新建"
        onCancel={() => setMode('list')}
        actions={<Button variant="primary" form="quote-form" type="submit">保存</Button>}
        side={<QuoteFormRules />}
      >
        {error && <div role="alert" className="error">{error}</div>}
        <form id="quote-form" className="actionBand" onSubmit={submit}>
          <FormSection title="报价基本信息">
            <div className="formFields">
              <TextField label="商机 ID" value={form.opportunityId} onChange={(event) => setForm({ ...form, opportunityId: event.target.value })} />
              <TextField label="客户 ID" value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              <TextField label="金额" value={form.amount} onChange={(event) => setForm({ ...form, amount: event.target.value })} />
              <TextField label="有效期截止日" type="date" value={form.validityEnd} onChange={(event) => setForm({ ...form, validityEnd: event.target.value })} />
              <TextField label="负责人 ID" value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
            </div>
          </FormSection>
          <div className="saveBar">
            <Button onClick={() => setMode('list')}>取消</Button>
            <Button variant="primary" type="submit">保存报价</Button>
          </div>
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    return <QuoteDetail quote={selected} onUpdated={updateSelected} onError={setError} onBack={() => setMode('list')} error={error} />;
  }

  return (
    <CrudListShell
      title="报价"
      description={`${scopeLabel} · 共 ${rows.length} 条 · 默认按更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={<><Button disabled={loading} aria-busy={loading || undefined} onClick={() => void refresh(search)}><RotateCcw size={16} aria-hidden="true" />{loading ? '刷新中' : '刷新'}</Button><Button variant="primary" onClick={() => setMode('create')}><Plus size={16} aria-hidden="true" />新建报价</Button></>}
      toolbar={
        <Toolbar
          searchValue={search}
          onSearchChange={setSearch}
          searchPlaceholder="搜索报价、商机或客户"
          filters={<label className="inlineCheckbox"><span>状态</span><select value={status} onChange={(event) => setStatus(event.target.value)}><option value="">全部状态</option>{statuses.map((item) => <option value={item} key={item}>{labelFor(quoteStatusLabel, item)}</option>)}</select></label>}
          actions={<Button disabled={loading} aria-busy={loading || undefined} onClick={() => void refresh(search)}>应用筛选</Button>}
        />
      }
      feedback={<ListAsyncFeedback error={error} loading={loading} selecting={Boolean(selectingId)} />}
      activeFilters={<ActiveFilterSummary onClear={() => { setSearch(''); setStatus(''); setSelectedIds([]); void refresh(''); }}><span className="chip">负责人：{scopeLabel}</span><span className="chip">状态：{status ? labelFor(quoteStatusLabel, status) : '全部状态'}</span></ActiveFilterSummary>}
      bulkBar={
        <BulkActionBar>
          <div className="bulkSummary"><span className="bulkCount">已选择 {selectedRows.length} 条</span><span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '报价批量操作保留导出与清除选择。')}</span></div>
          <div className="bulkActions">
            <ExportSelectedButton disabled={selectedRows.length === 0} onExport={exportSelected} />
            <Button className="bulkButton" variant="primary" disabled={selectedRows.length === 0} onClick={() => setSelectedIds([])}>清除选择</Button>
          </div>
        </BulkActionBar>
      }
      table={
        loading ? <ListTableLoading label="正在加载报价列表..." /> : <DataTable
          caption="报价结果表"
          rows={slice.items}
          rowKey={(quote) => quote.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(quote) => selected?.id === quote.id ? 'selected' : undefined}
          onRowClick={(quote) => void selectQuote(quote.id)}
          getRowAriaLabel={(quote) => `打开报价 ${quote.id}`}
          empty="没有符合当前筛选条件的报价。"
          columns={[
            {
              key: 'quote',
              header: '报价',
              width: '220px',
              render: (quote) => (
                <RecordIdentity
                  icon={<FileText size={17} aria-hidden="true" />}
                  title={quote.opportunityId}
                  titleAriaLabel={`打开报价 ${quote.id}`}
                  titleBusy={selectingId === quote.id}
                  subtitle={quote.id}
                  tone="purple"
                  onTitleClick={() => void selectQuote(quote.id)}
                />
              )
            },
            { key: 'customer', header: '客户', render: (quote) => quote.customerId },
            { key: 'status', header: '状态', render: (quote) => <StatusPill tone={quoteTone(quote.status)}>{labelFor(quoteStatusLabel, quote.status)}</StatusPill> },
            { key: 'amount', header: '金额', align: 'right', render: (quote) => <span className="amountText">{money(quote.amount)}</span> },
            { key: 'validity', header: '有效期', align: 'right', render: (quote) => quote.validityEnd },
            { key: 'updated', header: '更新时间', align: 'right', sortable: true, sortDirection: 'desc', render: (quote) => formatDate(quote.updatedAt) }
          ]}
          actions={(quote) => (
            <div className="rowActions">
              <ActionMenu
                label={`打开报价 ${quote.id} 的行操作菜单`}
                items={[
                  { label: '查看', onSelect: () => void selectQuote(quote.id) },
                  {
                    label: '发送',
                    onSelect: () => void changeRowStatus(quote, 'Sent'),
                    disabled: !canChangeQuoteTo(quote, 'Sent'),
                    reason: quote.status !== 'Draft' ? '仅草稿报价可发送。' : undefined
                  },
                  {
                    label: '接受',
                    onSelect: () => void changeRowStatus(quote, 'Accepted'),
                    disabled: !canChangeQuoteTo(quote, 'Accepted'),
                    reason: quote.status !== 'Sent' ? '仅已发送报价可接受。' : undefined
                  },
                  {
                    label: '拒绝',
                    onSelect: () => void changeRowStatus(quote, 'Rejected'),
                    disabled: !canChangeQuoteTo(quote, 'Rejected'),
                    reason: quote.status !== 'Sent' ? '仅已发送报价可拒绝。' : undefined,
                    tone: 'danger'
                  },
                  {
                    label: '标记过期',
                    onSelect: () => void changeRowStatus(quote, 'Expired'),
                    disabled: !canChangeQuoteTo(quote, 'Expired'),
                    reason: quote.status !== 'Draft' && quote.status !== 'Sent' ? '仅草稿或已发送报价可标记过期。' : undefined,
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

function canChangeQuoteTo(quote: Quote, toStatus: string) {
  if (toStatus === 'Sent') return quote.status === 'Draft';
  if (toStatus === 'Accepted' || toStatus === 'Rejected') return quote.status === 'Sent';
  if (toStatus === 'Expired') return quote.status === 'Draft' || quote.status === 'Sent';
  return false;
}

function QuoteFormRules() {
  return <div className="sideCard"><h3>字段状态</h3><div className="rule"><span>1</span><p>报价保存时状态固定为草稿。</p></div><div className="rule"><span>2</span><p>报价接受后才可进入合同创建流程。</p></div></div>;
}

function quoteTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Accepted') return 'success';
  if (status === 'Rejected' || status === 'Expired') return 'danger';
  if (status === 'Sent') return 'primary';
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
