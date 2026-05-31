# CRM System — Claude Entry (Planning + Audit)

This project inherits workspace governance. On Claude, this project is in the
**planning** phase (G1–G8) or the **independent audit** phase (G12). Per project
state, the current gate is G6 (MDA Modeling).

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

- Claude owns G1–G8 (planning) and G12 (independent audit).
- At G8, stop and produce a self-contained execution handoff package on disk;
  do not write implementation code.
- G9–G11 (implementation, QA execution, integration) belong to Codex.
- Follow the no-downgrade rule; P0/P1 items are only `Done`, `Blocked`, or
  `Formal Scope Change by User`.
- Project-specific rules may only strengthen workspace rules.

## Project-Specific Note

Architecture, MDA, task planning, and implementation artifacts were discarded by
user direction on 2026-05-29. G5 architecture was redesigned and passed
2026-05-30. New MDA (G6) must trace the accepted architecture, not redefine it.
Implementation remains blocked until G8.

## Editable Scope

Default editable scope: `projects/crm-system/`. Do not edit workspace-level files
unless the user explicitly requests it.
