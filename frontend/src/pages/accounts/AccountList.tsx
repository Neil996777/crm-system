import { FormEvent, useEffect, useState } from 'react';
import { Account, checkAccountDuplicate, createAccount, getAccount, listAccounts } from '../../api/accounts';
import { ApiError } from '../../api/client';
import { DuplicateWarningResult } from '../../api/duplicates';
import { DuplicateWarning } from '../../components/DuplicateWarning';
import { accountStatusLabel, archiveStatusLabel, labelFor } from '../../i18n/labels';
import { AccountDetail } from './AccountDetail';

export function AccountList() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [selected, setSelected] = useState<Account | null>(null);
  const [creating, setCreating] = useState(false);
  const [search, setSearch] = useState('');
  const [includeArchived, setIncludeArchived] = useState(false);
  const [error, setError] = useState('');
  const [duplicateWarning, setDuplicateWarning] = useState<DuplicateWarningResult | null>(null);
  const [form, setForm] = useState({ companyName: '', customerStatus: '', ownerId: '' });

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh(nextSearch = search) {
    const response = await listAccounts(nextSearch, includeArchived);
    setAccounts(response.items);
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
      setCreating(false);
      setForm({ companyName: '', customerStatus: '', ownerId: '' });
      setDuplicateWarning(null);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || '请求失败。');
    }
  }

  async function selectAccount(id: string) {
    setError('');
    setSelected(await getAccount(id));
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>公司/客户</h1>
          <p>管理客户记录及其联系人。</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => { setDuplicateWarning(null); setCreating((value) => !value); }}>新建客户</button>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={(event) => { event.preventDefault(); void refresh(search); }}>
            <label>
              搜索
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <label className="inlineCheckbox">
              <input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />
              包含已归档
            </label>
            <button className="secondaryButton" type="submit">搜索</button>
          </form>
          {creating && (
            <form className="createPanel" onSubmit={submit}>
              <label>
                公司名称
                <input value={form.companyName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, companyName: event.target.value }); }} />
              </label>
              <label>
                客户状态
                <input value={form.customerStatus} onChange={(event) => setForm({ ...form, customerStatus: event.target.value })} />
              </label>
              <label>
                负责人 ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">保存客户</button>
              {duplicateWarning ? (
                <DuplicateWarning
                  warning={duplicateWarning}
                  onProceed={() => void saveAccount(duplicateWarning.warningToken)}
                  onCancel={() => setDuplicateWarning(null)}
                />
              ) : null}
            </form>
          )}
          <div className="recordList" aria-label="客户记录">
            {accounts.length === 0 ? <p className="emptyState">暂无客户。</p> : accounts.map((account) => (
              <button className="recordRow" type="button" key={account.id} onClick={() => void selectAccount(account.id)}>
                <strong>{account.companyName}</strong>
                <span>{account.archived ? labelFor(archiveStatusLabel, 'Archived') : labelFor(accountStatusLabel, account.customerStatus)}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <AccountDetail account={selected} onArchived={() => { setSelected(null); void refresh(); }} onError={setError} /> : <p className="emptyState">选择客户以查看联系人。</p>}
        </div>
      </section>
    </main>
  );
}
