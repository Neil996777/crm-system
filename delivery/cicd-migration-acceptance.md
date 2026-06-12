# CI/CD 迁移 —— 验收矩阵 + 绑定约束(Yardstick)

> 角色:Claude 本职 = yardstick(需求/验收/绑定约束,G1–G3)。本文件是这次"CI/CD &
> 发布流水线迁移"跟进变更的**验收源头尺子**,也是 G8 handoff 审计与 G12 审计的判定依据。
> 权威:`../../standards/cicd-and-release-standard.md`(§1 核心规则)、`company/operating-model.md`
> (no-downgrade / 发布规则 / 证据规则)、政策变更 `../../company/policy-changes/...`(2026-06-06 标准采纳)。
> 配套:`cicd-migration-brief.md`、`cicd-migration-plan.md`(M1–M6 任务草案)、
> `cicd-current-state-gap-analysis.md`(现状差距 + D1–D4)。
> 状态:**Yardstick 已经 release owner 确认(2026-06-12)。** 交 Codex 从 G5 重入。

---

## 1. WHAT / WHY(需求)

**WHAT.** 把 CRM 的发布机制从"在生产机上构建"(on-host `npm build` + `docker compose up -d --build`、
`latest` 浮动标签、host `git checkout` 取源码)迁移为符合公司 CI/CD 标准的**离机构建 + digest-pinned
镜像 + image-only 部署**。

**WHY.** 公司 2026-06-06 采纳 CI/CD & 发布流水线标准;旧方式违反其 §1 全部核心规则,且历史上因在生产机
构建造成磁盘压力 / 陈旧构建事故。`/opt/crm-system` 已于 2026-06-06 清空,下一次 CRM 上线必须是
**首次标准合规部署**。

**边界:这是发布机制变更,不动应用代码。** 不改任何业务逻辑/API/数据模型/安全实现;只改"怎么构建、
怎么把已构建产物放上生产机运行"。

---

## 2. 验收矩阵(可测 Definition of Done)

| ID | P | 能力(Done 定义) | 标准依据 | 验证方法 |
|---|---|---|---|---|
| ACC-CICD-001 | P0 | **离机构建**:生产机不执行任何 build/compile(无 `npm build`/`docker build`/`compose build`/`up --build`/任何编译器·打包器);全部服务镜像 + 前端在生产机外(CI/构建工作站)构建 | §1.1 | 审 deploy runbook + `docker-compose.prod.yml` 无任何构建动作;CI 配置在离机环境构建 |
| ACC-CICD-002 | P0 | **CI 负责 build + test + 出镜像**:CI 构建产物、运行 QA 定义的测试套件(现有 61 e2e + 后端测试)、产出不可变镜像并使其对部署可用(push registry 或导出 digest-pinned 镜像工件) | §1.2 | CI 流水线运行记录:构建日志 + 测试结果 + 产出镜像清单 |
| ACC-CICD-003 | P0 | **CD 只拉取/加载运行**:生产机部署只 pull(registry)或 load(digest-pinned 导出镜像)**指定 tag/digest** 并运行;**不**在机上 build/compile/为构建 `git checkout` | §1.3 | 审 runbook:部署步骤只有 pull/load + run,无构建 |
| ACC-CICD-004 | P0 | **commit 可追溯**:每个部署镜像可溯到其精确源 commit(镜像 label 或发布证据里的 digest→commit 映射);**不**用 `latest` 等浮动标签做发布身份,发布身份 = digest | §1.4 | 发布证据含 digest→commit 表;镜像带 commit label;无 `latest` 作发布标识 |
| ACC-CICD-005 | P0 | **image-only 生产 compose**:`docker-compose.prod.yml` 全部服务引用镜像 tag/digest、**零 `build:` 键**、无机上前端构建步骤 | §1.1/§1.3/§1.5 | 审 compose:`grep build:` = 0;前端服务引用已构建镜像 |
| ACC-CICD-006 | P0 | **发布证据齐**:测试结果 + digest→commit 映射 + 部署 transcript + 部署后健康检查 + 回滚点(上一已知良好 digest) | §1.4 + 证据规则 | 发布证据工件可核 5 项齐全 |
| ACC-CICD-007 | P0 | **生产机不长期依赖源码**:生产机不保留为构建而存在的源码 checkout 作为运行依赖 | §1.5 | 审主机运行依赖:运行只依赖镜像 + 配置/secret,不依赖源码树 |
| ACC-CICD-008 | P0 | **零应用代码改动 / 不降级**:业务逻辑/API/数据模型/G12 安全修复 **0 diff**;发布内容 = 已过闸 commit,其任何 P0/P1 验收与 G12 结论一项不降 | operating-model no-downgrade | diff 审计:应用层 0 改动;发布内容 commit 已过 UI/UX G12 |

**说明:** ACC-CICD-001..007 把标准 §1 的 5 条 MUST + 证据规则逐条落成可测项;ACC-CICD-008 是
no-downgrade 红线的机制侧投影(机制变更不得借机改应用或降级任何已过闸结果)。

---

## 3. 绑定约束(C,红线)

- **C1 仅机制。** 不改任何应用/业务逻辑/API/数据模型/安全实现;前端/后端**应用源码 0 diff**。
  允许改动仅限:CI 配置、`docker-compose.prod.yml`、部署 runbook/脚本、镜像构建文件(Dockerfile)、
  发布证据模板。
- **C2 满足标准 §1 全部 MUST。** §1.1 离机构建、§1.2 CI build+test+image、§1.3 CD 只拉取/load、
  §1.4 digest 可追溯(禁 `latest` 作发布身份)、§1.5 生产机不依赖源码。绿流水线 ≠ 发布批准(标准明文)。
- **C3 发布内容 = 已过闸 commit。** 发布内容 = `66d2531`(已审后端 G12 + 已完成并过闸的 UI/UX:
  ACC-018 团队总览 / ACC-023 基础报表 + UI 表面)。**保 zh-CN;不改任何 enum/role 比较值。**
  (此条取代 brief/plan 里残留的"audited HEAD as-is"旧措辞——见 gap 分析 D3。)
- **C4 不跳 G10/G11/G12。** 机制迁移仍走 QA(G10)/集成(G11)/审计(G12);Codex 不自判 G12。
  Infrastructure Ops 是 G5/G8/G11/G12 必需评审人(项目层强化)。
- **C5 不超需放开 ingress。** 保 co-location 约束:只为 CRM `server_name` 放开必要 80/443,不接管主机
  超出 CRM 范围的入口(gap 分析 D4)。
- **C6 secret 不入库不入文档。** 任何密钥/口令/token/secret 值不写进项目文档或发布证据;只记位置与恢复
  方法(承 infra README / backup-recovery-plan)。

---

## 4. 已确认决策(release owner,2026-06-12)

- **D3 发布内容 commit = `66d2531`**(确认)。已过 UI/UX G12,符合 gap 分析 D3 的"过闸后再定 commit"前置。
- **D1 镜像分发形式:默认 export/load(`docker save`→`scp`→`docker load`)。** 依据:公司基础设施现无任何
  registry(infra register 全核实),单生产机 `srv-volcengine-sh-01`;export/load = 零新基础设施、完全合规
  §3;registry 仅在多主机/未来更优。**最终敲定属 Codex 的 G5 ADR(M1)**,但默认值钉在 export/load,除非
  release owner 决定为此新立 registry(当前不建议)。
- **D2 前端形态:nginx 容器镜像(推荐),不走 dist+scp。** 依据:标准把发布单元定义为带 digest 的不可变
  镜像;dist+scp 让前端以散文件落地、非 digest-pinned、commit 可追溯性弱且违背"镜像即发布单元";nginx 镜像
  与 10 个 Go 服务同构、同一条 CD 路径、单一 digest。**最终敲定属 G5 ADR**,推荐值 = nginx 镜像。
- **D4 ingress:** 见 C5,放开范围属 G5 ADR + Infrastructure Ops 评审。

---

## 5. 不在范围(Out of Scope)

- 任何应用/业务逻辑/API/数据模型/安全行为改动(C1)。
- 任何 P0/P1 验收或 G12 结论的降级(C2 no-downgrade)。
- 用绿流水线替代 G10/G11/G12 任一闸(C4;标准明文绿流水线≠发布批准)。
- 多主机/集群编排、蓝绿/金丝雀等超出"单机标准合规部署"的发布策略(可作未来项,不在本次)。

---

## 6. 闸路径

REGISTERED → **G5**(registry/channel ADR[M1] + image-only compose[M3] + digest runbook[M4])
→ **G7/G8**(M1–M6 任务包,Codex 产,Claude 审 G8 handoff)→ **G9–G11**(Codex 执行:CI 搭建、镜像构建、
离机测试、部署演练)→ **G12**(Claude 独立审计:逐条核 ACC-CICD-001..008 + C1–C6 + 发布证据)。

> 实现代码/配置在 G8 通过前不得开始(no implementation before G8)。
