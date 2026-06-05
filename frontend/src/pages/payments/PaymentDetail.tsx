import { useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract } from '../../api/contracts';
import { ActualPayment, PaymentPlan, createPaymentPlan, recordPayment } from '../../api/payments';
import { labelFor, localizeError, paymentStatusLabel } from '../../i18n/labels';

export function PaymentDetail({ contract, onError }: { contract: Contract; onError: (message: string) => void }) {
  const [plan, setPlan] = useState<PaymentPlan | null>(null);
  const [payment, setPayment] = useState<ActualPayment | null>(null);
  const [planForm, setPlanForm] = useState({ dueAmount: '', dueDate: '', currency: 'CNY' });
  const [paymentForm, setPaymentForm] = useState({ amount: '', paymentDate: '', idempotencyKey: '', note: '', currency: 'CNY' });
  const overdue = plan !== null && plan.status !== 'Paid' && plan.dueDate < today();

  async function submitPlan(event: React.FormEvent) {
    event.preventDefault();
    onError('');
    try {
      setPlan(await createPaymentPlan(contract.id, planForm));
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(localizeError(apiError));
    }
  }

  async function submitPayment(event: React.FormEvent) {
    event.preventDefault();
    onError('');
    try {
      const recorded = await recordPayment(contract.id, paymentForm);
      setPayment(recorded);
      setPaymentForm({ ...paymentForm, amount: '', idempotencyKey: '', note: '' });
      if (plan) {
        setPlan({ ...plan, status: recorded.paymentStatus });
      }
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(localizeError(apiError));
    }
  }

  return (
    <section className="detailPane" aria-label="回款详情">
      <div className="detailHeader">
        <div>
          <h2>{contract.opportunityId}</h2>
          <p>售后回款跟踪</p>
        </div>
        <span className="statusPill">{labelFor(paymentStatusLabel, payment?.paymentStatus ?? plan?.status ?? 'No plan')}</span>
      </div>
      {overdue && <div role="alert" className="error">回款计划已逾期。</div>}
      <dl className="detailGrid">
        <div>
          <dt>合同</dt>
          <dd>{contract.id}</dd>
        </div>
        <div>
          <dt>合同金额</dt>
          <dd>{contract.amount}</dd>
        </div>
        <div>
          <dt>计划状态</dt>
          <dd>{labelFor(paymentStatusLabel, plan?.status ?? 'No plan')}</dd>
        </div>
        <div>
          <dt>剩余金额</dt>
          <dd>{payment?.remainingAmount ?? contract.amount}</dd>
        </div>
      </dl>
      <form className="createPanel" onSubmit={submitPlan}>
        <label>
          计划金额
          <input value={planForm.dueAmount} onChange={(event) => setPlanForm({ ...planForm, dueAmount: event.target.value })} />
        </label>
        <label>
          计划到期日
          <input type="date" value={planForm.dueDate} onChange={(event) => setPlanForm({ ...planForm, dueDate: event.target.value })} />
        </label>
        <label>
          计划币种
          <input value={planForm.currency} onChange={(event) => setPlanForm({ ...planForm, currency: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">保存回款计划</button>
      </form>
      {plan && <p className="inlineNotice">计划状态：{labelFor(paymentStatusLabel, plan.status)}</p>}
      <form className="createPanel" onSubmit={submitPayment}>
        <label>
          回款金额
          <input value={paymentForm.amount} onChange={(event) => setPaymentForm({ ...paymentForm, amount: event.target.value })} />
        </label>
        <label>
          回款日期
          <input type="date" value={paymentForm.paymentDate} onChange={(event) => setPaymentForm({ ...paymentForm, paymentDate: event.target.value })} />
        </label>
        <label>
          幂等键
          <input value={paymentForm.idempotencyKey} onChange={(event) => setPaymentForm({ ...paymentForm, idempotencyKey: event.target.value })} />
        </label>
        <label>
          回款备注
          <input value={paymentForm.note} onChange={(event) => setPaymentForm({ ...paymentForm, note: event.target.value })} />
        </label>
        <label>
          回款币种
          <input value={paymentForm.currency} onChange={(event) => setPaymentForm({ ...paymentForm, currency: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">登记回款</button>
      </form>
      {payment && (
        <div className="inlineNotice">
          <p>回款状态：{labelFor(paymentStatusLabel, payment.paymentStatus)}</p>
          <p>剩余金额：{payment.remainingAmount}</p>
        </div>
      )}
    </section>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
