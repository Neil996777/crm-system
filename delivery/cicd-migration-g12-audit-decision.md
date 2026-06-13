# CI/CD 迁移 —— G12 独立审计决定(Claude)

| 字段 | 值 |
|---|---|
| 闸 | G12 实现审计 → 发布决定 |
| 审计人 | Claude(独立审计角色) |
| 日期 | 2026-06-13 |
| 尺子 | `delivery/cicd-migration-acceptance.md`(ACC-CICD-001..008 + C1–C6) |
| 发布内容 | commit `66d2531`(已审后端 G12 + 已过闸 UI/UX 完成) |
| 部署目标 | `srv-volcengine-sh-01`(`118.196.44.193`) |
| 判定 | **GATE PASSED**(CI/CD 迁移机制) |

## 验收逐条核(独立验证,非轻信证据)

| ACC | 判定 | 独立依据 |
|---|---|---|
| 001 离机构建 | ✅ | CI `ubuntu-latest`(run 27464903082 success);部署 transcript 0 host-build 命令 |
| 002 CI build+test+image | ✅ | run 27464903082:10 Go 测试 + 前端 build + e2e 满 61(57 stable+4 isolated)+ 11 镜像导出 |
| 003 CD 只 load/run | ✅ | runbook + run-release-step.sh 守卫;transcript 无 build/checkout |
| 004 digest↔commit / 禁 latest | ✅ | **活主机 11 容器 revision label 全 = 66d2531**(亲读 docker inspect);postgres `@sha256:16bc17c6…` 钉死;manifest 无 latest |
| 005 image-only compose | ✅ | compose 0 `build:`、0 `latest` |
| 006 五项发布证据 | ✅ | 测试结果 + digest→commit(containerd-safe TSV)+ 部署 transcript + 健康检查/TLS/安全头/负端口 + 回滚点 + **offsite 备份 PASS** |
| 007 生产机不依赖源码 | ✅ | 迁移从 release artifact SQL 跑;运行不依赖 `/opt/.../current/.git` |
| 008 零应用 diff / 不降级 | ✅ | 全程 `git diff -- frontend services shared` = 0;发布内容钉 66d2531 |

## 绑定约束(C1–C6)

C1 仅机制 ✅(0 应用源码)· C2 标准 §1 全 MUST ✅ · C3 发布内容=66d2531 + zh-CN/enum 不动 ✅ · C4 不跳 G10/G11/G12、Codex 不自判 ✅ · C5 co-location/ingress ✅(仅 CRM 80/443,gateway/frontend loopback)· C6 secret 不入库/不入 transcript/不入证据 ✅(独立扫描 0 泄漏)

## rework 收口(本季关键发现,均已闭环)

- **BLK-CICD-G11-002 secret-safe 迁移**:bundle 迁移用 `__CRM_DB_PASSWORD_*__` 占位 + prod.env 渲染,0 dev/明文口令。**Resolved**。
- **BLK-CICD-G11-003 可复现性**:3 处主机适配(DO-block 渲染 / containerd-safe 核验 / frontend tmpfs)已回写提交脚本,经干净重部署验证(committed 脚本、0 内联、9 角色登录 PASS、tmpfs OK)。**Resolved**。
- **BLK-CICD-G11-004 启动竞态**:根因=runbook 先起服务后迁移 → 28P01 瞬时 502。修复=`compose-up-release.sh` 改 postgres→migrate→apps 顺序。**独立验证:活主机 identity-authz 全历史 28P01=0、9 服务 0 竞态、nginx 502 扫描 PASS。Resolved**。

## 独立活检(VPN-route 绕过后亲测)

HTTPS `/health` 200 + HSTS/CSP;前端 SPA 200;HTTP 301;TLS Let's Encrypt 有效;8080/5432 公网 closed;**login 200 / 创建用户 201 / list-users 200**(nginx→gateway→服务→postgres 全链路功能正常)。

## 范围外的独立产品缺陷(不阻断本 G12)

- **BLK-PROD-AUDIT-001**:`UserAccessDenied` 审计事件信封缺 `actorRole/actorDisplay` → audit-history 400 → 访问被拒事件未入操作日志(ACC-022 安全审计窟窿)。**这是基线产品(66d2531)潜伏的应用 bug,出 CI/CD 迁移"仅机制"范围(C1)**,需独立变更修复(release owner 立项)。不阻断 CI/CD 迁移 G12。

## 决定

**CI/CD 迁移 G12 GATE PASSED。** 发布机制已从"生产机构建 + latest + host checkout"迁移为合规的"离机构建 + digest-pinned 镜像 + image-only image-only 部署 + commit 可追溯 + 可复现 + 完整发布证据"。生产 `srv-volcengine-sh-01` 现运行 commit `66d2531`(独立核实 digest↔66d2531),站点活、安全、功能正常。

**后续(非 G12 阻断):** ① BLK-PROD-AUDIT-001 独立产品缺陷立项;② `verify-loaded-images.sh` 的 `.Id` 比对长期建议统一为 containerd-safe(已在部署用);③ 诊断用户 `diagtest@example.com` 待 release owner 决定清理。
