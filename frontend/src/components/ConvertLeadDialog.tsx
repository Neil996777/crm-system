import { GitBranch } from 'lucide-react';
import { FormEvent, useState } from 'react';
import { ApiError } from '../api/client';
import { ConversionResult, Lead, convertLead } from '../api/leads';

type Props = {
  lead: Lead;
  onConverted: (result: ConversionResult) => void;
  onError: (message: string) => void;
};

export function ConvertLeadDialog({ lead, onConverted, onError }: Props) {
  const [open, setOpen] = useState(false);
  const [expectedAmount, setExpectedAmount] = useState('');
  const [expectedCloseDate, setExpectedCloseDate] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const disabled = lead.status !== 'Valid' || lead.ownerId === '';

  async function submit(event: FormEvent) {
    event.preventDefault();
    setSubmitting(true);
    onError('');
    try {
      onConverted(await convertLead(lead, expectedAmount, expectedCloseDate));
      setOpen(false);
    } catch (caught) {
      const error = caught as ApiError;
      onError(error.safeMessage || 'Request failed.');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="convertPanel">
      <button className="primaryButton iconButtonText" type="button" disabled={disabled} onClick={() => setOpen((value) => !value)}>
        <GitBranch size={16} />
        Convert lead
      </button>
      {open && (
        <form className="inlineForm convertForm" onSubmit={submit}>
          <label>
            Expected amount
            <input value={expectedAmount} onChange={(event) => setExpectedAmount(event.target.value)} />
          </label>
          <label>
            Expected close date
            <input type="date" value={expectedCloseDate} onChange={(event) => setExpectedCloseDate(event.target.value)} />
          </label>
          <button className="primaryButton" type="submit" disabled={submitting || expectedAmount.trim() === '' || expectedCloseDate === ''}>
            Convert
          </button>
        </form>
      )}
    </section>
  );
}
