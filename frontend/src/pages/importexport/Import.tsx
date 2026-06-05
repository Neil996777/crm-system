import { useState } from 'react';
import type { FormEvent } from 'react';
import { ImportRun, startImport } from '../../api/importexport';
import { labelFor, objectTypeLabel, runStatusLabel } from '../../i18n/labels';
import { ExportPanel } from './Export';

export function ImportExportPage() {
  const [objectType, setObjectType] = useState('lead');
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<ImportRun | null>(null);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

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
      const safe = err as { safeMessage?: string };
      setError(safe.safeMessage ?? '无法启动导入。');
    } finally {
      setBusy(false);
    }
  }

  return (
    <main className="page">
      <div className="pageHeader">
        <div>
          <h1>导入/导出</h1>
          <p>导入有权限访问记录的 CSV 文件</p>
        </div>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      <section className="listPanel">
        <form className="importForm" onSubmit={(event) => void submit(event)}>
          <label>
            对象类型
            <select value={objectType} onChange={(event) => setObjectType(event.target.value)}>
              <option value="lead">{labelFor(objectTypeLabel, 'lead')}</option>
            </select>
          </label>
          <label>
            CSV 文件
            <input type="file" accept=".csv,text/csv" onChange={(event) => setFile(event.target.files?.[0] ?? null)} />
          </label>
          <button className="primaryButton" type="submit" disabled={busy}>
            开始导入
          </button>
        </form>
      </section>
      {result ? (
        <section className="listPanel">
          <div className="sectionTitle">
            <h3>导入结果</h3>
            <span>{labelFor(runStatusLabel, result.status)}</span>
          </div>
          <p className="resultSummary">已导入 {result.successCount} / {result.totalRows} 行</p>
          {result.rowErrors.length === 0 ? (
            <p className="emptyState">没有行错误。</p>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>行</th>
                  <th>字段</th>
                  <th>规则</th>
                  <th>消息</th>
                </tr>
              </thead>
              <tbody>
                {result.rowErrors.map((row) => (
                  <tr key={`${row.rowNumber}-${row.field}-${row.code}`}>
                    <td>第 {row.rowNumber} 行</td>
                    <td>{row.field}</td>
                    <td>{row.code}</td>
                    <td>{row.safeMessage}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </section>
      ) : null}
      <ExportPanel />
    </main>
  );
}
