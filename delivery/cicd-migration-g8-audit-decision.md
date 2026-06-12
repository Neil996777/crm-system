# CI/CD 迁移 —— G8 Handoff 审计决定(Claude 独立审计)

| 字段 | 值 |
|---|---|
| 闸 | G8 设计/计划 → 实现 handoff |
| 审计人 | Claude(独立审计角色) |
| 日期 | 2026-06-12 |
| 尺子 | `delivery/cicd-migration-acceptance.md`(ACC-CICD-001..008 + C1–C6 + D1/D2/D3) |
| 被审包 | Codex G5/G8 产出 |
| 判定 | **GATE PASSED** — 可进 G9 实现(M2/M3/M4) |

## 被审交付物

- `docs/architecture/adr/ADR-CICD-001-image-channel-and-frontend-runtime.md`(M1 ADR:D1/D2)
- `docs/architecture/cicd-release-architecture.md`(M2–M6 G5 架构)
- `delivery/cicd-migration-g8-task-package.md`(G8 任务包 M1–M6)
- `delivery/cicd-release-evidence-template.md`(发布证据模板)
- `delivery/cicd-migration-acc-traceability.md`(ACC↔交付物追溯)
- `docs/product/decision-log.md`(DEC-023/024/025)
- `delivery/cicd-migration-brief.md` / `cicd-migration-plan.md` / `cicd-current-state-gap-analysis.md`(旧措辞对齐)

## 验收逐条核(ACC-CICD-001..008)

| ACC | 判定 | 依据 |
|---|---|---|
| 001 离机构建 | ✅ | ADR 拒绝 host build;arch M2 离机 CI;M3 compose 删 build;M4 runbook fail-fast 禁 host build |
| 002 CI build+test+image | ✅ | arch M2 阶段:10 Go 测试 + 前端 build + Playwright e2e + 镜像 + 导出 manifest |
| 003 CD 只 load/run | ✅ | M4 runbook = 校验 checksum + `docker load` + `up -d`(无 `--build`)+ 无 `git checkout` 取源构建 |
| 004 digest→commit / 禁 latest | ✅ | ADR digest 为发布身份 + OCI revision label;M3 删 `:-latest` 兜底;M6 运行镜像 digest 核验 |
| 005 image-only compose | ✅ | M3 差异表:10 个 `build:` 全删、`latest` 兜底改 commit tag、前端纳入 nginx 镜像 |
| 006 五项发布证据 | ✅ | 证据模板含 测试/digest→commit/部署transcript/健康检查/回滚点 + no-host-build 审计段 |
| 007 生产机不依赖源码 | ✅ | migrations 从 release artifact SQL 跑(非 `/opt/.../current/services` git checkout);M6 核验无 `.git` 依赖 |
| 008 零应用 diff / 不降级 | ✅ | scope 边界明确;发布内容钉 `66d2531`;**实测 `git status` 应用源码 0 diff**(仅 cicd 文档 + 治理) |

## 绑定约束逐条核(C1–C6)

- **C1 仅机制** ✅ —— 实测 `git status`:改动仅 cicd 交付物 + gate-status + decision-log;前端/后端/shared/e2e 源码 0 diff。任务包"允许文件类"限于 CI 配置/compose/runbook/Dockerfile/证据模板。
- **C2 标准 §1 全部 MUST** ✅ —— §1.1 离机、§1.2 CI build+test+image、§1.3 CD 只 load、§1.4 digest 身份(禁 latest)、§1.5 不依赖源码,逐条有设计落点。
- **C3 发布内容 = 66d2531 + 保 zh-CN + 不动 enum/role 值** ✅ —— 钉 66d2531;机制变更不触应用层故 enum/role 值天然不变。
- **C4 不跳 G10/G11/G12** ✅ —— 任务包明文 G9 后走 G10/G11 证据、G12 Claude 独立审;Codex 不自判闸。
- **C5 不超需放开 ingress** ✅ —— ADR + arch 保 loopback(gateway 8080 / frontend 8081)、仅 CRM `server_name` 80/443、不接管无关入口。
- **C6 secret 不入库/不入文档** ✅ —— ADR + 证据模板 + infra 评审点明文禁写 secret 值,只记位置与恢复法。

## 事实前提核验(不轻信描述)

- prod compose 真有 **10 个 `build:`**、10 服务真用 `${CRM_IMAGE_TAG:-latest}`、服务名与 ADR 镜像集一致、migration SQL 真从 `/opt/crm-system/current/services/.../migrations` 挂载 —— 设计的"现状差异"全部属实,迁移建在准确前提上。
- 旧 "as-is HEAD" 措辞已清(唯一残留 "Keep as-is in the new runbook" 指保留 backup 步骤,无关)。
- 本包**只含设计/计划,无实现代码**(nginx Dockerfile/CI 配置/compose 改动显式推迟到 G9)→ 满足"G8 前无实现代码"。

## 非阻断备注(带入 G9/G11,非 G8 失败项)

1. 证据模板 target host 写 `srv-volcengine-sh-01` / `118.196.44.193` —— 该公网 IP 由 **Infrastructure Ops 在 G9/G11 现场核**(infra register 我未独立确认 volcengine 公网 IP),非 G8 设计阻断。
2. ADR "base images pinned by digest where practical" 的 "where practical" 在设计阶段可接受;G9 落地时基础镜像应尽量 digest 固定,应用镜像 digest 由 manifest 钉(发布身份不受影响)。
3. Infrastructure Ops 是 G8/G11 必需评审人 —— 其 ingress/disk/secret/backup/monitoring 评审应在 G9 实现前补签(ADR/arch 已留 review hooks)。

## 决定

**G8 GATE PASSED。** 设计 + 任务包对齐 yardstick,8 项验收 + 6 条约束全部满足,事实前提核实,无实现代码先行。Codex 可进 **G9** 实现 M2/M3/M4,产 G10/G11 的 M5/M6 证据;Claude 在 G11 后做 **G12** 独立审计(逐条核运行镜像 digest↔66d2531 + 发布证据 + 无 host build 发生)再做发布决定。
