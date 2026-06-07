# Governance Sync Notice — Collaboration-Model Redirect

Issued: 2026-06-06 (company → project). Authority:
`../../company/policy-changes/2026-06-06-collaboration-model-redirect.md`,
`../../company/collaboration-model.md`.

The company collaboration model changed. **This project must align its own entry
files** (project agents do this; the governance change does not edit them).

## What changed
- **Claude** = yardstick + audit only: requirements + acceptance matrix (G1–G3),
  binding constraints (no-downgrade, P0/P1 status rules, security invariants,
  locked design system, real enums), Security Compliance, QA Test Design, Audit,
  Governance Auditor. Plus a **new G8 design/planning handoff audit**.
- **Codex** = design/production (G2–G8) + execution (G9–G11): Business Analyst,
  UX, UI, Domain Modeling, Architecture, Task Planner, FE, BE, QA Execution,
  Integration Owner, Infrastructure Ops.
- Unchanged: G1–G12 gate table, no-downgrade, release/evidence rules,
  "no implementation code before G8 passes", all prior G12 results.

## What this project must update (project scope)
- `CLAUDE.md`: the "Claude owns G1–G8 (planning)" wording → new yardstick + audit
  model (mirror `templates/project-CLAUDE.md`).
- `AGENTS.md` (project): mirror `templates/project-AGENTS.md` (Codex owns design
  G2–G8 + execution G9–G11).
- `planning/gate-status.md`: the "responsible platform" column for any future/open
  gate should reflect Codex for G2–G8 production and Claude for the G8 handoff
  audit + G12. Already-passed gates are historical and not rewritten.

No business/implementation code changes from this notice.

## Alignment Status

ALIGNED 2026-06-06. Project entry files updated to the new model (index of changes):
- `CLAUDE.md` — Claude = yardstick (G1–G3) + audit (G8 handoff audit + G12); stale
  "Claude owns G1–G8 (planning)" and "current gate G6" removed; current state
  recorded (G12 passed, live 2026-06-05, two follow-on changes in flight).
- `AGENTS.md` — Codex = design/production (G2–G8) + execution (G9–G11); added
  Production & Execution Rules (produce from on-disk G1 yardstick, do not change
  A/B/C, produce G8 handoff package, Kickback Protocol).
- `planning/gate-status.md` — added the 2026-06-06 redirect note (going-forward
  responsible-platform = new model; already-passed rows kept as historical);
  sync-point description updated; corrected a stale relative path.

Authority/index: `../../company/policy-changes/2026-06-06-collaboration-model-redirect.md`,
`../../company/collaboration-model.md`, this notice.
