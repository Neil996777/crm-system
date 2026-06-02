import { FormEvent, useEffect, useState } from 'react';
import { Account, checkAccountDuplicate, createAccount, getAccount, listAccounts } from '../../api/accounts';
import { ApiError } from '../../api/client';
import { DuplicateWarningResult } from '../../api/duplicates';
import { DuplicateWarning } from '../../components/DuplicateWarning';
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
      setError(apiError.safeMessage || 'Request failed.');
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
          <h1>Companies/Customers</h1>
          <p>Manage account records and their contacts.</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => { setDuplicateWarning(null); setCreating((value) => !value); }}>New customer</button>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={(event) => { event.preventDefault(); void refresh(search); }}>
            <label>
              Search
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <label className="inlineCheckbox">
              <input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />
              Include archived
            </label>
            <button className="secondaryButton" type="submit">Search</button>
          </form>
          {creating && (
            <form className="createPanel" onSubmit={submit}>
              <label>
                Company name
                <input value={form.companyName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, companyName: event.target.value }); }} />
              </label>
              <label>
                Customer status
                <input value={form.customerStatus} onChange={(event) => setForm({ ...form, customerStatus: event.target.value })} />
              </label>
              <label>
                Owner ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">Save customer</button>
              {duplicateWarning ? (
                <DuplicateWarning
                  warning={duplicateWarning}
                  onProceed={() => void saveAccount(duplicateWarning.warningToken)}
                  onCancel={() => setDuplicateWarning(null)}
                />
              ) : null}
            </form>
          )}
          <div className="recordList" aria-label="Customer records">
            {accounts.length === 0 ? <p className="emptyState">No customers found.</p> : accounts.map((account) => (
              <button className="recordRow" type="button" key={account.id} onClick={() => void selectAccount(account.id)}>
                <strong>{account.companyName}</strong>
                <span>{account.customerStatus}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <AccountDetail account={selected} onArchived={() => { setSelected(null); void refresh(); }} onError={setError} /> : <p className="emptyState">Select a customer to view contacts.</p>}
        </div>
      </section>
    </main>
  );
}
