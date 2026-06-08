# UI/UX A5 Text-token Exception Evidence

Date: 2026-06-07
Decision: DEC-UIUX-A5-001
Blocker: `planning/blockers.md` BLK-UIUX-G9-001
Scope: UIUX-001 token checkpoint before downstream UIUX-002..014 work

## Result

Token checkpoint ready for Claude token re-audit.

Codex added four strictly text-only contrast tokens to
`docs/ux-ui/design-system.md` Appendix A:

```css
--success-ink: #11803A;
--warning-ink: #A55A05;
--danger-ink: #CC2121;
--purple-ink: #7D4CE4;
```

No frontend implementation, backend, API, data model, enum/role comparison, or
business logic was changed in this checkpoint.

## Derivation Method

For each source color, Codex kept the locked token's HSL hue and saturation,
lowered only HSL lightness, rounded to 8-bit sRGB hex, and selected the first
distinct rounded color that reaches WCAG AA normal text contrast (4.5:1) on all
intended backgrounds:

- `--card` `#FFFFFF`
- `--section` `#F6F7FD`
- corresponding soft surface

The previous lighter rounded hex is shown below to prove the selected value is
the minimal 8-bit darkening under this method.

## Minimality Table

| Token | Source | Selected hex | Previous lighter rounded hex | Previous min contrast | Selected min contrast |
|---|---|---|---|---:|---:|
| `--success-ink` | `--success` `#16A34A` | `#11803A` | `#11813A` | 4.47 | 4.53 |
| `--warning-ink` | `--warning` `#D97706` | `#A55A05` | `#A55B05` | 4.49 | 4.53 |
| `--danger-ink` | `--danger` `#DC2626` | `#CC2121` | `#CD2121` | 4.48 | 4.52 |
| `--purple-ink` | `--purple` `#B79CF0` | `#7D4CE4` | `#7D4DE4` | 4.49 | 4.52 |

## WCAG Contrast Table

| Text token | Card `#FFFFFF` | Section `#F6F7FD` | Corresponding soft surface | Result |
|---|---:|---:|---:|---|
| `--success-ink` `#11803A` | 5.03 | 4.71 | 4.53 on `--mint-soft` `#E5F7F0` | Pass AA normal text |
| `--warning-ink` `#A55A05` | 5.17 | 4.83 | 4.53 on `--peach-soft` `#FDEDE5` | Pass AA normal text |
| `--danger-ink` `#CC2121` | 5.52 | 5.16 | 4.52 on danger soft `#FEE2E2` | Pass AA normal text |
| `--purple-ink` `#7D4CE4` | 5.22 | 4.88 | 4.52 on `--purple-soft` `#F2ECFD` | Pass AA normal text |
| readable secondary `--muted` `#475569` | 7.58 | 7.09 | N/A | Pass AA normal text |

## Token Diff Evidence

Additions only:

| Added token | Hex | Allowed use |
|---|---|---|
| `--success-ink` | `#11803A` | Readable success text only |
| `--warning-ink` | `#A55A05` | Readable warning text only |
| `--danger-ink` | `#CC2121` | Readable danger text only |
| `--purple-ink` | `#7D4CE4` | Readable purple/accent text only |

Existing locked values remain byte-for-byte unchanged:

| Locked token / literal | Hex |
|---|---|
| `--subtle` | `#94A3B8` |
| `--success` | `#16A34A` |
| `--warning` | `#D97706` |
| `--danger` | `#DC2626` |
| `--purple` | `#B79CF0` |
| `--mint-soft` | `#E5F7F0` |
| `--peach-soft` | `#FDEDE5` |
| danger soft literal | `#FEE2E2` |
| `--purple-soft` | `#F2ECFD` |
| `--card` | `#FFFFFF` |
| `--section` | `#F6F7FD` |
| `--primary` | `#2563EB` |
| `--primary-hover` | `#1D4ED8` |
| `--border` | `#EDF0F6` |

## Usage Update

- Readable colored status/accent text uses the matching `*-ink` token.
- Readable secondary, meta, caption, timestamp, and placeholder text uses
  `--muted` `#475569`.
- `--subtle` `#94A3B8` is decorative / non-text / disabled-only, not readable
  body text.
- Original solid semantic and support colors remain locked for swatches, icons,
  graph strokes, accent dots, and decorative marks.
- Backgrounds, soft tints, solid fills, brand colors, borders, and icon colors
  are not recolored by this decision.
