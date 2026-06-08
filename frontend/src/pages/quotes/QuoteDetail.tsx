import { CalendarDays, Coins, FileText, Hash } from 'lucide-react';
import { ApiError } from '../../api/client';
import { Quote, changeQuoteStatus } from '../../api/quotes';
import { DetailHero, DetailStat, StatusPill } from '../../components/CrudScaffold';
import { Panel } from '../../components/ui';
import { labelFor, localizeError, quoteStatusLabel } from '../../i18n/labels';

export function QuoteDetail({ quote, onUpdated, onError, onBack, error }: { quote: Quote; onUpdated: (quote: Quote) => Promise<void>; onError: (message: string) => void; onBack?: () => void; error?: string }) {
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
    <main className="content crudPage" aria-label="报价详情">
      <DetailHero
        eyebrow="返回报价列表"
        title={quote.opportunityId}
        subtitle={<><span>客户 {quote.customerId}</span><span>负责人 {quote.ownerId}</span><span>更新于 {formatDate(quote.updatedAt)}</span></>}
        icon={<FileText size={20} aria-hidden="true" />}
        status={<StatusPill tone={quoteTone(quote.status)}>{labelFor(quoteStatusLabel, quote.status)}</StatusPill>}
        onBack={onBack ?? (() => undefined)}
        stats={
          <>
            <DetailStat label="金额" value={money(quote.amount)} icon={<Coins size={17} aria-hidden="true" />} />
            <DetailStat label="有效期截止日" value={quote.validityEnd} icon={<CalendarDays size={17} aria-hidden="true" />} tone="peach" />
            <DetailStat label="版本" value={`v${quote.version}`} icon={<Hash size={17} aria-hidden="true" />} />
          </>
        }
      />
      {error ? <div role="alert" className="error">{error}</div> : null}
      {expired && <div role="alert" className="error">已过期报价不能关联合同。</div>}
      <p className="inlineNotice">状态：{labelFor(quoteStatusLabel, quote.status)}</p>
      <section className="detailContentGrid">
        <Panel>
          <div className="sectionHeader"><h2>报价字段</h2><StatusPill tone={quoteTone(quote.status)}>{labelFor(quoteStatusLabel, quote.status)}</StatusPill></div>
          <dl className="detailGrid">
            <div><dt>商机</dt><dd>{quote.opportunityId}</dd></div>
            <div><dt>客户</dt><dd>{quote.customerId}</dd></div>
            <div><dt>金额</dt><dd>{money(quote.amount)}</dd></div>
            <div><dt>有效期截止日</dt><dd>{quote.validityEnd}</dd></div>
          </dl>
          <div className="actionBand opportunityActions">
            <button className="secondaryButton" type="button" disabled={quote.status !== 'Draft'} onClick={() => void change('Sent')}>发送</button>
            <button className="secondaryButton" type="button" disabled={quote.status !== 'Sent'} onClick={() => void change('Accepted')}>接受</button>
            <button className="secondaryButton" type="button" disabled={quote.status !== 'Sent'} onClick={() => void change('Rejected')}>拒绝</button>
            <button className="secondaryButton" type="button" disabled={quote.status !== 'Draft' && quote.status !== 'Sent'} onClick={() => void change('Expired')}>标记过期</button>
          </div>
          <p className="inlineNotice">{accepted ? '可关联合同' : expired ? '禁止关联合同' : '待接受后关联合同'}</p>
        </Panel>
      </section>
    </main>
  );
}

function quoteTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Accepted') return 'success';
  if (status === 'Rejected' || status === 'Expired') return 'danger';
  if (status === 'Sent') return 'primary';
  return 'neutral';
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
