import { FormEvent, useEffect, useMemo, useState } from 'react';
import { ApiError } from '../../api/client';
import { DuplicateWarningResult } from '../../api/duplicates';
import { ConversionResult, Lead, checkLeadDuplicate, createLead, getLead, listLeads } from '../../api/leads';
import { DuplicateWarning } from '../../components/DuplicateWarning';
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
      setError(apiError.safeMessage || 'Request failed.');
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

  const selectedTitle = useMemo(() => selected ? selected.companyName || selected.leadName : 'Select a lead', [selected]);

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Leads</h1>
          <p>Capture, qualify, and convert sales leads.</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => { setDuplicateWarning(null); setCreating((value) => !value); }}>
          New lead
        </button>
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
                Lead name
                <input value={form.leadName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, leadName: event.target.value }); }} />
              </label>
              <label>
                Company name
                <input value={form.companyName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, companyName: event.target.value }); }} />
              </label>
              <label>
                Email
                <input value={form.email} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, email: event.target.value }); }} />
              </label>
              <label>
                Phone
                <input value={form.phone} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, phone: event.target.value }); }} />
              </label>
              <label>
                Source
                <input value={form.source} onChange={(event) => setForm({ ...form, source: event.target.value })} />
              </label>
              <label>
                Owner ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <label>
                Need summary
                <input value={form.needSummary} onChange={(event) => setForm({ ...form, needSummary: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">Save lead</button>
              {duplicateWarning ? (
                <DuplicateWarning
                  warning={duplicateWarning}
                  onProceed={() => void saveLead(duplicateWarning.warningToken)}
                  onCancel={() => setDuplicateWarning(null)}
                />
              ) : null}
            </form>
          )}
          <div className="recordList" aria-label="Lead records">
            {leads.length === 0 ? <p className="emptyState">No leads found.</p> : leads.map((lead) => (
              <button className="recordRow" type="button" key={lead.id} onClick={() => void selectLead(lead.id)}>
                <strong>{lead.companyName || lead.leadName}</strong>
                <span>{lead.status}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          <h2 className="srOnly">{selectedTitle}</h2>
          {selected ? <LeadDetail lead={selected} onUpdated={updateSelected} onConverted={converted} onError={setError} /> : <p className="emptyState">Select a lead to view qualification actions.</p>}
        </div>
      </section>
    </main>
  );
}
