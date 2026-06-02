import { ArchiveEligibility } from '../api/accounts';

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
        <strong>Active obligations block archive</strong>
        <ul>
          {eligibility.obligations.map((obligation) => (
            <li key={`${obligation.service}-${obligation.id}`}>{obligation.safeMessage || obligation.type}</li>
          ))}
        </ul>
      </div>
    );
  }
  return (
    <div className="archivePanel">
      <label>
        Archive reason
        <input value={reason} onChange={(event) => onReasonChange(event.target.value)} />
      </label>
      <div className="warningActions">
        <button className="primaryButton" type="button" onClick={onConfirm}>Confirm archive</button>
        <button className="secondaryButton" type="button" onClick={onCancel}>Cancel</button>
      </div>
    </div>
  );
}
