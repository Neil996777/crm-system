# ACC ↔ Test 覆盖追溯审计（2026-06-12）

> 目的：核实每个 P0/P1 产品验收（`docs/product/acceptance-matrix.md`，23 项 = 17 P0 + 6 P1）
> 是否有通过的测试覆盖，揪出"有验收、无测试"的 P0/P1 风险。
> 来源：58 个 e2e 用例（workers:2+retries:1 下 58/58 绿）。
> 发现：验收矩阵**没有 Verification 列引用 TEST-ID**（0 处）→ 覆盖此前是隐含未核实的；本表补上。

## 追溯矩阵

| ACC | P | 能力 | 覆盖测试 | 判定 |
|---|---|---|---|---|
| ACC-001 | P0 | 登录 + 按角色操作 | AUTH-LOGIN-001/005、002 | ✅ 覆盖 |
| ACC-002 | P0 | 三角色访问控制 | PERM-USERADMIN-002/003、OPLOG-004、A4-OPP-001、A4-REPORT-001、NAV-RETRIEVE-005、HISTORY-003 | ✅ 覆盖(多测试横切) |
| ACC-003 | P0 | 线索管理(增改查/搜筛/分配转移) | LEAD-CREATE-002、FUNC-ROWMENU-001、NAV-RETRIEVE-* | ⚠️ 部分:**转移负责人无专测**(矩阵明文要求 assign/transfer) |
| ACC-004 | P0 | 线索资格流转 + 历史 | LEAD-QUALIFY-003、004 | ✅ 覆盖 |
| ACC-005 | P0 | 公司/客户管理 | CUSTOMER-CRUD-002 | ✅ 覆盖(基本) |
| ACC-006 | P0 | 多联系人 | CONTACT-LINK-003 | ✅ 覆盖 |
| ACC-007 | P0 | 商机管理 | OPP-*、UIUX-P1-002 | ✅ 覆盖 |
| ACC-008 | P0 | 商机阶段流转 | OPP-STAGE-002 | ✅ 覆盖 |
| ACC-009 | P0 | 报价(每商机一报价) | QUOTE-LIFECYCLE-002(×2)、ACCEPT-001 | ✅ 覆盖 |
| ACC-010 | P0 | 合同(签约日期/差额原因) | CONTRACT-CREATE-002、LIFECYCLE-002 | ✅ 覆盖 |
| ACC-011 | P0 | 回款计划/实收 | PAYMENT-RECORD-002、GUARD-003 | ✅ 覆盖 |
| ACC-012 | P0 | 活动/笔记/任务(≥5 种记录上下文) | ACTIVITY-NOTE-002、TASK-LIFECYCLE-002 | ⚠️ 部分:**只测了 1 种记录上下文**;矩阵明文要求"对 lead/customer/opportunity/contract/payment 分别测" |
| ACC-013 | P0 | 关闭赢单/丢单 + 历史 | OPP-CLOSE-002、003 | ✅ 覆盖 |
| ACC-014 | P0 | 记录级历史 | HISTORY-001/004、003 | ✅ 覆盖 |
| ACC-015 | P0 | 列表/详情/搜索/筛选(happy/空/非法筛/权限) | NAV-RETRIEVE-001/003/004/005 | ✅ 覆盖(场景齐) |
| ACC-016 | P0 | 持久化(刷新/重登/重启) | PERSISTENCE-001..005 | ✅ 覆盖 |
| ACC-017 | P0 | 部署运行(真配置/真库) | 无 e2e(验证法=Integration/Manual/Audit) | ➖ 非 e2e 项;已于 2026-06-05 go-live 验证 |
| ACC-018 | P1 | 团队总览 | TEAM-OVERVIEW-003、A4-REPORT-002 | ✅ 覆盖 |
| ACC-019 | P1 | 查重告警(公司/联系人/线索) | DUPLICATE-WARN-001/004/005 | ⚠️ 部分:公司+线索已覆盖;**联系人电话/邮箱查重无专测** |
| ACC-020 | P1 | 导入/导出 | CSV-IMPORT-001/002、CSV-EXPORT-001 | ✅ 覆盖 |
| ACC-021 | P1 | 提醒(任务/合同/回款) | REMINDER-001/002/003、004 | ✅ 覆盖 |
| ACC-022 | P1 | 管理员全局操作日志 | OPLOG-001/002/005、004 | ✅ 覆盖 |
| ACC-023 | P1 | 基础报表 | BASIC-REPORT-002、A4-REPORT-001/002 | ✅ 覆盖 |

## 结论

- **无任何 P0/P1 验收完全无测试覆盖** —— 没有"裸奔"的红线。✅
- **ACC-017**：部署项,非 e2e（设计如此,验证法=Manual/Audit）,已于 go-live(2026-06-05)验证。非缺口。
- **3 个部分缺口（矩阵明文要求但测试未覆盖到的子场景）**:
  - **ACC-003**:线索"分配/转移负责人"无专门 e2e（产品验收明文要求 assign/transfer）。
  - **ACC-012**:活动/笔记/任务只测了 1 种记录上下文,矩阵要求对 lead/customer/opportunity/contract/payment **分别**测。
  - **ACC-019**:**联系人 电话/邮箱 查重**无专测（公司名、线索已覆盖）。
- **结构性问题**:`acceptance-matrix.md` 缺 Verification 列 → ACC↔TEST 追溯不显式。建议把本表的映射补进矩阵的 Verification 列(产品验收源头尺子应自带追溯)。

## 建议处理

1. **补 3 个部分缺口的 e2e**（ACC-003 转移、ACC-012 多上下文、ACC-019 联系人查重）→ 踢回 Codex 加测试。优先级 P1（这些是 P0/P1 验收的明文子场景）。
2. **给 acceptance-matrix.md 补 Verification 列**（每个 ACC 列出覆盖它的 TEST-ID）→ yardstick 增强,使追溯显式、可维护。
