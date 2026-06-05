import { CheckCircle2, RotateCcw, XCircle } from 'lucide-react';
import { useState } from 'react';
import { Lead, qualifyInvalid, qualifyValid, restoreInvalid } from '../api/leads';
import { ApiError } from '../api/client';
import { localizeError } from '../i18n/labels';

type Props = {
  lead: Lead;
  onUpdated: (lead: Lead) => void;
  onError: (message: string) => void;
};

export function QualificationActions({ lead, onUpdated, onError }: Props) {
  const [invalidReason, setInvalidReason] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const unassigned = lead.status === 'Unassigned';
  const converted = lead.status === 'Converted To Opportunity';

  async function run(action: () => Promise<Lead>) {
    setSubmitting(true);
    onError('');
    try {
      onUpdated(await action());
      setInvalidReason('');
    } catch (caught) {
      const error = caught as ApiError;
      onError(localizeError(error));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="actionBand" aria-label="确认线索">
      {unassigned && <p className="inlineNotice">未分配线索不能确认或转换。</p>}
      <button className="secondaryButton iconButtonText" type="button" disabled={submitting || unassigned || converted || lead.status !== 'Pending Qualification'} onClick={() => run(() => qualifyValid(lead))}>
        <CheckCircle2 size={16} />
        标记有效
      </button>
      <div className="inlineForm">
        <label>
          无效原因
          <input value={invalidReason} onChange={(event) => setInvalidReason(event.target.value)} disabled={submitting || unassigned || converted || lead.status !== 'Pending Qualification'} />
        </label>
        <button className="secondaryButton iconButtonText" type="button" disabled={submitting || unassigned || converted || lead.status !== 'Pending Qualification' || invalidReason.trim() === ''} onClick={() => run(() => qualifyInvalid(lead, invalidReason))}>
          <XCircle size={16} />
          标记无效
        </button>
      </div>
      <button className="secondaryButton iconButtonText" type="button" disabled={submitting || lead.status !== 'Invalid'} onClick={() => run(() => restoreInvalid(lead))}>
        <RotateCcw size={16} />
        恢复无效线索
      </button>
    </section>
  );
}
