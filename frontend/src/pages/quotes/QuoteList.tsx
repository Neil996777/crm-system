import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Quote, createQuote, getQuote, listQuotes } from '../../api/quotes';
import { QuoteDetail } from './QuoteDetail';

export function QuoteList() {
  const [quotes, setQuotes] = useState<Quote[]>([]);
  const [selected, setSelected] = useState<Quote | null>(null);
  const [creating, setCreating] = useState(false);
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');
  const [form, setForm] = useState({ opportunityId: '', customerId: '', amount: '', validityEnd: '', ownerId: '' });

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh(nextSearch = search) {
    const response = await listQuotes(nextSearch);
    setQuotes(response.items);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createQuote(form);
      setSelected(created);
      setCreating(false);
      setForm({ opportunityId: '', customerId: '', amount: '', validityEnd: '', ownerId: '' });
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
    }
  }

  async function selectQuote(id: string) {
    setError('');
    setSelected(await getQuote(id));
  }

  async function updateSelected(quote: Quote) {
    setSelected(quote);
    await refresh();
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Quotes</h1>
          <p>Create the single quote for an opportunity and manage acceptance.</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => setCreating((value) => !value)}>New quote</button>
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
                Validity end
                <input type="date" value={form.validityEnd} onChange={(event) => setForm({ ...form, validityEnd: event.target.value })} />
              </label>
              <label>
                Owner ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">Save quote</button>
            </form>
          )}
          <div className="recordList" aria-label="Quote records">
            {quotes.length === 0 ? <p className="emptyState">No quotes found.</p> : quotes.map((quote) => (
              <button className="recordRow" type="button" key={quote.id} onClick={() => void selectQuote(quote.id)}>
                <strong>{quote.opportunityId}</strong>
                <span>{quote.status} · {quote.amount}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <QuoteDetail quote={selected} onUpdated={updateSelected} onError={setError} /> : <p className="emptyState">Select a quote to manage status.</p>}
        </div>
      </section>
    </main>
  );
}
