# CRM System — Codex Entry (Design/Production + Execution)

This project inherits workspace-level agent rules and standards. On Codex, this
project covers the **design/production** phase (G2–G8: business analysis, UX, UI,
domain modeling, architecture, copy/IA, task planning) and the **execution** phase
(G9–G11: implementation, QA execution, integration). Claude owns the yardstick
(requirements + acceptance + binding constraints, G1–G3) and independent audit
(the G8 design/planning handoff audit + G12) — see `CLAUDE.md` and
`../../company/collaboration-model.md`.

> Aligned to the collaboration-model redirect on 2026-06-06. Authority/index:
> `../../company/policy-changes/2026-06-06-collaboration-model-redirect.md`,
> `../../company/collaboration-model.md`, project sync notice
> `COLLAB-MODEL-REDIRECT-SYNC-NOTICE.md`; mirrors `../../templates/project-AGENTS.md`.

## Required Workspace Context

Before working on this project, read:

- `../../AGENTS.md`
- `../../company/collaboration-model.md`
- `../../company/operating-model.md`
- `../../standards/acceptance-matrix-standard.md`
- `../../standards/status-and-priority-standard.md`
- `../../workflows/project-initialization.md`
- `../../workflows/software-delivery.md`

## Project Context

Before doing project work, read:

- `PROJECT_CONTEXT.md`
- `planning/gate-status.md` (current gate and responsible platform)

## Available Workspace Agents

- `../../agents/product-manager.md`
- `../../agents/business-analyst.md`
- `../../agents/ux-designer.md`
- `../../agents/ui-designer.md`
- `../../agents/domain-modeling.md`
- `../../agents/architecture.md`
- `../../agents/security-compliance.md`
- `../../agents/task-planner.md`
- `../../agents/qa-test-design.md`
- `../../agents/frontend-engineer.md`
- `../../agents/backend-engineer.md`
- `../../agents/qa-execution.md`
- `../../agents/integration-owner.md`
- `../../agents/infrastructure-ops.md`
- `../../agents/audit.md`

## Production & Execution Rules

- Produce design/architecture/modeling/copy/task planning (G2–G8) from the on-disk
  G1 yardstick (requirements + acceptance + binding constraints); do not rely on
  conversation context.
- Do not change the (A) requirements, (B) acceptance, or (C) binding constraints —
  those stay with Claude. Honor no-downgrade, the P0/P1 status rules, security
  invariants (IDOR, durable audit, optimistic concurrency, idempotency), the locked
  design system, and real enums (e.g. the six opportunity stages).
- Produce the self-contained G8 handoff package on disk; implementation (G9) may
  not start until Claude's G8 design/planning handoff audit passes.
- If the yardstick is wrong or insufficient, follow the Kickback Protocol in
  `../../company/collaboration-model.md`: record in `planning/blockers.md`, set the
  gate to `Gate Blocked` in `planning/gate-status.md`, and return to Claude.
- Do not use mock, stub, TODO, static UI, or non-persistent behavior for core P0/P1
  paths.
- After execution (G11), update `planning/gate-status.md` and hand back to Claude
  for the independent G12 audit.
- Project-specific rules may only strengthen workspace-level rules.

## Editable Scope

Default editable scope:

- `projects/crm-system/`

Do not edit workspace-level files unless the user explicitly requests it.
