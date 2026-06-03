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
| BLK-G11-001 | TASK-039 cannot be completed because production release requires a valid HTTPS endpoint with TLS certificate, but the deployment source says the domain is "Not specified yet" and absence of a valid HTTPS endpoint blocks production release. G11 also requires real runtime evidence on `srv-volcengine-sh-01` for endpoint, TLS status, security-group inbound rules, opened ports, health URL, deployment timestamp, operator, and smoke result. | Infrastructure Ops + Integration Owner | G11 handoff to G12; TASK-039 | ACC-017 (P0), ARCH-ACC-008/013/014/015, TEST-DEPLOY-SMOKE-001/002 | Resolved — `https://118.196.44.193` approved and provisioned with a Let's Encrypt IP certificate on 2026-06-03 | 2026-06-03 |
| BLK-G11-002 | TASK-039 has server-side HTTPS/TLS, health, redirect, service health, monitoring threshold, renewal, operator evidence, and Volcengine security-group API evidence. The API evidence binds instance `i-yemoz0an7kk36d2c9bp6` to security group `sg-1pm4k7f37z8xs643rg0fvk85e` through ENI `eni-13e8tbocd8f0g79iu5jer8idt` and confirms CRM gateway `8080` and PostgreSQL `5432` are not allowed from `0.0.0.0/0`. User approved releasing previous deployments; Codex stopped and removed host-network Hermes, so host-level `8642` no longer listens and CRM smoke still passes. The user removed old/non-CRM Volcengine security-group rules for TCP `8088`, TCP `8443`, and TCP `3389`; API post-cleanup verification confirms only public TCP `22`, `80`, and `443` remain. | Infrastructure Ops + Security Compliance | G11 handoff to G12; TASK-039 | ACC-017 (P0), ARCH-ACC-013, ARCH-ACC-014, TEST-DEPLOY-SMOKE-001/002 | Resolved | 2026-06-03 |

## Carried-forward Release Blockers (not gate blockers)

Recorded in `planning/gate-status.md` G12 row; repeated here for visibility. These were
release-time evidence items, not modeling blockers: encrypted off-server backup copy +
restore rehearsal, HTTPS/TLS endpoint, security-group, and monitoring evidence. As of
2026-06-03, Codex has recorded runtime evidence for these items under TASK-039 and
TASK-040; G12 still performs independent audit before any release decision. Refs:
OQ-001, RISK-002, `docs/architecture/deployment-notes.md`.

## Resolution Log

| Date | ID | Resolution | Deciding source |
|---|---|---|---|
| 2026-06-01 | BLK-001 | Opportunity `Status` field removed; Pipeline Stage is the sole lifecycle dimension. The "what are Status's values" gap is dissolved, not answered. | DEC-020 (Formal Scope Change by User) |
| 2026-06-01 | BLK-002 | Won decoupled from full payment (Won = contract Signed); multi-plan "fully paid" is no longer a Won precondition, so the aggregation gate is moot. Overpayment ceiling remains contract-level. | DEC-017 + DEC-019 (Formal Scope Change by User) |
| 2026-06-01 | BLK-003 | Each opportunity has exactly one quote; there is no second quote to accept, so the reject-vs-auto-demote question no longer exists. | DEC-018 (Formal Scope Change by User) |
| 2026-06-01 | BLK-A01 | Overdue-evaluation trigger resolved at PSM as on-read evaluation against the supplied `businessDate` (Asia/Shanghai); deterministic for G7 test design. | PSM "Resolved Mechanisms"; api-spec.md Reminder Query; FLOW-005 |
| 2026-06-03 | BLK-G11-001 | The user approved `https://118.196.44.193` as the production HTTPS endpoint. Codex installed Certbot 5.6.0, issued a Let's Encrypt IP certificate with SAN `IP Address:118.196.44.193`, enabled Nginx 443, verified HTTP→HTTPS redirect and server-side deploy smoke, and configured renewal timer/dry-run. | `docs/release/acc-017-evidence-template.md`; TASK-039 server evidence |
| 2026-06-03 | BLK-G11-002 | Volcengine API evidence confirms CRM `8080` and PostgreSQL `5432` are not publicly allowed; host-network Hermes `8642` was stopped/removed; old/non-CRM security-group rules TCP `8088`, TCP `8443`, and TCP `3389` were removed by the user and verified by API. | `docs/release/evidence/volcengine-security-group-post-cleanup-2026-06-03.json`; `docs/release/evidence/old-deployment-release-2026-06-03.json`; TASK-039 closure evidence |
| 2026-06-03 | TASK-040 release evidence | Encrypted PostgreSQL backup `crm-postgres-20260603T104620Z.sql.gz.enc` was produced on `srv-volcengine-sh-01`, copied to `srv-aliyun-bj-01`, verified by checksum, and restored in rehearsal run `20260603T104837Z`. `crm-backup.timer` is enabled and active for daily 02:00 backups. | `docs/release/acc-017-backup-evidence-template.md`; `docs/release/evidence/backup-restore-rehearsal-2026-06-03.json`; TASK-040 closure evidence |
