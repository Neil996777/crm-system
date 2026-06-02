import { Lead, ConversionResult } from '../../api/leads';
import { ConvertLeadDialog } from '../../components/ConvertLeadDialog';
import { HistoryTimeline } from '../../components/HistoryTimeline';
import { QualificationActions } from '../../components/QualificationActions';

type Props = {
  lead: Lead;
  onUpdated: (lead: Lead) => void;
  onConverted: (result: ConversionResult) => void;
  onError: (message: string) => void;
};

export function LeadDetail({ lead, onUpdated, onConverted, onError }: Props) {
  return (
    <section className="detailPane" aria-label="Lead detail" data-record-id={lead.id}>
      <div className="detailHeader">
        <div>
          <h2>{lead.companyName || lead.leadName}</h2>
          <p>{lead.source}</p>
        </div>
        <span className="statusPill">{lead.status}</span>
      </div>
      <dl className="detailGrid">
        <div>
          <dt>Owner</dt>
          <dd>{lead.ownerId || 'Unassigned'}</dd>
        </div>
        <div>
          <dt>Version</dt>
          <dd>{lead.version}</dd>
        </div>
        {lead.invalidReason && (
          <div>
            <dt>Invalid reason</dt>
            <dd>{lead.invalidReason}</dd>
          </div>
        )}
        {lead.convertedOpportunityId && (
          <div>
            <dt>Opportunity</dt>
            <dd>Opportunity: {lead.convertedOpportunityId}</dd>
          </div>
        )}
        {lead.convertedAccountId && (
          <div>
            <dt>Account</dt>
            <dd>Account: {lead.convertedAccountId}</dd>
          </div>
        )}
      </dl>
      {lead.status === 'Converted To Opportunity' ? <p className="inlineNotice">Converted leads are read-only for qualification.</p> : <QualificationActions lead={lead} onUpdated={onUpdated} onError={onError} />}
      <ConvertLeadDialog lead={lead} onConverted={onConverted} onError={onError} />
      <HistoryTimeline resource="leads" recordId={lead.id} reloadKey={lead.version} />
    </section>
  );
}
