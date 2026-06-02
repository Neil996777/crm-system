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
          <h1>Contacts</h1>
          <p>Find authorized contact records across customer accounts.</p>
        </div>
      </section>
      <section className="leadLayout">
        <div className="listPane">
          <form className="toolbar" onSubmit={submit}>
            <label>
              Search
              <input value={search} onChange={(event) => setSearch(event.target.value)} />
            </label>
            <button className="secondaryButton" type="submit">Search</button>
          </form>
          <div className="recordList" aria-label="Contact records">
            {contacts.length === 0 ? <p className="emptyState">No contacts found.</p> : contacts.map((contact) => (
              <button className="recordRow" type="button" key={contact.id} onClick={() => void selectContact(contact.id)}>
                <strong>{contact.contactName}</strong>
                <span>{contact.accountName || contact.accountId}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="detailShell">
          {selected ? (
            <section className="detailPane" aria-label="Contact detail">
              <div className="detailHeader">
                <div>
                  <h2>{selected.contactName}</h2>
                  <p>{selected.accountName || selected.accountId}</p>
                </div>
                <span className="statusPill">Contact</span>
              </div>
              <dl className="detailGrid">
                <div>
                  <dt>Account</dt>
                  <dd>{selected.accountName || selected.accountId}</dd>
                </div>
                <div>
                  <dt>Email</dt>
                  <dd>{selected.email || 'None'}</dd>
                </div>
                <div>
                  <dt>Phone</dt>
                  <dd>{selected.phone || 'None'}</dd>
                </div>
                <div>
                  <dt>Role note</dt>
                  <dd>{selected.roleNote || 'None'}</dd>
                </div>
              </dl>
            </section>
          ) : (
            <p className="emptyState">Select a contact to view details.</p>
          )}
        </div>
      </section>
    </main>
  );
}
