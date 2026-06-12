import { useState } from 'react';
import { Download } from 'lucide-react';
import { ExportRun, startExport } from '../../api/importexport';
import { Badge, Button, Panel, PanelHeader } from '../../components/ui';
import { fileSafetyLabel, labelFor, localizeError, objectTypeLabel, runStatusLabel } from '../../i18n/labels';

export function ExportPanel() {
  const [objectType, setObjectType] = useState('lead');
  const [includeArchived, setIncludeArchived] = useState(false);
  const [confirmed, setConfirmed] = useState(false);
  const [result, setResult] = useState<ExportRun | null>(null);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  async function runExport() {
    setError('');
    setResult(null);
    setBusy(true);
    try {
      setResult(await startExport({ objectType, confirmed, includeArchived }));
    } catch (err) {
      setError(localizeError(err as { safeMessage?: string }, '无法启动导出。'));
    } finally {
      setBusy(false);
    }
  }

  return (
    <Panel className="exportFlowPanel" aria-label="导出流程">
      <PanelHeader
        title="导出流程"
        description="确认范围并记录审计后执行"
        meta={result ? labelFor(runStatusLabel, result.status) : '待确认'}
        actions={<span className="panelIcon mint"><Download size={18} aria-hidden="true" /></span>}
      />
      {error && <div role="alert" className="errorBanner">{error}</div>}
      <div className="objectPills" aria-label="导出对象类型">
        <Badge tone="primary">{labelFor(objectTypeLabel, 'lead')}</Badge>
      </div>
      <div className="importForm">
        <label>
          对象类型
          <select value={objectType} onChange={(event) => setObjectType(event.target.value)}>
            <option value="lead">{labelFor(objectTypeLabel, 'lead')}</option>
          </select>
        </label>
        <label className="inlineCheckbox">
          <input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />
          包含归档：{includeArchived ? '是' : '否'}
        </label>
        <label className="inlineCheckbox exportConfirm">
          <input type="checkbox" checked={confirmed} onChange={(event) => setConfirmed(event.target.checked)} />
          确认导出范围并记录审计日志
        </label>
        <Button variant="primary" disabled={busy || !confirmed} onClick={() => void runExport()}>
          开始导出
        </Button>
      </div>
      {result ? (
        <section className="exportResult" aria-label="导出结果">
          <div className="sectionTitle">
            <h3>导出结果 {result.runId}</h3>
            <Badge tone="success">{labelFor(runStatusLabel, result.status)}</Badge>
          </div>
          <div className="runStatGrid">
            <RunStat label="导出行数" value={result.exportedCount} />
            <RunStat label="包含归档" value={result.archivedIncluded ? '是' : '否'} />
            <RunStat label="文件安全" value={fileSafetyText(result.fileSafety)} />
          </div>
          <div className="exportFileBox">
            <strong>文件：{result.filename}</strong>
            <p>内容：仅包含当前账号有权访问的{labelFor(objectTypeLabel, result.objectType)}记录；已去除不可访问记录。</p>
            <textarea readOnly value={result.content} aria-label="导出 CSV 内容" />
          </div>
        </section>
      ) : <p className="safeBox">导出前必须完成确认勾选。</p>}
      <AuditStatus
        operationLogStatus={result?.operationLogStatus}
        cleanupStatus={result?.cleanupStatus}
        retainedUntil={result?.retainedUntil}
      />
      <div className="panelFooterNote">
        <span>导出执行前必须完成确认勾选。</span>
        <span>失败态见源码注释。</span>
      </div>
    </Panel>
  );
}

function fileSafetyText(value: string) {
  return fileSafetyLabel[value] ?? fileSafetyLabel.unknown;
}

function RunStat({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="runStat">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function AuditStatus({ operationLogStatus, cleanupStatus, retainedUntil }: { operationLogStatus?: string; cleanupStatus?: string; retainedUntil?: string }) {
  return (
    <section className="auditCleanupBox" aria-label="审计与清理">
      <div className="sectionTitle">
        <h3>审计与清理</h3>
        <Badge tone={operationLogStatus ? 'success' : 'neutral'}>{operationLogStatus ? '已记录' : '待记录'}</Badge>
      </div>
      <p>
        审计记录状态：{operationLogStatus ? labelFor(runStatusLabel, operationLogStatus) : '待写入'} ·
        清理状态：{cleanupStatus ? labelFor(runStatusLabel, cleanupStatus) : '待处理'}
        {retainedUntil ? ` · 保留至 ${formatDate(retainedUntil)}` : null}
      </p>
    </section>
  );
}

function formatDate(value: string) {
  if (!value) return '';
  return value.length > 10 ? value.slice(0, 10) : value;
}
