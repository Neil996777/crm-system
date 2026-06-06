# Design System

## Document Control

- Project: CRM System
- Phase: G4 UI Design (G4c)
- Owner Agent: UI Designer
- Gate: This is the **G4c design-system deliverable**. It is the authoritative
  visual contract for every page archetype (overview/dashboard, list, detail,
  form, reports, admin).
- Source of truth (LOCKED, v4 "premium soft-modern" palette, APPROVED — do not
  recolor): `docs/ux-ui/mockups/_src/dashboard-v7-sales.html`,
  `dashboard-v7-manager.html`, `dashboard-v7-manager-focus.html`.
- Companion specs (do not contradict): `ui-spec.md`, `component-spec.md`,
  `responsive-spec.md`, and the UX-owned `interaction-spec.md` +
  `screen-state-spec.md`.
- Division of ownership: **UI Designer owns VISUAL treatment** (tokens, static
  look of every state). **UX Designer owns BEHAVIOR/MOTION** (timing,
  transitions, easing) in `interaction-spec.md` Part B. This document references
  interaction states by name and defers all motion to the UX spec.
- Status: Authored 2026-06-06; dashboard archetype proven. Pending acceptance as
  Architecture Input alongside the parallel UX interaction layer.

### Value provenance legend

- **[EXTRACTED]** — copied verbatim from the locked CSS. Authoritative.
- **[DERIVED]** — not present in source; defined here consistently with the
  token system. Flagged so it can be spot-checked. The dashboard mostly renders
  the **default** visual state, so loading/empty/error/disabled/selected/
  focused/hover/permission-denied/optimistic/success visual treatments are
  DERIVED from the token set unless a hover/selected/glow precedent exists in
  source.

---

## 1. Design Principles

- **Premium soft-modern.** Calm off-white app canvas, white cards, soft layered
  shadows, generous 16px radii, pastel-tinted accents. Never flat-harsh, never
  neon. The "科技感"/tech feel comes from **components, data-viz, and
  interaction — not from color**. The palette is LOCKED; do not add or
  re-saturate brand colors.
- **Icon-forward.** Every panel and KPI leads with a circular, softly-tinted
  icon badge. Navigation, flow rows, and strip-cards are icon-led. Icons orient
  the eye before text.
- **Soft layered shadows over hard borders.** Elevation is expressed with
  multi-layer low-opacity violet-tinted shadows (`--shadow`, `--badge-shadow`),
  paired with a near-invisible hairline border (`--border`). Cards float; they
  are not boxed.
- **Capsule + ghost-track geometry.** Pills (`border-radius:999px`) are the house
  shape for bars, badges, buttons, search, nav items, and the live segment. Bars
  ride inside a soft ghost track (`--section` fill) with a rounded fill on top.
- **Blob / liquid data-viz.** Donuts use round-capped stroke arcs on a soft
  track; trend lines use a smooth bezier with a gradient area fill ("liquid"
  underglow). Avoid bar-chart-grid heaviness.
- **Numeric discipline.** `font-variant-numeric:tabular-nums` is set globally so
  money, counts, deltas, and table figures align column-to-column. Money and
  counts are heavier weight (700–800).

---

## 2. Color Tokens

All values **[EXTRACTED]** verbatim from the `:root` block, identical across all
three locked files. No semantic label appears without a value. Do not recolor.

### Brand

| Token | Value | Usage |
|---|---|---|
| `--primary` | `#2563EB` | Primary brand. Primary buttons, active nav, links, fill gradients, live dot, focus ring base, selected accents. |
| `--primary-hover` | `#1D4ED8` | Primary hover/pressed; the darker stop in the "glow" funnel-fill gradient. |

There is a recurring **primary-tint surface `#EAF1FF`** [EXTRACTED] used for
active nav background, primary badge, avatar background, segment-on background,
hovered table row, and the `arrived`/selected row highlight. It is a literal hex
in source (not a `:root` var); treat it as the canonical **selected/active soft
fill** and reuse it — do not invent a different tint.

### Support pastels (each as solid + soft pair)

Solid = swatch/legend/accent dot; soft = icon-badge background and badge fill.

| Token | Value | Token (soft) | Value | Usage |
|---|---|---|---|---|
| `--lavender` | `#8B93F8` | `--lavender-soft` | `#ECEDFE` | Default/neutral icon badge; donut segment 1 (报价/quote); report accent dot; rail avatar bg. |
| `--sky` | `#5FB8F5` | `--sky-soft` | `#E6F4FE` | "sky" icon badge (funnel/pipeline/contract); lighter stop in funnel-fill gradient. |
| `--mint` | `#5BC8A0` | `--mint-soft` | `#E5F7F0` | "mint" icon badge (won/leaderboard); donut segment 2 (赢单/won); success badge fill; payment summary fill. |
| `--peach` | `#F6A98A` | `--peach-soft` | `#FDEDE5` | "peach" icon badge (tasks/payments/alerts); donut segment 3 (合同谈判); warning badge fill; pay flow-icon fill. |
| `--purple` | `#B79CF0` | `--purple-soft` | `#F2ECFD` | "purple" icon badge (stage composition); donut segment 4 (丢单风险/loss risk). |

### Neutrals / text / border / section

| Token | Value | Usage |
|---|---|---|
| `--text` | `#0F172A` | Primary text, h1, row titles, donut center number. |
| `--muted` | `#475569` | Secondary text, labels, button text, nav idle, eyebrows, chart axis labels. |
| `--subtle` | `#94A3B8` | Tertiary/de-emphasized text, meta, timestamps, placeholders, expand icon, donut caption. |
| `--border` | `#EDF0F6` | Hairline borders, dividers, chart gridlines, input/button outlines. |
| `--card` | `#FFFFFF` | Card/panel/table surface. |
| `--section` | `#F6F7FD` | Ghost track fill, idle rail icon bg, neutral badge bg, table header bg, donut track. |
| `--app` | `linear-gradient(135deg,#F3F2FC 0%,#F7F8FE 45%,#FBFCFF 100%)` | App canvas behind cards. Page `background:#FBFCFF` (the gradient's end stop) is the flat fallback. |

### Semantic

| Token | Value | Soft pairing | Usage |
|---|---|---|---|
| `--success` | `#16A34A` | `--mint-soft` `#E5F7F0` | Won/paid/positive delta text + money-arrived value; `badge.success`, `delta.up`. |
| `--warning` | `#D97706` | `--peach-soft` `#FDEDE5` | Overdue/expiring/negotiation/pending-signature; `badge.warning`, `delta.warn`. |
| `--danger` | `#DC2626` | `#FEE2E2` | Lost/overdue/invalid/denied; `badge.danger`, `delta.down`. **Danger soft fill `#FEE2E2`** is a literal hex in source (no `:root` var); treat as the canonical danger soft fill. |

### Gradients

| Name | Value | Usage |
|---|---|---|
| App canvas | `linear-gradient(135deg,#F3F2FC 0%,#F7F8FE 45%,#FBFCFF 100%)` (`--app`) | Behind all cards. |
| Funnel fill | `linear-gradient(90deg,var(--primary),var(--sky))` | Default capsule-bar fill. |
| Funnel fill (glow) | `linear-gradient(90deg,var(--primary-hover),var(--sky))` | Emphasized/active funnel stage. |
| Live report banner | `linear-gradient(90deg,#ECEDFE 0%,#F2ECFD 52%,#fff 100%)` | "今日实时战报" capsule banner. |
| Trend area fill | `#2563EB` at `.18` opacity → `#2563EB` at `0` opacity (vertical) | Line-chart "liquid" underglow (`linearGradient` defs). |

### Shadows

| Token | Value | Usage |
|---|---|---|
| `--shadow` | `0 1px 2px rgba(17,24,39,.04), 0 8px 24px rgba(82,72,156,.06), 0 24px 48px rgba(82,72,156,.05)` | Card/panel/stage/side-card elevation. Three-layer, violet-tinted (`82,72,156`). |
| `--badge-shadow` | `0 2px 6px rgba(82,72,156,.10)` | Small lifted elements: icon badges, badges, expand button, flow icons, swatches, rail icons, chips, ghost buttons. |
| Funnel fill drop | `0 8px 18px rgba(37,99,235,.14)` | Default funnel-fill underglow. |
| Funnel fill (glow) drop | `0 0 0 1px rgba(37,99,235,.16), 0 10px 22px rgba(37,99,235,.22)` | Emphasized funnel-fill ring + deeper underglow. |
| Funnel track inset | `inset 0 1px 2px rgba(82,72,156,.05)` | Ghost-track inner shadow (focus-size funnel only). |
| Live-dot halo | `0 0 0 4px color-mix(in srgb, var(--primary) 12%, #fff)` | Soft ring around the 7px live status dot. |

---

## 3. Typography Scale

- **Family** [EXTRACTED]: `Inter, "Noto Sans SC", ui-sans-serif, system-ui,
  -apple-system, BlinkMacSystemFont, "Segoe UI", "Microsoft YaHei", sans-serif`.
  Inter for Latin/numerals, Noto Sans SC for 中文.
- **Base** [EXTRACTED]: `font:14px/1.35` on `html,body`.
- **Numeric** [EXTRACTED]: `font-variant-numeric:tabular-nums` global. All
  money, counts, deltas, table figures use tabular numerals.

| Role | Size / line-height | Weight | Color | Source |
|---|---|---|---|---|
| Page title `h1` | 24px / 1.15 (focus) · 1 (overview) | 800 | `--text` | [EXTRACTED] |
| KPI metric | 30px / 1 | 700 | `--text` | [EXTRACTED] |
| Donut center number | 22–24px | 700 | `--text` | [EXTRACTED] |
| Report item value | 17px | 700 | `--primary` | [EXTRACTED] |
| Brand title | 16px | 700 | `--text` | [EXTRACTED] |
| Side-card value `strong` | 16px | 700 | `--primary` | [EXTRACTED] |
| Panel title (sales) | 16px | 600 | `--text` | [EXTRACTED] |
| Panel title (manager dense) | 15px | 600 | `--text` | [EXTRACTED] |
| Funnel label (focus) | 14px | 700 | `--muted` | [EXTRACTED] |
| Body / base | 14px / 1.35 | 400–500 | `--text`/`--muted` | [EXTRACTED] |
| Row title / value | 12–13px / 1.24–1.28 | 700 | `--text` | [EXTRACTED] |
| Page subtitle, quicklink, nav, stage-sub | 13px | 500–700 | `--muted`/`--primary` | [EXTRACTED] |
| Meta / panelMeta / updateMeta / timestamps | 12px | 500 | `--subtle` | [EXTRACTED] |
| Badge / delta / segment | 12px | 700 | semantic | [EXTRACTED] |
| Table header | 12px (sales 11px dense) | 600 | `--muted` | [EXTRACTED] |
| Micro (manager dense row span, table cells) | 11px | 500–700 | `--subtle`/`--text` | [EXTRACTED] |

Canonical scale: **24 / 17 / 16 / 15 / 14 / 13 / 12 / 11**. Weights in use:
**800 / 700 / 600 / 500 / 400**. Do not introduce sizes/weights outside this set
without amending this section.

---

## 4. Spacing, Radius, Elevation

- **Box model** [EXTRACTED]: `*{box-sizing:border-box}` globally. All width/
  height math (strip-card 92px, fixed table layout, padded tracks) assumes
  border-box. Frontend must preserve this — it is load-bearing for the
  fixed-height layouts.
- **Spacing rhythm** [DERIVED, observed]: multiples of 2 on a **4px base**.
  Observed values: 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 24, 28. Primary grid/card
  gap is **20px**; tight internal gaps 6–12px; section/card padding 16–28px.
- **Radius tokens** [EXTRACTED]:
  - `--radius: 16px` — cards, panels, stage, side-card, chart box, table wrap.
  - `--inner-radius: 12px` — declared for nested inner surfaces (reserved;
    apply to inner cards/insets when needed).
  - Brand mark: `10px`. Capsules (`999px`): bars, badges, buttons, search, nav
    items, segment, summary pills, chips. Circles (`50%`): all icon badges,
    avatars, rail icons, flow icons, swatches, live dot.
- **Elevation tiers**:
  1. Canvas (`--app` gradient, no shadow).
  2. Card/panel/stage/side-card → `--shadow` + 1px `--border`.
  3. Small lifted chips/badges/icons → `--badge-shadow`.
  4. Active data emphasis → funnel-fill glow ring/drop shadows.
- **Grid gaps & padding** [EXTRACTED]:
  - Role grid / KPI grid / focus columns gap: **20px**.
  - Content padding: `20px 24px` (overview), `24px` (focus).
  - Topbar padding: `0 24px`. Panel padding: 18px (sales) / 16px (manager dense)
    / 28px (focus stage) / `10px 16px` (strip-card).
- **Card recipe** [EXTRACTED]: `.card{background:var(--card); border:1px solid
  var(--border); border-radius:var(--radius); box-shadow:var(--shadow)}`.

---

## 5. Iconography

- **Form**: line/stroke SVG, no fill. Global rule [EXTRACTED]:
  `svg{stroke:currentColor; stroke-width:1.8; fill:none; stroke-linecap:round;
  stroke-linejoin:round}`. Icon color inherits the badge's `color`.
- **Stroke icon sizes** [EXTRACTED]: 18px inside icon badges / nav / rail /
  expand; 16px in side-cards and ghost buttons; 13px (sales) / 12px (manager) in
  flow icons. SVG `viewBox` is `0 0 16 16` for nav/rail/flow, `0 0 20 20` for
  panel icon badges.
- **Circular tinted icon badges** [EXTRACTED]:
  - Panel/KPI badge: **40px** circle (KPI uses 44px), `--badge-shadow`,
    soft-tint background per support color, icon color = the paired solid (or
    semantic) color.
  - Strip-card (focus) badge: **36px** circle, same tinting rules.
  - Flow-icon (payment/timeline rows): 24px (sales) / 22px (manager) circle.
- **Tint→color mapping** [EXTRACTED]:
  - default/lavender → bg `--lavender-soft`, icon `--primary`.
  - `.sky` → bg `--sky-soft`, icon `--primary`.
  - `.mint` → bg `--mint-soft`, icon `--success`.
  - `.peach` → bg `--peach-soft`, icon `--warning`.
  - `.purple` → bg `--purple-soft`, icon `--purple`.
  - `.flowIcon.pay` → bg `--peach-soft`, icon `--warning`.
- **When an icon leads a panel**: every panel header, KPI, strip-card, flow row,
  nav item, and rail item is icon-led. A panel without a leading icon is
  off-system. Icon-only controls (expand, rail icons) require an accessible name
  / tooltip per `component-spec.md` and `interaction-spec.md` accessibility
  rules.

---

## 6. Layout System

**Desktop-first posture** [EXTRACTED]: the locked mockups are authored at
`width:1440px` (canvas `1440px`, content column `1fr`). This is the reference
composition width. Breakpoint intent below aligns with `responsive-spec.md`
(Desktop ≥1200, Tablet 768–1199, Mobile <768); mobile/tablet reflow is governed
there and is out of fidelity scope this round but is **not omitted** — the
responsive contract still applies and must not be downgraded.

### 6.1 Navigation rule (mutually exclusive)

The shell has **exactly two navigation modes; never both at once**:

- **Expanded icon+text sidebar** — **248px** [EXTRACTED]. Used for overview/home
  and standard pages. Brand block (mark + title + subtitle) + vertical nav list
  of `navItem` (icon 16px + label, 44px tall, capsule, active = `#EAF1FF` bg +
  `--primary` text weight 700). App grid: `grid-template-columns:248px 1fr`.
- **Collapsed icon-only rail** — **72px** [EXTRACTED]. Used for focus/center-
  stage mode. `railBrand` (CRM mark) + vertical `railIcon` stack (44px circles,
  idle `--section`, active `--primary`) + `railAvatar` at bottom. App grid:
  `grid-template-columns:72px 1fr`.

Switching to focus collapses 248px→72px; returning expands back. Do not render a
half-expanded or dual sidebar.

#### Rail hover-expand flyout (NAV-01, ACCEPTED 2026-06-06) [DERIVED]

The collapsed 72px rail hover/focus-expands a **temporary label flyout overlay**
(DEC-UX-NAV-01). It is a transient overlay, not a layout mode: it **floats over
content and does not reflow the stage or change the 72px track**, so the
mutually-exclusive nav rule (§6.1) still holds. UI owns the static look; motion
(`motion-fast` `ease-decelerate` in, `ease-accelerate` out, snap under
reduced-motion) is owned by `interaction-spec.md` Part B / B-micro.

- **Surface**: a floating card anchored to the rail's right edge, vertically
  aligned to the rail icon stack. `background:var(--card)` `#FFFFFF`, 1px
  `--border` `#EDF0F6`, `border-radius:var(--radius)` 16px, `box-shadow:var(--shadow)`
  (the same three-layer violet-tinted card elevation — it floats like a card,
  not a hard popover). Inner padding 8px; small horizontal offset off the rail
  (~8px gap) so the shadow reads.
- **Items**: each row mirrors the **expanded nav item** token set (§7.9): flex
  row, 44px tall, `border-radius:999px`, `padding:0 12px`, gap 10px, icon 16px +
  label, idle `--muted`/500. The flyout shows the **label** the rail icon stands
  for — the icon may repeat or be omitted; the label is the point. **Active item**
  = bg `#EAF1FF`, text `--primary`/700 (matches the active rail icon, so the
  current location is unambiguous in the flyout).
- **No scrim.** The flyout is a lightweight label aid, not a modal; it does not
  dim the workspace (the focus-stage scrim `rgba(15,23,42,.06)` is a separate,
  unrelated overlay and is not re-used here). Hover-out closes after a small
  delay; focus-out closes.
- **z-index**: above the rail and the focus-stage scrim/stage so labels are never
  clipped by the stage; below toasts. (It must paint over the `1fr` content area
  without participating in its layout.)
- **focus-visible**: entering a rail icon by keyboard opens the flyout (parity
  with hover, per B6 "no hover-only"); the focused item shows the standard
  focus-visible ring (§8 `focused`: `0 0 0 3px color-mix(in srgb, var(--primary)
  35%, #fff)`), which must stay legible over both the white flyout surface and the
  `#EAF1FF` active-item fill. Each item keeps a ≥40px hit area (B6).
- **Reuse**: clicking a flyout item navigates identically to clicking the rail
  icon; the flyout adds no new action, only the label.

### 6.2 Topbar — 64px [EXTRACTED]

`grid-template-rows:64px 1fr`. Flex row, `padding:0 24px`, bottom 1px `--border`,
`background:rgba(255,255,255,.92)`. Contents left→right:

- **Search capsule**: 420px × 40px, `border-radius:999px`, 1px border, white,
  placeholder `--subtle` ("搜索线索、客户、合同或负责人").
- **Spacer** (`flex:1`).
- **Quick links** (`新建线索`, `导出报表`): `--primary`, 13px/700.
- **User menu**: capsule, 40px tall, 24px avatar circle (`#EAF1FF` bg) + name +
  email (`王敏 · sales@example.com`).

Page-level action cluster (`本月` button, `刷新数据` primary button, the
`实时更新 / 暂停` live segment toggle, `更新于 09:30` meta) lives in the **page
title row** under the topbar in the overview layout [EXTRACTED] — not in the
topbar itself. The topbar holds global search/quick-links/user; the page header
holds page-scoped controls.

### 6.3 Role home grids (equal-card)

- **Sales 我的工作台**: `roleGrid` = `grid-template-columns:repeat(3,1fr);
  grid-auto-rows:308px; gap:20px` → **6 panels, 3×2** [EXTRACTED].
- **Manager 团队工作台**: `repeat(4,1fr); grid-auto-rows:268px; gap:20px` →
  **8 panels, 4×2** [EXTRACTED].
- Equal-card grid: every panel is the same fixed-height cell (`grid-auto-rows`);
  panels never define their own height. Above the grid: page title row → live
  report banner → 4-up KPI row (`kpis: repeat(4,1fr); gap:20px`).

### 6.4 Overflow-safe panel recipe [EXTRACTED, manager file]

For panels inside a fixed-height grid cell whose content can overflow:

```
.panel { display:flex; flex-direction:column; overflow:hidden; }
.panelHeader { flex:0 0 auto; }            /* header never shrinks */
.list/.tableWrap/.timeline { flex:1; min-height:0; overflow:hidden; }
.footer { flex:0 0 auto; margin-top:auto; } /* footer pinned to bottom */
```

`min-height:0` on the flex child is mandatory or the panel will blow out its
cell. Footer uses `margin-top:auto` to pin to the bottom. This recipe is the
contract for any panel that lists rows in a height-capped card.

### 6.5 Focus layout [EXTRACTED, manager-focus file]

- `focus` = `grid-template-columns:1fr 300px; gap:20px; padding:24px`.
- **Center stage**: `1fr`, `padding:28px`, `min-height:1010px`, `--shadow`.
  Holds the expanded artifact (e.g., funnel + detail table). Stage head =
  icon + title + live sub on the left, ghost `返回` + `Esc 返回` chip on the right.
- **Right strip**: fixed **300px** column, `side` = `grid; gap:12px;
  align-content:start` of compact strip-cards.
- **Strip-card**: fixed **92px** height [EXTRACTED], `padding:10px 16px`,
  `flex-direction:column; justify-content:center; overflow:hidden`. Top row
  `grid-template-columns:36px 1fr auto` (icon / title / optional live dot);
  value row indented `padding-left:48px` (36px icon + 12px gap). Title + value
  truncate with ellipsis. Border-box math keeps the card at 92px regardless of
  content — do not let content change the height.
- Behind the focus content, the workspace dims (`rgba(15,23,42,.06)` overlay
  below the topbar) to spotlight the stage [EXTRACTED].

### 6.6 Breakpoint intent

| Name | Width | Layout posture |
|---|---|---|
| Desktop | ≥1200px (composition ref 1440px) | Full sidebar/rail, multi-column grids, full tables. |
| Tablet | 768–1199px | Collapse nav, reduce grid columns and table columns (see `responsive-spec.md`). |
| Mobile | <768px | Single-column stacked cards/rows, drawer nav (see `responsive-spec.md`). |

Mobile token values (font sizes, paddings) are **[DERIVED]** downstream; this
round delivers desktop fidelity. The responsive contract in `responsive-spec.md`
is not weakened.

---

## 7. Component Specs

Each maps to a `component-spec.md` CMP ID where one exists. Anatomy + variants +
exact tokens.

### 7.1 Card / Panel (→ container for CMP-002/004/011/012/014)

- **Anatomy**: white surface, 1px `--border`, `--radius` 16px, `--shadow`.
  `.panel` adds padding (18/16/28px) and, in fixed grids, the §6.4 flex recipe.
- **Variants**: standard panel; KPI card; stage (focus, 28px pad); side strip-
  card (92px); chart box (inner, `--radius`, height 184–218px).

### 7.2 Panel header (icon + title + meta + expand)

- **Anatomy** [EXTRACTED]: `panelHeader` flex row, `min-height:38–40px`,
  `margin-bottom:10–12px`. Left `titleGroup` = icon badge (40px) + stacked
  (`panelTitle` 15–16px/600, `panelMeta` 12px/500 `--subtle`, optional live dot).
  Right = `expand` button (30–32px circle, 1px border, white, `--subtle` icon,
  `--badge-shadow`, expand-corners SVG).
- **Behavior**: expand triggers Card→focus (§10). Motion deferred to UX spec.

### 7.3 Capsule / ghost-track bar (funnel fill) (→ CMP-008 funnel viz)

- **Anatomy** [EXTRACTED]: `funnelTrack` = pill (`999px`), `--section` fill,
  padding 4–5px, height 20px (manager) / 22px (sales) / 32px (focus). Inside,
  `funnelFill` = pill, gradient `primary→sky`, height 12/14/22px, drop shadow
  `0 8px 18px rgba(37,99,235,.14)`. Width = the metric percentage (inline %).
- **Variants**: default; `.glow` (gradient `primary-hover→sky` + ring +
  `0 10px 22px rgba(37,99,235,.22)`) for the emphasized stage.
- **Row layout**: `funnelRow` grid label / track / value (e.g.
  `86px 1fr 96px` sales, `118px 1fr 190px` focus with a `rate` subline).

### 7.4 Blob / liquid donut (→ CMP-014 share viz)

- **Anatomy** [EXTRACTED]: SVG `viewBox` 170 (manager) / 190 (sales). Base track
  `circle stroke:#F6F7FD` (=`--section`), `stroke-width:24–26`. Segments =
  round-capped (`stroke-linecap:round`) arcs via `stroke-dasharray` /
  `stroke-dashoffset`, `transform:rotate(-90 …)`, in segment order lavender
  `#8B93F8` → mint `#5BC8A0` → peach `#F6A98A` → purple `#B79CF0`. Center
  `donutCenter` text (22–24px/700 `--text`) + caption (`--subtle` 11–12px).
- **Legend** [EXTRACTED]: `donutWrap` grid `128–150px 1fr`. `legendItem` grid
  `swatch 22–24px / label / num`. Swatch = soft-tint circle (`*-soft`),
  `--badge-shadow`.

### 7.5 Data table (→ CMP-004 / CMP-013 / CMP-020)

- **Anatomy** [EXTRACTED]: `tableWrap` = 1px border, `--radius`, `overflow:hidden`.
  `table{width:100%; border-collapse:collapse; table-layout:fixed}`. `th` =
  `--section` bg, `--muted`, 11–12px/600, left-aligned, height 34px (manager) /
  48px (focus). `td` height 36/52px, top 1px `--border`. All cells
  `white-space:nowrap; overflow:hidden; text-overflow:ellipsis`. **Last column
  right-aligned** (`th:last-child, td:last-child{text-align:right}`) — numerics/
  money right-align. `.money` weight 800.
- **When to use a table vs a compact 2-line list**: use a **table** when ≥3
  comparable columns must align (leaderboard, focus detail, admin/logs). In a
  **narrow grid card** (overview panels), prefer the **compact 2-line `row`
  list** (§7.7): a wide column for `strong` title + `span` subline, and a right
  badge cluster. The locked manager overview uses the 3-column table only for the
  leaderboard; all other narrow panels use the 2-line list. Do not force a
  multi-column table into a narrow card.

### 7.6 Badge (→ CMP-007 Status Badge)

- **Anatomy** [EXTRACTED]: inline-flex pill, height 22px (manager) / 24px,
  padding `0 8–10px`, 12px/700, `--badge-shadow`, `white-space:nowrap`.
- **Variants** [EXTRACTED]:
  - `primary` → bg `#EAF1FF`, text `--primary`.
  - `success` → bg `--mint-soft`, text `--success`.
  - `warning` → bg `--peach-soft`, text `--warning`.
  - `danger` → bg `#FEE2E2`, text `--danger`.
  - `neutral` → bg `--section`, text `--muted`.
- **Rule** (from `component-spec.md` CMP-007): badge text is explicit; color
  alone is never sufficient.
- **Delta variant** (KPI): same pill; `up`→mint-soft/success, `down`→#FEE2E2/
  danger, `warn`→peach-soft/warning; arrow glyph ▲/▼ inline.

### 7.7 Payment / flow row + compact list row (→ CMP-011 Reminder Row)

- **List row** [EXTRACTED]: `row` grid `1fr auto`, `min-height:39–44px`, bottom
  1px `--border`. Left = `strong` title (12–13px/700 `--text`, ellipsis) + `span`
  subline (11–12px `--subtle`). Right = `badges` cluster (flex, gap 5–6px).
- **Payment row** [EXTRACTED]: grid `icon 22–24px / 1fr / auto`. `flowIcon`
  circle (`--lavender-soft`/`--primary`, `.pay`→`--peach-soft`/`--warning`).
  Right = `.money` value + status badge.
- **Arrived/selected highlight** [EXTRACTED]: `.arrived` bleeds to card edge
  (`margin:0 -16/-18px; padding-left/right:16/18px`), bg `#EAF1FF`, left 2px
  `--primary` border. This is the canonical **selected/new-item visual** for
  rows.
- **Payment summary pill** [EXTRACTED]: `paymentSummary` pill, height 24–26px,
  bg `--mint-soft`, `--muted` label, `.money` in `--success`. Pinned bottom-right
  via `footer{margin-top:auto}`.

### 7.8 Compact strip-card (focus right strip) (→ CMP-014 compact)

See §6.5 for full anatomy. 92px fixed, 36px icon, ellipsis title + value;
`sideValue strong` = `--primary` 16px, `.plain` = `--text` 14px, optional live
dot in top row.

### 7.9 Nav item (expanded) + rail icon (collapsed) (→ CMP-001 App Shell)

- **Nav item** [EXTRACTED]: flex row, 44px tall, capsule, `padding:0 12px`, gap
  10px, icon 16px + label, `--muted`/500. **Active** = bg `#EAF1FF`, text
  `--primary`/700.
- **Rail icon** [EXTRACTED]: 44px circle, idle bg `--section`/`--muted`,
  `--badge-shadow`. **Active** = bg `--primary`/white. `railBrand` (CRM,
  `--primary`/white) at top; `railAvatar` (`--lavender-soft`/`--primary`) at
  bottom.

### 7.10 Topbar live segment toggle

- **Anatomy** [EXTRACTED]: `liveSegment` = pill track, 38px tall, 3px inner pad,
  1px border, white. Two segments (`segmentOn` / `segmentOff`), each 30px,
  capsule, 12px/700. **On** = bg `#EAF1FF`, text `--primary`, leading 7px
  `liveDot` (`--primary` + halo). **Off** = `--muted`, transparent.
- States: this is a segmented toggle; selected segment = the on-treatment. Motion
  on switch deferred to UX spec.

### 7.11 Buttons (→ CMP-002 actions)

- **Primary** [EXTRACTED]: 38px, capsule, bg `--primary`, border `--primary`,
  white, weight 700, padding `0 14px`. Hover → `--primary-hover` [DERIVED].
- **Ghost / secondary** [EXTRACTED]: 38px (34px ghost in focus), capsule, white,
  1px `--border`, `--muted`/700. Focus-stage `ghost`/`chip` add `--badge-shadow`;
  `chip` = `--section` bg, `--subtle`, 12px (used for `Esc 返回` hint).

### 7.12 KPI stat card (→ CMP-014 Metric Tile)

- **Anatomy** [EXTRACTED]: `kpi` card grid `1fr 48px`, `min-height:100–104px`,
  `padding:18px`. Left column = `eyebrow` (12px/500 `--muted`) → `metric`
  (30px/1/700 `--text`) → `delta` pill. Right = 44px icon badge.
- **Use**: single headline number with trend. One KPI = one number; do not pack
  multiple metrics into one tile.

### 7.13 Live report banner

- **Anatomy** [EXTRACTED]: full-width capsule card, `min-height:48px`, gradient
  `#ECEDFE→#F2ECFD→#fff`, grid `150px repeat(4,1fr)`. Lead (`liveDot` + 今日实时
  战报, 13px/700) + 4 `reportItem`s (accent dot + label `--muted` 13px + `strong`
  value `--primary` 17px/700).

---

## 8. Visual States

The dashboard renders the **default** state. The treatments below are the
authoritative **static visual look** for each interactive component, using the
**exact state names** shared with the UX spec. **All timing/easing/motion is
deferred to `interaction-spec.md` (Part B).** Reconcile names 1:1 with
`screen-state-spec.md`.

| State | Visual treatment | Provenance |
|---|---|---|
| **default** | As specified per component in §7. | [EXTRACTED] |
| **hover** | Hover-lift: row/table `tr.hovered td` → bg `#EAF1FF` (extracted). Cards/buttons lift by deepening to `--shadow` with no color change; primary button → `--primary-hover`. Cursor pointer on actionable. | row/table [EXTRACTED]; card/button [DERIVED] |
| **focused** | Focus ring: `0 0 0 3px color-mix(in srgb, var(--primary) 35%, #fff)` outer ring + retain border. Mirrors the live-dot halo idiom (`box-shadow 0 0 0 Npx`). Always visible, never `outline:none` without a replacement. | [DERIVED] (idiom from live-dot halo) |
| **selected / active** | Soft fill `#EAF1FF` + accent: nav/segment active (extracted), row `.arrived` bleed + 2px `--primary` left border (extracted), funnel `.glow`. `#EAF1FF` is the canonical selected fill. | [EXTRACTED] |
| **disabled** | `opacity:0.45`; remove `--badge-shadow`/lift; `cursor:not-allowed`; keep label legible (no color inversion). Disabled actions still explain the missing requirement per `screen-state-spec.md` where safe. | [DERIVED] |
| **loading** | Skeleton: replace text/number with rounded (`999px` for pills, `--radius` for blocks) placeholder filled `--section` `#F6F7FD`; container keeps its size (layout-stable per `screen-state-spec.md`). Shimmer motion → UX spec. Inline action loading → spinner in `--primary` on white. | [DERIVED] |
| **empty** | Card body shows a centered empty message in `--muted` + optional 40px `--lavender-soft` icon badge + an allowed next-action ghost button. Names the missing data type; no unauthorized action (per CMP-018). | [DERIVED] |
| **error** | Inline/section alert: `#FEE2E2` bg, `--danger` text/icon, 1px ring `color-mix(in srgb, var(--danger) 30%, #fff)`, `--radius`. Field error: 1px `--danger` border + message below in `--danger` 12px. Safe summaries only (per CMP-010); never echo restricted values. | [DERIVED] (danger soft fill extracted) |
| **permission-denied** | `CMP-017` panel: neutral card, `--section` icon badge, `--muted` safe denial message, return ghost button. No restricted record name/existence, no danger styling (denial is not an error), no bypass action. | [DERIVED] |
| **optimistic-update** | Show the projected value immediately at reduced emphasis: value at `opacity:0.65` (or `--subtle`), no success badge yet; the new row uses the `.arrived` highlight idiom. On confirm → success; on failure → error + revert. Static look only; transition timing → UX spec. | [DERIVED] |
| **success** | Success pill/badge `--mint-soft` / `--success` (extracted); value resolves to full emphasis `--text`/`--success`; section success alert `--mint-soft` bg + `--success` text. Must identify the changed record (per `screen-state-spec.md`). | [EXTRACTED] base tokens |

Accessibility (carried from companion specs): focus indicator always visible;
target ≥44px where practical; status conveyed by text + color, never color
alone; contrast meets WCAG AA. The soft pastel-on-white badge fills already pair
dark semantic text with light tints for legibility.

---

## 9. Data-Visualization Decisions

Pick by metric type. **Defaulting a dashboard or report to a generic table is
forbidden** (UI Designer failure condition). Tables are for aligned multi-column
record comparison only.

| Metric type | Visualization | Component | Locked precedent |
|---|---|---|---|
| Single headline number + trend | **KPI stat card** | §7.12 | 本月新增线索, 商机金额, 赢单, 待办 |
| Stage progression / conversion | **Capsule funnel bars** (+ rate subline in focus) | §7.3 | 销售漏斗 (sales/manager/focus) |
| Part-to-whole composition | **Blob/liquid donut + legend** | §7.4 | 商机阶段构成 |
| Time series (≤ ~12 points) | **Liquid line chart** (bezier + gradient area) | chart box | 赢单金额趋势 近6月 |
| Ranked comparison across people/entities | **Table** (right-aligned numerics) | §7.5 | 销售业绩榜 |
| Live event/payment stream | **Flow/payment rows** (icon-led, `.arrived` for new) | §7.7 | 回款到账, 最近活动 |
| Task/reminder/risk list | **Compact 2-line list + status badges** | §7.7 | 待办与提醒/预警 |
| At-a-glance multi-metric summary | **Strip-cards** (focus) / **live report banner** | §7.8 / §7.13 | 团队工作台 focus strip, 今日实时战报 |

Reports archetype: build from the same vocabulary (KPI cards, donuts, funnels,
liquid line, then tables for grouped detail). A report screen that is only a
table is off-system.

---

## 10. Card → Focus (Visual Side Only)

Locked **visual end-states** only; the transition **motion is owned by the UX
interaction spec** (`interaction-spec.md` Part B) — defer timing/easing/sequence
there.

- **Before (overview)** [EXTRACTED]: equal-card grid (Sales 3×2 308px / Manager
  4×2 268px), full **248px** expanded sidebar, each panel has an `expand`
  control.
- **After (focus)** [EXTRACTED]:
  - The activated panel becomes the **center stage** (`1fr`, 28px pad,
    `min-height:1010px`, full artifact: e.g. larger funnel rows + a detail
    table).
  - The other panels collapse into the **300px right strip** of **92px
    strip-cards** (icon + title + value, ellipsis).
  - The sidebar collapses **248px → 72px icon rail**.
  - The workspace behind the stage **dims** (`rgba(15,23,42,.06)` overlay under
    the topbar) to spotlight the stage.
  - Exit affordances: `返回` ghost button + `Esc 返回` chip in the stage head.
- **Reverse**: strip-cards re-expand to the grid; rail re-expands to 248px; dim
  clears. Mutually-exclusive nav rule (§6.1) holds throughout.

---

## 11. Traceability / Status

- **Supersedes / feeds**: this document is the **token + visual authority** that
  `ui-spec.md` and `component-spec.md` defer to for concrete values. Where those
  files describe component structure/variants (CMP-001…CMP-022), this file
  supplies the **exact tokens and visual treatment**; they remain the structural
  index. No requirement in those files is weakened (no-downgrade).
- **Companion (not superseded)**: `interaction-spec.md` (UX-owned behavior/
  motion) and `screen-state-spec.md` (canonical state names). This file
  references their state names and defers motion to them.
- **Exemplar**: the **dashboard/overview archetype is the proven exemplar**.
  Every future archetype (list, detail, form, reports, admin) **must conform** to
  these tokens, the capsule/soft-shadow/icon-forward language, the navigation
  rule, the overflow-safe panel recipe, and the data-viz decision table.
- **Gate**: marked as the **G4c design-system deliverable**.
- **Locked palette**: v4 premium soft-modern, APPROVED. No recolor; no new brand
  colors. Any future color need must reuse an existing token or be raised as a
  formal change.

### Derived-value index (for spot-check)

Values **not** in the locked CSS, defined here consistently:

- 4px spacing base / rhythm enumeration (§4) — observed-pattern derivation.
- `--inner-radius` application targets (declared in source, usage derived).
- Focus ring `0 0 0 3px color-mix(... 35% ...)` (§8) — idiom borrowed from the
  live-dot halo.
- `disabled opacity:0.45`, loading skeleton fill = `--section`, empty/error/
  permission-denied/optimistic panel treatments (§8).
- Card/button hover-lift and primary-button `--primary-hover` on hover (§8) —
  table-row hover `#EAF1FF` is extracted; card/button hover is derived.
- All mobile/tablet token specifics (§6.6) — desktop fidelity delivered;
  responsive contract lives in `responsive-spec.md`.

### States with motion deferred to the UX interaction spec (for reconciliation)

hover, focused, selected/active, loading (skeleton shimmer + inline spinner
spin), optimistic-update (projected→confirmed/revert transition), success
(appearance), and the entire **Card→Focus** transition (collapse/expand/dim
sequence). UI owns their static look; UX owns their timing/easing/sequence.
