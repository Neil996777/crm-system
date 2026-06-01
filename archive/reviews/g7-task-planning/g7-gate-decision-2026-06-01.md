# G7 (Task Planning) — Gate Decision

## Document Control

- Project: CRM System
- Gate: G7 — MDA + Test Model → Task Planning
- Owner: Domain Modeling + QA Test Design (Claude, planning platform)
- Reviewers: Product Manager, Architecture, Security Compliance, QA Test Design,
  Infrastructure Ops, Task Planner (handoff-readiness)
- Date: 2026-06-01
- Decision: **Gate Passed**
- Archive note: Gate evidence only. Not design authority.

## Scope of the Gate

The G7 task-planning package under `delivery/`:

- `delivery/tasks.md` — 40 tasks (TASK-001..040), each with the full 17-field
  schema (objective, capability, ACC, reference docs, concrete file changes, owner
  agent, prerequisites, DoD, acceptance method, automated tests, manual verification,
  MDA traceability, TDD, no-downgrade items, blocker).
- `delivery/acceptance-task-map.md` — ACC-001..023 → tasks (23/23).
- `delivery/task-dependencies.md` — acyclic dependency DAG / build order.
- `delivery/delivery-plan.md` — committed stack, repo layout, conventions, phasing.

It is an acceptance-driven, Codex-executable plan built to the post-2026-06-01
model (DEC-017..022): each task = one user capability + one ACC + real code changes
+ tests + a reproducible verification; zero TBD; service boundaries, no-downgrade,
and TDD encoded per task; carried release blockers staged as explicit G11/G12
evidence tasks.

## How produced and verified

Authored by Domain Modeling + QA Test Design (Task Planner role), then multi-agent
audited (author ≠ reviewer) across four lenses (coverage/traceability,
executability/DAG, no-downgrade/service-boundary/fidelity, TDD/test-mapping), a fix
round, and a consistency re-audit — all PASS; 66/66 delivery TEST families resolve
1:1 to `modeling/test-model.md`.

## Reviewer Sign-Off

| Reviewer | Verdict | Notes |
|---|---|---|
| Product Manager | SIGN-OFF | ACC 23/23; no-downgrade; DEC-017..022 reflected; no out-of-scope; each task acceptance-checkable. |
| Architecture | SIGN-OFF | Service decomposition/data-ownership/S2S boundaries respected; stack matches ADR-ARCH-001 + DEC-021/022; contracts-before-consumers build order. |
| Security Compliance | SIGN-OFF | Authz/abuse/audit/retention controls tasked server-side with backend negative tests; PM-*/ABUSE-* covered. |
| QA Test Design | SIGN-OFF | Tests 1:1 to test-model; P0 positive+negative; TDD fail-first; EDGE 37/37 + ABUSE 22/22 homed; retired tests excluded. |
| Infrastructure Ops | SIGN-OFF (concern resolved) | Runtime/persistence/release-evidence tasked and staged for G11/G12. Concern: operator-access provisioning was cited-not-tasked → resolved (TASK-039 now tasks the least-privilege deploy/ops user + records the pre-G8 Security Compliance review of SSH/sudo). |
| Task Planner (handoff) | SIGN-OFF (concern resolved) | Self-contained, zero-TBD, acyclic DAG, 17/17 fields. Concern: `infrastructure-ops` owner label → resolved (recognized project role per `AGENTS.md` + `process/process-gap-register.md`; added to `delivery-plan.md` owner enumeration). |

All concerns were non-blocking and have been resolved in the plan. Remaining
non-blocking hygiene notes (explicit ABUSE-0NN tokens on some reference lines;
range-citation expansion) are deferred to G8/G9 and do not affect gate passage.

## Pre-G8 entry condition (recorded, not a G7 blocker)

Per `deployment-notes.md` "Operator Access", Security Compliance must review the
operator-access design (SSH access, key ownership, sudo boundary) before G8
implementation tasks are approved. Recorded on TASK-039 (field 17).

## Decision

**Gate Passed.** The G7 task-planning package is acceptance-driven, fully traceable
(ACC-001..023), Codex-executable (zero TBD, concrete file changes, acyclic DAG),
audited (author ≠ reviewer), and signed off by all six reviewer roles. Proceed to
**G8 (Task Planning → Implementation — Claude→Codex handoff)**: assemble the
self-contained execution handoff package; Claude writes no implementation code.

Carried release blockers (off-server backup + restore rehearsal, HTTPS/TLS,
security-group, monitoring) remain release-time evidence for G11/G12.
