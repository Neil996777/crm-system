import { useEffect, useState } from 'react';
import { getManagerOverview, ManagerOverview as ManagerOverviewData } from '../../api/reports';
import { BasicReports } from './BasicReports';

export function ManagerOverview() {
  const [overview, setOverview] = useState<ManagerOverviewData | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh() {
    setError('');
    try {
      setOverview(await getManagerOverview());
    } catch (err) {
      const safe = err as { safeMessage?: string };
      setError(safe.safeMessage ?? 'Unable to load manager overview.');
    }
  }

  return (
    <main className="page">
      <div className="pageHeader">
        <div>
          <h1>Manager Team Overview</h1>
          <p>{overview ? `${overview.scope} scope · ${overview.filters.archived}` : 'Team records'}</p>
        </div>
        <button className="secondaryButton" type="button" onClick={() => void refresh()}>
          Refresh
        </button>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      {overview?.emptyState ? <p className="emptyState">No team records yet.</p> : null}
      {overview ? (
        <>
          <section className="metricGrid" aria-label="Team metrics">
            <Metric label="Leads" value={overview.metrics.leadCount} />
            <Metric label="Opportunities" value={overview.metrics.opportunityCount} />
            <Metric label="Quotes" value={overview.metrics.quoteAmount} />
            <Metric label="Contracts" value={overview.metrics.contractAmount} />
            <Metric label="Paid" value={overview.metrics.paidAmount} />
            <Metric label="Receivable" value={overview.metrics.receivableAmount} />
            <Metric label="Open tasks" value={overview.metrics.taskCount} />
            <Metric label="Won / Lost" value={`${overview.metrics.wonCount} / ${overview.metrics.lostCount}`} />
          </section>
          <section className="listPanel">
            <div className="sectionTitle">Pipeline Status</div>
            {overview.pipeline.length === 0 ? (
              <p className="emptyState">No opportunities in the team pipeline.</p>
            ) : (
              <table>
                <thead>
                  <tr>
                    <th>Stage</th>
                    <th>Count</th>
                    <th>Amount</th>
                  </tr>
                </thead>
                <tbody>
                  {overview.pipeline.map((row) => (
                    <tr key={row.key}>
                      <td>{row.label}</td>
                      <td>{row.count}</td>
                      <td>{row.amount} {overview.currency}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </section>
        </>
      ) : null}
      <BasicReports />
    </main>
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
