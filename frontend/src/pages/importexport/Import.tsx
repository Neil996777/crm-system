import { useState } from 'react';
import type { FormEvent } from 'react';
import { ImportRun, startImport } from '../../api/importexport';

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
      setError('Choose a CSV file.');
      return;
    }
    setBusy(true);
    try {
      const content = await file.text();
      setResult(await startImport({ objectType, filename: file.name, content }));
    } catch (err) {
      const safe = err as { safeMessage?: string };
      setError(safe.safeMessage ?? 'Unable to start import.');
    } finally {
      setBusy(false);
    }
  }

  return (
    <main className="page">
      <div className="pageHeader">
        <div>
          <h1>Import/Export</h1>
          <p>CSV import for authorized records</p>
        </div>
      </div>
      {error && <div role="alert" className="errorBanner">{error}</div>}
      <section className="listPanel">
        <form className="importForm" onSubmit={(event) => void submit(event)}>
          <label>
            Object type
            <select value={objectType} onChange={(event) => setObjectType(event.target.value)}>
              <option value="lead">Lead</option>
            </select>
          </label>
          <label>
            CSV file
            <input type="file" accept=".csv,text/csv" onChange={(event) => setFile(event.target.files?.[0] ?? null)} />
          </label>
          <button className="primaryButton" type="submit" disabled={busy}>
            Start import
          </button>
        </form>
      </section>
      {result ? (
        <section className="listPanel">
          <div className="sectionTitle">
            <h3>Import Result</h3>
            <span>{result.status}</span>
          </div>
          <p className="resultSummary">Imported {result.successCount} of {result.totalRows} rows</p>
          {result.rowErrors.length === 0 ? (
            <p className="emptyState">No row errors.</p>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Row</th>
                  <th>Field</th>
                  <th>Rule</th>
                  <th>Message</th>
                </tr>
              </thead>
              <tbody>
                {result.rowErrors.map((row) => (
                  <tr key={`${row.rowNumber}-${row.field}-${row.code}`}>
                    <td>Row {row.rowNumber}</td>
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
    </main>
  );
}
