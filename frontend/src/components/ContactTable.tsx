import { Contact } from '../api/accounts';

export function ContactTable({ contacts }: { contacts: Contact[] }) {
  if (contacts.length === 0) {
    return <p className="emptyState">暂无联系人。</p>;
  }

  return (
    <div className="tableWrap">
      <table className="dataTable" aria-label="联系人">
      <thead>
        <tr>
          <th>姓名</th>
          <th>邮箱</th>
          <th>电话</th>
          <th>角色备注</th>
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
    </div>
  );
}
