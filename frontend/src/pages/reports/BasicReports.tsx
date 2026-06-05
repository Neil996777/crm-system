import { useEffect, useState } from 'react';
import { BasicReport, getBasicReport, GroupRow, PaymentGroupRow } from '../../api/reports';
import { contractStatusLabel, labelFor, leadStatusLabel, localizeError, opportunityStageLabel, paymentStatusLabel, quoteStatusLabel, reportArchiveFilterLabel, reportScopeLabel } from '../../i18n/labels';

export function BasicReports() {
  const [report, setReport] = useState<BasicReport | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh() {
    setError('');
    try {
      setReport(await getBasicReport());
    } catch (err) {
      setError(localizeError(err as { safeMessage?: string }, '无法加载基础报表。'));
    }
  }

  return (
    <section className="reportSection" aria-label="基础销售报表">
      <div className="sectionHeader">
        <div>
          <h2>基础销售报表</h2>
          <p>{report ? `${labelFor(reportScopeLabel, report.scope)}范围 · ${labelFor(reportArchiveFilterLabel, report.filters.archived)}` : '已持久化的授权记录'}</p>
        </div>
        <button className="secondaryButton" type="button" onClick={() => void refresh()}>
          刷新
        </button>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      {report?.emptyState ? <p className="emptyState">暂无报表记录。</p> : null}
      {report ? (
        <>
          <section className="metricGrid" aria-label="基础报表指标">
            <Metric label="线索" value={report.metrics.leadCount} />
            <Metric label="商机" value={report.metrics.opportunityCount} />
            <Metric label="报价" value={`${report.metrics.quoteAmount} ${report.currency}`} />
            <Metric label="合同" value={`${report.metrics.contractAmount} ${report.currency}`} />
            <Metric label="已回款" value={`${report.metrics.paidAmount} ${report.currency}`} />
            <Metric label="应收" value={`${report.metrics.receivableAmount} ${report.currency}`} />
            <Metric label="赢单" value={report.metrics.wonCount} />
            <Metric label="丢单" value={report.metrics.lostCount} />
          </section>
          <div className="reportTables">
            <GroupedTable title="按状态统计线索" firstColumn="状态" rows={report.breakdowns.leadsByStatus} currency={report.currency} labels={leadStatusLabel} />
            <GroupedTable title="按阶段统计商机" firstColumn="阶段" rows={report.breakdowns.opportunitiesByStage} currency={report.currency} labels={opportunityStageLabel} />
            <GroupedTable title="按状态统计报价" firstColumn="状态" rows={report.breakdowns.quotesByStatus} currency={report.currency} labels={quoteStatusLabel} />
            <GroupedTable title="按状态统计合同" firstColumn="状态" rows={report.breakdowns.contractsByStatus} currency={report.currency} labels={contractStatusLabel} />
            <PaymentTable rows={report.breakdowns.paymentsByStatus} currency={report.currency} />
          </div>
        </>
      ) : null}
    </section>
  );
}

function Metric({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="metricTile">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function GroupedTable({ title, firstColumn, rows, currency, labels }: { title: string; firstColumn: string; rows: GroupRow[]; currency: string; labels: Record<string, string> }) {
  return (
    <section className="listPanel">
      <div className="sectionTitle">{title}</div>
      {rows.length === 0 ? (
        <p className="emptyState">暂无数据行。</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>{firstColumn}</th>
              <th>数量</th>
              <th>金额</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.key}>
                <td>{labelFor(labels, row.label)}</td>
                <td>{row.count}</td>
                <td>{row.amount} {currency}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}

function PaymentTable({ rows, currency }: { rows: PaymentGroupRow[]; currency: string }) {
  return (
    <section className="listPanel">
      <div className="sectionTitle">按状态统计回款</div>
      {rows.length === 0 ? (
        <p className="emptyState">暂无数据行。</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>状态</th>
              <th>数量</th>
              <th>金额</th>
              <th>应收</th>
              <th>已回款</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.key}>
                <td>{labelFor(paymentStatusLabel, row.label)}</td>
                <td>{row.count}</td>
                <td>{row.amount} {currency}</td>
                <td>{row.dueAmount} {currency}</td>
                <td>{row.paidAmount} {currency}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}
