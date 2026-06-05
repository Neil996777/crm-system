import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Quote, createQuote, getQuote, listQuotes } from '../../api/quotes';
import { labelFor, quoteStatusLabel } from '../../i18n/labels';
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
      setError(apiError.safeMessage || '请求失败。');
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
          <h1>报价</h1>
          <p>为商机创建唯一报价并管理接受状态。</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => setCreating((value) => !value)}>新建报价</button>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={(event) => { event.preventDefault(); void refresh(search); }}>
            <label>
              搜索
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <button className="secondaryButton" type="submit">搜索</button>
          </form>
          {creating && (
            <form className="createPanel" onSubmit={submit}>
              <label>
                商机 ID
                <input value={form.opportunityId} onChange={(event) => setForm({ ...form, opportunityId: event.target.value })} />
              </label>
              <label>
                客户 ID
                <input value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              </label>
              <label>
                金额
                <input value={form.amount} onChange={(event) => setForm({ ...form, amount: event.target.value })} />
              </label>
              <label>
                有效期截止日
                <input type="date" value={form.validityEnd} onChange={(event) => setForm({ ...form, validityEnd: event.target.value })} />
              </label>
              <label>
                负责人 ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">保存报价</button>
            </form>
          )}
          <div className="recordList" aria-label="报价记录">
            {quotes.length === 0 ? <p className="emptyState">暂无报价。</p> : quotes.map((quote) => (
              <button className="recordRow" type="button" key={quote.id} onClick={() => void selectQuote(quote.id)}>
                <strong>{quote.opportunityId}</strong>
                <span>{labelFor(quoteStatusLabel, quote.status)} · {quote.amount}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <QuoteDetail quote={selected} onUpdated={updateSelected} onError={setError} /> : <p className="emptyState">选择报价以管理状态。</p>}
        </div>
      </section>
    </main>
  );
}
