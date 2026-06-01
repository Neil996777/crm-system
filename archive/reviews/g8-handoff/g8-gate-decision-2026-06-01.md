# G8 (Task Planning → Implementation) — Gate Decision [HANDOFF: Claude → Codex]

## Document Control

- Project: CRM System
- Gate: G8 — Task Planning → Implementation **[HANDOFF: Claude → Codex]**
- Owner: Task Planner (Claude, planning platform)
- Reviewers: Task Planner (handoff completeness), Infrastructure Ops (required),
  Security Compliance, independent Audit
- Date: 2026-06-01
- Decision: **Gate Passed — handed off to Codex**
- Archive note: Gate evidence only. Not design authority.

## Scope of the Gate

G8 is the planning→execution handoff. The deliverable is a self-contained execution
handoff package on disk; Claude writes no implementation code at or after G8.

Handoff package (all on disk):
- `delivery/G8-handoff.md` — single entry document: platform handoff, read order,
  committed stack/layout, per-task execution contract, build-to-current-model rules,
  G9/G10/G11 expectations, standing conditions/blockers, governance.
- `delivery/delivery-plan.md`, `delivery/tasks.md` (40 tasks), `delivery/task-dependencies.md`, `delivery/acceptance-task-map.md`.
- Authority: `modeling/*`, `docs/*`, `planning/*`, `process/*`; gate evidence under `archive/reviews/*`.

## Pre-G8 condition satisfied

The G7-recorded pre-G8 condition — Security Compliance review of the operator-access
design — is satisfied: **Approved with conditions**
(`archive/reviews/g8-handoff/security-operator-access-review-2026-06-01.md`); its four
conditions are carried into G9/G11 evidence via TASK-039.

## Reviewer Sign-Off

| Reviewer | Verdict | Notes |
|---|---|---|
| Task Planner (handoff) | SIGN-OFF | Package self-contained, zero-TBD, 40 tasks/17 fields, acyclic DAG, ACC 23/23; no implementation code in repo; re-planning forbidden; no open product decisions. |
| Infrastructure Ops | SIGN-OFF | Deployment target/ownership carried; operator-access tasked + Security review recorded; all carried release blockers staged as G11/G12 evidence under owning tasks. |
| Security Compliance | SIGN-OFF-WITH-CONCERNS | Operator-access pre-condition met; full server-side authz / S2S / append-only audit / authenticated-actor / classification-retention contract carried with backend negative tests + no-downgrade; HTTPS/TLS/cookie/exposure/backup correctly staged as G11/G12 evidence. |
| Audit (independent) | SIGN-OFF-WITH-CONCERNS | Platform discipline (code-free, stop/return-at-G12), self-containment, traceability spine (ACC 23/23 → 40 tasks → test-model, retired tests excluded), no-downgrade, and consistency all hold; minor doc-hygiene only. |

All concerns are non-blocking (release-time evidence correctly deferred to G11/G12;
minor path/hygiene notes addressed). No REJECT.

## Decision

**Gate Passed — handed off to Codex.** The G8 execution handoff package is complete,
self-contained, governance-compliant, and code-free. Execution proceeds on the Codex
platform:
- G9 Implementation (Frontend/Backend Engineer), G10 QA Execution, G11 Integration —
  per `delivery/G8-handoff.md` and the delivery plan.
- Claude performs no implementation. Claude resumes for the **independent G12 audit**
  after G11 passes, before any release decision.

Carried into execution (G9/G11/G12 evidence, not gate blockers): operator-access
conditions; encrypted off-server backup + restore rehearsal; HTTPS/TLS; security
group; monitoring/health; restart-survival persistence. Real-world dependency: a
valid HTTPS endpoint / TLS certificate (and domain, if used) for production release.
