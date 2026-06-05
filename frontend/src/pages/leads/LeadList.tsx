import { FormEvent, useEffect, useMemo, useState } from 'react';
import { ApiError } from '../../api/client';
import { DuplicateWarningResult } from '../../api/duplicates';
import { ConversionResult, Lead, checkLeadDuplicate, createLead, getLead, listLeads } from '../../api/leads';
import { DuplicateWarning } from '../../components/DuplicateWarning';
import { labelFor, leadStatusLabel, localizeError } from '../../i18n/labels';
import { LeadDetail } from './LeadDetail';

export function LeadList() {
  const [leads, setLeads] = useState<Lead[]>([]);
  const [selected, setSelected] = useState<Lead | null>(null);
  const [search, setSearch] = useState('');
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState('');
  const [duplicateWarning, setDuplicateWarning] = useState<DuplicateWarningResult | null>(null);
  const [form, setForm] = useState({ leadName: '', companyName: '', email: '', phone: '', source: '', ownerId: '', needSummary: '' });

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh(nextSearch = search) {
    const response = await listLeads(nextSearch);
    setLeads(response.items);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    await saveLead();
  }

  async function saveLead(proceedWarningToken?: string) {
    setError('');
    try {
      if (!proceedWarningToken) {
        const warning = await checkLeadDuplicate({ companyName: form.companyName, email: form.email, phone: form.phone });
        if (warning.result === 'PossibleDuplicate' && warning.warningToken) {
          setDuplicateWarning(warning);
          return;
        }
      }
      const created = await createLead({
        leadName: form.leadName,
        companyName: form.companyName,
        email: form.email,
        phone: form.phone,
        source: form.source,
        ownerId: form.ownerId,
        needSummary: form.needSummary,
        proceedWarningToken
      });
      setCreating(false);
      setForm({ leadName: '', companyName: '', email: '', phone: '', source: '', ownerId: '', needSummary: '' });
      setDuplicateWarning(null);
      setSelected(created);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectLead(id: string) {
    setError('');
    setSelected(await getLead(id));
  }

  async function updateSelected(lead: Lead) {
    setSelected(lead);
    await refresh();
  }

  async function converted(result: ConversionResult) {
    const lead = await getLead(result.leadId);
    setSelected(lead);
    await refresh();
  }

  const selectedTitle = useMemo(() => selected ? selected.companyName || selected.leadName : '选择线索', [selected]);

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>线索</h1>
          <p>录入、确认并转换销售线索。</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => { setDuplicateWarning(null); setCreating((value) => !value); }}>
          新建线索
        </button>
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
                线索名称
                <input value={form.leadName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, leadName: event.target.value }); }} />
              </label>
              <label>
                公司名称
                <input value={form.companyName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, companyName: event.target.value }); }} />
              </label>
              <label>
                邮箱
                <input value={form.email} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, email: event.target.value }); }} />
              </label>
              <label>
                电话
                <input value={form.phone} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, phone: event.target.value }); }} />
              </label>
              <label>
                来源
                <input value={form.source} onChange={(event) => setForm({ ...form, source: event.target.value })} />
              </label>
              <label>
                负责人 ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <label>
                需求摘要
                <input value={form.needSummary} onChange={(event) => setForm({ ...form, needSummary: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">保存线索</button>
              {duplicateWarning ? (
                <DuplicateWarning
                  warning={duplicateWarning}
                  onProceed={() => void saveLead(duplicateWarning.warningToken)}
                  onCancel={() => setDuplicateWarning(null)}
                />
              ) : null}
            </form>
          )}
          <div className="recordList" aria-label="线索记录">
            {leads.length === 0 ? <p className="emptyState">暂无线索。</p> : leads.map((lead) => (
              <button className="recordRow" type="button" key={lead.id} onClick={() => void selectLead(lead.id)}>
                <strong>{lead.companyName || lead.leadName}</strong>
                <span>{labelFor(leadStatusLabel, lead.status)}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          <h2 className="srOnly">{selectedTitle}</h2>
          {selected ? <LeadDetail lead={selected} onUpdated={updateSelected} onConverted={converted} onError={setError} /> : <p className="emptyState">选择线索以查看确认操作。</p>}
        </div>
      </section>
    </main>
  );
}
