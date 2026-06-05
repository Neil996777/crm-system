import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Opportunity, createOpportunity, getOpportunity, listOpportunities } from '../../api/opportunities';
import { labelFor, opportunityStageLabel } from '../../i18n/labels';
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
      setError(apiError.safeMessage || '请求失败。');
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
          <h1>商机</h1>
          <p>推进商机阶段并处理赢单/丢单结果。</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => setCreating((value) => !value)}>新建商机</button>
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
                标题
                <input value={form.title} onChange={(event) => setForm({ ...form, title: event.target.value })} />
              </label>
              <label>
                客户 ID
                <input value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              </label>
              <label>
                负责人 ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <label>
                预计金额
                <input value={form.expectedAmount} onChange={(event) => setForm({ ...form, expectedAmount: event.target.value })} />
              </label>
              <label>
                预计关闭日期
                <input type="date" value={form.expectedCloseDate} onChange={(event) => setForm({ ...form, expectedCloseDate: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">保存商机</button>
            </form>
          )}
          <div className="recordList" aria-label="商机记录">
            {opportunities.length === 0 ? <p className="emptyState">暂无商机。</p> : opportunities.map((opportunity) => (
              <button className="recordRow" type="button" key={opportunity.id} onClick={() => void selectOpportunity(opportunity.id)}>
                <strong>{opportunity.title || opportunity.id}</strong>
                <span>{labelFor(opportunityStageLabel, opportunity.stage)}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <OpportunityDetail opportunity={selected} onUpdated={updateSelected} onError={setError} /> : <p className="emptyState">选择商机以查看阶段和关闭操作。</p>}
        </div>
      </section>
    </main>
  );
}
