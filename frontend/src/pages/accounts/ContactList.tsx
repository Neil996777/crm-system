import { FormEvent, useEffect, useState } from 'react';
import { Contact, getContact, listAllContacts } from '../../api/accounts';

export function ContactList() {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [selected, setSelected] = useState<Contact | null>(null);
  const [search, setSearch] = useState('');

  useEffect(() => {
    void refresh();
  }, []);

  async function refresh(nextSearch = search) {
    const response = await listAllContacts(nextSearch);
    setContacts(response.items);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    await refresh(search);
  }

  async function selectContact(id: string) {
    setSelected(await getContact(id));
  }

  return (
    <main className="content">
      <section className="pageHeader">
        <div>
          <h1>联系人</h1>
          <p>查找有权限访问的客户联系人记录。</p>
        </div>
      </section>
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={submit}>
            <label>
              搜索
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <button className="secondaryButton" type="submit">搜索</button>
          </form>
          <div className="recordList" aria-label="联系人记录">
            {contacts.length === 0 ? <p className="emptyState">暂无联系人。</p> : contacts.map((contact) => (
              <button className="recordRow" type="button" key={contact.id} onClick={() => void selectContact(contact.id)}>
                <strong>{contact.contactName}</strong>
                <span>{contact.accountName || contact.accountId}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? (
            <section className="detailPane" aria-label="联系人详情">
              <div className="detailHeader">
                <div>
                  <h2>{selected.contactName}</h2>
                  <p>{selected.accountName || selected.accountId}</p>
                </div>
                <span className="statusPill">联系人</span>
              </div>
              <dl className="detailGrid">
                <div>
                  <dt>客户</dt>
                  <dd>{selected.accountName || selected.accountId}</dd>
                </div>
                <div>
                  <dt>邮箱</dt>
                  <dd>{selected.email || '无'}</dd>
                </div>
                <div>
                  <dt>电话</dt>
                  <dd>{selected.phone || '无'}</dd>
                </div>
                <div>
                  <dt>角色备注</dt>
                  <dd>{selected.roleNote || '无'}</dd>
                </div>
              </dl>
            </section>
          ) : (
            <p className="emptyState">选择联系人以查看详情。</p>
          )}
        </div>
      </section>
    </main>
  );
}
