import { ManagedUser, UserStatus } from '../api/users';
import { UserRole } from '../api/auth';
import { labelFor, roleLabel, userStatusLabel } from '../i18n/labels';

type Props = {
  user: ManagedUser;
  nextRole: UserRole;
  nextStatus: UserStatus;
  onCancel: () => void;
  onConfirm: () => void;
};

export function RoleStatusChangeDialog({ user, nextRole, nextStatus, onCancel, onConfirm }: Props) {
  return (
    <section className="dialogPanel" role="dialog" aria-label="角色/状态变更确认">
      <h3>确认角色/状态变更</h3>
      <p>目标用户：{user.displayName}</p>
      <p>原角色：{labelFor(roleLabel, user.role)}</p>
      <p>新角色：{labelFor(roleLabel, nextRole)}</p>
      <p>原状态：{labelFor(userStatusLabel, user.status)}</p>
      <p>新状态：{labelFor(userStatusLabel, nextStatus)}</p>
      <p>访问影响：下一次受保护请求会重新评估该角色和状态。</p>
      <p>操作日志：本次确认变更会记录用于审计复核。</p>
      <div className="opportunityActions">
        <button className="secondaryButton" type="button" onClick={onCancel}>取消</button>
        <button className="primaryButton" type="button" onClick={onConfirm}>确认变更</button>
      </div>
    </section>
  );
}
