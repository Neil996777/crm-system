import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Opportunity, createOpportunity, getOpportunity, listOpportunities } from '../../api/opportunities';
import { OpportunityDetail } from './OpportunityDetail';

export function OpportunityList() {
  const [opportunities, setOpportunities] = useState<Opportunity[]>([]);
  const [selected, setSelected] = useState<Opportunity | null>(null);
  const [creating, setCreating] = useState(false);
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');
  const [form, setForm] = useState({
    title: '',
    customerId: '',
    ownerId: '',
    expectedAmount: '',
    expectedCloseDate: ''
  });

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh(nextSearch = search) {
    const response = await listOpportunities(nextSearch);
    setOpportunities(response.items);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createOpportunity({ ...form, stage: 'New Opportunity' });
      setSelected(created);
      setCreating(false);
      setForm({ title: '', customerId: '', ownerId: '', expectedAmount: '', expectedCloseDate: '' });
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
    }
  }

  async function selectOpportunity(id: string) {
    setError('');
    setSelected(await getOpportunity(id));
  }

  async function updateSelected(opportunity: Opportunity) {
    setSelected(opportunity);
    await refresh();
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Opportunities</h1>
          <p>Move deals through pipeline stages and terminal close outcomes.</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => setCreating((value) => !value)}>New opportunity</button>
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
                Title
                <input value={form.title} onChange={(event) => setForm({ ...form, title: event.target.value })} />
              </label>
              <label>
                Customer ID
                <input value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              </label>
              <label>
                Owner ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <label>
                Expected amount
                <input value={form.expectedAmount} onChange={(event) => setForm({ ...form, expectedAmount: event.target.value })} />
              </label>
              <label>
                Expected close date
                <input type="date" value={form.expectedCloseDate} onChange={(event) => setForm({ ...form, expectedCloseDate: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">Save opportunity</button>
            </form>
          )}
          <div className="recordList" aria-label="Opportunity records">
            {opportunities.length === 0 ? <p className="emptyState">No opportunities found.</p> : opportunities.map((opportunity) => (
              <button className="recordRow" type="button" key={opportunity.id} onClick={() => void selectOpportunity(opportunity.id)}>
                <strong>{opportunity.title || opportunity.id}</strong>
                <span>{opportunity.stage}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <OpportunityDetail opportunity={selected} onUpdated={updateSelected} onError={setError} /> : <p className="emptyState">Select an opportunity to view stage and close actions.</p>}
        </div>
      </section>
    </main>
  );
}
