# 标准应用审查报告

## 1. 审查范围

本报告基于当前有效文档进行项目内 agent 审查：

- `../../AGENTS.md`
- `AGENTS.md`
- `PROJECT_CONTEXT.md`
- `STANDARD-APPLICATION-INDEX.md`
- 当前 `docs/product/`、`docs/business/`、`docs/ux-ui/`、`docs/security/`
  下的有效文档
- 公司标准、工作流、默认微服务治理策略和相关模板

`archive/` 只作为历史参考，不作为当前架构、MDA、任务、测试、集成、
审计或发布权威。本报告不写实现代码，不决定最终服务拆分，不修改
P0/P1 等级或完成标准。

## 2. 当前 Gate 状态（更新于 2026-06-01）

- 当前阶段：G6 已通过，进入 G7 任务规划
- G5 Architecture Design：Gate Passed（2026-05-30）
- G6 MDA Modeling：Gate Passed（2026-06-01，六角色会签；
  `archive/reviews/g6-mda/g6-mda-gate-decision-2026-06-01.md`）
- 当前 Gate：G7 Task Planning（Domain Modeling + QA Test Design）
- G8 状态：未进入
- 实现状态：Blocked，G8 通过前不得实现

历史（2026-05-29 重置）：项目曾保留产品、业务、UX/UI、安全输入并废弃旧架构/
MDA/任务/实现/测试/部署/QA-集成-审计材料。此后已**重建并通过** G5 架构
(service-boundary-first，SVC-001..010) 与 G6 MDA 包（`modeling/`：CIM/PIM/PSM/
traceability-matrix/test-model），服务治理链路已补齐。下方各"补齐项/GAP"按
2026-06-01 进度核对。

## 3. Agent 后续补齐项

### product-manager

需要补齐：

- 在不改变 P0/P1 等级和完成标准的前提下，维护产品验收矩阵。
- 为 P0/P1 验收项补齐新标准要求的字段：
  - Business Capability
  - Related Services / Service Candidates
  - Service Owner Agent
  - Related Contracts
- 与 business-analyst 一起产出 P0/P1 business capability map。
- 将无法映射到能力、服务候选、合同或 owner 的 P0/P1 项标记为阻塞或
  待责任 agent 处理，不能标记 Done。
- 重新确认 OQ-001 与 ACC-017 的关系；旧架构结论已废弃。

不能做：

- 不能单独决定最终服务拆分。
- 不能通过改写验收标准来适配架构或实现。

### business-analyst

当前业务文档已有流程、规则、边界条件、权限场景和术语表，但还需要
面向 service-boundary-first 补齐：

- Business capability map：把 BP/BR/ACC 映射为业务能力。
- Cross-capability flow：标出跨能力流程，例如线索转商机、报价到合同、
  合同到回款、回款到赢单、任务提醒、导入导出、审计日志。
- Business event list：列出业务事件、触发条件、接收方、失败行为。
- Data responsibility notes：说明哪些业务能力负责哪些核心数据的
  业务含义和生命周期。
- 检查现有业务流程是否足够支持架构定义服务边界，但不直接定义技术服务。
- 将与 P0/P1 相关的未决业务规则反馈给 product-manager。

### ux-designer

当前 UX 文档覆盖用户旅程、流程、屏幕流、交互和状态，但还需要补齐：

- UX path 到 business capability 的映射。
- 跨能力用户路径的服务链证据需求，例如：
  - 线索到商机
  - 报价到合同
  - 合同到回款
  - 回款到赢单
  - 导入导出
  - 报表和日志
- 明确哪些 UX 状态需要后端合同支持：
  - loading
  - empty
  - validation error
  - permission denied
  - conflict / stale data
  - partial failure
  - long-running import/export
- 输出 UX 对 architecture、qa-tdd、integration-owner 的补充要求。

### ui-designer

当前 UI 文档覆盖页面、组件、响应式、视觉状态，但还需要补齐：

- UI state 到 business capability / acceptance item 的映射。
- 标出哪些组件状态依赖 API、权限、错误合同或长任务状态合同。
- 补齐服务支持相关的 UI 状态要求：
  - 授权失败和权限隐藏的差异
  - 导入导出进度、部分失败、下载过期
  - 报表空态、权限过滤、聚合不可见
  - 审计日志只读和敏感字段展示
- 明确 UI 不替代后端权限和业务规则。

### security-compliance

当前安全文档覆盖认证、权限、隐私、审计、滥用场景和合规风险，但还需要
按服务治理补齐：

- Service trust boundary notes。
- Service-to-service authentication / authorization 要求。
- 敏感跨服务数据流：
  - 用户/角色
  - 客户与联系人
  - 报价、合同、回款金额
  - 审计日志
  - 导入导出文件
  - 报表聚合
- 服务间调用的审计、correlation ID、最小载荷、错误安全摘要要求。
- 公网暴露、密钥、备份、日志、权限提升相关的阻塞条件。
- 将缺失服务信任边界或跨服务权限规则的 P0/P1 范围标为安全阻塞。

### architecture

这是当前下一阶段主责 agent。需要从当前有效输入重新设计，不得复用旧架构
作为权威。

需要补齐：

- `docs/architecture/architecture.md`
- `docs/architecture/module-boundaries.md`
- `docs/architecture/api-spec.md`
- `docs/architecture/data-design.md`
- `docs/architecture/integration-design.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/frontend-backend-contract.md`
- `docs/architecture/deployment-notes.md`
- `docs/architecture/risk-register.md`
- service architecture ADR
- service list with exactly one `Service Owner Agent` per service candidate
- service data ownership map
- API/event/error/permission contracts
- deployment boundaries
- observability strategy
- cross-service reliability strategy

必须处理：

- OQ-001：生产部署目标、域名、数据库、备份位置、环境所有权。
- 如果延迟物理微服务拆分，必须出 ADR；ADR 只能延迟物理拆分，不能移除
  服务边界、owner、合同、数据归属、测试、集成证据和审计追踪。

不能做：

- 不能替 product-manager 改产品范围。
- 不能替 business-analyst 改业务规则。
- 不能替 security-compliance 弱化安全要求。
- 不能在本报告中直接决定最终 CRM 服务拆分。

### domain-modeling

当前旧 MDA 已废弃。domain-modeling 必须等待新架构通过 G5 后重建。

需要补齐：

- `modeling/CIM.md`
- `modeling/PIM.md`
- `modeling/PSM.md`
- `modeling/domain-model.md`
- `modeling/state-machines.md`
- `modeling/domain-events.md`
- `modeling/service-mapping.md` 或 PSM 中等价服务映射章节
- `modeling/traceability-matrix.md`
- `modeling/test-model.md`

必须覆盖：

- P0/P1 acceptance -> business capability -> service/service candidate
- service owner agent
- aggregate ownership
- data ownership
- API/event/error/permission contracts
- state machines and failure paths
- service-to-service permission and reliability
- task/test/integration/audit traceability

不能复用旧 MDA 作为当前权威。

### task-planner

当前不能进入 task planning。task-planner 只能在新架构、MDA、traceability、
test model 准备好后工作。

需要补齐：

- `delivery/tasks.md`
- `delivery/task-dependencies.md`
- `delivery/delivery-plan.md`
- `delivery/acceptance-task-map.md`
- `planning/blockers.md` （治理，已存在于 `planning/`）

> 2026-06-01 决定：执行类任务产物放顶层 `delivery/`（与 `docs/` 设计、`modeling/`
> MDA、`planning/` 治理、`process/` 过程分离）；`blockers.md` 属治理，留在 `planning/`。

每个 service-backed 任务必须包含：

- business capability
- service / service candidate
- Service Owner Agent
- Primary Flow Owner Agent where applicable
- acceptance ID
- contract reference
- data ownership
- forbidden boundary access
- test requirements
- integration evidence plan
- MDA trace
- no-downgrade rule
- blocker record

### qa-tdd

当前 QA 文档已废弃，需要在新 MDA/test model 阶段重建。

需要补齐：

- P0/P1 acceptance test model。
- Service contract tests。
- Boundary tests。
- Permission and abuse-case tests。
- Persistence and recovery tests。
- Cross-service failure tests：
  - idempotency
  - retry
  - timeout
  - compensation
  - event duplication / ordering
  - correlation ID
- 测试 ID 到 ACC、service、contract、PSM 的映射。
- G8 前指出无法测试或合同缺失的 P0/P1 阻塞。

### integration-owner

当前不能做集成验证，但需要在 G8 前参与集成计划。

需要补齐：

- P0/P1 service-chain evidence plan。
- 每条集成链路的：
  - acceptance ID
  - business capability
  - services involved
  - contracts/events
  - persisted data evidence
  - permission checks
  - failure recovery
  - correlation ID
  - environment prerequisites
- ACC-017 的环境可达性、部署、配置、持久化、备份、恢复、烟测证据计划。

### infrastructure-ops

当前项目 `AGENTS.md` 未列出 `../../agents/infrastructure-ops.md`，但工作区
根入口和 `STANDARD-APPLICATION-INDEX.md` 已要求 infrastructure-ops 参与。
这需要在后续流程中纳入执行。

需要补齐：

- 部署环境需求记录。
- 服务器、数据库、域名、端口、公网入口、TLS、反向代理需求。
- secrets metadata 要求，不记录 secret 值。
- 备份与恢复要求。
- 监控与健康检查要求。
- 基础设施阻塞项与风险。
- 与 architecture 协作解决 OQ-001。

不能做：

- 不能决定 CRM 业务架构。
- 不能决定最终服务拆分。
- 不能绕过 security-compliance 做公网暴露或密钥相关决策。

### audit

当前 audit 不能基于旧 archive 作当前结论。需要准备新的反向审查计划。

需要补齐：

- 从 ACC-001 到 ACC-023 反向追踪：
  - business capability
  - service / service candidate
  - service owner agent
  - contract
  - architecture
  - PSM
  - task
  - test
  - integration evidence
  - infrastructure evidence where applicable
- 审查服务 owner 缺失、合同缺失、数据归属缺失、跨服务直接访问风险。
- 审查 no mock / no stub / no TODO / no static-only / no non-persistent。
- 为 G5/G6/G7/G8 准备阻塞判定标准。

## 4. G8 前必须关闭的 P0/P1 缺口清单

| ID | 严重级别 | 缺口 | 当前状态 | 责任 agent |
|---|---|---|---|---|
| GAP-G8-001 | P0 | P0/P1 business capability map 缺失 | 未产出独立能力映射 | product-manager, business-analyst |
| GAP-G8-002 | P0 | `acceptance-matrix.md` 缺少新标准字段：Business Capability、Related Services、Service Owner Agent、Related Contracts | 当前矩阵仍是旧字段结构 | product-manager |
| GAP-G8-003 | P0 | 每个软件支持的 P0/P1 验收项尚未映射 service/service candidate | 未映射 | product-manager, architecture |
| GAP-G8-004 | P0 | 每个 service candidate 尚未指定唯一 `Service Owner Agent` | 未定义服务候选和 owner | architecture |
| GAP-G8-005 | P0 | API/event/error/permission contracts 缺失 | 架构未重建 | architecture, security-compliance |
| GAP-G8-006 | P0 | Data ownership 和 forbidden cross-service access rules 缺失 | 架构未重建 | architecture, domain-modeling |
| GAP-G8-007 | P0 | Service trust boundaries 和 service-to-service 权限规则缺失 | 安全文档尚未服务化补齐 | security-compliance |
| GAP-G8-008 | P0 | OQ-001 重新打开，ACC-017 的生产部署目标和环境所有权未定 | Open / Reopened | architecture, infrastructure-ops |
| GAP-G8-009 | P0 | 新架构文档不存在 | `docs/architecture/` 当前不存在 | architecture |
| GAP-G8-010 | P0 | 新 MDA/PSM/service mapping/traceability 不存在 | `modeling/` 当前不存在 | domain-modeling |
| GAP-G8-011 | P0 | 新 test model 不存在 | 当前无有效 test model | qa-tdd, domain-modeling |
| GAP-G8-012 | P0 | 新 task plan 不存在，任务无法映射 service、owner、contract、acceptance、tests、boundaries | 当前无有效 task plan | task-planner |
| GAP-G8-013 | P1 | P0/P1 service-backed flows 的 integration evidence plan 缺失 | 当前无有效集成计划 | integration-owner |
| GAP-G8-014 | P1 | 基础设施环境需求、server/database/domain/port/backup/monitoring ownership 记录缺失 | OQ-001 未关闭 | infrastructure-ops, architecture |
| GAP-G8-015 | P1 | OQ-016 初始数据/迁移需求未定 | Open | product-manager, business-analyst |
| GAP-G8-016 | P0 | project `AGENTS.md` 未列出 infrastructure-ops，但当前标准要求该 agent 参与基础设施审查 | 项目入口不完整 | audit / project maintainer |

> **GAP 状态更新（2026-06-01）**：上表为 2026-05-29 重置期快照。随 G5（2026-05-30）
> 与 G6（2026-06-01）通过，以下已解决：
> - GAP-G8-007（服务信任边界/S2S 权限）— `docs/security/service-boundary-security.md`，G5 通过。
> - GAP-G8-008 / 部分 GAP-G8-014（OQ-001 生产部署目标）— `docs/architecture/deployment-notes.md` 已定；剩余为发布期证据（备份/恢复、TLS、安全组、监控），G11/G12 验证。
> - GAP-G8-009（架构文档）— `docs/architecture/` 已建，G5 通过。
> - GAP-G8-010（MDA/PSM/service mapping/traceability）— `modeling/` 已建（CIM/PIM/PSM/traceability-matrix），G6 通过。
> - GAP-G8-011（test model）— `modeling/test-model.md`，G6 通过。
>
> 仍 Open：GAP-G8-012（task plan，G7 产出，落 `delivery/`）、GAP-G8-013（集成证据计划，G7/G8）、GAP-G8-015（OQ-016 初始数据/迁移，发布规划）、GAP-G8-016（`AGENTS.md` 列入 infrastructure-ops）。

## 5. 建议的下一步执行顺序

1. product-manager + business-analyst：补 business capability map，并标出
   ACC-001 到 ACC-023 对应能力。
2. product-manager：更新 acceptance matrix 结构，加入服务治理字段；字段可先
   标为待 architecture 确认，但不能缺列。
3. business-analyst：补 cross-capability flow、business event list、data
   responsibility notes。
4. ux-designer：补 UX path 到 capability 的映射和服务链状态需求。
5. ui-designer：补 UI state 到 capability/contract support 的映射。
6. security-compliance：补 service trust boundary、跨服务权限、敏感数据流、
   secrets/public exposure/审计约束。
7. architecture：基于当前有效输入重新做 G5 架构设计；如延迟物理拆分，
   必须出 ADR。
8. infrastructure-ops：与 architecture 协作补环境需求和 OQ-001，不决定业务
   服务拆分。
9. G5 review：由 architecture 作为 Gate owner，PM/BA/UX/UI/Security 复审。
10. domain-modeling：G5 通过后重建 CIM/PIM/PSM/service mapping/traceability。
11. qa-tdd：与 domain-modeling 产出新 test model。
12. G7 review：确认 MDA、traceability、test model 覆盖 P0/P1 和服务治理。
13. task-planner + qa-tdd + integration-owner：生成 G8 任务、依赖、交付计划、
    acceptance-task map、blockers、integration evidence plan。
14. audit：对 G8 前链路做反向审查；缺服务 owner、合同、数据归属、测试或
    集成计划时阻塞 G8。

## 6. 审查结论

当前 CRM 项目不应进入 MDA、任务规划或实现。下一步不是写代码，也不是直接
决定最终服务拆分，而是按新标准补齐产品/业务/UX/UI/安全对服务治理的输入，
再由 architecture 重新完成 G5 架构设计。

## 7. 补充执行状态

2026-05-29 已完成一轮 G5 前置输入补强：

- `docs/product/business-capability-map.md` 已补 business capability map。
- `docs/product/acceptance-matrix.md` 已补 Business Capability、Related
  Services / Service Candidates、Service Owner Agent、Related Contracts 字段。
- `docs/business/service-governance-inputs.md` 已补跨能力流程、业务事件、
  数据责任输入。
- `docs/ux-ui/service-state-mapping.md` 已补 UX/UI 流程、状态与合同支持需求。
- `docs/security/service-boundary-security.md` 已补服务信任边界、跨服务安全
  输入和安全阻塞条件。
- `AGENTS.md` 已纳入 `infrastructure-ops`。

仍未关闭的 G8 前关键缺口：

- 最终服务边界、服务 owner、合同、数据归属仍需 Architecture 在 G5 设计中
  确认。
- OQ-001 仍需 Architecture + Infrastructure Ops 关闭。
- 新 MDA、PSM、traceability、test model、task plan、integration evidence
  plan 仍需在后续 Gate 中重建。
