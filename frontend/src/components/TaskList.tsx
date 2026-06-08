import { FormEvent, useEffect, useMemo, useState } from 'react';
import { ArrowRight, CalendarDays, CheckCircle2, ClipboardList, MoreHorizontal, Plus, RotateCcw } from 'lucide-react';
import { ApiError } from '../api/client';
import { WorkTask, changeTaskStatus, createTask, listTasks } from '../api/work';
import { useSession } from '../auth/SessionProvider';
import {
  ActiveFilterSummary,
  CrudListShell,
  CrudPagination,
  DEFAULT_PAGE_SIZE,
  DetailHero,
  DetailStat,
  ExportSelectedButton,
  FormSection,
  FormShell,
  RecordIdentity,
  StatusPill,
  exportRows,
  paginate
} from './CrudScaffold';
import { BulkActionBar, Button, DataTable, Panel, SelectField, TextField, Toolbar } from './ui';
import { labelFor, localizeError, objectTypeLabel, taskStatusLabel } from '../i18n/labels';

type Mode = 'list' | 'create' | 'detail';

const statuses = Object.keys(taskStatusLabel);
const relatedTypes = ['Opportunity', 'Lead', 'Account', 'Contact', 'Quote', 'Contract', 'Payment'];

export function TaskList() {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [tasks, setTasks] = useState<WorkTask[]>([]);
  const [selected, setSelected] = useState<WorkTask | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [activeOnly, setActiveOnly] = useState(false);
  const [asOfToday, setAsOfToday] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [form, setForm] = useState({
    relatedType: 'Opportunity',
    relatedId: '',
    title: '',
    dueDate: '',
    ownerId: user?.id ?? ''
  });

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (user?.role === 'Sales') {
      setForm((current) => ({ ...current, ownerId: user.id }));
    }
  }, [user?.id, user?.role]);

  async function refresh(nextActiveOnly = activeOnly, nextAsOfToday = asOfToday) {
    try {
      const response = await listTasks({
        activeOnly: nextActiveOnly,
        businessDate: nextAsOfToday ? today() : undefined
      });
      setTasks(response.items);
      setPage(1);
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createTask({
        relatedType: form.relatedType,
        relatedId: form.relatedId,
        ownerId: user?.role === 'Sales' && user.id ? user.id : form.ownerId,
        title: form.title,
        dueDate: form.dueDate
      });
      setSelected(created);
      setForm({ relatedType: 'Opportunity', relatedId: '', title: '', dueDate: '', ownerId: user?.id ?? '' });
      setMode('detail');
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function complete(task = selected) {
    if (!task || task.status === 'Completed' || task.status === 'Cancelled') return;
    setError('');
    try {
      const updated = await changeTaskStatus(task.id, 'Completed', task.version);
      if (selected?.id === task.id) setSelected(updated);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function completeSelected() {
    if (selectedRows.length === 0) return;
    setError('');
    try {
      let completed = 0;
      for (const task of selectedRows) {
        if (task.status === 'Completed' || task.status === 'Cancelled') continue;
        await changeTaskStatus(task.id, 'Completed', task.version);
        completed += 1;
      }
      setSelectedIds([]);
      setNotice(`已完成 ${completed} 条任务。`);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  const rows = useMemo(() => tasks.filter((task) => {
    if (status && task.status !== status) return false;
    const needle = search.trim().toLowerCase();
    if (!needle) return true;
    return [task.title, task.relatedId, task.relatedType, labelFor(objectTypeLabel, task.relatedType), task.ownerId]
      .some((value) => value.toLowerCase().includes(needle));
  }), [tasks, search, status]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((task) => selectedIds.includes(task.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function clearFilters() {
    setSearch('');
    setStatus('');
    setActiveOnly(false);
    setAsOfToday(false);
    setSelectedIds([]);
    void refresh(false, false);
  }

  function toggleRow(task: WorkTask, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, task.id])] : value.filter((id) => id !== task.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((task) => task.id) : []);
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('tasks-selected.csv', selectedRows.map((task) => ({
      任务: task.title,
      状态: labelFor(taskStatusLabel, task.status),
      关联记录: `${labelFor(objectTypeLabel, task.relatedType)} ${task.relatedId}`,
      负责人: task.ownerId,
      到期日: task.dueDate
    })));
    setNotice(`已导出 ${selectedRows.length} 条任务。`);
  }

  if (mode === 'create') {
    return (
      <FormShell
        title="新建任务"
        description="创建跟进任务并绑定已有业务记录。"
        badge="新建"
        onCancel={() => setMode('list')}
        actions={<Button variant="primary" form="task-form" type="submit">保存</Button>}
        side={<TaskFormRules salesLocked={user?.role === 'Sales'} />}
      >
        {error && <div role="alert" className="error">{error}</div>}
        <form id="task-form" className="actionBand" onSubmit={submit}>
          <FormSection title="任务基本信息">
            <div className="formFields">
              <SelectField label="关联类型" value={form.relatedType} onChange={(event) => setForm({ ...form, relatedType: event.currentTarget.value })}>
                {relatedTypes.map((type) => <option value={type} key={type}>{labelFor(objectTypeLabel, type)}</option>)}
              </SelectField>
              <TextField label="关联记录 ID" value={form.relatedId} onChange={(event) => setForm({ ...form, relatedId: event.currentTarget.value })} />
              <TextField className="full" label="任务标题" value={form.title} onChange={(event) => setForm({ ...form, title: event.currentTarget.value })} />
              <TextField label="任务到期日" type="date" value={form.dueDate} onChange={(event) => setForm({ ...form, dueDate: event.currentTarget.value })} />
              <TextField label="负责人 ID" value={form.ownerId} disabled={user?.role === 'Sales'} onChange={(event) => setForm({ ...form, ownerId: event.currentTarget.value })} />
            </div>
          </FormSection>
          <div className="saveBar">
            <Button onClick={() => setMode('list')}>取消</Button>
            <Button variant="primary" type="submit">保存任务</Button>
          </div>
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    const terminal = selected.status === 'Completed' || selected.status === 'Cancelled';
    return (
      <main className="content crudPage" aria-label="任务详情">
        <DetailHero
          eyebrow="返回任务列表"
          title={selected.title}
          subtitle={<><span>{labelFor(objectTypeLabel, selected.relatedType)} {selected.relatedId}</span><span>负责人 {selected.ownerId}</span><span>到期 {selected.dueDate}</span></>}
          icon={<ClipboardList size={20} aria-hidden="true" />}
          status={<StatusPill tone={taskTone(selected)}>{labelFor(taskStatusLabel, selected.status)}</StatusPill>}
          onBack={() => setMode('list')}
          actions={<Button variant="primary" disabled={terminal} onClick={() => void complete()}>完成任务</Button>}
          stats={
            <>
              <DetailStat label="关联记录" value={`${labelFor(objectTypeLabel, selected.relatedType)} ${selected.relatedId}`} icon={<ClipboardList size={17} aria-hidden="true" />} />
              <DetailStat label="到期日" value={selected.dueDate} icon={<CalendarDays size={17} aria-hidden="true" />} tone={isOverdue(selected) ? 'peach' : 'mint'} />
              <DetailStat label="状态" value={labelFor(taskStatusLabel, selected.status)} icon={<CheckCircle2 size={17} aria-hidden="true" />} />
              <DetailStat label="版本" value={`v${selected.version}`} icon={<ClipboardList size={17} aria-hidden="true" />} tone="peach" />
            </>
          }
        />
        {isOverdue(selected) ? <div role="alert" className="error">任务已逾期，请优先处理。</div> : null}
        <section className="detailContentGrid">
          <Panel>
            <div className="sectionHeader"><h2>任务字段</h2><StatusPill tone={taskTone(selected)}>{labelFor(taskStatusLabel, selected.status)}</StatusPill></div>
            <dl className="detailGrid">
              <div><dt>任务标题</dt><dd>{selected.title}</dd></div>
              <div><dt>关联记录</dt><dd>{labelFor(objectTypeLabel, selected.relatedType)} {selected.relatedId}</dd></div>
              <div><dt>负责人</dt><dd>{selected.ownerId}</dd></div>
              <div><dt>到期日</dt><dd>{selected.dueDate}</dd></div>
              <div><dt>版本</dt><dd>v{selected.version}</dd></div>
            </dl>
            <div className="actionBand opportunityActions">
              <Button disabled title="任务当前无批量归档/取消页面入口；按 A3 禁用。">取消任务</Button>
            </div>
          </Panel>
        </section>
      </main>
    );
  }

  return (
    <CrudListShell
      title="任务"
      description={`${scopeLabel} · 共 ${rows.length} 条任务 · 默认按到期日和更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={
        <>
          <Button onClick={() => void refresh(activeOnly, asOfToday)}><RotateCcw size={16} aria-hidden="true" />刷新</Button>
          <Button variant="primary" onClick={() => setMode('create')}><Plus size={16} aria-hidden="true" />新建任务</Button>
        </>
      }
      toolbar={
        <Toolbar
          searchValue={search}
          onSearchChange={setSearch}
          searchPlaceholder="搜索任务、关联记录或负责人"
          filters={
            <>
              <label className="inlineCheckbox">
                <span>状态</span>
                <select value={status} onChange={(event) => setStatus(event.currentTarget.value)}>
                  <option value="">全部状态</option>
                  {statuses.map((statusOption) => <option value={statusOption} key={statusOption}>{labelFor(taskStatusLabel, statusOption)}</option>)}
                </select>
              </label>
              <label className="inlineCheckbox">
                <input type="checkbox" checked={activeOnly} onChange={(event) => setActiveOnly(event.currentTarget.checked)} />
                仅活动任务
              </label>
              <label className="inlineCheckbox">
                <input type="checkbox" checked={asOfToday} onChange={(event) => setAsOfToday(event.currentTarget.checked)} />
                截至今日
              </label>
            </>
          }
          actions={<Button onClick={() => void refresh(activeOnly, asOfToday)}>应用筛选</Button>}
        />
      }
      activeFilters={
        <ActiveFilterSummary onClear={clearFilters}>
          <span className="chip">负责人：{scopeLabel}</span>
          <span className="chip">状态：{status ? labelFor(taskStatusLabel, status) : '全部状态'}</span>
          <span className="chip">活动：{activeOnly ? '仅活动任务' : '全部任务'}</span>
          <span className="chip">日期：{asOfToday ? '截至今日' : '不限到期日'}</span>
        </ActiveFilterSummary>
      }
      bulkBar={
        <BulkActionBar>
          <div className="bulkSummary">
            <span className="bulkCount">已选择 {selectedRows.length} 条</span>
            <span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '任务批量完成通过现有单条状态接口逐条执行；转移/归档无任务端点。')}</span>
          </div>
          <div className="bulkActions">
            {user?.role !== 'Sales' ? (
              <>
                <Button className="bulkButton" disabled title="任务当前无负责人转移接口；按 A3 禁用。">批量转移负责人</Button>
                <Button className="bulkButton" disabled title="任务当前无归档接口；按 A3 禁用。">批量归档</Button>
              </>
            ) : null}
            <Button className="bulkButton" disabled={selectedRows.length === 0} onClick={() => void completeSelected()}>批量完成</Button>
            <ExportSelectedButton disabled={selectedRows.length === 0} onExport={exportSelected} />
            <Button className="bulkButton" variant="primary" disabled={selectedRows.length === 0} onClick={() => setSelectedIds([])}>清除选择</Button>
          </div>
        </BulkActionBar>
      }
      table={
        <DataTable
          caption="任务结果表"
          rows={slice.items}
          rowKey={(task) => task.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(task) => selected?.id === task.id ? 'selected' : undefined}
          empty="没有符合当前筛选条件的任务。"
          columns={[
            {
              key: 'task',
              header: '任务',
              width: '280px',
              render: (task) => (
                <button className="recordLinkButton" type="button" onClick={() => { setSelected(task); setMode('detail'); }}>
                  <RecordIdentity icon={<ClipboardList size={17} aria-hidden="true" />} title={task.title} subtitle={isOverdue(task) ? '已逾期 · 需优先处理' : `到期 ${task.dueDate}`} tone={isOverdue(task) ? 'peach' : 'sky'} />
                </button>
              )
            },
            { key: 'related', header: '关联记录', render: (task) => `${labelFor(objectTypeLabel, task.relatedType)} ${task.relatedId}` },
            { key: 'status', header: '状态', render: (task) => <StatusPill tone={taskTone(task)}>{labelFor(taskStatusLabel, task.status)}</StatusPill> },
            { key: 'owner', header: '负责人', render: (task) => task.ownerId || '未分配' },
            { key: 'dueDate', header: '到期日', align: 'right', sortable: true, sortDirection: 'desc', render: (task) => task.dueDate || '未设置' }
          ]}
          actions={(task) => (
            <div className="rowActions">
              <button className="rowAction" type="button" aria-label={`查看任务 ${task.title}`} onClick={() => { setSelected(task); setMode('detail'); }}><ArrowRight size={16} aria-hidden="true" /></button>
              <span className="rowAction" aria-hidden="true"><MoreHorizontal size={16} /></span>
            </div>
          )}
        />
      }
      pagination={<CrudPagination slice={slice} onPageChange={setPage} onPageSizeChange={(next) => { setPageSize(next); setPage(1); }} />}
    />
  );
}

function TaskFormRules({ salesLocked }: { salesLocked: boolean }) {
  return (
    <div className="sideCard">
      <h3>字段状态</h3>
      <div className="rule"><span>1</span><p>任务必须绑定已有业务记录。</p></div>
      <div className="rule"><span>2</span><p>{salesLocked ? '销售新建任务负责人锁定为本人。' : '负责人可按授权范围填写。'}</p></div>
      <div className="rule"><span>3</span><p>完成状态通过现有单条状态接口持久化。</p></div>
    </div>
  );
}

function taskTone(task: WorkTask): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (task.status === 'Completed') return 'success';
  if (task.status === 'Cancelled') return 'neutral';
  if (task.status === 'Overdue' || isOverdue(task)) return 'danger';
  return 'primary';
}

function isOverdue(task: WorkTask) {
  return task.status !== 'Completed' && task.status !== 'Cancelled' && Boolean(task.dueDate) && task.dueDate < today();
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
