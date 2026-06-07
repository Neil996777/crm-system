# UI/UX Contrast Audit — G9 Kickback Evidence

Date: 2026-06-07
Scope: Pre-implementation check for UIUX-001 against
`docs/ux-ui/requirements/uiux-implementation.requirements.md` A5 and C6.

## Result

**Blocked.** The locked design-system palette contains common text/background
pairs that do not meet WCAG AA for normal text. The 2026-06-07 yardstick
supplement says this A5(AA) vs C6(no recolor) conflict must kick back to Claude;
Codex must not silently recolor or continue with known non-AA text pairings.

Blocker: `planning/blockers.md` BLK-UIUX-G9-001.
Gate: `planning/gate-status.md` UI/UX G9 = Gate Blocked.

## Method

Contrast was calculated using the WCAG relative luminance formula:

```text
(L1 + 0.05) / (L2 + 0.05)
```

AA thresholds:

- Normal text: 4.5:1
- Large text: 3:1

## Measured Pairs

| Pair | Foreground | Background | Contrast | Result |
|---|---|---|---:|---|
| text on card | `#0F172A` | `#FFFFFF` | 17.85 | Pass normal AA |
| muted on card | `#475569` | `#FFFFFF` | 7.58 | Pass normal AA |
| subtle on card | `#94A3B8` | `#FFFFFF` | 2.56 | **Fail normal/large AA** |
| primary on card | `#2563EB` | `#FFFFFF` | 5.17 | Pass normal AA |
| primary on tint | `#2563EB` | `#EAF1FF` | 4.56 | Pass normal AA |
| muted on section | `#475569` | `#F6F7FD` | 7.09 | Pass normal AA |
| subtle on section | `#94A3B8` | `#F6F7FD` | 2.40 | **Fail normal/large AA** |
| success on mint-soft | `#16A34A` | `#E5F7F0` | 2.97 | **Fail normal/large AA** |
| warning on peach-soft | `#D97706` | `#FDEDE5` | 2.79 | **Fail normal/large AA** |
| danger on danger-soft | `#DC2626` | `#FEE2E2` | 3.95 | **Fail normal AA** |
| white on primary | `#FFFFFF` | `#2563EB` | 5.17 | Pass normal AA |
| white on primary-hover | `#FFFFFF` | `#1D4ED8` | 6.70 | Pass normal AA |
| purple on purple-soft | `#B79CF0` | `#F2ECFD` | 2.02 | **Fail normal/large AA** |

## Required Claude Decision

One of the following is needed before G9 implementation can continue:

- Approve minimal contrast-only token exceptions and record them as a formal
  design-system/yardstick decision.
- Clarify which failing token pairings are decorative/non-text only and provide
  an AA-safe text pairing rule that still satisfies C6.
- Revise the locked design-system or A5/C6 yardstick to remove the contradiction.

Until then, UIUX-001 and all downstream UIUX tasks remain blocked.

