import { useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract, changeContractStatus } from '../../api/contracts';
import { contractStatusLabel, labelFor } from '../../i18n/labels';

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
      onError(apiError.safeMessage || '请求失败。');
    }
  }

  return (
    <section className="detailPane" aria-label="合同详情">
      <div className="detailHeader">
        <div>
          <h2>{contract.id}</h2>
          <p>状态：{labelFor(contractStatusLabel, contract.status)}</p>
        </div>
        <span className="statusPill">{labelFor(contractStatusLabel, contract.status)}</span>
      </div>
      {pendingReminder && <div role="alert" className="error">待签署合同的预计签署日期已过。</div>}
      <dl className="detailGrid">
        <div>
          <dt>报价</dt>
          <dd>{contract.quoteId}</dd>
        </div>
        <div>
          <dt>商机</dt>
          <dd>{contract.opportunityId}</dd>
        </div>
        <div>
          <dt>客户</dt>
          <dd>{contract.customerId}</dd>
        </div>
        <div>
          <dt>金额</dt>
          <dd>{contract.amount}</dd>
        </div>
        <div>
          <dt>预计签署日期</dt>
          <dd>{contract.expectedSignedDate}</dd>
        </div>
        <div>
          <dt>签署/生效日期</dt>
          <dd>{contract.signedEffectiveDate || '未签署'}</dd>
        </div>
        <div>
          <dt>合同备注</dt>
          <dd>{contract.contractNote}</dd>
        </div>
        <div>
          <dt>金额差异原因</dt>
          <dd>{contract.amountDifferenceReason || '无'}</dd>
        </div>
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
    </section>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
