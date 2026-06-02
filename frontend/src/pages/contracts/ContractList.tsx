import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract, createContract, getContract, listContracts } from '../../api/contracts';
import { ContractDetail } from './ContractDetail';

export function ContractList() {
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [selected, setSelected] = useState<Contract | null>(null);
  const [creating, setCreating] = useState(false);
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');
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

  async function refresh(nextSearch = search) {
    const response = await listContracts(nextSearch);
    setContracts(response.items);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createContract(form);
      setSelected(created);
      setCreating(false);
      setForm({ quoteId: '', opportunityId: '', customerId: '', amount: '', expectedSignedDate: '', contractNote: '', amountDifferenceReason: '', ownerId: '' });
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
    }
  }

  async function selectContract(id: string) {
    setError('');
    setSelected(await getContract(id));
  }

  async function updateSelected(contract: Contract) {
    setSelected(contract);
    await refresh();
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Contracts</h1>
          <p>Create Pending Signature contracts from accepted quotes and manage signing.</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => setCreating((value) => !value)}>New contract</button>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={(event) => { event.preventDefault(); void refresh(search); }}>
            <label>
              Search
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <button className="secondaryButton" type="submit">Search</button>
          </form>
          {creating && (
            <form className="createPanel" onSubmit={submit}>
              <label>
                Quote ID
                <input value={form.quoteId} onChange={(event) => setForm({ ...form, quoteId: event.target.value })} />
              </label>
              <label>
                Opportunity ID
                <input value={form.opportunityId} onChange={(event) => setForm({ ...form, opportunityId: event.target.value })} />
              </label>
              <label>
                Customer ID
                <input value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              </label>
              <label>
                Amount
                <input value={form.amount} onChange={(event) => setForm({ ...form, amount: event.target.value })} />
              </label>
              <label>
                Expected signed date
                <input type="date" value={form.expectedSignedDate} onChange={(event) => setForm({ ...form, expectedSignedDate: event.target.value })} />
              </label>
              <label>
                Contract note
                <textarea value={form.contractNote} onChange={(event) => setForm({ ...form, contractNote: event.target.value })} />
              </label>
              <label>
                Amount difference reason
                <input value={form.amountDifferenceReason} onChange={(event) => setForm({ ...form, amountDifferenceReason: event.target.value })} />
              </label>
              <label>
                Owner ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">Save contract</button>
            </form>
          )}
          <div className="recordList" aria-label="Contract records">
            {contracts.length === 0 ? <p className="emptyState">No contracts found.</p> : contracts.map((contract) => (
              <button className="recordRow" type="button" key={contract.id} onClick={() => void selectContract(contract.id)}>
                <strong>{contract.opportunityId}</strong>
                <span>{contract.status} · {contract.amount}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <ContractDetail contract={selected} onUpdated={updateSelected} onError={setError} /> : <p className="emptyState">Select a contract to manage status.</p>}
        </div>
      </section>
    </main>
  );
}
