import { useState } from 'react';
import { ApiError } from '../../api/client';
import { Opportunity, changeOpportunityStage, closeOpportunityLost, closeOpportunityWon, getOpportunity } from '../../api/opportunities';
import { ActivityNoteTaskPanel } from '../../components/ActivityNoteTaskPanel';
import { CloseOpportunityDialog } from '../../components/CloseOpportunityDialog';
import { StageStepper } from '../../components/StageStepper';
import { labelFor, localizeError, lostReasonLabel, opportunityStageLabel } from '../../i18n/labels';

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
      onError(localizeError(apiError));
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
      onError(localizeError(apiError));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="detailPane" aria-label="商机详情">
      <div className="detailHeader">
        <div>
          <h2>{opportunity.title || opportunity.id}</h2>
          <p>当前阶段：{labelFor(opportunityStageLabel, opportunity.stage)}</p>
        </div>
        <span className="statusPill">{terminal ? '已关闭记录' : labelFor(opportunityStageLabel, opportunity.stage)}</span>
      </div>

      <StageStepper currentStage={opportunity.stage} terminal={terminal || busy} onSelectStage={(stage) => void selectStage(stage)} />

      <dl className="detailGrid">
        <div>
          <dt>客户</dt>
          <dd>{opportunity.customerId}</dd>
        </div>
        <div>
          <dt>负责人</dt>
          <dd>{opportunity.ownerId}</dd>
        </div>
        <div>
          <dt>预计金额</dt>
          <dd>{opportunity.expectedAmount}</dd>
        </div>
        <div>
          <dt>预计关闭日期</dt>
          <dd>{opportunity.expectedCloseDate}</dd>
        </div>
        {opportunity.wonContractId && (
          <div>
            <dt>赢单合同</dt>
            <dd>{opportunity.wonContractId}</dd>
          </div>
        )}
        {opportunity.lostReasonCode && (
          <div>
            <dt>丢单原因</dt>
            <dd>{labelFor(lostReasonLabel, opportunity.lostReasonCode)}</dd>
          </div>
        )}
      </dl>

      <div className="actionBand opportunityActions">
        <button className="secondaryButton" type="button" onClick={() => setCloseMode('Won')} disabled={terminal || busy}>关闭为赢单</button>
        <button className="secondaryButton" type="button" onClick={() => setCloseMode('Lost')} disabled={terminal || busy}>关闭为丢单</button>
      </div>

      {closeMode && <CloseOpportunityDialog mode={closeMode} onCancel={() => setCloseMode(null)} onConfirm={close} />}
      <ActivityNoteTaskPanel relatedType="Opportunity" relatedId={opportunity.id} ownerId={opportunity.ownerId} onError={onError} />
    </section>
  );
}
