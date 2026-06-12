import { FormEvent, useEffect, useMemo, useState } from 'react';
import { Plus, RotateCcw, UserRound } from 'lucide-react';
import { Contact, createContact, getContact, listAllContacts } from '../../api/accounts';
import { ApiError } from '../../api/client';
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
} from '../../components/CrudScaffold';
import { BulkActionBar, Button, DataTable, TextField, Toolbar } from '../../components/ui';
import { useSession } from '../../auth/SessionProvider';
import { localizeError } from '../../i18n/labels';

type Mode = 'list' | 'create' | 'detail';

export function ContactList({ targetRecordId, onTargetHandled }: { targetRecordId?: string; onTargetHandled?: () => void }) {
  const { user } = useSession();
  const [mode, setMode] = useState<Mode>('list');
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [selected, setSelected] = useState<Contact | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [form, setForm] = useState({ accountId: '', contactName: '', email: '', phone: '', roleNote: '' });

  useEffect(() => {
    void refresh();
  }, []);

  useEffect(() => {
    if (!targetRecordId) return;
    void selectContact(targetRecordId).finally(() => onTargetHandled?.());
  }, [targetRecordId, onTargetHandled]);

  async function refresh(nextSearch = search) {
    const response = await listAllContacts(nextSearch);
    setContacts(response.items);
    setPage(1);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createContact(form.accountId, {
        contactName: form.contactName,
        email: form.email,
        phone: form.phone,
        roleNote: form.roleNote
      });
      setSelected(created);
      setForm({ accountId: '', contactName: '', email: '', phone: '', roleNote: '' });
      setMode('detail');
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectContact(id: string) {
    setSelected(await getContact(id));
    setMode('detail');
  }

  const rows = useMemo(() => contacts, [contacts]);
  const slice = paginate(rows, page, pageSize);
  const selectedRows = rows.filter((contact) => selectedIds.includes(contact.id));
  const scopeLabel = user?.role === 'Sales' ? '本人范围' : user?.role === 'Sales Manager' ? '团队范围' : '全部范围';

  function toggleRow(contact: Contact, checked: boolean) {
    setSelectedIds((value) => checked ? [...new Set([...value, contact.id])] : value.filter((id) => id !== contact.id));
  }

  function toggleAll(checked: boolean) {
    setSelectedIds(checked ? slice.items.map((contact) => contact.id) : []);
  }

  function exportSelected() {
    if (selectedRows.length === 0) return;
    exportRows('contacts-selected.csv', selectedRows.map((contact) => ({
      联系人: contact.contactName,
      客户: contact.accountName || contact.accountId,
      邮箱: contact.email,
      电话: contact.phone,
      角色备注: contact.roleNote
    })));
    setNotice(`已导出 ${selectedRows.length} 条联系人。`);
  }

  if (mode === 'create') {
    return (
      <FormShell
        title="新建联系人"
        description="联系人必须归属一个有权限访问的客户。"
        badge="新建"
        onCancel={() => setMode('list')}
        actions={<Button variant="primary" form="contact-form" type="submit">保存</Button>}
        side={<ContactFormRules />}
      >
        {error && <div role="alert" className="error">{error}</div>}
        <form id="contact-form" className="actionBand" onSubmit={submit}>
          <FormSection title="联系人基本信息">
            <div className="formFields">
              <TextField label="客户 ID" value={form.accountId} onChange={(event) => setForm({ ...form, accountId: event.target.value })} />
              <TextField label="联系人姓名" value={form.contactName} onChange={(event) => setForm({ ...form, contactName: event.target.value })} />
              <TextField label="邮箱" value={form.email} onChange={(event) => setForm({ ...form, email: event.target.value })} />
              <TextField label="电话" value={form.phone} onChange={(event) => setForm({ ...form, phone: event.target.value })} />
              <TextField className="full" label="角色备注" value={form.roleNote} onChange={(event) => setForm({ ...form, roleNote: event.target.value })} />
            </div>
          </FormSection>
          <div className="saveBar">
            <Button onClick={() => setMode('list')}>取消</Button>
            <Button variant="primary" type="submit">保存联系人</Button>
          </div>
        </form>
      </FormShell>
    );
  }

  if (mode === 'detail' && selected) {
    return (
      <main className="content crudPage" aria-label="联系人详情">
        <DetailHero
          eyebrow="返回联系人列表"
          title={selected.contactName}
          subtitle={<><span>{selected.accountName || selected.accountId}</span><span>更新于 {formatDate(selected.updatedAt)}</span></>}
          icon={<UserRound size={20} aria-hidden="true" />}
          status={<StatusPill>联系人</StatusPill>}
          onBack={() => setMode('list')}
          stats={
            <>
              <DetailStat label="客户" value={selected.accountName || selected.accountId} icon={<UserRound size={17} aria-hidden="true" />} />
              <DetailStat label="邮箱" value={selected.email || '无'} icon={<UserRound size={17} aria-hidden="true" />} tone="peach" />
              <DetailStat label="电话" value={selected.phone || '无'} icon={<UserRound size={17} aria-hidden="true" />} />
            </>
          }
        />
        <section className="detailContentGrid">
          <div className="panel">
            <div className="sectionHeader"><h2>联系人字段</h2><span className="badge">详情</span></div>
            <dl className="detailGrid">
              <div><dt>客户</dt><dd>{selected.accountName || selected.accountId}</dd></div>
              <div><dt>邮箱</dt><dd>{selected.email || '无'}</dd></div>
              <div><dt>电话</dt><dd>{selected.phone || '无'}</dd></div>
              <div><dt>角色备注</dt><dd>{selected.roleNote || '无'}</dd></div>
            </dl>
          </div>
        </section>
      </main>
    );
  }

  return (
    <CrudListShell
      title="联系人"
      description={`${scopeLabel} · 共 ${rows.length} 条 · 默认按更新时间倒序`}
      scope={`${scopeLabel} · 第 ${slice.page} / ${slice.totalPages} 页`}
      actions={
        <>
          <Button onClick={() => void refresh(search)}><RotateCcw size={16} aria-hidden="true" />刷新</Button>
          <Button variant="primary" onClick={() => setMode('create')}><Plus size={16} aria-hidden="true" />新建联系人</Button>
        </>
      }
      toolbar={<Toolbar searchValue={search} onSearchChange={setSearch} searchPlaceholder="搜索联系人或客户" actions={<Button onClick={() => void refresh(search)}>应用筛选</Button>} />}
      activeFilters={<ActiveFilterSummary onClear={() => { setSearch(''); setSelectedIds([]); void refresh(''); }}><span className="chip">负责人：{scopeLabel}</span><span className="chip">关键词：{search || '全部'}</span></ActiveFilterSummary>}
      bulkBar={
        <BulkActionBar>
          <div className="bulkSummary">
            <span className="bulkCount">已选择 {selectedRows.length} 条</span>
            <span className="bulkHint">{notice || (user?.role === 'Sales' ? '销售角色仅保留导出和清除选择。' : '联系人当前无转移/归档接口，按 A3 仅提供导出与清除选择。')}</span>
          </div>
          <div className="bulkActions">
            {user?.role !== 'Sales' ? (
              <>
                <Button className="bulkButton" disabled title="联系人当前无负责人转移接口；按 A3 禁用。">批量转移负责人</Button>
                <Button className="bulkButton" disabled title="联系人当前无归档接口；按 A3 禁用。">批量归档</Button>
              </>
            ) : null}
            <ExportSelectedButton disabled={selectedRows.length === 0} onExport={exportSelected} />
            <Button className="bulkButton" variant="primary" onClick={() => setSelectedIds([])} disabled={selectedRows.length === 0}>清除选择</Button>
          </div>
        </BulkActionBar>
      }
      table={
        <DataTable
          caption="联系人结果表"
          rows={slice.items}
          rowKey={(contact) => contact.id}
          selectedRowKeys={selectedIds}
          onToggleRow={toggleRow}
          onToggleAll={toggleAll}
          getRowClassName={(contact) => selected?.id === contact.id ? 'selected' : undefined}
          onRowClick={(contact) => void selectContact(contact.id)}
          getRowAriaLabel={(contact) => `打开联系人 ${contact.contactName}`}
          empty="没有符合当前筛选条件的联系人。"
          columns={[
            {
              key: 'name',
              header: '联系人',
              width: '220px',
              render: (contact) => (
                <RecordIdentity
                  icon={<UserRound size={17} aria-hidden="true" />}
                  title={contact.contactName}
                  titleAriaLabel={`打开联系人 ${contact.contactName}`}
                  subtitle={contact.accountName || contact.accountId}
                  tone="mint"
                  onTitleClick={() => void selectContact(contact.id)}
                />
              )
            },
            { key: 'email', header: '邮箱', render: (contact) => contact.email || '无' },
            { key: 'phone', header: '电话', render: (contact) => contact.phone || '无' },
            { key: 'role', header: '角色备注', render: (contact) => contact.roleNote || '无' },
            { key: 'updated', header: '更新时间', align: 'right', sortable: true, sortDirection: 'desc', render: (contact) => formatDate(contact.updatedAt) }
          ]}
        />
      }
      pagination={<CrudPagination slice={slice} onPageChange={setPage} onPageSizeChange={(next) => { setPageSize(next); setPage(1); }} />}
    />
  );
}

function ContactFormRules() {
  return (
    <div className="sideCard">
      <h3>字段状态</h3>
      <div className="rule"><span>1</span><p>客户 ID 指向联系人所属公司客户。</p></div>
      <div className="rule"><span>2</span><p>邮箱或电话用于跨客户检索和重复检查。</p></div>
    </div>
  );
}

function formatDate(value: string) {
  if (!value) return '未更新';
  return value.length > 10 ? value.slice(0, 10) : value;
}
