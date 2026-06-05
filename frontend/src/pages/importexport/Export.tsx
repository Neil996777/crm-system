import { useState } from 'react';
import { ExportRun, startExport } from '../../api/importexport';
import { labelFor, objectTypeLabel } from '../../i18n/labels';

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
      const safe = err as { safeMessage?: string };
      setError(safe.safeMessage ?? '无法启动导出。');
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="listPanel">
      <div className="sectionTitle">
        <h3>CSV 导出</h3>
        <span>{includeArchived ? '包含已归档' : '不包含已归档'}</span>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      <div className="importForm">
        <label>
          对象类型
          <select value={objectType} onChange={(event) => setObjectType(event.target.value)}>
            <option value="lead">{labelFor(objectTypeLabel, 'lead')}</option>
          </select>
        </label>
        <label className="inlineCheckbox">
          <input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />
          包含已归档
        </label>
        <label className="inlineCheckbox">
          <input type="checkbox" checked={confirmed} onChange={(event) => setConfirmed(event.target.checked)} />
          确认导出范围并记录审计日志
        </label>
        <button className="primaryButton" type="button" disabled={busy || !confirmed} onClick={() => void runExport()}>
          开始导出
        </button>
      </div>
      {result ? (
        <div className="exportResult">
          <p className="resultSummary">已导出 {result.exportedCount} 行线索</p>
          <p>{result.archivedIncluded ? '包含已归档' : '不包含已归档'}</p>
          <textarea readOnly value={result.content} aria-label="导出 CSV 内容" />
        </div>
      ) : null}
    </section>
  );
}
