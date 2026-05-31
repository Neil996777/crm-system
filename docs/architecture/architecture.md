# Architecture

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30
- Source Inputs:
  - `PROJECT_CONTEXT.md`
  - `STANDARD-APPLICATION-INDEX.md`
  - `docs/product/prd.md`
  - `docs/product/acceptance-matrix.md`
  - `docs/product/business-capability-map.md`
  - `docs/business/service-governance-inputs.md`
  - `docs/ux-ui/service-state-mapping.md`
  - `docs/security/service-boundary-security.md`

## Overview

The CRM System will use a physical multi-service architecture for the committed release.

The deployment runtime host is one Volcengine ECS server (`srv-volcengine-sh-01`)
running Docker Compose, with a second Alibaba Cloud ECS server (`srv-aliyun-bj-01`)
as the off-server backup target only. Each Go backend service runs in its own
Docker container. PostgreSQL is self-hosted in a Docker container on the runtime
host. A reverse proxy exposes only the web/API entrypoint. Backend services
communicate only on the internal Docker network. See `deployment-notes.md` for
the registered assets and runtime-host co-location constraints.

The architecture follows service-boundary-first governance and physically
separates service runtime containers from the start. Service boundaries are
based on business capabilities and DDD bounded-context candidates, not pages,
controllers, tables, or technical layers.

## Architecture Goals

- Support the full P0/P1 ToB CRM loop without mock, static-only, TODO,
  in-memory-only, or non-persistent core paths.
- Preserve high cohesion by keeping each service responsible for one business
  capability area.
- Preserve low coupling by requiring public API or event contracts for all
  cross-service collaboration.
- Enforce database ownership through independent PostgreSQL database or schema
  permissions per service.
- Provide architecture artifacts that Domain Modeling can represent in PSM.
- Provide contract, data, permission, integration, and deployment constraints
  for Task Planner, QA TDD, Integration Owner, and Audit.

## Constraints

- Product: P0/P1 acceptance items in `docs/product/acceptance-matrix.md` are the
  source of truth and cannot be downgraded, deleted, merged away, weakened, or
  accepted as partial work.
- Business: the system must support the complete committed ToB CRM loop from lead to
  customer/contact, opportunity, quote, contract, payment, closure, activity,
  reminders, history, reporting, and import/export.
- UX/UI: API and error contracts must support loading, empty, validation,
  permission denied, blocked transition, conflict, partial failure, long-running
  operation, read-only history, sensitive display, and archived context states.
- Security: frontend hiding or disabling actions is not authorization. Backend
  services must enforce role, record scope, service-to-service authorization,
  safe errors, auditability, and sensitive data handling.
- Technical: backend services use Go; runtime isolation uses Docker containers;
  orchestration uses Docker Compose for the committed release; persistence uses PostgreSQL.
- Operations: production runtime host is Volcengine ECS (`srv-volcengine-sh-01`)
  with Alibaba Cloud ECS (`srv-aliyun-bj-01`) as the off-server backup target;
  PostgreSQL local backups are automatic, timestamped, encrypted, and retained
  for 7 days. Production release also requires off-server backup copy evidence to
  the Alibaba host; same-host-only backup is a release-blocking gap, not an
  accepted P0/P1 completion state.
- Public access: production login/session traffic is HTTPS-only. IP-based
  internal validation is allowed before release, but production ACC-017 evidence
  must record the final domain or approved endpoint, TLS certificate, security
  group rules, health checks, monitoring target, and backup evidence.

## Service Architecture Strategy

| Decision | Value | Rationale |
|---|---|---|
| Service strategy | Physical multi-service from the committed release | User selected multiple Go microservices with Docker isolation. |
| Runtime orchestration | Docker Compose on Volcengine ECS runtime host `srv-volcengine-sh-01` | Fits current infrastructure and team scale; avoids Kubernetes complexity in the committed scope. |
| Database deployment | Self-hosted PostgreSQL Docker container on the runtime host | User selected self-hosted database on the runtime host resource. |
| Data isolation | One PostgreSQL instance, service-isolated database or schema plus service account | Balances operational simplicity with service ownership and forbidden cross-service data access. |
| External access | Reverse proxy exposes only web/API entrypoint | Reduces public attack surface and keeps service network internal. |
| Backup | Local automatic encrypted daily backup, timestamped, retain 7 days, plus off-server production backup requirement | Local backup matches user-selected baseline; same-host-only backup remains a production release blocker. |
| Production entry | HTTPS-only through reverse proxy | Login/session traffic cannot use plaintext HTTP in production. HTTP may only redirect to HTTPS. |
| Internal trust | Authenticated service-to-service calls | Docker internal network is a transport boundary, not a trust boundary. |

## Service List

| Service ID | Service | Bounded Context / Capability | Service Owner Agent | Deployment Boundary | Primary Acceptance IDs |
|---|---|---|---|---|---|
| SVC-001 | gateway-bff | API entry, frontend aggregation, request routing | backend-engineer | Independent container | ACC-001 to ACC-023 |
| SVC-002 | identity-authz-service | Identity and role access | backend-engineer | Independent container | ACC-001, ACC-002, ACC-022 |
| SVC-003 | lead-service | Lead intake and qualification | backend-engineer | Independent container | ACC-003, ACC-004, ACC-019 |
| SVC-004 | account-service | Account and contact management | backend-engineer | Independent container | ACC-005, ACC-006, ACC-019 |
| SVC-005 | opportunity-service | Opportunity pipeline | backend-engineer | Independent container | ACC-007, ACC-008, ACC-013 |
| SVC-006 | commercial-service | Quote, contract, and payment lifecycle | backend-engineer | Independent container | ACC-009, ACC-010, ACC-011, ACC-013, ACC-021 |
| SVC-007 | work-service | Activities, notes, tasks, and reminders | backend-engineer | Independent container | ACC-012, ACC-021 |
| SVC-008 | audit-history-service | Record-local history and operation logs | backend-engineer | Independent container | ACC-014, ACC-022 |
| SVC-009 | reporting-service | Team overview, reports, and read models | backend-engineer | Independent container | ACC-018, ACC-023 |
| SVC-010 | import-export-service | CSV import/export runs and row results | backend-engineer | Independent container | ACC-020, ACC-022 |

Every service currently has exactly one `Service Owner Agent`. G8 task planning
may introduce more specialized project agents later, but cannot remove service
ownership.

## Topology

```mermaid
flowchart LR
  Browser[Browser 浏览器]
  Proxy[Reverse Proxy 反向代理<br/>HTTPS only 仅HTTPS<br/>Nginx or Caddy]
  BFF[gateway-bff 网关]

  subgraph ECS[Volcengine ECS 运行机 srv-volcengine-sh-01 / Docker Compose]
    Proxy
    BFF
    Auth[identity-authz-service 身份鉴权服务]
    Lead[lead-service 线索服务]
    Account[account-service 客户服务]
    Opp[opportunity-service 商机服务]
    Comm[commercial-service 商务服务]
    Work[work-service 工作服务]
    Audit[audit-history-service 审计历史服务]
    Report[reporting-service 报表服务]
    Import[import-export-service 导入导出服务]
    PG[(PostgreSQL<br/>service-owned database/schema 各服务独占库/schema)]
  Backup[backup job 备份任务<br/>encrypted daily local 每日本地加密<br/>retain 7 days 留存7天]
  end

  OffHost[(off-server backup target 异机备份目标<br/>srv-aliyun-bj-01 Alibaba Cloud Beijing 阿里云北京<br/>required before production release 发布前必需)]

  Browser --> Proxy --> BFF
  BFF --> Auth
  BFF --> Lead
  BFF --> Account
  BFF --> Opp
  BFF --> Comm
  BFF --> Work
  BFF --> Audit
  BFF --> Report
  BFF --> Import

  Auth --> PG
  Lead --> PG
  Account --> PG
  Opp --> PG
  Comm --> PG
  Work --> PG
  Audit --> PG
  Report --> PG
  Import --> PG
  Backup --> PG
  Backup -. production release evidence 生产发布证据 .-> OffHost
```

Only the reverse proxy and web/API entrypoint are exposed outside the host.
Database and internal service ports are not publicly exposed.

## Service Data Boundary Diagram

```mermaid
flowchart TB
  Auth[identity-authz-service 身份鉴权服务] --> AuthDB[(identity_authz db/schema 库)]
  Lead[lead-service 线索服务] --> LeadDB[(lead db/schema 库)]
  Account[account-service 客户服务] --> AccountDB[(account db/schema 库)]
  Opp[opportunity-service 商机服务] --> OppDB[(opportunity db/schema 库)]
  Comm[commercial-service 商务服务] --> CommDB[(commercial db/schema 库)]
  Work[work-service 工作服务] --> WorkDB[(work db/schema 库)]
  Audit[audit-history-service 审计历史服务] --> AuditDB[(audit_history db/schema 库)]
  Report[reporting-service 报表服务] --> ReportDB[(reporting db/schema 库)]
  Import[import-export-service 导入导出服务] --> ImportDB[(import_export db/schema 库)]

  Lead -. Query/Command API 查询/命令API .-> Account
  Lead -. Query/Command API 查询/命令API .-> Opp
  Opp -. Query API 查询API .-> Comm
  Work -. Query API 查询API .-> Account
  Work -. Query API 查询API .-> Opp
  Work -. Query API 查询API .-> Comm
  Import -. Command/Query API 命令/查询API .-> Lead
  Import -. Command/Query API 命令/查询API .-> Account
  Import -. Command/Query API 命令/查询API .-> Opp
  Import -. Command/Query API 命令/查询API .-> Comm
  Report -. Event/API projection 事件/API投影 .-> Lead
  Report -. Event/API projection 事件/API投影 .-> Account
  Report -. Event/API projection 事件/API投影 .-> Opp
  Report -. Event/API projection 事件/API投影 .-> Comm
```

Each solid arrow is a service accessing only its own database/schema. Each
dotted arrow is a public business API, event projection, or approved read-model
interaction. Dotted arrows are not database access.

## Core Service Flow Matrix

| Flow ID | Flow | Primary Services | Supporting Services | Acceptance IDs | Detail |
|---|---|---|---|---|---|
| ARCH-FLOW-001 | Sign in and protected work | gateway-bff, identity-authz-service | target service, audit-history-service | ACC-001, ACC-002, ACC-022 | See sequence below. |
| ARCH-FLOW-002 | Lead to opportunity | lead-service, account-service, opportunity-service | gateway-bff, identity-authz-service, audit-history-service, reporting-service | ACC-003 to ACC-007, ACC-014, ACC-019 | See sequence below. |
| ARCH-FLOW-003 | Opportunity to quote and contract | opportunity-service, commercial-service | identity-authz-service, audit-history-service, reporting-service | ACC-007 to ACC-010, ACC-014 | See sequence below. |
| ARCH-FLOW-004 | Payment to Won | commercial-service, opportunity-service | identity-authz-service, audit-history-service, reporting-service | ACC-011, ACC-013, ACC-014 | See sequence below. |
| ARCH-FLOW-005 | Work reminders | work-service, commercial-service | record-owning services, identity-authz-service | ACC-012, ACC-021 | See sequence below. |
| ARCH-FLOW-006 | Import/export | import-export-service, target domain services | identity-authz-service, audit-history-service | ACC-016, ACC-020, ACC-022 | See sequence below. |
| ARCH-FLOW-007 | Reports and overview | reporting-service | source service events/APIs, identity-authz-service | ACC-018, ACC-023 | See sequence below. |
| ARCH-FLOW-008 | Backup and restore | PostgreSQL, backup job | runtime services, infrastructure-ops | ACC-016, ACC-017, ACC-022 | See sequence below. |
| ARCH-FLOW-009 | Archive eligibility and active obligations | record-owning services | work-service, commercial-service, audit-history-service | ACC-005, ACC-007, ACC-010, ACC-012, ACC-014 | See sequence below. |
| ARCH-FLOW-010 | Owner transfer and open work transfer | lead/account/opportunity services, work-service | identity-authz-service, audit-history-service | ACC-003, ACC-005, ACC-007, ACC-012, ACC-014 | See sequence below. |
| ARCH-FLOW-011 | Close Lost and terminal edit protection | opportunity-service | work-service, audit-history-service, reporting-service | ACC-013, ACC-014, ACC-021 | See sequence below. |
| ARCH-FLOW-012 | Duplicate warning and proceed-after-warning | lead-service, account-service | gateway-bff, identity-authz-service | ACC-019 | See sequence below. |

## Sequence: Sign In And Protected Work

```mermaid
sequenceDiagram
  participant U as User 用户
  participant G as gateway-bff 网关
  participant A as identity-authz-service 身份鉴权服务
  participant T as target domain service 目标领域服务
  participant H as audit-history-service 审计历史服务

  U->>G: Sign in 登录
  G->>A: Authenticate credentials 校验凭证
  A->>A: Validate user, role, active status 校验用户/角色/启用状态
  A-->>H: UserSignedIn or UserAccessDenied event 事件:登录成功或拒绝访问
  A-->>G: Session and role context 会话与角色上下文
  U->>G: Protected action 受保护操作
  G->>A: Check permission(actor, action, resource) 校验权限(操作者,动作,资源)
  A-->>G: Allowed or denied 允许或拒绝
  alt allowed 允许时
    G->>T: Forward command/query with actor and correlationId 转发命令/查询(带操作者与correlationId)
    T->>T: Enforce domain permission and business rule 校验领域权限与业务规则
    T-->>G: Result 结果
  else denied 拒绝时
    G-->>U: Safe permission denial 安全的权限拒绝
  end
```

## Sequence: Lead To Opportunity

```mermaid
sequenceDiagram
  participant U as Sales 销售
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant L as lead-service 线索服务
  participant A as account-service 客户服务
  participant O as opportunity-service 商机服务
  participant H as audit-history-service 审计历史服务
  participant R as reporting-service 报表服务

  U->>G: Convert lead 转换线索
  G->>P: Check permission 校验权限
  P-->>G: Allowed 允许
  G->>L: ConvertLead(command, idempotencyKey) 转换线索(命令,幂等key)
  L->>L: Validate state and conversion-once guard 校验状态与"仅转换一次"约束
  L->>A: CreateOrLinkAccountContact 创建或关联客户/联系人
  A-->>L: Account/contact references 客户/联系人引用
  L->>O: CreateOpportunity 创建商机
  O-->>L: Opportunity reference 商机引用
  L->>L: Persist converted state 持久化已转换状态
  L-->>H: LeadConverted event 事件:线索已转换
  O-->>H: OpportunityCreated event 事件:商机已创建
  L-->>R: LeadConverted event 事件:线索已转换
  O-->>R: OpportunityCreated event 事件:商机已创建
  L-->>G: Conversion result 转换结果
  G-->>U: Converted lead and opportunity link 已转换线索与商机关联
```

## Sequence: Opportunity To Quote And Contract

```mermaid
sequenceDiagram
  participant U as Sales 销售
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant O as opportunity-service 商机服务
  participant C as commercial-service 商务服务(报价/合同/回款)
  participant H as audit-history-service 审计历史服务
  participant R as reporting-service 报表服务

  U->>G: Advance opportunity to quote work 推进商机到报价阶段
  G->>P: Check opportunity permission 校验商机权限
  P-->>G: Allowed 允许
  G->>O: ChangeStage(command) 变更阶段(命令)
  O->>O: Validate allowed transition 校验允许的状态转移
  O-->>H: OpportunityStageChanged event 事件:商机阶段已变更
  O-->>R: OpportunityStageChanged event 事件:商机阶段已变更

  U->>G: Create or accept quote 创建或接受报价
  G->>P: Check commercial permission 校验商务权限
  P-->>G: Allowed 允许
  G->>C: Quote command 报价命令
  C->>O: GetOpportunitySummary 获取商机摘要
  O-->>C: Authorized opportunity summary 已授权的商机摘要
  C->>C: Validate quote status and accepted uniqueness 校验报价状态与"唯一已接受"约束
  C-->>H: QuoteAccepted event where applicable 事件:报价已接受(如适用)
  C-->>R: QuoteAccepted event where applicable 事件:报价已接受(如适用)

  U->>G: Create contract from accepted quote 由已接受报价创建合同
  G->>C: Contract command 合同命令
  C->>C: Validate accepted quote, dates, note, amount difference reason 校验已接受报价/日期/备注/金额差异原因
  C-->>H: ContractStatusChanged event 事件:合同状态已变更
  C-->>G: Contract result 合同结果
```

## Sequence: Payment To Won

```mermaid
sequenceDiagram
  participant U as Sales 销售
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant C as commercial-service 商务服务(报价/合同/回款)
  participant O as opportunity-service 商机服务
  participant H as audit-history-service 审计历史服务
  participant R as reporting-service 报表服务

  U->>G: Record payment 记录回款
  G->>P: Check commercial permission 校验商务权限
  P-->>G: Allowed 允许
  G->>C: RecordPayment(command, idempotencyKey) 记录回款(命令,幂等key)
  C->>C: Validate amount and overpayment 校验金额与超额支付
  C-->>H: PaymentRecorded event 事件:回款已记录
  C-->>R: PaymentRecorded event 事件:回款已记录
  C-->>G: Payment status 回款状态

  U->>G: Close opportunity Won 关闭商机为赢单
  G->>P: Check opportunity close permission 校验商机关闭权限
  P-->>G: Allowed 允许
  G->>O: CloseWon(command, idempotencyKey) 赢单关闭(命令,幂等key)
  O->>C: GetPaymentStatusSummary 获取回款状态摘要
  C-->>O: Paid / not paid 已全额/未全额
  O->>O: Persist terminal Won if fully paid 全额则持久化终态"赢单"
  O-->>H: OpportunityClosed event 事件:商机已关闭
  O-->>R: OpportunityClosed event 事件:商机已关闭
  O-->>G: Won result 赢单结果
  G-->>U: Won state 赢单状态
```

## Sequence: Work Reminders

```mermaid
sequenceDiagram
  participant U as Sales 销售
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant W as work-service 工作服务(活动/任务/提醒)
  participant C as commercial-service 商务服务(报价/合同/回款)
  participant T as target record service 目标记录服务

  U->>G: Open reminder area 打开提醒区
  G->>P: Check reminder permission 校验提醒权限
  P-->>G: Allowed 允许
  G->>W: Query reminders(actor scope) 查询提醒(按操作者范围)
  W->>T: Get safe related-record summaries 获取安全的关联记录摘要
  T-->>W: Authorized summaries or hidden 已授权摘要或隐藏
  W->>C: Get due contract/payment reminder eligibility 获取到期合同/回款的提醒资格
  C-->>W: Authorized due/overdue summary 已授权的到期/逾期摘要
  W->>W: Exclude completed, cancelled, signed, terminated, fully paid, archived, or unauthorized items 排除已完成/已取消/已签署/已终止/已全额/已归档/未授权项
  W-->>G: Reminder list 提醒列表
  G-->>U: Authorized reminders 已授权的提醒
```

## Sequence: Import / Export

```mermaid
sequenceDiagram
  participant U as Admin/Manager 管理员/经理
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant I as import-export-service 导入导出服务
  participant T as target domain service 目标领域服务
  participant H as audit-history-service 审计历史服务

  U->>G: Start CSV import/export 启动CSV导入/导出
  G->>P: Check import/export permission 校验导入导出权限
  P-->>G: Allowed 允许
  G->>I: Start run(file or export criteria) 启动运行(文件或导出条件)
  I->>I: Validate format and scope 校验格式与范围
  loop each import row or export page 每个导入行或导出分页
    I->>T: Domain command/query with idempotency key or scope 领域命令/查询(带幂等key或范围)
    T-->>I: Success, authorized data, or safe row error 成功/已授权数据/安全的行级错误
  end
  I-->>H: ImportRunCompleted or ExportRunCompleted event 事件:导入或导出运行完成
  I-->>G: Run result summary 运行结果摘要
  G-->>U: Row results or export metadata 行级结果或导出元数据
```

## Sequence: Reports And Overview

```mermaid
sequenceDiagram
  participant S as Source domain services 源领域服务
  participant R as reporting-service 报表服务
  participant U as Administrator/Manager 管理员/经理
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务

  S-->>R: Domain events for report projection 用于报表投影的领域事件
  R->>R: Update owned read model 更新自有读模型
  U->>G: Open overview/report 打开总览/报表
  G->>P: Check report permission 校验报表权限
  P-->>G: Allowed with scope 允许(带范围)
  G->>R: Query report(scope, filters) 查询报表(范围,过滤)
  R->>R: Apply authorization before aggregate response 聚合返回前先做授权
  R-->>G: Authorized metrics 已授权的指标
  G-->>U: Overview/report result 总览/报表结果
```

## Sequence: Backup And Restore Evidence

```mermaid
sequenceDiagram
  participant B as backup job 备份任务
  participant DB as PostgreSQL 数据库
  participant FS as ECS local encrypted backup directory 运行机本地加密备份目录
  participant X as Off-server backup target 异机备份目标
  participant O as Operator/Infrastructure Ops 运维/基础设施
  participant S as Runtime services 运行期服务

  B->>DB: Create timestamped backup 创建带时间戳的备份
  DB-->>B: Backup stream 备份数据流
  B->>FS: Write new backup file 写入新备份文件
  B->>FS: Delete files older than 7 days 删除7天前的旧文件
  alt production release evidence 生产发布证据
    B->>X: Copy encrypted backup off-server 加密备份异机拷贝
    X-->>B: Stored with timestamp and checksum 已存(带时间戳与校验和)
  else pre-release local-only 预发布仅本地
    B-->>O: Mark release blocker BACKUP_OFF_SERVER_MISSING 标记发布阻塞:缺异机备份
  end
  O->>FS: Select backup for restore rehearsal 选取备份做恢复演练
  O->>DB: Restore into controlled target/procedure 恢复到受控目标/流程
  O->>S: Verify health, data, history, logs, and DB permissions 校验健康/数据/历史/日志/库权限
  S-->>O: Restore evidence result 恢复证据结果
```

## Sequence: Archive Eligibility And Active Obligations

```mermaid
sequenceDiagram
  participant U as Admin/Manager 管理员/经理
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant R as record-owning service 记录所属服务
  participant W as work-service 工作服务(活动/任务/提醒)
  participant C as commercial-service 商务服务(报价/合同/回款)
  participant H as audit-history-service 审计历史服务

  U->>G: Request archive record 请求归档记录
  G->>P: Check archive permission 校验归档权限
  P-->>G: Allowed 允许
  G->>R: GetArchiveEligibility(recordId) 获取归档资格(记录ID)
  R->>W: Query active tasks/follow-ups by record 按记录查在途任务/跟进
  R->>C: Query active quote/contract/payment obligations where relevant 查相关在途报价/合同/回款义务
  W-->>R: Active work obligations 在途工作义务
  C-->>R: Active commercial obligations 在途商务义务
  alt no active obligations 无在途义务
    G->>R: ArchiveRecord(expectedVersion, reason) 归档记录(期望版本,原因)
    R->>R: Persist archived state 持久化归档状态
    R-->>H: RecordArchived event 事件:记录已归档
    R-->>G: Archived result 归档结果
  else active obligations exist 存在在途义务
    R-->>G: ARCHIVE_BLOCKED with obligation DTOs 归档被阻塞(含义务DTO)
    G-->>U: Blocked archive state with refresh/retry option 阻塞态(带刷新/重试选项)
  end
```

## Sequence: Owner Transfer And Open Work Transfer

```mermaid
sequenceDiagram
  participant U as Admin/Manager 管理员/经理
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant R as record-owning service 记录所属服务
  participant W as work-service 工作服务(活动/任务/提醒)
  participant H as audit-history-service 审计历史服务

  U->>G: Change owner 变更负责人
  G->>P: Check owner transfer permission 校验负责人转移权限
  P-->>G: Allowed 允许
  G->>R: ChangeOwner(expectedVersion, newOwnerId, reason) 变更负责人(期望版本,新负责人ID,原因)
  R->>R: Persist owner change and outbox OwnerChanged 持久化负责人变更并发件箱记录OwnerChanged
  R-->>H: OwnerChanged history event 历史事件:负责人已变更
  R-->>W: OwnerChanged event or TransferOpenWork command 事件OwnerChanged或转移在途工作命令
  W->>W: Transfer open tasks/follow-ups unless manual exception exists 转移在途任务/跟进(除非有手动例外)
  W-->>H: OpenWorkTransferred event 事件:在途工作已转移
  W-->>R: Transfer completed or retryable failure 转移完成或可重试失败
  R-->>G: Owner change result with workTransferStatus 负责人变更结果(含工作转移状态)
```

If open work transfer fails, the owning service must expose
`workTransferStatus = PendingRetry | Failed` and the work-service must retry by
idempotency key. A manual reassignment exception requires Administrator or Sales
Manager permission and a required reason.

## Sequence: Close Lost And Terminal Edit Protection

```mermaid
sequenceDiagram
  participant U as Sales/Manager 销售/经理
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant O as opportunity-service 商机服务
  participant W as work-service 工作服务(活动/任务/提醒)
  participant H as audit-history-service 审计历史服务
  participant R as reporting-service 报表服务

  U->>G: Close opportunity Lost 关闭商机为输单
  G->>P: Check close permission 校验关闭权限
  P-->>G: Allowed 允许
  G->>O: CloseLost(expectedVersion, lostReason, closeDate) 输单关闭(期望版本,输单原因,关闭日期)
  O->>O: Validate required lostReason and non-terminal state 校验必填输单原因与非终态
  O->>O: Persist terminal Lost state 持久化终态"输单"
  O-->>H: OpportunityClosedLost event 事件:商机已输单关闭
  O-->>R: OpportunityClosedLost event 事件:商机已输单关闭
  O-->>G: Lost result 输单结果
  U->>G: Edit closed opportunity 编辑已关闭商机
  G->>O: UpdateOpportunity(expectedVersion) 更新商机(期望版本)
  O-->>G: TERMINAL_RECORD_READ_ONLY 终态记录只读
  U->>G: Add post-close note/task 关闭后追加笔记/任务
  G->>W: CreateNoteOrTask(relatedOpportunityId) 创建笔记或任务(关联商机ID)
  W-->>H: WorkItemCreated event 事件:工作项已创建
```

Won/Lost are terminal opportunity states in the committed scope. Post-close notes and follow-up
tasks are allowed through work-service only; they do not reopen or edit the
closed opportunity.

## Sequence: Duplicate Warning And Proceed-After-Warning

```mermaid
sequenceDiagram
  participant U as User 用户
  participant G as gateway-bff 网关
  participant P as identity-authz-service 身份鉴权服务
  participant L as lead-service 线索服务
  participant A as account-service 客户服务

  U->>G: Create lead/account/contact 创建线索/客户/联系人
  G->>P: Check create permission 校验创建权限
  P-->>G: Allowed 允许
  G->>L: CreateLead or duplicate probe 创建线索或重复探测
  L->>L: Normalize company, contact, phone, email, province 归一化公司/联系人/电话/邮箱/省份
  L->>A: Query safe duplicate candidates 查询安全的重复候选
  A-->>L: Safe match summaries only 仅返回安全匹配摘要
  alt probable duplicate 疑似重复
    L-->>G: DUPLICATE_WARNING with warningToken 重复预警(带warningToken)
    U->>G: Proceed after warning 预警后继续
    G->>L: CreateLead(proceedWarningToken) 创建线索(继续令牌)
    L->>L: Create new record without merge/overwrite 创建新记录(不合并/不覆盖)
  else no duplicate 无重复
    L->>L: Create record 创建记录
  end
  L-->>G: Created result 创建结果
```

## Architecture Principles

- A service owns its own data and writes its own data.
- No service may receive, request, or use another service's database
  credentials.
- Cross-service reads must use target service Query API, events, or an approved
  read model.
- Cross-service writes must use target service Command API.
- Shared packages may contain contracts, DTO schemas, constants, and generated
  clients only. Shared business implementation is prohibited.
- Core CRM records are never hard-deleted in the committed scope. They are closed, terminated,
  archived, or moved through explicit lifecycle states.
- Reporting is based on an owned read model or target service APIs, not direct
  cross-service table reads.
- Import/export must call domain services and cannot bypass validation,
  authorization, history, or audit rules.
- Audit/history must be durable and append-only through normal CRM workflows.
- Editable P0 records must expose a concurrency token. Mutating commands must
  include `expectedVersion` and return `VERSION_CONFLICT` when stale.
- Archive, owner transfer, close Won/Lost, quote acceptance, contract status,
  payment, import/export, and user lifecycle changes must create history or
  operation log events.

## Key Design Decisions

| Decision ID | Decision | Status | Notes |
|---|---|---|---|
| ADR-ARCH-001 | Use physical multi-service Go backend with Docker Compose on one runtime host (`srv-volcengine-sh-01`, Volcengine ECS); `srv-aliyun-bj-01` is the off-server backup target only. | Accepted for G5 Re-review | Details in `service-architecture-adr.md`. |
| ADR-ARCH-002 | Use one self-hosted PostgreSQL instance with service-isolated database/schema and users. | Accepted for G5 Re-review | Direct cross-service database access is forbidden. |
| ADR-ARCH-003 | Use local automatic encrypted PostgreSQL backups with 7-day retention as baseline. | Accepted for pre-release only | Same-host-only backup is a production release blocker until off-server backup evidence exists. |
| ADR-ARCH-004 | Do not create a unified database CRUD service. | Accepted | Services expose business APIs, not database operation APIs. |
| ADR-ARCH-005 | Enforce HTTPS-only production ingress and authenticated service-to-service calls. | Accepted for G5 Re-review | Details in `authz-architecture.md` and `deployment-notes.md`. |

## G5 Handoff To MDA

Domain Modeling must represent these architecture decisions in PSM:

- service mapping and bounded contexts
- aggregate ownership and data ownership
- public API, event, error, permission, and DTO contracts
- service-to-service permission rules
- state machines and failure paths
- idempotency, timeout, retry, compensation, and correlation ID rules
- backup, restore, deployment, and observability constraints
- forbidden cross-service imports and forbidden database access

## Gate Status

This architecture package is revised for G5 re-review. G5 is not passed until
Product Manager, Business Analyst, UX Designer, UI Designer, and Security
Compliance review it and no P0/P1 blocker remains.
