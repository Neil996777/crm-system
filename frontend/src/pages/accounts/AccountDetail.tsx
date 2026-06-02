import { useEffect, useState } from 'react';
import { Account, ArchiveEligibility, archiveAccount, Contact, getAccountArchiveEligibility, listContacts } from '../../api/accounts';
import { ApiError } from '../../api/client';
import { AddContactDialog } from '../../components/AddContactDialog';
import { ArchiveConfirmation } from '../../components/ArchiveConfirmation';
import { ContactTable } from '../../components/ContactTable';

type Props = {
  account: Account;
  onArchived: () => void;
  onError: (message: string) => void;
};

export function AccountDetail({ account, onArchived, onError }: Props) {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [eligibility, setEligibility] = useState<ArchiveEligibility | null>(null);
  const [reason, setReason] = useState('');

  useEffect(() => {
    void refreshContacts();
  }, [account.id]);

  async function refreshContacts() {
    const response = await listContacts(account.id);
    setContacts(response.items);
  }

  async function created() {
    await refreshContacts();
  }

  async function startArchive() {
    onError('');
    try {
      const next = await getAccountArchiveEligibility(account.id);
      setEligibility(next);
      setReason('Archive inactive customer record');
    } catch (caught) {
      const error = caught as ApiError;
      onError(error.safeMessage || 'Request failed.');
    }
  }

  async function confirmArchive() {
    onError('');
    try {
      await archiveAccount(account.id, account.version, reason);
      setEligibility(null);
      onArchived();
    } catch (caught) {
      const error = caught as ApiError;
      onError(error.safeMessage || 'Request failed.');
    }
  }

  return (
    <section className="detailPane" aria-label="Customer detail">
      <div className="detailHeader">
        <div>
          <h2>{account.companyName}</h2>
          <p>{account.ownerId}</p>
        </div>
        <div className="headerActions">
          <span className="statusPill">{account.archived ? 'Archived' : account.customerStatus}</span>
          {!account.archived && <button className="secondaryButton" type="button" onClick={() => void startArchive()}>Archive</button>}
        </div>
      </div>
      <ArchiveConfirmation
        eligibility={eligibility}
        reason={reason}
        onReasonChange={setReason}
        onConfirm={() => void confirmArchive()}
        onCancel={() => setEligibility(null)}
      />
      <dl className="detailGrid">
        <div>
          <dt>Owner</dt>
          <dd>{account.ownerId}</dd>
        </div>
        <div>
          <dt>Version</dt>
          <dd>{account.version}</dd>
        </div>
      </dl>
      <section className="relatedSection">
        <div className="sectionTitle">
          <h3>Contacts</h3>
          <AddContactDialog accountId={account.id} onCreated={created} onError={onError} />
        </div>
        <ContactTable contacts={contacts} />
      </section>
      <section className="relatedSection">
        <h3>Related records</h3>
        <p className="emptyState">No related opportunities, contracts, payments, or history loaded yet.</p>
      </section>
    </section>
  );
}
