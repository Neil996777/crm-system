import { useEffect, useState } from 'react';
import { Building2, Hash, UserRound, UsersRound } from 'lucide-react';
import { Account, ArchiveEligibility, archiveAccount, Contact, getAccountArchiveEligibility, listContacts } from '../../api/accounts';
import { ApiError } from '../../api/client';
import { AddContactDialog } from '../../components/AddContactDialog';
import { ArchiveConfirmation } from '../../components/ArchiveConfirmation';
import { ContactTable } from '../../components/ContactTable';
import { DetailHero, DetailStat, StatusPill } from '../../components/CrudScaffold';
import { Panel, SkeletonBlock } from '../../components/ui';
import { accountStatusLabel, archiveStatusLabel, labelFor, localizeError } from '../../i18n/labels';

type Props = {
  account: Account;
  onArchived: () => void;
  onError: (message: string) => void;
  onBack?: () => void;
  error?: string;
};

export function AccountDetail({ account, onArchived, onError, onBack, error }: Props) {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [contactsLoading, setContactsLoading] = useState(true);
  const [eligibility, setEligibility] = useState<ArchiveEligibility | null>(null);
  const [reason, setReason] = useState('');

  useEffect(() => {
    void refreshContacts();
  }, [account.id]);

  async function refreshContacts() {
    setContactsLoading(true);
    onError('');
    try {
      const response = await listContacts(account.id);
      setContacts(response.items);
    } catch (caught) {
      const apiError = caught as ApiError;
      onError(localizeError(apiError));
    } finally {
      setContactsLoading(false);
    }
  }

  async function created() {
    await refreshContacts();
  }

  async function startArchive() {
    onError('');
    try {
      const next = await getAccountArchiveEligibility(account.id);
      setEligibility(next);
      setReason('归档非活跃客户记录');
    } catch (caught) {
      const error = caught as ApiError;
      onError(localizeError(error));
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
      onError(localizeError(error));
    }
  }

  return (
    <main className="content crudPage" aria-label="客户详情">
      <DetailHero
        eyebrow="返回客户列表"
        title={account.companyName}
        subtitle={
          <>
            <span>负责人 {account.ownerId}</span>
            <span>更新于 {formatDate(account.updatedAt)}</span>
          </>
        }
        icon={<Building2 size={20} aria-hidden="true" />}
        status={<StatusPill tone={account.archived ? 'warning' : accountTone(account.customerStatus)}>{account.archived ? labelFor(archiveStatusLabel, 'Archived') : labelFor(accountStatusLabel, account.customerStatus)}</StatusPill>}
        onBack={onBack ?? (() => undefined)}
        actions={!account.archived ? <button className="secondaryButton" type="button" onClick={() => void startArchive()}>归档</button> : null}
        stats={
          <>
            <DetailStat label="负责人" value={account.ownerId} icon={<UserRound size={17} aria-hidden="true" />} />
            <DetailStat label="联系人" value={`${contacts.length} 个`} icon={<UsersRound size={17} aria-hidden="true" />} tone="mint" />
            <DetailStat label="版本" value={`v${account.version}`} icon={<Hash size={17} aria-hidden="true" />} />
          </>
        }
      />
      {error ? <div role="alert" className="error">{error}</div> : null}
      <ArchiveConfirmation
        eligibility={eligibility}
        reason={reason}
        onReasonChange={setReason}
        onConfirm={() => void confirmArchive()}
        onCancel={() => setEligibility(null)}
      />
      <section className="detailContentGrid">
        <Panel>
          <div className="sectionHeader">
            <h2>联系人</h2>
            <AddContactDialog accountId={account.id} onCreated={created} onError={onError} />
          </div>
          {contactsLoading ? <SkeletonBlock lines={3} label="正在加载联系人..." /> : <ContactTable contacts={contacts} />}
        </Panel>
        <Panel>
          <div className="sectionHeader">
            <h2>客户字段</h2>
            <span className="badge">详情</span>
          </div>
          <dl className="detailGrid">
            <div>
              <dt>客户状态</dt>
              <dd>{labelFor(accountStatusLabel, account.customerStatus)}</dd>
            </div>
            <div>
              <dt>归档状态</dt>
              <dd>{account.archived ? labelFor(archiveStatusLabel, 'Archived') : '活动记录'}</dd>
            </div>
            <div>
              <dt>归档原因</dt>
              <dd>{account.archiveReason || '无'}</dd>
            </div>
          </dl>
        </Panel>
      </section>
    </main>
  );
}

function accountTone(status: string): 'primary' | 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'Active') return 'success';
  if (status === 'Inactive') return 'warning';
  return 'neutral';
}

function formatDate(value: string) {
  if (!value) return '未更新';
  return value.length > 10 ? value.slice(0, 10) : value;
}
