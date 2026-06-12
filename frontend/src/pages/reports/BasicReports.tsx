import { useEffect, useState } from 'react';
import { BarChart3, BriefcaseBusiness, FileText, Landmark, ListChecks, ReceiptText, RefreshCcw, Target, Users } from 'lucide-react';
import { BasicReport, getBasicReport, GroupRow, PaymentGroupRow } from '../../api/reports';
import { Badge, Button, DataTable, FunnelBars, MetricCard, PageHeader, Panel, PanelHeader } from '../../components/ui';
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

  const scope = report ? labelFor(reportScopeLabel, report.scope) : '授权';
  const archiveFilter = report ? labelFor(reportArchiveFilterLabel, report.filters.archived) : '默认仅活动记录';

  return (
    <main className="content reportDashboard" data-uiux="reports-team">
      <PageHeader
        title="团队报表"
        description={`${scope}范围 · ${archiveFilter}${report?.filters.from ? ` · ${report.filters.from} 至 ${report.filters.to}` : ''}`}
        actions={(
          <>
            <Button variant="primary" onClick={() => void refresh()}>
              <RefreshCcw size={16} aria-hidden="true" />
              刷新
            </Button>
          </>
        )}
      />
      {error && <div role="alert" className="errorBanner">{error}</div>}
      {report?.emptyState ? <p className="emptyState">暂无报表记录。</p> : null}
      {report ? (
        <>
          <section className="reportsKpiStrip" aria-label="基础报表指标">
            <MetricCard label="线索数" value={report.metrics.leadCount} icon={<ListChecks size={18} aria-hidden="true" />} delta="线索总数" />
            <MetricCard label="商机数" value={report.metrics.opportunityCount} tone="mint" icon={<Landmark size={18} aria-hidden="true" />} delta="商机总数" />
            <MetricCard label="任务数" value={report.metrics.taskCount} tone="peach" icon={<BriefcaseBusiness size={18} aria-hidden="true" />} delta="任务总数" />
            <MetricCard label="赢单数" value={report.metrics.wonCount} tone="mint" icon={<Target size={18} aria-hidden="true" />} delta="已赢商机" />
            <MetricCard label="丢单数" value={report.metrics.lostCount} tone="purple" icon={<Target size={18} aria-hidden="true" />} delta="已丢商机" />
            <MetricCard className="currencyMetric" label="报价额" value={compactMoney(report.metrics.quoteAmount, report.currency)} tone="peach" icon={<FileText size={18} aria-hidden="true" />} delta="报价金额" />
            <MetricCard className="currencyMetric" label="合同额" value={compactMoney(report.metrics.contractAmount, report.currency)} tone="purple" icon={<ReceiptText size={18} aria-hidden="true" />} delta="合同金额" />
            <MetricCard className="currencyMetric" label="已回款额" value={compactMoney(report.metrics.paidAmount, report.currency)} tone="mint" icon={<ReceiptText size={18} aria-hidden="true" />} delta="实收金额" />
            <MetricCard className="currencyMetric" label="应收额" value={compactMoney(report.metrics.receivableAmount, report.currency)} tone="sky" icon={<BriefcaseBusiness size={18} aria-hidden="true" />} delta="待收金额" />
          </section>

          <section className="reportsMainGrid">
            <Panel className="reportVisualPanel" aria-label="管道分析">
              <PanelHeader
                title="管道分析"
                description="商机按阶段分布"
                meta="六阶段"
                actions={<span className="panelIcon sky"><BarChart3 size={18} aria-hidden="true" /></span>}
              />
              <FunnelBars
                rows={report.breakdowns.opportunitiesByStage.map((row) => ({
                  label: labelFor(opportunityStageLabel, row.label),
                  value: row.count,
                  suffix: <> · {money(row.amount, report.currency)}</>
                }))}
              />
            </Panel>

            <Panel className="reportVisualPanel" aria-label="负责人分组" data-report-owner-group="true" tabIndex={-1}>
              <PanelHeader
                title="负责人分组"
                description="按负责人汇总"
                meta="团队范围"
                actions={<span className="panelIcon mint"><Users size={18} aria-hidden="true" /></span>}
              />
              <DataTable
                caption="负责人分组表"
                rows={report.groups}
                rowKey={(row) => row.key}
                empty="暂无负责人分组数据。"
                columns={[
                  { key: 'owner', header: '负责人', render: (row) => row.label || row.key },
                  { key: 'count', header: '商机数', align: 'right', render: (row) => row.count },
                  { key: 'amount', header: '金额', align: 'right', render: (row) => <span className="money">{money(row.amount, report.currency)}</span> }
                ]}
              />
            </Panel>
          </section>

          <section className="reportBreakdownGrid" aria-label="状态阶段分解">
            <BreakdownCard title="线索按状态" rows={report.breakdowns.leadsByStatus} labels={leadStatusLabel} currency={report.currency} />
            <BreakdownCard title="商机按阶段" rows={report.breakdowns.opportunitiesByStage} labels={opportunityStageLabel} currency={report.currency} />
            <BreakdownCard title="报价按状态" rows={report.breakdowns.quotesByStatus} labels={quoteStatusLabel} currency={report.currency} />
            <BreakdownCard title="合同按状态" rows={report.breakdowns.contractsByStatus} labels={contractStatusLabel} currency={report.currency} />
            <PaymentBreakdownCard rows={report.breakdowns.paymentsByStatus} currency={report.currency} />
          </section>
        </>
      ) : null}
    </main>
  );
}

function BreakdownCard({ title, rows, labels, currency }: { title: string; rows: GroupRow[]; labels: Record<string, string>; currency: string }) {
  return (
    <Panel className="breakdownCard">
      <PanelHeader title={title} />
      <div className="breakdownList">
        {rows.length === 0 ? <p className="emptyState compact">暂无数据行。</p> : rows.map((row) => (
          <div className="breakdownItem" key={row.key}>
            <span>{labelFor(labels, row.label)}</span>
            <strong>{row.count}</strong>
            <Badge>{money(row.amount, currency)}</Badge>
          </div>
        ))}
      </div>
    </Panel>
  );
}

function PaymentBreakdownCard({ rows, currency }: { rows: PaymentGroupRow[]; currency: string }) {
  return (
    <Panel className="breakdownCard">
      <PanelHeader title="回款按状态" />
      <div className="breakdownList">
        {rows.length === 0 ? <p className="emptyState compact">暂无数据行。</p> : rows.map((row) => (
          <div className="breakdownItem" key={row.key}>
            <span>{labelFor(paymentStatusLabel, row.label)}</span>
            <strong>{row.count}</strong>
            <Badge>{money(row.paidAmount || row.amount, currency)} / {money(row.dueAmount, currency)}</Badge>
          </div>
        ))}
      </div>
    </Panel>
  );
}

function numberValue(value: string | number | undefined) {
  const parsed = typeof value === 'number' ? value : Number(value ?? 0);
  return Number.isFinite(parsed) ? parsed : 0;
}

function money(value: string | number, currency: string) {
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency, maximumFractionDigits: 0 }).format(numberValue(value));
}

function compactMoney(value: string | number, currency: string) {
  const amount = numberValue(value);
  const absolute = Math.abs(amount);
  if (absolute >= 100_000_000) return `${currencySymbol(currency)}${formatCompactNumber(amount / 100_000_000)}亿`;
  if (absolute >= 10_000) return `${currencySymbol(currency)}${formatCompactNumber(amount / 10_000)}万`;
  return money(amount, currency);
}

function formatCompactNumber(value: number) {
  return new Intl.NumberFormat('zh-CN', { maximumFractionDigits: Math.abs(value) >= 10 ? 0 : 1 }).format(value);
}

function currencySymbol(currency: string) {
  const symbol = new Intl.NumberFormat('zh-CN', { style: 'currency', currency, maximumFractionDigits: 0 })
    .formatToParts(0)
    .find((part) => part.type === 'currency')?.value;
  return symbol ?? `${currency} `;
}
