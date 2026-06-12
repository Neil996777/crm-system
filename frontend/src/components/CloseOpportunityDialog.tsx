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
  const canSubmit = mode === 'Won'
    ? contractId.trim() !== ''
    : reasonCode.trim() !== '' && reasonDetail.trim() !== '';

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!canSubmit) return;
    setSaving(true);
    try {
      await onConfirm({ contractId: contractId.trim(), closeDate, lostReason: { code: reasonCode.trim(), detail: reasonDetail.trim() } });
    } finally {
      setSaving(false);
    }
  }

  return (
    <form className="dialogPanel" onSubmit={submit}>
      <div className="sectionTitle">
        <h3>{mode === 'Won' ? '确认赢单' : '确认丢单'}</h3>
        <button className="secondaryButton" type="button" onClick={onCancel}>取消</button>
      </div>
      {mode === 'Won' && (
        <label>
          合同 ID
          <input required value={contractId} onChange={(event) => setContractId(event.target.value)} />
        </label>
      )}
      <label>
        关闭日期
        <input type="date" value={closeDate} onChange={(event) => setCloseDate(event.target.value)} />
      </label>
      {mode === 'Lost' && (
        <>
          <label>
            丢单原因
            <select required value={reasonCode} onChange={(event) => setReasonCode(event.target.value)}>
              <option value="">选择原因</option>
              <option value="PRICE">价格</option>
              <option value="COMPETITOR">竞争对手</option>
              <option value="NO_BUDGET">无预算</option>
              <option value="NO_DECISION">未决策</option>
              <option value="OTHER">其他</option>
            </select>
          </label>
          <label>
            原因详情
            <input required value={reasonDetail} onChange={(event) => setReasonDetail(event.target.value)} />
          </label>
        </>
      )}
      <button className="primaryButton" type="submit" disabled={saving || !canSubmit}>
        {mode === 'Won' ? '确认赢单' : '确认丢单'}
      </button>
    </form>
  );
}
