import { useState } from 'react';
import { CalendarDays, Coins, CreditCard, Hash } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Contract } from '../../api/contracts';
import { ActualPayment, PaymentPlan, createPaymentPlan, recordPayment } from '../../api/payments';
import { DetailHero, DetailStat, StatusPill } from '../../components/CrudScaffold';
import { Panel } from '../../components/ui';
import { labelFor, localizeError, paymentStatusLabel } from '../../i18n/labels';

export function PaymentDetail({ contract, onError, onBack, error }: { contract: Contract; onError: (message: string) => void; onBack?: () => void; error?: string }) {
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
    <main className="content crudPage" aria-label="回款详情">
      <DetailHero
        eyebrow="返回回款列表"
        title={contract.opportunityId}
        subtitle={<><span>合同 {contract.id}</span><span>客户 {contract.customerId}</span><span>负责人 {contract.ownerId}</span></>}
        icon={<CreditCard size={20} aria-hidden="true" />}
        status={<StatusPill tone={paymentTone(payment?.paymentStatus ?? plan?.status ?? 'No plan')}>{labelFor(paymentStatusLabel, payment?.paymentStatus ?? plan?.status ?? 'No plan')}</StatusPill>}
        onBack={onBack ?? (() => undefined)}
        stats={
          <>
            <DetailStat label="合同金额" value={money(contract.amount)} icon={<Coins size={17} aria-hidden="true" />} />
            <DetailStat label="计划到期" value={plan?.dueDate ?? '未设置'} icon={<CalendarDays size={17} aria-hidden="true" />} tone="peach" />
            <DetailStat label="剩余金额" value={money(payment?.remainingAmount ?? contract.amount)} icon={<Coins size={17} aria-hidden="true" />} tone="mint" />
            <DetailStat label="合同版本" value={`v${contract.version}`} icon={<Hash size={17} aria-hidden="true" />} />
          </>
        }
      />
      {error ? <div role="alert" className="error">{error}</div> : null}
      {overdue && <div role="alert" className="error">回款计划已逾期。</div>}
      <section className="detailContentGrid">
        <Panel>
          <div className="sectionHeader"><h2>回款字段</h2><StatusPill tone={paymentTone(plan?.status ?? 'No plan')}>{labelFor(paymentStatusLabel, plan?.status ?? 'No plan')}</StatusPill></div>
          <dl className="detailGrid">
            <div><dt>合同</dt><dd>{contract.id}</dd></div>
            <div><dt>合同金额</dt><dd>{money(contract.amount)}</dd></div>
            <div><dt>计划状态</dt><dd>{labelFor(paymentStatusLabel, plan?.status ?? 'No plan')}</dd></div>
            <div><dt>剩余金额</dt><dd>{money(payment?.remainingAmount ?? contract.amount)}</dd></div>
          </dl>
          <form className="actionBand formFields" onSubmit={submitPlan}>
            <label>计划金额<input value={planForm.dueAmount} onChange={(event) => setPlanForm({ ...planForm, dueAmount: event.target.value })} /></label>
            <label>计划到期日<input type="date" value={planForm.dueDate} onChange={(event) => setPlanForm({ ...planForm, dueDate: event.target.value })} /></label>
            <label>计划币种<input value={planForm.currency} onChange={(event) => setPlanForm({ ...planForm, currency: event.target.value })} /></label>
            <button className="primaryButton" type="submit">保存回款计划</button>
          </form>
          {plan && <p className="inlineNotice">计划状态：{labelFor(paymentStatusLabel, plan.status)}</p>}
        </Panel>
        <Panel>
          <div className="sectionHeader"><h2>登记回款</h2><StatusPill tone={paymentTone(payment?.paymentStatus ?? 'No plan')}>{labelFor(paymentStatusLabel, payment?.paymentStatus ?? 'No plan')}</StatusPill></div>
          <form className="actionBand formFields" onSubmit={submitPayment}>
            <label>回款金额<input value={paymentForm.amount} onChange={(event) => setPaymentForm({ ...paymentForm, amount: event.target.value })} /></label>
            <label>回款日期<input type="date" value={paymentForm.paymentDate} onChange={(event) => setPaymentForm({ ...paymentForm, paymentDate: event.target.value })} /></label>
            <label>幂等键<input value={paymentForm.idempotencyKey} onChange={(event) => setPaymentForm({ ...paymentForm, idempotencyKey: event.target.value })} /></label>
            <label>回款备注<input value={paymentForm.note} onChange={(event) => setPaymentForm({ ...paymentForm, note: event.target.value })} /></label>
            <label>回款币种<input value={paymentForm.currency} onChange={(event) => setPaymentForm({ ...paymentForm, currency: event.target.value })} /></label>
            <button className="primaryButton" type="submit">登记回款</button>
          </form>
          {payment && (
            <div className="inlineNotice">
              <p>回款状态：{labelFor(paymentStatusLabel, payment.paymentStatus)}</p>
              <p>剩余金额：{payment.remainingAmount}</p>
            </div>
          )}
        </Panel>
      </section>
    </main>
  );
}

function paymentTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Paid') return 'success';
  if (status === 'Overdue') return 'danger';
  if (status === 'PartiallyPaid' || status === 'Unpaid' || status === 'Pending') return 'warning';
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
