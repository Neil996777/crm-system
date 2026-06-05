import { Check, Circle, CircleDot } from 'lucide-react';
import { labelFor, opportunityStageLabel } from '../i18n/labels';

const stages = ['New Opportunity', 'Needs Confirmed', 'Quote', 'Contract Negotiation', 'Won', 'Lost'];

export function StageStepper({ currentStage, terminal, onSelectStage }: { currentStage: string; terminal: boolean; onSelectStage: (stage: string) => void }) {
  const currentIndex = stages.indexOf(currentStage);
  return (
    <div className="stageStepper" aria-label="商机阶段">
      {stages.map((stage, index) => {
        const complete = currentStage === 'Lost' ? false : index < currentIndex;
        const current = stage === currentStage;
        const next = !terminal && index === currentIndex + 1 && stage !== 'Won' && stage !== 'Lost';
        const blocked = !current && !next;
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
    </div>
  );
}
