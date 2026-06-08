import { Building2, GitBranch, Hash, ListChecks, Mail, UserRound } from 'lucide-react';
import { Lead, ConversionResult } from '../../api/leads';
import { ConvertLeadDialog } from '../../components/ConvertLeadDialog';
import { DetailHero, DetailStat, StatusPill } from '../../components/CrudScaffold';
import { HistoryTimeline } from '../../components/HistoryTimeline';
import { QualificationActions } from '../../components/QualificationActions';
import { Panel } from '../../components/ui';
import { labelFor, leadStatusLabel } from '../../i18n/labels';

type Props = {
  lead: Lead;
  onUpdated: (lead: Lead) => void;
  onConverted: (result: ConversionResult) => void;
  onError: (message: string) => void;
  onBack?: () => void;
  error?: string;
};

export function LeadDetail({ lead, onUpdated, onConverted, onError, onBack, error }: Props) {
  return (
    <section aria-label="线索详情" data-record-id={lead.id}>
      <DetailHero
        eyebrow="返回线索列表"
        title={lead.companyName || lead.leadName || lead.id}
        subtitle={
          <>
            <span>{lead.source || '未填写来源'}</span>
            <span>负责人 {lead.ownerId || '未分配'}</span>
            <span>更新于 {formatDate(lead.updatedAt)}</span>
          </>
        }
        icon={<ListChecks size={20} aria-hidden="true" />}
        status={<StatusPill tone={leadTone(lead.status)}>{labelFor(leadStatusLabel, lead.status)}</StatusPill>}
        onBack={onBack ?? (() => undefined)}
        stats={
          <>
            <DetailStat label="负责人" value={lead.ownerId || labelFor(leadStatusLabel, 'Unassigned')} icon={<UserRound size={17} aria-hidden="true" />} />
            <DetailStat label="联系方式" value={lead.email || lead.phone || '未填写'} icon={<Mail size={17} aria-hidden="true" />} tone="peach" />
            <DetailStat label="转化客户" value={lead.convertedAccountId || '未转化'} icon={<Building2 size={17} aria-hidden="true" />} tone="mint" />
            <DetailStat label="转化商机" value={lead.convertedOpportunityId || '未转化'} icon={<GitBranch size={17} aria-hidden="true" />} />
            <DetailStat label="版本" value={`v${lead.version}`} icon={<Hash size={17} aria-hidden="true" />} />
          </>
        }
      />
      {error ? <div role="alert" className="error">{error}</div> : null}
      <section className="detailContentGrid">
        <Panel>
          <div className="sectionHeader">
            <h2>确认与转换</h2>
            <StatusPill tone={leadTone(lead.status)}>{labelFor(leadStatusLabel, lead.status)}</StatusPill>
          </div>
          <dl className="detailGrid">
            <div>
              <dt>线索名称</dt>
              <dd>{lead.leadName || '未填写'}</dd>
            </div>
            <div>
              <dt>公司名称</dt>
              <dd>{lead.companyName || '未填写'}</dd>
            </div>
            <div>
              <dt>需求摘要</dt>
              <dd>{lead.needSummary || '无'}</dd>
            </div>
            {lead.invalidReason ? (
              <div>
                <dt>无效原因</dt>
                <dd>{lead.invalidReason}</dd>
              </div>
            ) : null}
          </dl>
          {lead.status === 'Converted To Opportunity' ? <p className="inlineNotice">已转换线索在确认环节为只读。</p> : <QualificationActions lead={lead} onUpdated={onUpdated} onError={onError} />}
          <ConvertLeadDialog lead={lead} onConverted={onConverted} onError={onError} />
        </Panel>
        <Panel>
          <div className="sectionHeader">
            <h2>历史记录</h2>
            <span className="badge primary">安全摘要</span>
          </div>
          <HistoryTimeline resource="leads" recordId={lead.id} reloadKey={lead.version} />
        </Panel>
      </section>
    </section>
  );
}

function leadTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Valid' || status === 'Converted To Opportunity') return 'success';
  if (status === 'Invalid') return 'danger';
  if (status === 'Pending Qualification') return 'warning';
  return 'neutral';
}

function formatDate(value: string) {
  if (!value) return '未更新';
  return value.length > 10 ? value.slice(0, 10) : value;
}
