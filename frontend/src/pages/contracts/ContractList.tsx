import { FormEvent, useEffect, useState } from 'react';
import { ApiError } from '../../api/client';
import { Contract, createContract, getContract, listContracts } from '../../api/contracts';
import { contractStatusLabel, labelFor, localizeError } from '../../i18n/labels';
import { ContractDetail } from './ContractDetail';

export function ContractList() {
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [selected, setSelected] = useState<Contract | null>(null);
  const [creating, setCreating] = useState(false);
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');
  const [form, setForm] = useState({
    quoteId: '',
    opportunityId: '',
    customerId: '',
    amount: '',
    expectedSignedDate: '',
    contractNote: '',
    amountDifferenceReason: '',
    ownerId: ''
  });

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh(nextSearch = search) {
    const response = await listContracts(nextSearch);
    setContracts(response.items);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');
    try {
      const created = await createContract(form);
      setSelected(created);
      setCreating(false);
      setForm({ quoteId: '', opportunityId: '', customerId: '', amount: '', expectedSignedDate: '', contractNote: '', amountDifferenceReason: '', ownerId: '' });
      await refresh();
    } catch (caught) {
      const apiError = caught as ApiError;
      setError(localizeError(apiError));
    }
  }

  async function selectContract(id: string) {
    setError('');
    setSelected(await getContract(id));
  }

  async function updateSelected(contract: Contract) {
    setSelected(contract);
    await refresh();
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>合同</h1>
          <p>基于已接受报价创建待签署合同并管理签署。</p>
        </div>
        <button className="primaryButton" type="button" onClick={() => setCreating((value) => !value)}>新建合同</button>
      </section>
      {error && <div role="alert" className="error">{error}</div>}
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={(event) => { event.preventDefault(); void refresh(search); }}>
            <label>
              搜索
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <button className="secondaryButton" type="submit">搜索</button>
          </form>
          {creating && (
            <form className="createPanel" onSubmit={submit}>
              <label>
                报价 ID
                <input value={form.quoteId} onChange={(event) => setForm({ ...form, quoteId: event.target.value })} />
              </label>
              <label>
                商机 ID
                <input value={form.opportunityId} onChange={(event) => setForm({ ...form, opportunityId: event.target.value })} />
              </label>
              <label>
                客户 ID
                <input value={form.customerId} onChange={(event) => setForm({ ...form, customerId: event.target.value })} />
              </label>
              <label>
                金额
                <input value={form.amount} onChange={(event) => setForm({ ...form, amount: event.target.value })} />
              </label>
              <label>
                预计签署日期
                <input type="date" value={form.expectedSignedDate} onChange={(event) => setForm({ ...form, expectedSignedDate: event.target.value })} />
              </label>
              <label>
                合同备注
                <textarea value={form.contractNote} onChange={(event) => setForm({ ...form, contractNote: event.target.value })} />
              </label>
              <label>
                金额差异原因
                <input value={form.amountDifferenceReason} onChange={(event) => setForm({ ...form, amountDifferenceReason: event.target.value })} />
              </label>
              <label>
                负责人 ID
                <input value={form.ownerId} onChange={(event) => setForm({ ...form, ownerId: event.target.value })} />
              </label>
              <button className="primaryButton" type="submit">保存合同</button>
            </form>
          )}
          <div className="recordList" aria-label="合同记录">
            {contracts.length === 0 ? <p className="emptyState">暂无合同。</p> : contracts.map((contract) => (
              <button className="recordRow" type="button" key={contract.id} onClick={() => void selectContract(contract.id)}>
                <strong>{contract.opportunityId}</strong>
                <span>{labelFor(contractStatusLabel, contract.status)} · {contract.amount}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? <ContractDetail contract={selected} onUpdated={updateSelected} onError={setError} /> : <p className="emptyState">选择合同以管理状态。</p>}
        </div>
      </section>
    </main>
  );
}
