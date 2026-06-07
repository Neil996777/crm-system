# CRM System — Claude Entry (Requirements/Standards + Audit)

This project inherits workspace governance. On Claude, this project is in the
**yardstick** phase (requirements + acceptance + binding constraints, G1–G3) or
the **independent audit** phase (the G8 design/planning handoff audit + G12).

> Aligned to the collaboration-model redirect on 2026-06-06. Authority/index:
> `../../company/policy-changes/2026-06-06-collaboration-model-redirect.md`,
> `../../company/collaboration-model.md`, project sync notice
> `COLLAB-MODEL-REDIRECT-SYNC-NOTICE.md`; mirrors `../../templates/project-CLAUDE.md`.

## Required Workspace Context

Before working, read:

- `../../CLAUDE.md`
- `../../company/collaboration-model.md`
- `../../company/operating-model.md`
- `../../standards/acceptance-matrix-standard.md`
- `../../standards/status-and-priority-standard.md`
- `../../workflows/project-initialization.md`
- `../../workflows/software-delivery.md`

## Project Context

- `PROJECT_CONTEXT.md`
- `planning/gate-status.md` (current gate and responsible platform)

## Platform Rules

- Claude owns the yardstick (requirements + acceptance + binding constraints,
  G1–G3) and independent audit (the G8 design/planning handoff audit + G12).
- Codex produces design/architecture/modeling/copy/task planning (G2–G8) and does
  implementation/QA/integration (G9–G11). Do not produce or implement those on
  Claude.
- At the G8 handoff, Codex produces the self-contained handoff package on disk and
  Claude audits it against the yardstick; implementation may not start until that
  audit passes (no implementation code before G8).
- The (A) requirements / (B) acceptance / (C) binding constraints (no-downgrade,
  P0/P1 status rules, security invariants, locked design system, real enums) are
  red lines and are never moved to Codex.
- Follow the no-downgrade rule; P0/P1 items are only `Done`, `Blocked`, or
  `Formal Scope Change by User`.
- Project-specific rules may only strengthen workspace rules.

## Project-Specific Note

Architecture, MDA, task planning, and implementation artifacts were discarded by
user direction on 2026-05-29; G5 architecture was redesigned and passed 2026-05-30.
The base product passed the G12 audit and went live on 2026-06-05 (zh-CN phase 1+2
localized and live). Two follow-on changes are in flight: **UI/UX completion**
(`delivery/uiux-completion-charter.md`) and **CI/CD migration**
(`delivery/cicd-migration-plan.md`). Under the new model their design / architecture
/ task production (G2–G8) is Codex's; Claude defines requirements / acceptance /
binding constraints (e.g. `docs/ux-ui/requirements/`, the locked design system,
the real opportunity-stage enum) and performs the G8 handoff audit + G12. Current
gate/responsible-platform state is tracked in `planning/gate-status.md`.

## Editable Scope

Default editable scope: `projects/crm-system/`. Do not edit workspace-level files
unless the user explicitly requests it.
