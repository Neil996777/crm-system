import { useState } from 'react';
import { ExportRun, startExport } from '../../api/importexport';

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
      setError(safe.safeMessage ?? 'Unable to start export.');
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="listPanel">
      <div className="sectionTitle">
        <h3>CSV Export</h3>
        <span>{includeArchived ? 'Archived included' : 'Archived excluded'}</span>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      <div className="importForm">
        <label>
          Object type
          <select value={objectType} onChange={(event) => setObjectType(event.target.value)}>
            <option value="lead">Lead</option>
          </select>
        </label>
        <label className="inlineCheckbox">
          <input type="checkbox" checked={includeArchived} onChange={(event) => setIncludeArchived(event.target.checked)} />
          Include archived
        </label>
        <label className="inlineCheckbox">
          <input type="checkbox" checked={confirmed} onChange={(event) => setConfirmed(event.target.checked)} />
          Confirm export scope and audit log
        </label>
        <button className="primaryButton" type="button" disabled={busy || !confirmed} onClick={() => void runExport()}>
          Start export
        </button>
      </div>
      {result ? (
        <div className="exportResult">
          <p className="resultSummary">Exported {result.exportedCount} lead rows</p>
          <p>{result.archivedIncluded ? 'Archived included' : 'Archived excluded'}</p>
          <textarea readOnly value={result.content} aria-label="Export CSV content" />
        </div>
      ) : null}
    </section>
  );
}
