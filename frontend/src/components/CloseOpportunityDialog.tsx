import { FormEvent, useState } from 'react';

type Mode = 'Won' | 'Lost';

export function CloseOpportunityDialog({
  mode,
  onCancel,
  onConfirm
}: {
  mode: Mode;
  onCancel: () => void;
  onConfirm: (input: { contractId: string; closeDate: string; lostReason: { code: string; detail: string } }) => Promise<void>;
}) {
  const [contractId, setContractId] = useState('');
  const [closeDate, setCloseDate] = useState('');
  const [reasonCode, setReasonCode] = useState('');
  const [reasonDetail, setReasonDetail] = useState('');
  const [saving, setSaving] = useState(false);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setSaving(true);
    try {
      await onConfirm({ contractId, closeDate, lostReason: { code: reasonCode, detail: reasonDetail } });
    } finally {
      setSaving(false);
    }
  }

  return (
    <form className="dialogPanel" onSubmit={submit}>
      <div className="sectionTitle">
        <h3>{mode === 'Won' ? 'Close Won' : 'Close Lost'}</h3>
        <button className="secondaryButton" type="button" onClick={onCancel}>Cancel</button>
      </div>
      {mode === 'Won' && (
        <label>
          Contract ID
          <input value={contractId} onChange={(event) => setContractId(event.target.value)} />
        </label>
      )}
      <label>
        Close date
        <input type="date" value={closeDate} onChange={(event) => setCloseDate(event.target.value)} />
      </label>
      {mode === 'Lost' && (
        <>
          <label>
            Lost reason
            <select value={reasonCode} onChange={(event) => setReasonCode(event.target.value)}>
              <option value="">Select reason</option>
              <option value="PRICE">PRICE</option>
              <option value="COMPETITOR">COMPETITOR</option>
              <option value="NO_BUDGET">NO_BUDGET</option>
              <option value="NO_DECISION">NO_DECISION</option>
              <option value="OTHER">OTHER</option>
            </select>
          </label>
          <label>
            Reason detail
            <input value={reasonDetail} onChange={(event) => setReasonDetail(event.target.value)} />
          </label>
        </>
      )}
      <button className="primaryButton" type="submit" disabled={saving}>
        {mode === 'Won' ? 'Confirm Won' : 'Confirm Lost'}
      </button>
    </form>
  );
}
