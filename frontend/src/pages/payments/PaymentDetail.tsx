import { useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract } from '../../api/contracts';
import { ActualPayment, PaymentPlan, createPaymentPlan, recordPayment } from '../../api/payments';

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
      onError(apiError.safeMessage || 'Request failed.');
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
      onError(apiError.safeMessage || 'Request failed.');
    }
  }

  return (
    <section className="detailPane" aria-label="Payment detail">
      <div className="detailHeader">
        <div>
          <h2>{contract.opportunityId}</h2>
          <p>Post-sale payment tracking</p>
        </div>
        <span className="statusPill">{payment?.paymentStatus ?? plan?.status ?? 'No plan'}</span>
      </div>
      {overdue && <div role="alert" className="error">Payment plan is overdue.</div>}
      <dl className="detailGrid">
        <div>
          <dt>Contract</dt>
          <dd>{contract.id}</dd>
        </div>
        <div>
          <dt>Contract amount</dt>
          <dd>{contract.amount}</dd>
        </div>
        <div>
          <dt>Plan status</dt>
          <dd>{plan?.status ?? 'No plan'}</dd>
        </div>
        <div>
          <dt>Remaining amount</dt>
          <dd>{payment?.remainingAmount ?? contract.amount}</dd>
        </div>
      </dl>
      <form className="createPanel" onSubmit={submitPlan}>
        <label>
          Plan amount
          <input value={planForm.dueAmount} onChange={(event) => setPlanForm({ ...planForm, dueAmount: event.target.value })} />
        </label>
        <label>
          Plan due date
          <input type="date" value={planForm.dueDate} onChange={(event) => setPlanForm({ ...planForm, dueDate: event.target.value })} />
        </label>
        <label>
          Plan currency
          <input value={planForm.currency} onChange={(event) => setPlanForm({ ...planForm, currency: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">Save payment plan</button>
      </form>
      {plan && <p className="inlineNotice">Plan status: {plan.status}</p>}
      <form className="createPanel" onSubmit={submitPayment}>
        <label>
          Payment amount
          <input value={paymentForm.amount} onChange={(event) => setPaymentForm({ ...paymentForm, amount: event.target.value })} />
        </label>
        <label>
          Payment date
          <input type="date" value={paymentForm.paymentDate} onChange={(event) => setPaymentForm({ ...paymentForm, paymentDate: event.target.value })} />
        </label>
        <label>
          Idempotency key
          <input value={paymentForm.idempotencyKey} onChange={(event) => setPaymentForm({ ...paymentForm, idempotencyKey: event.target.value })} />
        </label>
        <label>
          Payment note
          <input value={paymentForm.note} onChange={(event) => setPaymentForm({ ...paymentForm, note: event.target.value })} />
        </label>
        <label>
          Payment currency
          <input value={paymentForm.currency} onChange={(event) => setPaymentForm({ ...paymentForm, currency: event.target.value })} />
        </label>
        <button className="primaryButton" type="submit">Record payment</button>
      </form>
      {payment && (
        <div className="inlineNotice">
          <p>Payment status: {payment.paymentStatus}</p>
          <p>Remaining amount: {payment.remainingAmount}</p>
        </div>
      )}
    </section>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}
