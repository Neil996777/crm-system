import { ArchiveEligibility } from '../api/accounts';
import { localizeMessage } from '../i18n/labels';

type Props = {
  eligibility: ArchiveEligibility | null;
  reason: string;
  onReasonChange: (reason: string) => void;
  onConfirm: () => void;
  onCancel: () => void;
};

export function ArchiveConfirmation({ eligibility, reason, onReasonChange, onConfirm, onCancel }: Props) {
  if (!eligibility) return null;
  if (!eligibility.canArchive) {
    return (
      <div role="alert" className="archivePanel">
        <strong>仍有未完成事项，不能归档</strong>
        <ul>
          {eligibility.obligations.map((obligation) => (
            <li key={`${obligation.service}-${obligation.id}`}>{localizeMessage(obligation.safeMessage, obligation.type)}</li>
          ))}
        </ul>
      </div>
    );
  }
  return (
    <div className="archivePanel">
      <label>
        归档原因
        <input value={reason} onChange={(event) => onReasonChange(event.target.value)} />
      </label>
      <div className="warningActions">
        <button className="primaryButton" type="button" onClick={onConfirm}>确认归档</button>
        <button className="secondaryButton" type="button" onClick={onCancel}>取消</button>
      </div>
    </div>
  );
}
