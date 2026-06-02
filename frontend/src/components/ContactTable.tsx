import { Contact } from '../api/accounts';

export function ContactTable({ contacts }: { contacts: Contact[] }) {
  if (contacts.length === 0) {
    return <p className="emptyState">No contacts.</p>;
  }

  return (
    <table className="dataTable" aria-label="Contacts">
      <thead>
        <tr>
          <th>Name</th>
          <th>Email</th>
          <th>Phone</th>
          <th>Role note</th>
        </tr>
      </thead>
      <tbody>
        {contacts.map((contact) => (
          <tr key={contact.id}>
            <td>{contact.contactName}</td>
            <td>{contact.email || '-'}</td>
            <td>{contact.phone || '-'}</td>
            <td>{contact.roleNote || '-'}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
