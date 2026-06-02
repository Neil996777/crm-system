import { useState } from 'react';
import { ApiError } from '../../api/client';
import { Opportunity, changeOpportunityStage, closeOpportunityLost, closeOpportunityWon, getOpportunity } from '../../api/opportunities';
import { ActivityNoteTaskPanel } from '../../components/ActivityNoteTaskPanel';
import { CloseOpportunityDialog } from '../../components/CloseOpportunityDialog';
import { StageStepper } from '../../components/StageStepper';

type CloseMode = 'Won' | 'Lost' | null;

export function OpportunityDetail({
  opportunity,
  onUpdated,
  onError
}: {
  opportunity: Opportunity;
  onUpdated: (opportunity: Opportunity) => Promise<void>;
  onError: (message: string) => void;
}) {
  const [closeMode, setCloseMode] = useState<CloseMode>(null);
  const [busy, setBusy] = useState(false);
  const terminal = opportunity.stage === 'Won' || opportunity.stage === 'Lost';

  async function refresh() {
    await onUpdated(await getOpportunity(opportunity.id));
  }

  async function selectStage(stage: string) {
    if (terminal || stage === opportunity.stage) return;
    setBusy(true);
    onError('');
    try {
      const updated = await changeOpportunityStage(opportunity.id, opportunity.version, stage);
      await onUpdated(updated);
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(apiError.safeMessage || 'Request failed.');
      await refresh();
    } finally {
      setBusy(false);
    }
  }

  async function close(input: { contractId: string; closeDate: string; lostReason: { code: string; detail: string } }) {
    if (!closeMode) return;
    setBusy(true);
    onError('');
    try {
      if (closeMode === 'Won') {
        await closeOpportunityWon(opportunity.id, { expectedVersion: opportunity.version, contractId: input.contractId, closeDate: input.closeDate });
      } else {
        await closeOpportunityLost(opportunity.id, { expectedVersion: opportunity.version, closeDate: input.closeDate, lostReason: input.lostReason });
      }
      setCloseMode(null);
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(apiError.safeMessage || 'Request failed.');
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="detailPane" aria-label="Opportunity detail">
      <div className="detailHeader">
        <div>
          <h2>{opportunity.title || opportunity.id}</h2>
          <p>Current stage: {opportunity.stage}</p>
        </div>
        <span className="statusPill">{terminal ? 'Terminal record' : opportunity.stage}</span>
      </div>

      <StageStepper currentStage={opportunity.stage} terminal={terminal || busy} onSelectStage={(stage) => void selectStage(stage)} />

      <dl className="detailGrid">
        <div>
          <dt>Customer</dt>
          <dd>{opportunity.customerId}</dd>
        </div>
        <div>
          <dt>Owner</dt>
          <dd>{opportunity.ownerId}</dd>
        </div>
        <div>
          <dt>Expected amount</dt>
          <dd>{opportunity.expectedAmount}</dd>
        </div>
        <div>
          <dt>Expected close</dt>
          <dd>{opportunity.expectedCloseDate}</dd>
        </div>
        {opportunity.wonContractId && (
          <div>
            <dt>Won contract</dt>
            <dd>{opportunity.wonContractId}</dd>
          </div>
        )}
        {opportunity.lostReasonCode && (
          <div>
            <dt>Lost reason</dt>
            <dd>{opportunity.lostReasonCode}</dd>
          </div>
        )}
      </dl>

      <div className="actionBand opportunityActions">
        <button className="secondaryButton" type="button" onClick={() => setCloseMode('Won')} disabled={terminal || busy}>Close Won</button>
        <button className="secondaryButton" type="button" onClick={() => setCloseMode('Lost')} disabled={terminal || busy}>Close Lost</button>
      </div>

      {closeMode && <CloseOpportunityDialog mode={closeMode} onCancel={() => setCloseMode(null)} onConfirm={close} />}
      <ActivityNoteTaskPanel relatedType="Opportunity" relatedId={opportunity.id} ownerId={opportunity.ownerId} onError={onError} />
    </section>
  );
}
