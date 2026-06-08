import { Check, Circle, CircleDot } from 'lucide-react';
import { labelFor, opportunityStageLabel } from '../i18n/labels';

const stages = ['New Opportunity', 'Needs Confirmed', 'Quote', 'Contract Negotiation'];

export function StageStepper({ currentStage, terminal, onSelectStage }: { currentStage: string; terminal: boolean; onSelectStage: (stage: string) => void }) {
  const currentIndex = stages.indexOf(currentStage);
  const displayIndex = currentIndex >= 0 ? currentIndex : stages.length - 1;
  return (
    <div className="stageStepper stageStepperBranched" aria-label="商机阶段">
      {stages.map((stage, index) => {
        const complete = currentStage !== 'Lost' && index < displayIndex;
        const current = stage === currentStage;
        const blocked = terminal || (!current && index > displayIndex + 1);
        const Icon = current ? CircleDot : complete ? Check : Circle;
        return (
          <button
            type="button"
            key={stage}
            className={`stageStep ${current ? 'current' : ''} ${complete ? 'complete' : ''} ${blocked ? 'blocked' : ''}`}
            onClick={() => onSelectStage(stage)}
            disabled={terminal && !current}
            aria-current={current ? 'step' : undefined}
          >
            <Icon size={16} />
            <span>{labelFor(opportunityStageLabel, stage)}</span>
          </button>
        );
      })}
      <div className="stageOutcomes" aria-label="终态分支">
        <span className="outcomesLabel">终态分支</span>
        <span className={`outcome success ${currentStage === 'Won' ? 'current' : ''}`}>
          {labelFor(opportunityStageLabel, 'Won')}
          <small>合同谈判 + 已签合同</small>
        </span>
        <span className={`outcome danger ${currentStage === 'Lost' ? 'current' : ''}`}>
          {labelFor(opportunityStageLabel, 'Lost')}
          <small>需填写丢单原因</small>
        </span>
      </div>
    </div>
  );
}
