import { ApiError } from '../../api/client';
import { Quote, changeQuoteStatus } from '../../api/quotes';

export function QuoteDetail({ quote, onUpdated, onError }: { quote: Quote; onUpdated: (quote: Quote) => Promise<void>; onError: (message: string) => void }) {
  const expired = quote.status === 'Expired';
  const accepted = quote.status === 'Accepted';

  async function change(toStatus: string) {
    onError('');
    try {
      await onUpdated(await changeQuoteStatus(quote.id, quote.version, toStatus));
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(apiError.safeMessage || 'Request failed.');
    }
  }

  return (
    <section className="detailPane" aria-label="Quote detail">
      <div className="detailHeader">
        <div>
          <h2>{quote.opportunityId}</h2>
          <p>Status: {quote.status}</p>
        </div>
        <span className="statusPill">{quote.status}</span>
      </div>
      {expired && <div role="alert" className="error">Expired quote cannot be linked to a contract.</div>}
      <dl className="detailGrid">
        <div>
          <dt>Opportunity</dt>
          <dd>{quote.opportunityId}</dd>
        </div>
        <div>
          <dt>Customer</dt>
          <dd>{quote.customerId}</dd>
        </div>
        <div>
          <dt>Amount</dt>
          <dd>{quote.amount}</dd>
        </div>
        <div>
          <dt>Validity end</dt>
          <dd>{quote.validityEnd}</dd>
        </div>
      </dl>
      <div className="actionBand opportunityActions">
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Draft'} onClick={() => void change('Sent')}>Send</button>
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Sent'} onClick={() => void change('Accepted')}>Accept</button>
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Sent'} onClick={() => void change('Rejected')}>Reject</button>
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Draft' && quote.status !== 'Sent'} onClick={() => void change('Expired')}>Expire</button>
      </div>
      <p className="inlineNotice">{accepted ? 'Contract link available' : expired ? 'Contract link blocked' : 'Contract link pending acceptance'}</p>
    </section>
  );
}
