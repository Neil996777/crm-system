import { useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract, changeContractStatus } from '../../api/contracts';

export function ContractDetail({ contract, onUpdated, onError }: { contract: Contract; onUpdated: (contract: Contract) => Promise<void>; onError: (message: string) => void }) {
  const [signedEffectiveDate, setSignedEffectiveDate] = useState(contract.signedEffectiveDate ?? '');
  const pendingReminder = contract.status === 'Pending Signature' && contract.expectedSignedDate < today();

  async function change(toStatus: string) {
    onError('');
    try {
      const next = await changeContractStatus(contract.id, contract.version, toStatus, signedEffectiveDate);
      setSignedEffectiveDate(next.signedEffectiveDate ?? signedEffectiveDate);
      await onUpdated(next);
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(apiError.safeMessage || 'Request failed.');
    }
  }

  return (
    <section className="detailPane" aria-label="Contract detail">
      <div className="detailHeader">
        <div>
          <h2>{contract.id}</h2>
          <p>Status: {contract.status}</p>
        </div>
        <span className="statusPill">{contract.status}</span>
      </div>
      {pendingReminder && <div role="alert" className="error">Pending signature expected date has passed.</div>}
      <dl className="detailGrid">
        <div>
          <dt>Quote</dt>
          <dd>{contract.quoteId}</dd>
        </div>
        <div>
          <dt>Opportunity</dt>
          <dd>{contract.opportunityId}</dd>
        </div>
        <div>
          <dt>Customer</dt>
          <dd>{contract.customerId}</dd>
        </div>
        <div>
          <dt>Amount</dt>
          <dd>{contract.amount}</dd>
        </div>
        <div>
          <dt>Expected signed date</dt>
          <dd>{contract.expectedSignedDate}</dd>
        </div>
        <div>
          <dt>Signed/effective date</dt>
          <dd>{contract.signedEffectiveDate || 'Not signed'}</dd>
        </div>
        <div>
          <dt>Contract note</dt>
          <dd>{contract.contractNote}</dd>
        </div>
        <div>
          <dt>Amount difference reason</dt>
          <dd>{contract.amountDifferenceReason || 'None'}</dd>
        </div>
      </dl>
      <label className="singleField">
        Signed/effective date
        <input type="date" value={signedEffectiveDate} onChange={(event) => setSignedEffectiveDate(event.target.value)} />
      </label>
      <div className="actionBand opportunityActions">
        <button className="secondaryButton" type="button" disabled={contract.status !== 'Pending Signature'} onClick={() => void change('Signed')}>Sign</button>
        <button className="secondaryButton" type="button" disabled={contract.status !== 'Signed'} onClick={() => void change('Active')}>Activate</button>
        <button className="secondaryButton" type="button" disabled={contract.status !== 'Active'} onClick={() => void change('Completed')}>Complete</button>
        <button className="secondaryButton" type="button" disabled={contract.status === 'Completed' || contract.status === 'Terminated'} onClick={() => void change('Terminated')}>Terminate</button>
      </div>
    </section>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
