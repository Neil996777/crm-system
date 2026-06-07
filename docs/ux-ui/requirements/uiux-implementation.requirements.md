# 需求规格 —— UI/UX 设计实现（前端落地）

状态：Requirements（Claude 出，作为 Codex G7/G8 任务包的依据 + G8 交接审计/G12 审计基准）。
日期：2026-06-06
模型：新协作模型（Claude 需求+审计 / Codex 设计产出 + G7/G8 任务包 + G9–G11 实现）。
所属变更：UI/UX 完成（`delivery/uiux-completion-charter.md`，定义 release CONTENT）。
范围：`projects/crm-system/` 前端/设计；**不动后端/API/数据模型/业务逻辑**。

## 0. 目标
把已锁定的 UI/UX 设计落到真实 React 前端（`frontend/src/pages/*.tsx` 当前为旧朴素 UI）。
实现目标 = 9 个已 G8 审过的页型 mockup + `docs/ux-ui/design-system.md`（已锁定，设计闭环
2026-06-06）。视觉/交互呈现层改造，行为/逻辑/枚举值不变。

## 1. 实现范围（14 导航页 → 9 页型）
- 概览=工作台；列表(+角色变体)·详情·表单= 线索/客户/联系人/商机/报价/合同/回款/任务（复用范式）；
  报表=报表；管理=用户与角色；提醒中心；导入/导出；操作日志。
- 参考基准帧：`docs/ux-ui/mockups/*.png`（9 页型，含角色变体）。

## 2. "设计已实现"验收标准（Definition of Done，= G8/G12 审计标尺）
### 全局
- A1 设计系统在 React 中落地：颜色/字阶/间距/圆角/阴影 token 与组件（卡片/面板/表格/徽/按钮/
  导航/步进器等）**逐字派生自 `design-system.md`**——不得 recolor、不得新增 token。
- A2 14 个导航页各自换肤到对应页型 mockup；8 个 CRUD 实体复用 列表/详情/表单 范式。
- A3 交互状态按 `screen-state-spec.md` / design-system §8 全实现：loading / empty / error /
  disabled / selected / focused / hover / permission-denied / optimistic-update / success。
- A4 角色/权限门控按已审规则实现（硬性，不得弱化）：
  - 销售列表隐藏 批量转移负责人 / 批量归档（CanEdit/CanArchive）；
  - 详情终态(赢单/丢单)只读；赢单需"合同谈判+已签合同"(DEC-017)、丢单需原因；阶段线性单向；
  - 表单新建/编辑阶段排除终态；销售负责人锁定为本人(CanCreate)；
  - 用户与角色页、操作日志页仅管理员；末位管理员保护置灰；
  - 数据范围：销售本人 / 经理团队 / 管理员全部（CanRead）。
- A5 **无障碍基线**（cross-cutting 标准强制）：键盘可达可操作、可见 focus 环（design-system
  focus-visible）、文本对比度 ≥ AA、交互组件语义/ARIA、表单有 label。
  - **A5↔C6 冲突裁决**（G8 审计补充 2026-06-07）：若某锁定配色实测文本对比度不达 AA，
    属 A5 与 C6(不准 recolor) 的标尺冲突——实现者**不得默默改色**，须按 Kickback 协议回退给
    Claude；Claude 可批准"仅为达标的最小对比度 token 例外"（记入决策），否则维持锁定色。
    此情形应作为一条阻塞/回退条件对待。
  - **DEC-UIUX-A5-001（裁决已下，2026-06-07，关闭 BLK-UIUX-G9-001；用户选定方案 A）**：
    Codex 实测锁定调色板，确认多对"可读文字/背景"不达 AA（`--subtle`#94A3B8 灰字 2.40–2.56；
    success/warning/danger/purple 彩字压各自 *-soft 浅片 2.02–3.95）。裁决——
    **批准一个边界严格的"仅对比度"文字色 token 例外**：
    1) 为作"可读文字"用的状态/强调色，新增最小化的**更深文字变体**（如 `--success-ink`/
       `--warning-ink`/`--danger-ink`/`--purple-ink`，命名 Codex 定），每个=锁定色相**为达 AA
       所需的最小加深**（正文 ≥4.5:1；大/粗体 ≥3:1），其余不变。
    2) 次要/弱化可读文字一律改用已达标的 `--muted`#475569；`--subtle`#94A3B8 **重新归类为
       仅装饰/非文本**（分隔线、边框、装饰图标、禁用态——WCAG 1.4.3 豁免），不得作可读正文。
    **仍锁定不动（C6 继续约束）**：所有背景浅片(`*-soft`/`--card`/`--section`)、实色填充、
    品牌主色、边框、图形/图标色一律**逐字节锁定**（图标/图形对比按 1.4.11 仅需 3:1，可装饰）。
    例外**仅限文字易读性**，不得用于改任何填充/背景/按钮/徽章底色；每个新 token 必须是 text-only。
    **Codex 在 UIUX-001 须交付（回 Claude 复核后方可向下游推进，并并入 G12）**：每个新 `*-ink`
    文字 token 的精确 hex + 重算 WCAG 表证明在其所用背景(白/section/对应 soft)上 ≥4.5:1；
    设计系统 token diff 证明背景/填充/边框/图标 hex **仅有新增、无修改**；用法更新（彩字走 *-ink、
    装饰走原色）。Codex 把此决策作无障碍附录追加进 `design-system.md`（设计系统产出归 Codex）。
- A6 **桌面优先 + 响应式不破版**：1440px 为主锐；在合理宽度区间不错位/不溢出（本轮不要求完整
  移动端断点，design-system §6.6 断点意图作"优雅降级"处理）。
- A7 动效按 `interaction-spec.md` Part B（motion token + prefers-reduced-motion 必须 snap），
  采保守档（已认可决策：导航 hover 展开浮层 NAV-01、计数 count-up 保守 MOTION-02、
  实时合并默认 LIVE-03）。

### 折入实现的 2 个 G8 观察
- A8 提醒中心 status/priority 的显示值对齐后端真实枚举（labels.ts 暂无映射→实现时确认真值，
  必要时补 labels 映射，仍不得改后端比较值）。
- A9 导入/导出样本一致性（archivedIncluded 等）在真实数据下自洽。

## 3. 绑定约束（硬性，不可外移/弱化）
- C1 **仅前端/设计**：不改后端、API、数据模型、业务逻辑、服务边界。
- C2 **无降级**：不弱化任何 P0/P1 功能验收或既往 G12 安全修复（IDOR、durable audit、乐观并发、
  幂等）。换肤**不得改变授权行为或数据暴露面**。
- C3 **保 zh-CN**（phase 1+2）：不得回退英文；**不得改枚举/角色的比较值**（仅经 labels.ts 显示映射）。
- C4 **真实枚举**：六阶段等只用真值；显示走 labels.ts，比较值不动。
- C5 **E2E 保绿**：仅在 DOM 变化处更新选择器；断言/覆盖不减、不 skip；选择器因结构变化失效时
  修选择器而非弱化断言。
- C6 与锁定 mockup + design-system 一致；科技感来自组件/交互而非改色。

## 4. Out of scope
- 后端/新功能；CI/CD 迁移（独立变更，定义 release MECHANISM）；完整移动端断点；部署（另行 gated）。

## 5. 交给 Codex 的产出（G7/G8 任务包，Codex 决定"怎么做"）
- 技术地基方案：把 design-system 抽成 React 样式系统 + 组件库（CSS 变量/组件，技术选型 Codex 定，
  须派生自锁定 token）。
- 任务分解与依赖：地基任务 → 各页型组件 → 各实体应用 → 状态/角色门控 → a11y → e2e 更新；
  每任务绑定受影响 UI 面 + ACC-018/023 + 明确"设计已实现"验收。
- `delivery/` 下的可执行任务文件 + `planning/gate-status.md` 状态更新（G7/G8）。

## 6. G8 交接审计 / G12 审计我会核什么
- G8（任务包）：是否覆盖全部 14 页/9 页型、A1–A9 是否都有任务与验收、C1–C6 是否写入约束、
  可 Codex 执行、无 P0/P1 遗漏。
- G12（实现后）：逐页对照 mockup（视觉）；A3 状态 / A4 角色门控 / A5 a11y 实测；e2e 绿且无 skip；
  zh-CN live 抽检；**安全/授权无回归**（G12 不变量）；枚举/角色比较值 diff 无变化；无后端改动。
