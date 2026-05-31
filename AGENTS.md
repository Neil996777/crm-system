# CRM System — Codex Entry (Execution)

This project inherits workspace-level agent rules and standards. On Codex, this
project is in the execution phase (G9–G11). Planning (G1–G8) and audit (G12) run
on Claude — see `CLAUDE.md` and `../../company/collaboration-model.md`.

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

## Editable Scope

Default editable scope:

- `projects/crm-system/`

Do not edit workspace-level files unless the user explicitly requests it.

## Project Rules

- Follow the workspace no-downgrade rule.
- P0/P1 items cannot be downgraded, deleted, merged away, weakened, or accepted as partial work.
- Do not write implementation code before Gate G8 passes.
- Do not use mock, stub, TODO, static UI, or non-persistent behavior to satisfy core CRM paths.
- Project-specific rules may only strengthen workspace-level rules.
- Current architecture, MDA, task planning, and implementation artifacts were discarded by user direction on 2026-05-29. New architecture must be designed from the retained product, business, UX/UI, and security inputs before any new MDA, task planning, or implementation work.
