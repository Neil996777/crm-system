# Blockers — CRM System

This register tracks open blockers and pre-gate decisions surfaced during planning.
It is a companion to `planning/gate-status.md` (the gate sync point). On a kickback,
set the affected gate to `Gate Blocked` in gate-status, record it here, and return
to Claude. See `../../../company/collaboration-model.md`.

Status values: `Open` / `In Review` / `Resolved` / `Formal Scope Change by User`.
No-downgrade rule applies: a blocker touching a P0/P1 acceptance may only move to
`Resolved` (with the deciding source recorded) or `Formal Scope Change by User`.

## Open Modeling Decisions — must close before G8

These were surfaced by the G6 MDA modeling (CIM/PIM) and its independent multi-agent
audits. Each is a genuine upstream-source gap that MDA correctly declined to invent
(per the no-invention discipline).

> **Update 2026-06-01: ALL THREE RESOLVED by a Formal Scope Change by User**
> (decision-log.md DEC-017..020). BLK-001 → DEC-020 (Opportunity Status field
> removed); BLK-002 → DEC-017/019 (Won = contract Signed; payment decoupled, so no
> multi-plan "fully paid" Won-gate); BLK-003 → DEC-018 (one quote per opportunity;
> no second quote to accept). Cascade applied across baseline → architecture → MDA
> and re-audited. See Resolution Log below.

| ID | Blocker | Owner | Blocks | Touches | Status | Opened |
|---|---|---|---|---|---|---|
| BLK-001 | Opportunity Status enumerated value set distinct from Pipeline Stage is not enumerated in the accepted sources (CIM-016 commits Status as a persisted dimension; PRD-007/ACC-007 require status as a persisted field but enumerate no values). Decide whether a distinct Status enumeration exists or whether Status is realized through Stage. | Product Manager + Business Analyst | G8 handoff (and G7 test design determinism) | ACC-007 (P0) | Resolved — Formal Scope Change by User (DEC-020) | 2026-06-01 |
| BLK-002 | Multi-plan contract "fully paid" aggregation is undefined. Overpayment ceiling is fixed at contract level (EDGE-019) and Won requires contract-level full payment (BR-008, DEC-012), but the accepted sources do not define how multiple Payment Plans under one contract aggregate into "contract fully paid." Confirm the aggregation rule. | Product Manager + Business Analyst | G8 handoff (and G7 deterministic Won/overpayment test design) | ACC-011, ACC-013 (P0) | Resolved — Formal Scope Change by User (DEC-017/019) | 2026-06-01 |
| BLK-003 | Second-quote-accept observable outcome is unspecified. The "at most one Accepted quote per opportunity" invariant is fixed, but EDGE-012 does not specify the observable business result of accepting a second quote (reject the second accept vs auto-demote the current Accepted quote). A product decision is required before deterministic test design. | Product Manager | G7 test design / G8 handoff | ACC-009 (P0) | Resolved — Formal Scope Change by User (DEC-018) | 2026-06-01 |

### Source pointers
- BLK-001: `modeling/CIM.md` CIM-016; `modeling/PIM.md` PIM-SM-003, PIM-OPEN-003; `docs/product/acceptance-matrix.md` ACC-007; PRD-007.
- BLK-002: `modeling/PIM.md` PIM-SM-006, PIM-INV-007/023/025, PIM-OPEN-005; `docs/business/business-rules.md` BR-008; `docs/product/decision-log.md` DEC-012; `docs/business/edge-cases.md` EDGE-019.
- BLK-003: `modeling/CIM.md` CIM-PROC-008 Open/Blocked; `modeling/PIM.md` PIM-SM-004, PIM-OPEN-001; `docs/business/edge-cases.md` EDGE-012; ACC-009.

## Architecture / PSM-deferred (not PM/BA blockers; tracked for G7/PSM)

These are correctly deferred to PSM/Architecture at PIM altitude and are recorded for
the PSM artifact and G7 test design; they are not pre-G8 PM/BA decisions.

| ID | Item | Owner | Resolves at | Refs |
|---|---|---|---|---|
| BLK-A01 | Overdue-evaluation trigger (on-read vs scheduled) — needed for deterministic overdue test design; modeled at PIM as a Business-Date guard, mechanism deferred. | Architecture (in PSM) | PSM / G7 | PIM-OPEN-002, CIM-034, BR-021 |

## Open G11 Release Evidence Blockers

These are active execution blockers surfaced during Codex G9-G11 execution. Per the
no-invention discipline, Codex may not invent production endpoint, TLS, security-group,
monitoring, or runtime evidence.

| ID | Blocker | Owner | Blocks | Touches | Status | Opened |
|---|---|---|---|---|---|---|
| BLK-G11-001 | TASK-039 cannot be completed because production release requires a valid HTTPS endpoint with TLS certificate, but the deployment source says the domain is "Not specified yet" and absence of a valid HTTPS endpoint blocks production release. G11 also requires real runtime evidence on `srv-volcengine-sh-01` for endpoint, TLS status, security-group inbound rules, opened ports, health URL, deployment timestamp, operator, and smoke result. | Infrastructure Ops + Integration Owner | G11 handoff to G12; TASK-039 | ACC-017 (P0), ARCH-ACC-008/013/014/015, TEST-DEPLOY-SMOKE-001/002 | Open | 2026-06-03 |

## Carried-forward Release Blockers (not gate blockers)

Recorded in `planning/gate-status.md` G12 row; repeated here for visibility. These are
release-time evidence items, not modeling or current-gate blockers:
encrypted off-server backup copy + restore rehearsal, HTTPS/TLS endpoint, security-group,
and monitoring evidence. Refs: OQ-001, RISK-002, `docs/architecture/deployment-notes.md`.

## Resolution Log

| Date | ID | Resolution | Deciding source |
|---|---|---|---|
| 2026-06-01 | BLK-001 | Opportunity `Status` field removed; Pipeline Stage is the sole lifecycle dimension. The "what are Status's values" gap is dissolved, not answered. | DEC-020 (Formal Scope Change by User) |
| 2026-06-01 | BLK-002 | Won decoupled from full payment (Won = contract Signed); multi-plan "fully paid" is no longer a Won precondition, so the aggregation gate is moot. Overpayment ceiling remains contract-level. | DEC-017 + DEC-019 (Formal Scope Change by User) |
| 2026-06-01 | BLK-003 | Each opportunity has exactly one quote; there is no second quote to accept, so the reject-vs-auto-demote question no longer exists. | DEC-018 (Formal Scope Change by User) |
| 2026-06-01 | BLK-A01 | Overdue-evaluation trigger resolved at PSM as on-read evaluation against the supplied `businessDate` (Asia/Shanghai); deterministic for G7 test design. | PSM "Resolved Mechanisms"; api-spec.md Reminder Query; FLOW-005 |
