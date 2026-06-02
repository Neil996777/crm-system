import { DuplicateWarningResult } from '../api/duplicates';

type Props = {
  warning: DuplicateWarningResult;
  onProceed: () => void;
  onCancel: () => void;
};

export function DuplicateWarning({ warning, onProceed, onCancel }: Props) {
  return (
    <div role="alert" className="duplicateWarning">
      <div>
        <strong>Possible duplicate</strong>
        <p>{warning.matches.length > 0 ? warning.matches.map((match) => match.safeSummary).join(', ') : 'A similar record may already exist.'}</p>
      </div>
      <div className="warningActions">
        <button className="primaryButton" type="button" onClick={onProceed}>Create anyway</button>
        <button className="secondaryButton" type="button" onClick={onCancel}>Review input</button>
      </div>
    </div>
  );
}
