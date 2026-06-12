# UI/UX 多维审计缺口清单（2026-06-12）

> 尺子：`docs/ux-ui/requirements/uiux-audit-matrix.md`。方法：5 个并行代码审计 agent（D2/D3/D5/D8/D10/D12）
> + Claude 实测（D4/D6 待跑）。已消重、按优先级排。OpportunityDetail 关闭赢单/丢单对活跃商机可用（非缺口）。
> 说明：本清单是发现，不是已修；P0/P1 进 blockers 踢回，P2 由 release owner 定批次。

## P1 —— 真功能/反馈缺口（前端可修）

| # | 页面/位置 | D# | 现状 | 期望 |
|---|---|---|---|---|
| F1 | **全部列表页**（leads/accounts/contacts/opportunities/quotes/contracts/payments）`select*()` + `refresh()` | D10/D2 | 无 try-catch → 点名字/行打开记录、或刷新列表**失败时静默无提示**（getX 失败=啥也没发生）。可能就是 release owner 之前"点名字没反应"的真因之一 | 包 try-catch + `setError(localizeError)` 显示错误 |
| F2 | 全部列表页 | D5 | 刷新/选择记录时**无 loading 态**（按钮无 busy、无骨架） | async 期间显 loading/禁用 |
| F3 | QuoteDetail / ContractDetail | D10 | 状态变更 catch 了错误但**组件不渲染 error** → 失败静默 | 渲染 `{error && <alert>}` |
| F4 | CloseOpportunityDialog（关闭赢单/丢单） | D2/D5 | 提交未按必填门控:赢单缺 contractId、丢单缺 reasonCode 也能提交 | 必填齐才允许提交 |
| F5 | Import.tsx:83 | D2 | 未选文件时"开始导入"仍可点（只在提交后报错） | `disabled={busy \|\| !file}` |
| F6 | BasicReports `本月` / `按负责人分组` | D2 | `本月` 只弹"已切换"提示、**不真正筛选**;aria-pressed 写死 true | 真筛选或移除/正确反映状态 |
| F7 | OpportunityList 行菜单项（疑似） | D2/D3 | "转移负责人"等菜单项 `onSelect` 只是 `selectOpportunity()` 打开详情,并非行内执行 → 菜单"承诺"了行内动作却只跳详情 | 行内执行,或标注"（去详情页）",或精简菜单 |

## 决策点（release owner）—— "无后端接口"的禁用按钮

多页有**硬编码 disabled + title「…无接口;按 A3 禁用」**的批量按钮(后端无端点,C1 不许动后端):
PaymentList 批量转移/归档、TaskList 取消任务/批量转移/批量归档、AccountList 批量转移、ContactList 批量转移/归档、
OpportunityList/QuoteList/ContractList 批量转移(部分)。
release owner 之前嫌"半成品感"。**二选一**:(a) 移除这些按钮(更干净);(b) 保留 disabled+原因提示。

## P2 —— 打磨/一致性

| # | 页面 | D# | 现状 → 期望 |
|---|---|---|---|
| P2-1 | PaymentList 标题 | D7/D2 | 标题显示 `opportunityId`,但 aria-label/点击=`contract.id` → 标题应显示主记录标识、与动作一致 |
| P2-2 | ContactList | D7 | 唯一没有 ··· 行菜单(只有名字点击);其余 6 页都有 → 确认是否补一个最小菜单(查看/可用动作)或接受 |
| P2-3 | OperationLogs 筛选 | D8 | "操作人：全部" 实为动作筛选 → 应为"操作：全部" |
| P2-4 | UserManagement | D8 | "唯一启用管理员"标签对任意管理员都显示 → 应仅 `isLastActiveAdministrator` 时显示 |
| P2-5 | Import/Export 类型徽标 | D2 | 客户/联系人/商机/报价/合同 徽标看着可选,实际只有"线索"可选 → 去掉或标注仅展示 |
| P2-6 | Export fileSafety | D8 | 未映射值走硬编码中文兜底 → 补全 fileSafetyLabel 映射 |
| P2-7 | 各列表 StatusPill 色调 | D7 | 每页各自实现 tone 逻辑(重复) → 抽共享 helper |
| P2-8 | 各列表列数/顺序 | D7 | 同类实体列结构不统一 → 按实体类型标准化(认知负荷,非功能) |

## 实测维度结果（Claude 已跑 2026-06-12）

- **D4 响应式 — 基本 PASS**:14 页 × 1280/1366/1440 程序化测,**13 页 0 溢出**;唯一发现:
  **导入/导出页**「开始导出」按钮(+1280px 下确认勾选)略超视口右缘(PAGE_HOVERFLOW=0,无横向滚动)
  → **P2-9 D4**:导出页操作行在窄桌面宽度略溢出,收一下。
- **D6 a11y — focus 环可见 PASS**;主导航键盘可达(TEST-UIUX-A5-001 绿);行菜单 aria-haspopup/Escape/键盘可达;
  行菜单 portal 经早先实测 `elementFromPoint`=菜单(最上层,BLK-018 portal 修复生效)。
  深度键盘遍历(每个表单逐字段)未穷尽,框架到位;如需逐页穷举可再开一轮。
- **D1 视觉**:G12 已逐页 1440px vs mockup 过;本轮新增控件(行菜单/名字链接)视觉随 design-system token,
  无新色(grep clean)。维持 G12 结论。

## 关联在飞的修复
- **BLK-UIUX-G12-018**(行菜单 portal + 名字可点 + 删 →):Codex 已回传(56/56),菜单最上层经实测确认;
  名字点击经 Codex e2e(FUNC-ROWMENU 创建线索→点名字→开详情)验证。**待 Claude 正式签**(数据churn 导致本地探针撞空列表,
  非产品问题;以 e2e + 早先 elementFromPoint 实测为准)。
- **BLK-UIUX-G12-019**(P1 包 + 删无后端禁用按钮):已踢回 Codex。

## 整体读数

- **D12 权限**:全页通过(三角色门控、最后管理员守卫、操作日志只显 safeSummary 无 raw before/after,C2 满足)。✅
- **D8 文案**:基本全 zh-CN、走 labels.ts;少数 P2 标签误写。✅大体
- **D3 导航**:名字点击 + 深链 + 返回 主路径通(F1 的静默失败是边界)。
- **主问题集中在 D2(功能完整性)/D5(状态反馈)/D10(错误处理)** —— 即"控件在但反馈/错误/边界没做全"。
