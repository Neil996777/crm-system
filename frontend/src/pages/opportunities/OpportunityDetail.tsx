import { FormEvent, useState } from 'react';
import { Archive, BriefcaseBusiness, CalendarDays, Coins, GitBranch, Hash, ReceiptText, UserRoundPen } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Opportunity, archiveOpportunity, changeOpportunityStage, closeOpportunityLost, closeOpportunityWon, getOpportunity, updateOpportunity } from '../../api/opportunities';
import { useSession } from '../../auth/SessionProvider';
import { ActivityNoteTaskPanel } from '../../components/ActivityNoteTaskPanel';
import { CloseOpportunityDialog } from '../../components/CloseOpportunityDialog';
import { DetailHero, DetailStat, StatusPill } from '../../components/CrudScaffold';
import { Button, DataTable, Panel, TextField } from '../../components/ui';
import { StageStepper } from '../../components/StageStepper';
import { labelFor, localizeError, lostReasonLabel, opportunityStageLabel } from '../../i18n/labels';

type CloseMode = 'Won' | 'Lost' | null;

export function OpportunityDetail({
  opportunity,
  onUpdated,
  onError,
  onBack,
  onEdit,
  error
}: {
  opportunity: Opportunity;
  onUpdated: (opportunity: Opportunity) => Promise<void>;
  onError: (message: string) => void;
  onBack?: () => void;
  onEdit: (opportunity: Opportunity) => void;
  error?: string;
}) {
  const { user } = useSession();
  const [closeMode, setCloseMode] = useState<CloseMode>(null);
  const [transferOpen, setTransferOpen] = useState(false);
  const [transferOwnerId, setTransferOwnerId] = useState('');
  const [busy, setBusy] = useState(false);
  const terminal = opportunity.stage === 'Won' || opportunity.stage === 'Lost';
  const canManage = user?.role === 'Administrator' || user?.role === 'Sales Manager';

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

  async function archive() {
    if (terminal || !canManage || opportunity.archived) return;
    setBusy(true);
    onError('');
    try {
      await onUpdated(await archiveOpportunity(opportunity.id, opportunity.version, '详情归档商机记录'));
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(localizeError(apiError));
    } finally {
      setBusy(false);
    }
  }

  async function transferOwner(event: FormEvent) {
    event.preventDefault();
    if (terminal || !canManage || transferOwnerId.trim() === '') return;
    setBusy(true);
    onError('');
    try {
      const updated = await updateOpportunity(opportunity.id, {
        customerId: opportunity.customerId,
        ownerId: transferOwnerId.trim(),
        stage: opportunity.stage,
        expectedAmount: opportunity.expectedAmount,
        expectedCloseDate: opportunity.expectedCloseDate,
        title: opportunity.title,
        expectedVersion: opportunity.version
      });
      setTransferOpen(false);
      setTransferOwnerId('');
      await onUpdated(updated);
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(localizeError(apiError));
    } finally {
      setBusy(false);
    }
  }

  return (
    <main className="content crudPage" aria-label="商机详情">
      <DetailHero
        eyebrow="返回商机列表"
        title={opportunity.title || opportunity.id}
        subtitle={
          <>
            <span>当前阶段：{labelFor(opportunityStageLabel, opportunity.stage)}</span>
            <span>客户 {opportunity.customerId}</span>
            <span>负责人 {opportunity.ownerId}</span>
            <span>{opportunity.archived ? '已归档' : '活动记录'}</span>
            <span>更新于 {formatDate(opportunity.updatedAt)}</span>
          </>
        }
        icon={<BriefcaseBusiness size={20} aria-hidden="true" />}
        status={<StatusPill tone={stageTone(opportunity.stage)}>{terminal ? '已关闭记录' : labelFor(opportunityStageLabel, opportunity.stage)}</StatusPill>}
        onBack={onBack ?? (() => undefined)}
        actions={terminal ? (
          <>
            <button className="primaryButton" type="button" disabled>关闭为赢单</button>
            <button className="secondaryButton dangerText" type="button" disabled>关闭为丢单</button>
          </>
        ) : (
          <>
            {canManage ? (
              <>
                <button className="secondaryButton" type="button" onClick={() => onEdit(opportunity)} disabled={busy}>编辑</button>
                <button className="secondaryButton" type="button" onClick={() => { setTransferOpen(true); setTransferOwnerId(opportunity.ownerId); }} disabled={busy}>
                  <UserRoundPen size={15} aria-hidden="true" />
                  转移负责人
                </button>
                <button className="secondaryButton" type="button" onClick={() => void archive()} disabled={busy || opportunity.archived}>
                  <Archive size={15} aria-hidden="true" />
                  归档
                </button>
              </>
            ) : null}
            <button className="primaryButton" type="button" onClick={() => setCloseMode('Won')} disabled={busy}>关闭为赢单</button>
            <button className="secondaryButton dangerText" type="button" onClick={() => setCloseMode('Lost')} disabled={busy}>关闭为丢单</button>
          </>
        )}
        stats={
          <>
            <DetailStat label="预计金额" value={money(opportunity.expectedAmount)} icon={<Coins size={17} aria-hidden="true" />} />
            <DetailStat label="预计签约" value={opportunity.expectedCloseDate || '未填写'} icon={<CalendarDays size={17} aria-hidden="true" />} tone="peach" />
            <DetailStat label="下一合法动作" value={terminal ? '只读' : nextAction(opportunity.stage)} icon={<GitBranch size={17} aria-hidden="true" />} />
            <DetailStat label="赢单合同" value={opportunity.wonContractId || '未关联'} icon={<ReceiptText size={17} aria-hidden="true" />} tone="mint" />
            <DetailStat label="版本" value={`v${opportunity.version}`} icon={<Hash size={17} aria-hidden="true" />} />
          </>
        }
      />
      {error ? <div role="alert" className="error">{error}</div> : null}
      {transferOpen ? (
        <Panel aria-label="转移负责人">
          <form className="actionBand formFields" onSubmit={transferOwner}>
            <TextField
              className="full"
              label="新负责人 ID"
              aria-label="新负责人 ID"
              value={transferOwnerId}
              onChange={(event) => setTransferOwnerId(event.currentTarget.value)}
              hint="经理/管理员可将商机负责人转移给团队成员。"
            />
            <div className="saveBar full">
              <Button onClick={() => setTransferOpen(false)}>取消</Button>
              <Button variant="primary" type="submit" disabled={busy || transferOwnerId.trim() === ''}>确认转移负责人</Button>
            </div>
          </form>
        </Panel>
      ) : null}

      <Panel className="stageCard">
        <div className="sectionHeader">
          <div>
            <h2>阶段步进器</h2>
            <p>线性单向：新商机 → 需求已确认 → 报价 → 合同谈判；终态为二选一结果。</p>
          </div>
          <StatusPill tone={terminal ? stageTone(opportunity.stage) : 'primary'}>
            {terminal ? labelFor(opportunityStageLabel, opportunity.stage) : '进行中'}
          </StatusPill>
        </div>
        <StageStepper currentStage={opportunity.stage} terminal={terminal || busy} onSelectStage={(stage) => void selectStage(stage)} />
        {terminal ? <p className="inlineNotice">终态商机只读，阶段、关闭操作和工作记录新增均已停用。</p> : null}
      </Panel>

      <section className="detailContentGrid">
        <Panel aria-label="关联报价 / 合同 / 回款">
          <div className="sectionHeader">
            <h2>关联报价 / 合同 / 回款</h2>
            <span className="badge">DEC-018</span>
          </div>
          <DataTable
            caption="商机关联记录"
            rows={relatedRows(opportunity)}
            rowKey={(row) => row.id}
            empty="暂无关联报价、合同或回款记录。"
            columns={[
              { key: 'type', header: '类型', render: (row) => row.type },
              { key: 'id', header: '编号', render: (row) => row.id },
              { key: 'status', header: '状态', render: (row) => <StatusPill tone={row.tone}>{row.status}</StatusPill> },
              { key: 'amount', header: '金额', align: 'right', render: (row) => row.amount }
            ]}
          />
          <div className="detailGrid">
            <div>
              <dt>客户</dt>
              <dd>{opportunity.customerId}</dd>
            </div>
            <div>
              <dt>负责人</dt>
              <dd>{opportunity.ownerId}</dd>
            </div>
            {opportunity.lostReasonCode ? (
              <div>
                <dt>丢单原因</dt>
                <dd>{labelFor(lostReasonLabel, opportunity.lostReasonCode)}</dd>
              </div>
            ) : null}
          </div>
        </Panel>
        <Panel aria-label="活动、备注、任务">
          <div className="sectionHeader">
            <h2>活动 / 备注 / 任务</h2>
            <span className="badge primary">实时</span>
          </div>
          <ActivityNoteTaskPanel relatedType="Opportunity" relatedId={opportunity.id} ownerId={opportunity.ownerId} readOnly={terminal} onError={onError} />
        </Panel>
      </section>

      {closeMode && <CloseOpportunityDialog mode={closeMode} onCancel={() => setCloseMode(null)} onConfirm={close} />}
    </main>
  );
}

function relatedRows(opportunity: Opportunity) {
  const rows: Array<{ type: string; id: string; status: string; amount: string; tone: 'primary' | 'success' | 'warning' | 'danger' | 'neutral' }> = [];
  if (opportunity.wonContractId) {
    rows.push({ type: '合同', id: opportunity.wonContractId, status: '已签署', amount: money(opportunity.expectedAmount), tone: 'success' });
  }
  if (opportunity.stage === 'Won') {
    rows.push({ type: '回款计划', id: '待创建', status: '待回款', amount: money(opportunity.expectedAmount), tone: 'warning' });
  }
  return rows;
}

function stageTone(stage: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (stage === 'Won') return 'success';
  if (stage === 'Lost') return 'danger';
  if (stage === 'Contract Negotiation') return 'warning';
  if (stage === 'Quote') return 'primary';
  return 'neutral';
}

function nextAction(stage: string) {
  if (stage === 'New Opportunity') return '需求已确认';
  if (stage === 'Needs Confirmed') return '报价';
  if (stage === 'Quote') return '合同谈判';
  if (stage === 'Contract Negotiation') return '关闭为赢单/丢单';
  return '只读';
}

function money(value: string) {
  const number = Number(value);
  if (!Number.isFinite(number)) return value || '未填写';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY', maximumFractionDigits: 0 }).format(number);
}

function formatDate(value: string) {
  if (!value) return '未更新';
  return value.length > 10 ? value.slice(0, 10) : value;
}
