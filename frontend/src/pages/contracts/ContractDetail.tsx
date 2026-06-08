import { useState } from 'react';
import { CalendarDays, Coins, FileSignature, Hash } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Contract, changeContractStatus } from '../../api/contracts';
import { DetailHero, DetailStat, StatusPill } from '../../components/CrudScaffold';
import { Panel } from '../../components/ui';
import { contractStatusLabel, labelFor, localizeError } from '../../i18n/labels';

export function ContractDetail({ contract, onUpdated, onError, onBack, error }: { contract: Contract; onUpdated: (contract: Contract) => Promise<void>; onError: (message: string) => void; onBack?: () => void; error?: string }) {
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
      onError(localizeError(apiError));
    }
  }

  return (
    <main className="content crudPage" aria-label="合同详情">
      <DetailHero
        eyebrow="返回合同列表"
        title={contract.id}
        subtitle={<><span>商机 {contract.opportunityId}</span><span>客户 {contract.customerId}</span><span>负责人 {contract.ownerId}</span></>}
        icon={<FileSignature size={20} aria-hidden="true" />}
        status={<StatusPill tone={contractTone(contract.status)}>{labelFor(contractStatusLabel, contract.status)}</StatusPill>}
        onBack={onBack ?? (() => undefined)}
        stats={
          <>
            <DetailStat label="金额" value={money(contract.amount)} icon={<Coins size={17} aria-hidden="true" />} />
            <DetailStat label="预计签署" value={contract.expectedSignedDate} icon={<CalendarDays size={17} aria-hidden="true" />} tone="peach" />
            <DetailStat label="签署/生效" value={contract.signedEffectiveDate || '未签署'} icon={<CalendarDays size={17} aria-hidden="true" />} tone="mint" />
            <DetailStat label="版本" value={`v${contract.version}`} icon={<Hash size={17} aria-hidden="true" />} />
          </>
        }
      />
      {error ? <div role="alert" className="error">{error}</div> : null}
      {pendingReminder && <div role="alert" className="error">待签署合同的预计签署日期已过。</div>}
      <p className="inlineNotice">状态：{labelFor(contractStatusLabel, contract.status)}</p>
      <section className="detailContentGrid">
        <Panel>
          <div className="sectionHeader"><h2>合同字段</h2><StatusPill tone={contractTone(contract.status)}>{labelFor(contractStatusLabel, contract.status)}</StatusPill></div>
          <dl className="detailGrid">
            <div><dt>报价</dt><dd>{contract.quoteId}</dd></div>
            <div><dt>商机</dt><dd>{contract.opportunityId}</dd></div>
            <div><dt>客户</dt><dd>{contract.customerId}</dd></div>
            <div><dt>金额</dt><dd>{money(contract.amount)}</dd></div>
            <div><dt>预计签署日期</dt><dd>{contract.expectedSignedDate}</dd></div>
            <div><dt>签署/生效日期</dt><dd>{contract.signedEffectiveDate || '未签署'}</dd></div>
            <div><dt>合同备注</dt><dd>{contract.contractNote || '无'}</dd></div>
            <div><dt>金额差异原因</dt><dd>{contract.amountDifferenceReason || '无'}</dd></div>
          </dl>
          <label className="singleField">
            签署/生效日期
            <input type="date" value={signedEffectiveDate} onChange={(event) => setSignedEffectiveDate(event.target.value)} />
          </label>
          <div className="actionBand opportunityActions">
            <button className="secondaryButton" type="button" disabled={contract.status !== 'Pending Signature'} onClick={() => void change('Signed')}>签署</button>
            <button className="secondaryButton" type="button" disabled={contract.status !== 'Signed'} onClick={() => void change('Active')}>启用</button>
            <button className="secondaryButton" type="button" disabled={contract.status !== 'Active'} onClick={() => void change('Completed')}>完成</button>
            <button className="secondaryButton" type="button" disabled={contract.status === 'Completed' || contract.status === 'Terminated'} onClick={() => void change('Terminated')}>终止</button>
          </div>
        </Panel>
      </section>
    </main>
  );
}

function contractTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Signed' || status === 'Active' || status === 'Completed') return 'success';
  if (status === 'Terminated') return 'danger';
  if (status === 'Pending Signature') return 'warning';
  return 'neutral';
}

function money(value: string) {
  const number = Number(value);
  if (!Number.isFinite(number)) return value || '未填写';
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY', maximumFractionDigits: 0 }).format(number);
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
