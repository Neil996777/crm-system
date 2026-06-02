import { ManagedUser, UserStatus } from '../api/users';
import { UserRole } from '../api/auth';

type Props = {
  user: ManagedUser;
  nextRole: UserRole;
  nextStatus: UserStatus;
  onCancel: () => void;
  onConfirm: () => void;
};

export function RoleStatusChangeDialog({ user, nextRole, nextStatus, onCancel, onConfirm }: Props) {
  return (
    <section className="dialogPanel" role="dialog" aria-label="Role/status change confirmation">
      <h3>Confirm role/status change</h3>
      <p>Target user: {user.displayName}</p>
      <p>Old role: {user.role}</p>
      <p>New role: {nextRole}</p>
      <p>Old status: {user.status}</p>
      <p>New status: {nextStatus}</p>
      <p>Access impact: the next protected request re-evaluates this role and status.</p>
      <p>Operation log: this confirmed change is recorded for audit review.</p>
      <div className="opportunityActions">
        <button className="secondaryButton" type="button" onClick={onCancel}>Cancel</button>
        <button className="primaryButton" type="button" onClick={onConfirm}>Confirm change</button>
      </div>
    </section>
  );
}
