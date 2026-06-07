# 需求规格 —— 特殊页型原型（提醒中心 / 导入导出 / 操作日志）

状态：Requirements（Claude 出，作为 G8 审计基准）。定"做什么/怎么判定/约束"，不定视觉。
日期：2026-06-06
模型：新协作模型（Claude 需求+审计 / Codex 设计+实现）。
范围：`projects/crm-system/`，PREVIEW-ONLY（静态 mock 图）。

## 0. 共享视觉契约（已锁定，直接套用）
- 锁定设计系统：`docs/ux-ui/design-system.md`。参考帧：已定稿的 dashboard/list/detail/form/
  reports/admin 各帧。外壳(侧栏248px/顶栏64px)、组件类、表格样式逐字复用；不得新增/改色或 token；
  全 zh-CN；金额 tabular+¥；破坏性/确认动作注释标明走确认流程。
- 这三页都不是标准 CRUD 范式，但视觉语言必须与已定稿页一致。

---

## 1. 提醒中心
产物：`mockups/_src/reminders-center.html` + `mockups/reminders-center.png`
角色上下文：销售视图（本人提醒，owner 数据范围）。

### 必须呈现
- 外壳（"提醒中心"active）+ 顶栏 + 页头（提醒中心 + 业务日期 businessDate + 时区）。
- 提醒列表/分组：按提醒类型或 到期/逾期 分组皆可。每条提醒展示：
  类型徽、关联记录(relatedRecord.display + 类型)、负责人(ownerDisplay)、到期日(dueDate)、
  优先级(priority)、状态(status)。
- 每条提供跳转到关联记录的入口（如"查看"）。
- 空态："今日无待处理提醒"（注释定义 loading/empty）。

### 领域标尺（硬性，G8 审，源：frontend/src/api/reminders.ts）
- 提醒类型只用 5 个真实枚举（reminderTypeLabel）：任务到期 / 任务逾期 / 合同待签署 /
  回款到期 / 回款逾期。不得杜撰类型。
- 字段对齐 ReminderRow：type / relatedRecord{type,id,display} / ownerDisplay / dueDate /
  status / priority / version。
- **不要发明"已读/未读"开关**（接口无此语义）；可呈现 status/priority，但取真实字段、不臆造取值。
- 数据范围：销售=本人提醒；若做经理视图则为团队（owner 维度）。

---

## 2. 导入/导出（单页，双流程）
产物：`mockups/_src/import-export.html` + `mockups/import-export.png`
角色上下文：经理视图（数据范围=actor 有权访问的记录）。

### 必须呈现（一页含 导入 与 导出 两区，或选项卡）
- 外壳（"导入/导出"active）+ 顶栏 + 页头（导入/导出 + 副标题"导入/导出有权限访问的记录"）。
- **导入区**：对象类型选择（实体下拉）→ 选择 CSV 文件（文件名）→ 开始导入；
  导入结果卡：状态 + 总行数/成功数/失败数 + **逐行错误表**(行号/字段/错误信息) +
  审计记录状态 + 保留期(retainedUntil)。须展示一个"部分失败"结果示例（含 rowErrors）。
- **导出区**：对象类型选择 → 包含归档(includeArchived 勾选) → **确认勾选"确认导出范围并记录审计日志"** →
  开始导出；导出结果卡：导出行数 + 文件(可下载/只读内容) + 文件安全(fileSafety) + 保留期。

### 领域标尺（硬性，源：frontend/src/api/importexport.ts + labels.ts objectTypeLabel）
- 对象类型只用真实实体（objectTypeLabel）：线索/客户/联系人/商机/报价/合同。不得杜撰实体。
- 导入结果字段对齐 ImportRun：totalRows/successCount/failureCount/rowErrors[{rowNumber,field,code,
  safeMessage}]/operationLogStatus/cleanupStatus/retainedUntil。
- 导出字段对齐 ExportRun：exportedCount/archivedIncluded/content/operationLogStatus/fileSafety/
  retainedUntil。
- **导出须经"确认+记审计"**才执行（confirmed 必勾）——体现为显式确认动作，注释标明记审计。
- 数据范围：仅导出/导入 actor 有权访问的记录（注释标明），不越权全量。
- 状态覆盖：导入"部分失败"+导出"成功"两个结果态；注释定义 进行中/失败 态。

---

## 3. 操作日志
产物：`mockups/_src/operation-log.html` + `mockups/operation-log.png`
角色上下文：**管理员视图（仅管理员）**。

### 必须呈现
- 外壳（"操作日志"active）+ 顶栏 + 页头（操作日志 + 副标题"仅管理员可访问"+ 时间/操作人/类型筛选）。
- 只读审计列表/时间线，每条：操作人(actorDisplay + 角色徽 actorRole)、动作(action)、
  对象(resourceType + resourceId)、结果(result)、摘要(safeSummary)、时间(occurredAt)。
- 防篡改提示：每条带 eventHash（可折叠/弱化展示），页面标注"审计链不可篡改"。
- 空态/分页：注释定义 loading/empty；可分页。

### 领域标尺（硬性，源：frontend/src/api/history.ts + oplog.ts）
- 字段对齐 HistoryEvent：actorUserId/actorRole/actorDisplay/action/resourceType/resourceId/
  result/safeSummary/occurredAt/eventHash（**用 safeSummary，不展示 before/after 原始内容**）。
- 角色徽只用三枚举（管理员/销售经理/销售）。
- **整页仅管理员可访问/可见**（operation_log.read 仅 Administrator）——页面级权限门控，注释标明；
  非管理员不应进入。
- 只读：无编辑/删除任何日志的入口（审计不可改）。

---

## 4. 移交给 Codex 决定的设计项
分组/排序方式、卡片 vs 表格、徽变体到枚举的映射(仅既有变体)、各态视觉、是否单独出空/载入态图、
样本数据具体内容（须满足各页领域标尺）。

## 5. Out of scope / 约束
- 不改后端/API/数据模型/业务逻辑；不部署；不弱化任何 P0/P1 或既往 G12 修复。

## 6. G8 审计我会核什么
每页对照领域标尺：真实枚举/字段名、权限门控（提醒 owner 范围、导出确认+记审计、操作日志仅管理员只读）、
不杜撰字段/取值（尤其提醒不臆造已读、操作日志用 safeSummary）、状态覆盖、无 recolor、与参考帧视觉一致、
PREVIEW-ONLY 合规。
