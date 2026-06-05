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
        <strong>可能重复</strong>
        <p>{warning.matches.length > 0 ? warning.matches.map((match) => match.safeSummary).join(', ') : '可能已存在相似记录。'}</p>
      </div>
      <div className="warningActions">
        <button className="primaryButton" type="button" onClick={onProceed}>仍然创建</button>
        <button className="secondaryButton" type="button" onClick={onCancel}>返回检查</button>
      </div>
    </div>
  );
}
