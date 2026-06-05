import { Lead, ConversionResult } from '../../api/leads';
import { ConvertLeadDialog } from '../../components/ConvertLeadDialog';
import { HistoryTimeline } from '../../components/HistoryTimeline';
import { QualificationActions } from '../../components/QualificationActions';
import { labelFor, leadStatusLabel } from '../../i18n/labels';

type Props = {
  lead: Lead;
  onUpdated: (lead: Lead) => void;
  onConverted: (result: ConversionResult) => void;
  onError: (message: string) => void;
};

export function LeadDetail({ lead, onUpdated, onConverted, onError }: Props) {
  return (
    <section className="detailPane" aria-label="线索详情" data-record-id={lead.id}>
      <div className="detailHeader">
        <div>
          <h2>{lead.companyName || lead.leadName}</h2>
          <p>{lead.source}</p>
        </div>
        <span className="statusPill">{labelFor(leadStatusLabel, lead.status)}</span>
      </div>
      <dl className="detailGrid">
        <div>
          <dt>负责人</dt>
          <dd>{lead.ownerId || labelFor(leadStatusLabel, 'Unassigned')}</dd>
        </div>
        <div>
          <dt>版本</dt>
          <dd>{lead.version}</dd>
        </div>
        {lead.invalidReason && (
          <div>
            <dt>无效原因</dt>
            <dd>{lead.invalidReason}</dd>
          </div>
        )}
        {lead.convertedOpportunityId && (
          <div>
            <dt>商机</dt>
            <dd>商机：{lead.convertedOpportunityId}</dd>
          </div>
        )}
        {lead.convertedAccountId && (
          <div>
            <dt>客户</dt>
            <dd>客户：{lead.convertedAccountId}</dd>
          </div>
        )}
      </dl>
      {lead.status === 'Converted To Opportunity' ? <p className="inlineNotice">已转换线索在确认环节为只读。</p> : <QualificationActions lead={lead} onUpdated={onUpdated} onError={onError} />}
      <ConvertLeadDialog lead={lead} onConverted={onConverted} onError={onError} />
      <HistoryTimeline resource="leads" recordId={lead.id} reloadKey={lead.version} />
    </section>
  );
}
