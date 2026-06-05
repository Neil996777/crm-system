import { UserPlus } from 'lucide-react';
import { FormEvent, useState } from 'react';
import { checkContactDuplicate, Contact, createContact } from '../api/accounts';
import { ApiError } from '../api/client';
import { DuplicateWarningResult } from '../api/duplicates';
import { localizeError } from '../i18n/labels';
import { DuplicateWarning } from './DuplicateWarning';

type Props = {
  accountId: string;
  onCreated: (contact: Contact) => void;
  onError: (message: string) => void;
};

export function AddContactDialog({ accountId, onCreated, onError }: Props) {
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({ contactName: '', email: '', phone: '', roleNote: '' });
  const [duplicateWarning, setDuplicateWarning] = useState<DuplicateWarningResult | null>(null);

  async function submit(event: FormEvent) {
    event.preventDefault();
    await saveContact();
  }

  async function saveContact(proceedWarningToken?: string) {
    onError('');
    try {
      if (!proceedWarningToken) {
        const warning = await checkContactDuplicate({ email: form.email, phone: form.phone });
        if (warning.result === 'PossibleDuplicate' && warning.warningToken) {
          setDuplicateWarning(warning);
          return;
        }
      }
      const contact = await createContact(accountId, { ...form, proceedWarningToken });
      onCreated(contact);
      setForm({ contactName: '', email: '', phone: '', roleNote: '' });
      setDuplicateWarning(null);
      setOpen(false);
    } catch (caught) {
      const error = caught as ApiError;
      onError(localizeError(error));
    }
  }

  return (
    <section className="convertPanel">
      <button className="primaryButton iconButtonText" type="button" onClick={() => setOpen((value) => !value)}>
        <UserPlus size={16} />
        添加联系人
      </button>
      {open && (
        <form className="createPanel" onSubmit={submit}>
          <label>
            联系人姓名
            <input value={form.contactName} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, contactName: event.target.value }); }} />
          </label>
          <label>
            邮箱
            <input value={form.email} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, email: event.target.value }); }} />
          </label>
          <label>
            电话
            <input value={form.phone} onChange={(event) => { setDuplicateWarning(null); setForm({ ...form, phone: event.target.value }); }} />
          </label>
          <label>
            角色备注
            <input value={form.roleNote} onChange={(event) => setForm({ ...form, roleNote: event.target.value })} />
          </label>
          <button className="primaryButton" type="submit">保存联系人</button>
          {duplicateWarning ? (
            <DuplicateWarning
              warning={duplicateWarning}
              onProceed={() => void saveContact(duplicateWarning.warningToken)}
              onCancel={() => setDuplicateWarning(null)}
            />
          ) : null}
        </form>
      )}
    </section>
  );
}
