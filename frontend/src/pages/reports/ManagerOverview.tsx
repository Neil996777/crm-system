import { useEffect, useState } from 'react';
import { getManagerOverview, ManagerOverview as ManagerOverviewData } from '../../api/reports';
import { labelFor, opportunityStageLabel, reportArchiveFilterLabel, reportScopeLabel } from '../../i18n/labels';
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
      setError(safe.safeMessage ?? '无法加载经理总览。');
    }
  }

  return (
    <main className="page">
      <div className="pageHeader">
        <div>
          <h1>经理团队总览</h1>
          <p>{overview ? `${labelFor(reportScopeLabel, overview.scope)}范围 · ${labelFor(reportArchiveFilterLabel, overview.filters.archived)}` : '团队记录'}</p>
        </div>
        <button className="secondaryButton" type="button" onClick={() => void refresh()}>
          刷新
        </button>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      {overview?.emptyState ? <p className="emptyState">暂无团队记录。</p> : null}
      {overview ? (
        <>
          <section className="metricGrid" aria-label="团队指标">
            <Metric label="线索" value={overview.metrics.leadCount} />
            <Metric label="商机" value={overview.metrics.opportunityCount} />
            <Metric label="报价" value={overview.metrics.quoteAmount} />
            <Metric label="合同" value={overview.metrics.contractAmount} />
            <Metric label="已回款" value={overview.metrics.paidAmount} />
            <Metric label="应收" value={overview.metrics.receivableAmount} />
            <Metric label="待处理任务" value={overview.metrics.taskCount} />
            <Metric label="赢单 / 丢单" value={`${overview.metrics.wonCount} / ${overview.metrics.lostCount}`} />
          </section>
          <section className="listPanel">
            <div className="sectionTitle">销售管道状态</div>
            {overview.pipeline.length === 0 ? (
              <p className="emptyState">团队销售管道暂无商机。</p>
            ) : (
              <table>
                <thead>
                  <tr>
                    <th>阶段</th>
                    <th>数量</th>
                    <th>金额</th>
                  </tr>
                </thead>
                <tbody>
                  {overview.pipeline.map((row) => (
                    <tr key={row.key}>
                      <td>{labelFor(opportunityStageLabel, row.label)}</td>
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
