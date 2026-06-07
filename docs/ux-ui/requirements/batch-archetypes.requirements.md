# 需求规格 —— 批量页型原型（详情 / 表单 / 报表 / 管理）

状态：Requirements（Claude 出，作为 G8 审计基准）。本文件定"做什么 / 怎么判定合格 / 约束"，
**不定**视觉（字体/字号/版面/间距/调性/组件）——视觉地基已锁定，见下。
日期：2026-06-06
模型：新协作模型（Claude 需求+审计 / Codex 设计+实现）。
范围：`projects/crm-system/`，PREVIEW-ONLY（静态 mock 图，本轮不做可点击原型）。
范例实体：商机系（详情=商机、表单=商机），报表=团队报表，管理=用户与角色。

## 0. 共享视觉契约（已锁定，所有页直接套用，不得偏离）
- 锁定设计系统：`docs/ux-ui/design-system.md`（§3 字阶、§4 间距/圆角/阴影、§5 图标、
  §6 版面/导航规则/64px 顶栏/网格/断点、§7 组件、§8 状态、§9 数据可视化）。
- 参考帧（已定稿）：`mockups/dashboard-v7-sales|manager|manager-focus.png`、
  `mockups/list-opportunities.png`、`mockups/list-opportunities-sales.png`。
- 不得新增/改动颜色或 token；外壳（侧栏 248px / 顶栏 64px）、组件类、表格样式逐字复用。
- 全 zh-CN；金额 tabular + ¥；新建/编辑等破坏性/终态动作注释标明走确认流程。

## 0b. 全局领域真源（避免批量复制错误）
- 角色（3）：管理员 / 销售经理 / 销售（`identity-authz/.../role.go`）。
- 商机阶段（6，唯一枚举）：新商机 / 需求已确认 / 报价 / 合同谈判 / 赢单 / 丢单
  （`frontend/src/i18n/labels.ts`）。**禁止**把回款/合同状态当阶段。
- 其它状态枚举一律取自 labels.ts：报价 草稿/已发送/已接受/已拒绝/已过期；合同 待签署/
  已签署/启用/已完成/已终止；回款 无计划/未回款/待回款/部分回款/已回款/已逾期/已取消。

---

## 1. 详情页 —— 商机详情
产物：`mockups/_src/detail-opportunity.html` + `mockups/detail-opportunity.png`
角色上下文：销售经理视图（manager@example.com）。

### 必须呈现
- 外壳（展开侧栏，"商机"active）+ 顶栏；面包屑/返回到商机列表。
- 头部：商机名 + 当前阶段徽 + 客户 + 负责人 + 金额 + 预计签约 + 更新时间；主操作区。
- **阶段步进器（StageStepper）**：线性展示 新商机→需求已确认→报价→合同谈判→(赢单/丢单)，
  高亮当前阶段。
- 关联区块（面板）：关联报价（一条，DEC-018）、关联合同、回款计划/记录、活动/备注/任务时间线。
- 终态信息：若赢单显示「赢单合同」；若丢单显示「丢单原因」。

### 领域标尺（硬性，G8 审）
- 阶段步进只能体现**线性单向**推进（不可跳级/回退）；"推进到下一阶段"动作只指向 `AdvanceStage`
  的下一个合法阶段。
- **赢单**动作仅在「合同谈判」且**存在已签署合同**时可用（否则禁用/不可点，注释标明 DEC-017）；
  **丢单**动作须走"填原因码"确认流程（原因取 labels.ts `lostReasonLabel`）。
- **终态（赢单/丢单）记录只读**：不出现改阶段/编辑入口（注释标明 TERMINAL_RECORD_READ_ONLY）。
- 操作按角色门控：编辑/转移负责人/归档仅经理·管理员；阶段枚举值真实。
- 状态覆盖：默认；空关联（如"暂无关联合同"）；终态只读态——至少注释定义 loading/empty。

---

## 2. 表单页 —— 新建/编辑商机
产物：`mockups/_src/form-opportunity.html` + `mockups/form-opportunity.png`
角色上下文：销售视图（sales@example.com，新建场景）。

### 必须呈现
- 外壳 + 顶栏 + 页头（新建商机 / 取消·保存）。
- 表单字段（对齐 `CreateOpportunity`/`UpdateOpportunity` 必填项）：
  商机名、客户（选择）、负责人、阶段（下拉，六枚举）、预计金额、预计签约日期。
- 字段校验示例：至少展示一个**校验错误态**（如金额为空/非法）+ 一个**正常态**。
- 保存区：主按钮"保存"、次按钮"取消"。

### 领域标尺（硬性，G8 审）
- 必填项与后端一致：customerId/ownerId/stage/expectedAmount/expectedCloseDate/title 均必填
  （`domain.UpdateOpportunity` 校验）。
- **销售视图：负责人锁定为本人**（`CanCreateOpportunity` 强制 ownerID==actor）——负责人字段
  对销售为只读/预填自己，不可选他人。经理/管理员才可选负责人。
- 阶段下拉：**新建时只列非终态四项**（新商机/需求已确认/报价/合同谈判），默认"新商机"——
  禁止直接新建为终态（赢单/丢单），因终态须走 CloseWon(需已签合同)/CloseLost(需原因)。编辑场景同理
  不经关闭动作不可置为终态。[REFINED 2026-06-06 per G8 observation]
- 编辑场景体现乐观并发（隐含 version/expectedVersion；静态可在注释说明）。
- 状态覆盖：正常输入态 + 字段校验错误态；注释定义 saving（提交中）/ 提交失败 态。

---

## 3. 报表页 —— 团队报表
产物：`mockups/_src/reports-team.html` + `mockups/reports-team.png`
角色上下文：销售经理视图。对应 ACC-018（团队概览）/ ACC-023（基础销售报表）。

### 必须呈现
- 外壳（"报表"active）+ 顶栏 + 页头（团队报表 + 时间范围/分组筛选 + 导出）。
- **概览指标卡**（OverviewMetrics）：线索数、商机数、任务数、赢单数、丢单数、报价额、
  合同额、已回款额、应收额（金额 tabular+¥）。
- **管道分布**（pipeline / opportunitiesByStage）：按六阶段的数量+金额（复用漏斗或条形）。
- **分组明细**（breakdowns）：线索按状态、商机按阶段、报价按状态、合同按状态、回款按状态
  （回款行含 应收/已回款）——各用真实枚举。
- 空态：emptyState（"所选范围暂无数据"）至少注释定义。

### 领域标尺（硬性，G8 审）
- 指标字段名/口径对齐 `frontend/src/api/reports.ts`（OverviewMetrics/ManagerOverview/
  BasicReport）；不杜撰指标。
- 各分组维度的枚举值取自 labels.ts（阶段六枚举、各状态枚举），不混用。
- 数据范围：经理=团队（scope/teamId 体现）；不展示越权的全公司数据除非管理员。

---

## 4. 管理页 —— 用户与角色
产物：`mockups/_src/admin-users.html` + `mockups/admin-users.png`
角色上下文：管理员视图（admin/Administrator）。

### 必须呈现
- 外壳（"管理：用户与角色"active）+ 顶栏 + 页头（用户与角色 + 新建用户）。
- 用户表格：显示名、邮箱、角色（徽：管理员/销售经理/销售）、状态（启用/停用）、操作
  （编辑/停用/改角色）。
- 至少一行体现**末位管理员保护**：唯一在用管理员的"停用/降级"操作禁用 + 提示文案
  （注释标明 last-admin guard）。

### 领域标尺（硬性，G8 审）
- 字段对齐 `identity-authz/.../user.go`：email、displayName、role、status；**不展示密码/哈希**。
- 角色仅三枚举（管理员/销售经理/销售）；状态仅 启用/停用。
- **整页仅管理员可见/可操作**（user.*/role.* 仅 Administrator）——这是页面级权限门控，
  注释标明；非管理员不应进入此页。
- **末位管理员保护**：不能停用/降级/删除最后一个在用管理员（`WouldRemoveLastActiveAdministrator`）——
  对应操作置灰 + 解释。
- 状态覆盖：默认列表；注释定义 loading/empty。

---

## 5. 移交给 Codex 决定的设计项
列/字段的精确布局与宽度、步进器/卡片/图表的具体视觉、徽变体到枚举的映射（仅用既有变体）、
各态的视觉处理、是否单独出空/载入态图、样本数据具体内容（须满足各页领域标尺）。

## 6. Out of scope / 约束
- 其它实体（客户/合同/线索）的详情/表单本轮不做（范式复用，实现时再铺）。
- 不改后端/API/数据模型/业务逻辑；不部署；不弱化任何 P0/P1 或既往 G12 修复。

## 7. G8 审计我会核什么
每页对照其领域标尺：真实枚举、真实字段名、权限/角色门控（含销售负责人锁定、管理页仅管理员、
末位管理员保护、终态只读、赢单需已签合同、阶段线性单向）、状态覆盖、无 recolor、与参考帧
视觉一致、PREVIEW-ONLY 合规。
