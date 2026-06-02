import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract } from '../../api/contracts';
import { listPaymentContracts } from '../../api/payments';
import { PaymentDetail } from './PaymentDetail';

export function PaymentList() {
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [selected, setSelected] = useState<Contract | null>(null);
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh(nextSearch = search) {
    try {
      const response = await listPaymentContracts(nextSearch);
      setContracts(response.items);
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(apiError.safeMessage || 'Request failed.');
    }
  }

  function searchContracts(event: FormEvent) {
    event.preventDefault();
    setError('');
    void refresh(search);
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>Payments</h1>
          <p>Post-sale payment tracking for signed and active commercial work.</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={searchContracts}>
            <label>
              Search
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <button className="secondaryButton" type="submit">Search</button>
          </form>
          <div className="recordList" aria-label="Payment contract records">
            {contracts.length === 0 ? <p className="emptyState">No contracts found.</p> : contracts.map((contract) => (
              <button className="recordRow" type="button" key={contract.id} onClick={() => { setError(''); setSelected(contract); }}>
                <strong>{contract.opportunityId}</strong>
                <span>{contract.status} · {contract.amount}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <PaymentDetail key={selected.id} contract={selected} onError={setError} /> : <p className="emptyState">Select a contract to manage payment plan and actual payment.</p>}
        </div>
      </section>
    </main>
  );
}
