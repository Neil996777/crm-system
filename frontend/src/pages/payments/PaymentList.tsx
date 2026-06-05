import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract } from '../../api/contracts';
import { listPaymentContracts } from '../../api/payments';
import { contractStatusLabel, labelFor } from '../../i18n/labels';
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
      setError(apiError.safeMessage || '请求失败。');
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
          <h1>回款</h1>
          <p>跟踪已签署和启用合同的售后回款。</p>
        </div>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={searchContracts}>
            <label>
              搜索
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <button className="secondaryButton" type="submit">搜索</button>
          </form>
          <div className="recordList" aria-label="回款合同记录">
            {contracts.length === 0 ? <p className="emptyState">暂无合同。</p> : contracts.map((contract) => (
              <button className="recordRow" type="button" key={contract.id} onClick={() => { setError(''); setSelected(contract); }}>
                <strong>{contract.opportunityId}</strong>
                <span>{labelFor(contractStatusLabel, contract.status)} · {contract.amount}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <PaymentDetail key={selected.id} contract={selected} onError={setError} /> : <p className="emptyState">选择合同以管理回款计划和实际回款。</p>}
        </div>
      </section>
    </main>
  );
}
