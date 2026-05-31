# Gate Status — CRM System

> ONBOARDING NOTE (2026-05-30): Added when the project was brought under the
> platform collaboration model. Statuses below reflect PROJECT_CONTEXT.md as of
> 2026-05-30 (G5 passed, G6 current). Confirm and maintain going forward.

This file is the single synchronization point between Claude (planning + audit)
and Codex (execution). Update it on every gate change, handoff, and kickback.
See `../../../../company/collaboration-model.md`.

Status values: `Gate Draft` / `Gate Review` / `Gate Blocked` / `Gate Passed`.

| Gate | Transition | Owner | Platform | Status | Reviewer sign-off | Date | Current blocker |
|---|---|---|---|---|---|---|---|
| G1 | Idea -> Requirement Discussion | Product Manager | Claude | Gate Passed | | | |
| G2 | Requirement Discussion -> PRD | Product Manager | Claude | Gate Passed | | | |
| G3 | PRD -> Acceptance Matrix | Product Manager | Claude | Gate Passed | | | |
| G4a | Acceptance Matrix -> Business Design | Business Analyst | Claude | Gate Passed | | | |
| G4b | Business Design -> UX Design | UX Designer | Claude | Gate Passed | | | |
| G4c | UX Design -> UI Design | UI Designer | Claude | Gate Passed | | | |
| G4d | UI Design -> Security Design (Design Closure) | Security Compliance | Claude | Gate Passed | PM, BA, UX, UI (retroactive closure) | 2026-05-30 | |
| G5 | Design Closure -> Architecture Design | Architecture | Claude | Gate Passed | all required reviewers (incl. Infrastructure Ops) | 2026-05-30 | |
| G6 | Architecture Design -> MDA Modeling | Domain Modeling + Architecture | Claude | Gate Draft | PM, BA, UX, UI, Security, QA Test Design | | Not started — next gate; MDA package not yet created |
| G7 | MDA + Test Model -> Task Planning | Domain Modeling + QA Test Design | Claude | Gate Draft | | | |
| G8 | Task Planning -> Implementation **[HANDOFF: Claude -> Codex]** | Task Planner | Claude -> Codex | Gate Draft | | | |
| G9 | Implementation -> QA | Frontend / Backend Engineer | Codex | Gate Draft | | | |
| G10 | QA -> Integration | QA Execution | Codex | Gate Draft | | | |
| G11 | Integration -> Audit **[RETURN: Codex -> Claude]** | Integration Owner | Codex -> Claude | Gate Draft | | | |
| G12 | Audit -> Release/Rework | Audit | Claude | Gate Draft | | | release blockers carried: off-server backup+restore, HTTPS/TLS, security-group, monitoring evidence |

## Handoff Log

| Date | From platform | To platform | Gate | Note |
|---|---|---|---|---|
| | | | | |

## Notes

- G8 may not pass until the self-contained G8 execution handoff package exists on
  disk (see `company/collaboration-model.md`).
- On a kickback, set the affected gate back to `Gate Blocked`, record it here and
  in `planning/blockers.md`, and return to Claude.
- G12 audit is mandatory and performed on Claude before any release decision.
- Carried-forward release blockers (not current gate blockers): encrypted
  off-server backup copy + restore rehearsal, HTTPS/TLS endpoint, security-group,
  and monitoring evidence.
- G1–G4c are the retained design baseline: produced/passed under the pre-split
  combined G4, retained through the 2026-05-29 reset, supplemented for
  service-boundary governance, and re-validated by the G4d retroactive closure
  check and the pre-G6 design audit.
- Evidence pointers:
  - G4d closure: `archive/reviews/g4-design-closure/g4d-design-closure-decision-2026-05-30.md`
  - G5 final decision: `archive/reviews/g5-architecture/g5-architecture-final-decision-2026-05-30.md`
  - Pre-G6 design audit: `archive/reviews/g5-to-g6-handoff/pre-g6-design-audit-2026-05-30.md`
- Project-layer strengthening (`process/process-gap-register.md`):
  Infrastructure Ops is a required reviewer at G5/G8/G11/G12; every CRM
  cross-capability flow must name a `Primary Flow Owner Agent`.
- OPEN (decide at G8): where the G8 task/delivery artifacts live —
  `tasks.md`, `task-dependencies.md`, `delivery-plan.md`,
  `acceptance-task-map.md`, `blockers.md`. On 2026-05-31, `planning/` and
  `process/` were moved out of `docs/` (design-only). `STANDARD-APPLICATION-REVIEW.md`
  currently lists these task files under `planning/`. The discarded 2026-05-29
  G8 cycle had placed them in a top-level `delivery/` folder (referenced only in
  `archive/reviews/g8-task-planning/`; that folder no longer exists). Decide then
  whether G8 rebuilds them under `planning/` or a new `delivery/`, keeping
  governance (gate-status, blockers) separate from execution artifacts.
