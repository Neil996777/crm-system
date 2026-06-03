import { useEffect, useState } from 'react';
import { BasicReport, getBasicReport, GroupRow, PaymentGroupRow } from '../../api/reports';

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
      const safe = err as { safeMessage?: string };
      setError(safe.safeMessage ?? 'Unable to load basic reports.');
    }
  }

  return (
    <section className="reportSection" aria-label="Basic sales reports">
      <div className="sectionHeader">
        <div>
          <h2>Basic Sales Reports</h2>
          <p>{report ? `${report.scope} scope · ${report.filters.archived}` : 'Persisted authorized records'}</p>
        </div>
        <button className="secondaryButton" type="button" onClick={() => void refresh()}>
          Refresh
        </button>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      {report?.emptyState ? <p className="emptyState">No report records yet.</p> : null}
      {report ? (
        <>
          <section className="metricGrid" aria-label="Basic report metrics">
            <Metric label="Leads" value={report.metrics.leadCount} />
            <Metric label="Opportunities" value={report.metrics.opportunityCount} />
            <Metric label="Quotes" value={`${report.metrics.quoteAmount} ${report.currency}`} />
            <Metric label="Contracts" value={`${report.metrics.contractAmount} ${report.currency}`} />
            <Metric label="Paid" value={`${report.metrics.paidAmount} ${report.currency}`} />
            <Metric label="Receivable" value={`${report.metrics.receivableAmount} ${report.currency}`} />
            <Metric label="Won" value={report.metrics.wonCount} />
            <Metric label="Lost" value={report.metrics.lostCount} />
          </section>
          <div className="reportTables">
            <GroupedTable title="Leads by Status" firstColumn="Status" rows={report.breakdowns.leadsByStatus} currency={report.currency} />
            <GroupedTable title="Opportunities by Stage" firstColumn="Stage" rows={report.breakdowns.opportunitiesByStage} currency={report.currency} />
            <GroupedTable title="Quotes by Status" firstColumn="Status" rows={report.breakdowns.quotesByStatus} currency={report.currency} />
            <GroupedTable title="Contracts by Status" firstColumn="Status" rows={report.breakdowns.contractsByStatus} currency={report.currency} />
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

function GroupedTable({ title, firstColumn, rows, currency }: { title: string; firstColumn: string; rows: GroupRow[]; currency: string }) {
  return (
    <section className="listPanel">
      <div className="sectionTitle">{title}</div>
      {rows.length === 0 ? (
        <p className="emptyState">No rows.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>{firstColumn}</th>
              <th>Count</th>
              <th>Amount</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.key}>
                <td>{row.label}</td>
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
      <div className="sectionTitle">Payments by Status</div>
      {rows.length === 0 ? (
        <p className="emptyState">No rows.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>Status</th>
              <th>Count</th>
              <th>Amount</th>
              <th>Due</th>
              <th>Paid</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.key}>
                <td>{row.label}</td>
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
