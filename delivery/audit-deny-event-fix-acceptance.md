# 审计「访问被拒」事件修复 —— 验收矩阵 + 绑定约束(Yardstick)

> 角色:Claude 本职 = yardstick(需求/验收/绑定约束,G1–G3)。本文件是 `BLK-PROD-AUDIT-001`
> 修复(独立变更)的验收源头尺子,也是 G8 handoff 审计与 G12 审计的判定依据。
> 缘起:2026-06-13 在测试服务器 `srv-volcengine-sh-01` 上,管理员创建用户时偶发 502,排查发现
> `UserAccessDenied` 审计事件被 audit-history 拒绝(HTTP 400),"访问被拒"事件**从未进入操作日志**。
> 权威:`docs/product/acceptance-matrix.md`(ACC-022 全局操作日志)、基线 G12 安全不变量(safe-summary、
> 审计哈希链)、`company/operating-model.md`(no-downgrade)。
> 范围声明:**这是应用层缺陷,独立于 CI/CD 迁移**,不复用其"仅机制/0 应用源码"约束;走完整 G 闸。

---

## 1. WHAT / WHY(需求)

**WHAT.** 修复审计事件管线,使 identity-authz 发出的**所有**审计相关事件(尤其 `UserAccessDenied`)
携带 audit-history 入库所需的**完整信封**,从而"访问被拒"等安全事件能**可靠落入操作日志**。

**WHY.** "访问被拒"是最该审计的安全类别。现状:`UserAccessDenied` 事件 payload 仅 `{reason, result, actorId}`,
缺 `actorRole`/`actorDisplay` 等信封必填字段(对比 `UserSignedIn` 携带这些字段、可正常入库)→ audit-history
校验返回 **400** → 这些事件永不入库,outbox 每 2 秒重试刷屏。后果:**操作日志(ACC-022)对"访问被拒"存在静默
丢失** = 安全/合规审计窟窿;且 audit-history 拒绝时不记原因(可诊断性缺口)。这是基线产品(`66d2531`)潜伏缺陷,
本次部署 + 拒绝流量将其暴露。

**边界**:本变更修审计**完整性 + 契约校验/可诊断性**;不改业务功能、不放宽任何权限门控。

---

## 2. 验收矩阵(可测 Definition of Done)

| ID | P | 能力(Done 定义) | 验证方法 |
|---|---|---|---|
| ACC-AUDIT-001 | P0 | **拒绝事件可入库**:任一"访问被拒"动作(如 Sales 试图执行管理员专属操作)产生的 `UserAccessDenied` 审计事件被 audit-history **接受(2xx)** 并出现在操作日志中 | e2e:触发拒绝 → 管理员操作日志可见该事件 |
| ACC-AUDIT-002 | P0 | **信封完整**:落库的拒绝事件含 actor 标识 + **actor 角色** + 动作 + 结果 + 原因(reason_code),字段齐全且语义正确 | 审 audit_history.events 行 + 操作日志展示 |
| ACC-AUDIT-003 | P0 | **无 400 / outbox 可排空**:良构域事件的审计 append 不再返回 400;identity-authz outbox 无"毒丸"积压(待发计数随时间归零) | 集成测试 + 运行期 outbox 待发计数 |
| ACC-AUDIT-004 | P1 | **拒绝可诊断**:audit-history 拒绝一个事件时记录**结构化原因**(哪个字段/为何),不再静默 400 | 审日志:故意发不合规事件 → 有原因记录 |
| ACC-AUDIT-005 | P0 | **覆盖测试补齐**:新增 deny→audit 路径的 e2e/集成测试(此前缺口:"访问被拒"无审计测试),纳入套件且稳定绿 | 测试存在 + 套件通过(workers:2/retries:1,0 fail) |
| ACC-AUDIT-006 | P0 | **不降级既有审计**:`UserSignedIn`/`UserRoleStatusChanged` 等已正常入库的事件继续正常,审计哈希链(prev_hash/event_hash)完整、可验证 | 回归:既有审计 e2e 全绿 + 链完整性校验 |

---

## 3. 绑定约束(C,红线)

- **C1 保 safe-summary 安全不变量(承基线 G12)**:操作日志只渲染 `safeSummary`,**不得**出现 raw before/after。
  拒绝事件落库同样只带安全摘要,不泄敏感数据。
- **C2 保审计哈希链**:prev_hash/event_hash 防篡改链不得破坏;修复后链仍可端到端验证。
- **C3 不放宽权限**:本变更**只补审计字段/契约**,**不得**改变任何角色门控、enum/role 比较值或访问决策本身
  (访问仍照常被拒,只是这次"拒绝"被正确审计)。
- **C4 保 zh-CN**:操作日志的拒绝事件展示走 `labels.ts`,中文显示正确。
- **C5 no-downgrade**:不弱化任何既有审计覆盖或 P0/P1 验收;本变更是**增强**。
- **C6 走完整闸**:这是应用变更,**不复用 CI/CD 迁移的"仅机制"豁免**;G2–G8 设计/任务由 Codex 产、Claude 审 G8,
  G9–G11 Codex 执行,Claude G12;Codex 不自判闸。

---

## 4. 不在范围(Out of Scope)

- CI/CD 发布机制(已 G12 PASS,本变更只改应用代码,经同一流水线发布)。
- 重新设计审计事件模型/操作日志 UI(只补缺失字段 + 校验诊断,不重构)。
- 历史已丢失的拒绝事件追溯(测试环境,无需回填;真生产上线前修好即可)。

---

## 5. 闸路径(建议)

G1–G3 Yardstick(本文件,Claude,**待 release owner 确认**)→ G5/G7/G8 设计 + 任务包(Codex 产,Claude 审 G8 handoff)
→ G9–G11(Codex 实现 + QA + 集成)→ G12(Claude 独立审计:逐条核 ACC-AUDIT-001..006 + C1–C6,含触发拒绝看操作日志、
outbox 排空、链完整)。实现代码在 G8 通过前不得开始。

> 关联:`planning/blockers.md` BLK-PROD-AUDIT-001;根因诊断证据见该 blocker(`UserAccessDenied` payload 缺
> `actorRole`/`actorDisplay`,audit-history 400)。
