import { useRef, useState } from 'react';
import type { FormEvent } from 'react';
import { Upload } from 'lucide-react';
import { ImportRun, startImport } from '../../api/importexport';
import { Badge, Button, DataTable, PageHeader, Panel, PanelHeader } from '../../components/ui';
import { labelFor, localizeError, localizeMessage, objectTypeLabel, runStatusLabel } from '../../i18n/labels';
import { ExportPanel } from './Export';

export function ImportExportPage() {
  const [objectType, setObjectType] = useState('lead');
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<ImportRun | null>(null);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError('');
    setResult(null);
    if (!file) {
      setError('请选择 CSV 文件。');
      return;
    }
    setBusy(true);
    try {
      const content = await file.text();
      setResult(await startImport({ objectType, filename: file.name, content }));
    } catch (err) {
      setError(localizeError(err as { safeMessage?: string }, '无法启动导入。'));
    } finally {
      setBusy(false);
    }
  }

  function focusImportForm() {
    fileInputRef.current?.scrollIntoView({ block: 'center', behavior: 'smooth' });
    fileInputRef.current?.focus();
  }

  return (
    <main className="content importExportPage" data-uiux="import-export">
      <PageHeader
        title="导入/导出"
        description="导入/导出有权限访问的记录 · 经理视图 · 仅处理授权范围"
        actions={(
          <>
            <Button variant="primary" onClick={focusImportForm}>
              <Upload size={16} aria-hidden="true" />
              新建导入
            </Button>
          </>
        )}
      />
      {error && <div role="alert" className="errorBanner">{error}</div>}
      <section className="importExportGrid">
        <Panel className="importFlowPanel" aria-label="导入流程">
          <PanelHeader
            title="导入流程"
            description="CSV 文件 · 部分失败示例"
            meta={result ? labelFor(runStatusLabel, result.status) : '待开始'}
            actions={<span className="panelIcon sky"><Upload size={18} aria-hidden="true" /></span>}
          />
          <div className="objectPills" aria-label="导入对象类型">
            <Badge tone="primary">{labelFor(objectTypeLabel, 'lead')}</Badge>
            <Badge>客户</Badge>
            <Badge>联系人</Badge>
            <Badge>商机</Badge>
            <Badge>报价</Badge>
            <Badge>合同</Badge>
          </div>
          <form className="importForm" onSubmit={(event) => void submit(event)}>
            <label>
              对象类型
              <select value={objectType} onChange={(event) => setObjectType(event.target.value)}>
                <option value="lead">{labelFor(objectTypeLabel, 'lead')}</option>
              </select>
            </label>
            <label>
              CSV 文件
              <input ref={fileInputRef} type="file" accept=".csv,text/csv" onChange={(event) => setFile(event.target.files?.[0] ?? null)} />
            </label>
            <button className="primaryButton" type="submit" disabled={busy}>
              开始导入
            </button>
          </form>
          {result ? <ImportResult result={result} /> : <p className="safeBox">导入只写入当前角色可访问范围内的数据。</p>}
          <AuditStatus
            operationLogStatus={result?.operationLogStatus}
            cleanupStatus={result?.cleanupStatus}
            retainedUntil={result?.retainedUntil}
          />
          <div className="panelFooterNote">
            <span>导入只写入当前角色可访问范围内的数据。</span>
            <span>进行中/失败态见源码注释。</span>
          </div>
        </Panel>
        <ExportPanel />
      </section>
    </main>
  );
}

function ImportResult({ result }: { result: ImportRun }) {
  return (
    <section className="runResult" aria-label="导入结果">
      <div className="sectionTitle">
        <h3>导入结果 {result.runId}</h3>
        <Badge tone={result.failureCount > 0 ? 'warning' : 'success'}>{labelFor(runStatusLabel, result.status)}</Badge>
      </div>
      <div className="runStatGrid" aria-label="导入结果字段">
        <RunStat label="总行数" value={result.totalRows} />
        <RunStat label="成功数" value={result.successCount} />
        <RunStat label="失败数" value={result.failureCount} />
      </div>
      <DataTable
        caption="导入逐行错误表"
        rows={result.rowErrors}
        rowKey={(row) => `${row.rowNumber}-${row.field}-${row.code}`}
        empty="没有行错误。"
        columns={[
          { key: 'row', header: '行号', render: (row) => `第 ${row.rowNumber} 行` },
          { key: 'field', header: '字段', render: (row) => row.field || '整行' },
          { key: 'message', header: '错误信息', render: (row) => localizeMessage(row.safeMessage, row.code) },
          { key: 'code', header: '代码', render: (row) => row.code }
        ]}
      />
    </section>
  );
}

export function AuditStatus({ operationLogStatus, cleanupStatus, retainedUntil }: { operationLogStatus?: string; cleanupStatus?: string; retainedUntil?: string }) {
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

export function RunStat({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="runStat">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function formatDate(value: string) {
  if (!value) return '';
  return value.length > 10 ? value.slice(0, 10) : value;
}
