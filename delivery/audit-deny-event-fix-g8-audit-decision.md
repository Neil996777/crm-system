# 审计「访问被拒」修复 —— G8 Handoff 审计决定(Claude)

| 字段 | 值 |
|---|---|
| 闸 | G8 设计/计划 → 实现 handoff |
| 审计人 | Claude(独立审计角色) |
| 日期 | 2026-06-13 |
| 尺子 | `delivery/audit-deny-event-fix-acceptance.md`(ACC-AUDIT-001..006 + C1–C6) |
| 被审包 | Codex G5/G8 产出(`bd9a81e`) |
| 判定 | **GATE PASSED** — 可进 G9 实现(M1–M6） |

## 被审交付物
- `docs/architecture/audit-deny-event-fix-design.md`(G5 设计:根因 + 目标契约 + 架构 + non-goals）
- `delivery/audit-deny-event-fix-g8-task-package.md`(M1–M6 任务包 + ACC/约束映射）
- `delivery/audit-deny-event-fix-traceability.md`(根因↔代码↔设计↔ACC 追溯）

## 根因核验(对真实代码,不轻信设计）
- ✅ `auth.go:262 appendAccessDenied` 实发 `{actorId,reason,result}`,缺 actorRole/actorDisplay;6 处调用属实。
- ✅ `auth.go:153-156 UserSignedIn` 带 actorId/actorRole/actorDisplay/role(健康参照)。
- ✅ `outbox.go:152` 有 `X-Actor-Role` fallback `System`(会伪装角色 → 设计正确要求修生产者而非靠 fallback)。
- ✅ `audit-history server.go` 校验失败 400 + 从 S2S 头取 actor(空则拒)。
- ✅ `user_admin.go:235 requireAdministrator` 调 `appendAccessDenied(...,"user_admin_denied")` 后 403。
设计建在准确根因上。

## 验收逐条核(ACC-AUDIT-001..006）
6 项验收均映射到 M1–M6 任务且有 fail-first 测试要求:001/002 完整信封入库+可见(M1/M2/M5)、003 无 400+outbox 排空(M1/M2)、004 结构化拒绝诊断(M3)、005 **真实 deny→操作日志 e2e、明禁 seeded fixture/skip/弱化**(M5)、006 哈希链 + 既有审计不降级(M4)。

## 绑定约束逐条核(C1–C6）
- **C1** 拒绝事件 `{}` before/after + safe summary,操作日志只渲染 safeSummary ✅
- **C2** 走 `EventRepo.Append`、不改/删已入行、测试验 `prevHash==上条 eventHash` ✅
- **C3(红线)** 403/401/`allowed=false` 全保留、不改角色门控/enum/role 比较值、**audit-only actor label 绝不进域授权比较**、forbidden 清单明确 ✅
- **C4** 新标签走 `labels.ts` ✅ · **C5** 既有事件继续正常入库、no-downgrade ✅ · **C6** 走完整闸、G8 前无实现代码、Codex 不自判 ✅

## 范围核
保持 audit-history 校验**严格**(修生产者、不弱化契约)——方向正确;无应用实现代码先行(G8 前);改动面限于 identity-authz/audit-history/测试/labels.ts/可选 shared 契约,无授权语义变更。

## 决定
**G8 GATE PASSED。** 设计 + 任务包对齐尺子,根因经真实代码核实,6 验收 + 6 约束全覆盖、C3 红线守得严,无实现代码先行。Codex 可进 **G9** 实现 M1–M4 → G10/G11 产 M5/M6 证据 → Claude **G12** 独立审计(触发真实拒绝看操作日志、outbox 排空、链完整、无授权 diff)再决定 BLK-PROD-AUDIT-001 关闭。
