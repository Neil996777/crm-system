import { ApiError } from '../../api/client';
import { Quote, changeQuoteStatus } from '../../api/quotes';
import { labelFor, localizeError, quoteStatusLabel } from '../../i18n/labels';

export function QuoteDetail({ quote, onUpdated, onError }: { quote: Quote; onUpdated: (quote: Quote) => Promise<void>; onError: (message: string) => void }) {
  const expired = quote.status === 'Expired';
  const accepted = quote.status === 'Accepted';

  async function change(toStatus: string) {
    onError('');
    try {
      await onUpdated(await changeQuoteStatus(quote.id, quote.version, toStatus));
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(localizeError(apiError));
    }
  }

  return (
    <section className="detailPane" aria-label="报价详情">
      <div className="detailHeader">
        <div>
          <h2>{quote.opportunityId}</h2>
          <p>状态：{labelFor(quoteStatusLabel, quote.status)}</p>
        </div>
        <span className="statusPill">{labelFor(quoteStatusLabel, quote.status)}</span>
      </div>
      {expired && <div role="alert" className="error">已过期报价不能关联合同。</div>}
      <dl className="detailGrid">
        <div>
          <dt>商机</dt>
          <dd>{quote.opportunityId}</dd>
        </div>
        <div>
          <dt>客户</dt>
          <dd>{quote.customerId}</dd>
        </div>
        <div>
          <dt>金额</dt>
          <dd>{quote.amount}</dd>
        </div>
        <div>
          <dt>有效期截止日</dt>
          <dd>{quote.validityEnd}</dd>
        </div>
      </dl>
      <div className="actionBand opportunityActions">
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Draft'} onClick={() => void change('Sent')}>发送</button>
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Sent'} onClick={() => void change('Accepted')}>接受</button>
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Sent'} onClick={() => void change('Rejected')}>拒绝</button>
        <button className="secondaryButton" type="button" disabled={quote.status !== 'Draft' && quote.status !== 'Sent'} onClick={() => void change('Expired')}>标记过期</button>
      </div>
      <p className="inlineNotice">{accepted ? '可关联合同' : expired ? '禁止关联合同' : '待接受后关联合同'}</p>
    </section>
  );
}
