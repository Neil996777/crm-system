# G5 Architecture Re-Review Decision

## Document Control

- Project: CRM System
- Review Date: 2026-05-30
- Gate: G5 Architecture Design (re-review after repair)
- Gate Owner: Architecture
- Required Reviewers: Product Manager, Business Analyst, UX Designer, UI Designer, Security Compliance
- Additional Required Reviewer (project strengthening GAP-PROC-001): Infrastructure Ops
- Decision: Block
- Method: Review delegated to each reviewer agent. Each verified its assigned
  blockers against the ACTUAL active architecture documents under
  `docs/architecture/`, not against the repair note's claims.
- Archive Note: This file records G5 re-review evidence. It is not design
  authority and does not replace active architecture documents.

## Inputs Reviewed

- Repaired architecture documents under `docs/architecture/`
- `g5-architecture-review-decision-2026-05-29.md` (original blockers)
- `g5-architecture-repair-note-2026-05-30.md` (repair claims)
- Upstream product / business / UX-UI / security source documents
- `company/infrastructure/` registers (for Infrastructure Ops)

No implementation code was written. No P0/P1 item was downgraded, deleted,
merged away, weakened, or accepted as partial work.

## Per-Reviewer Decision

| Reviewer | Decision | Basis |
|---|---|---|
| Product Manager | Pass | G5-BLK-001, G5-ISS-001 verified closed in active docs; ACC-017/ACC-019 scope intact. |
| UX Designer | Pass | G5-BLK-007, G5-BLK-008, G5-ISS-001 backed by archive-eligibility, concurrency, and duplicate-warning contracts. |
| UI Designer | Pass | All six UI-domain blockers bind to concrete DTOs in api-spec / frontend-backend-contract. |
| Security Compliance | Pass | G5-BLK-003/004/005/006, G5-ISS-004/005 implemented with concrete mechanisms. |
| Business Analyst | Block | G5-ISS-002 report metric contract does not represent BR-017 status-based groupings. |
| Infrastructure Ops | Block | Production host/provider not reconciled with company server inventory; no off-server backup target exists. |

## Repaired Items Confirmed Closed

Verified present and concrete in the active architecture documents (evidence
quoted in each reviewer's return):

- G5-BLK-001 (endpoint strategy, environment ownership split, security-group
  evidence requirement, monitoring target, release evidence ownership) — design-level closed.
- G5-BLK-003 (service-to-service auth: service-token + signed headers, rotation,
  rejection, caller verification; internal network not a trust boundary).
- G5-BLK-004 (opaque server-side session, expiry, logout, disabled-user
  handling, authz-version recheck/revocation).
- G5-BLK-005 (backup directory permissions, encryption, key handling, restore
  logging, restore privacy).
- G5-BLK-006 (HTTPS-only ingress, HTTP→HTTPS redirect, TLS, secure cookies,
  security headers, no public admin/debug).
- G5-BLK-007 (archive eligibility API, obligation DTO, blocked response,
  retry/refresh, history events).
- G5-BLK-008 (version/expectedVersion concurrency contract, VERSION_CONFLICT).
- G5-BLK-009 (owner transfer contract, OwnerChanged/transfer flow, manual
  exception, retry/failure, history events).
- G5-BLK-010 (Close Lost contract, required lostReason, terminal edit
  protection, post-close work path, audit events).
- G5-ISS-001 (duplicate normalization, cross-object safe lookup, warning token,
  proceed-after-warning, no merge/overwrite).
- G5-ISS-003 (reminder timezone Asia/Shanghai, business-date rules, row DTO).
- G5-ISS-004 (import/export object scope, routing, validation, CSV formula
  safety, temp retention/deletion, auth, operation log).
- G5-ISS-005 (append-only audit API, no update/delete, DB permission
  constraints, prevHash/eventHash tamper evidence).

## Remaining Open Blockers (G5 Re-Review)

| ID | Priority | Owner | Blocker | Required Fix |
|---|---|---|---|---|
| G5-RE-001 | P1 | Architecture, Business Analyst, UI Designer | Extends G5-ISS-002. The report metric contract (`GET /reports/sales-overview`) provides a flat metrics block plus `groupBy=owner\|stage\|province`, but does not represent the BR-017 status-based groupings (leads by status, quotes by status/amount, contracts by status/amount, payments by status/amount), and the `groups[]` row DTO shape is undefined. A P1 business rule is not faithfully represented. | Add status-based groupings required by BR-017 to the report contract, define the `groups[]` element DTO (group key, label, count, amount), and confirm ACC-018/ACC-023 coverage. |
| G5-RE-002 | P1 | Architecture, Infrastructure Ops | Extends G5-BLK-001 environment ownership. Architecture names "Alibaba Cloud ECS" as the production target, but `company/infrastructure/server-inventory.md` lists the Alibaba asset as a 2 vCPU / 2 GiB test/ops candidate while the production-candidate runtime is the Volcengine asset (4 vCPU / 8 GiB). The production host/provider is not reconciled with the company server inventory, and environment-ownership closure does not name a specific registered asset. | Reconcile the production target with the company server inventory, name the specific registered asset, and confirm the capacity profile against the chosen asset (not the pre-release/pilot profile). |

## Recorded Release Gaps (Not G5 Design Blockers)

These are correctly fenced as release-blocking gaps in the architecture and are
NOT treated as new G5 design blockers, but they must close before production
release and are recorded here for traceability:

- G5-BLK-002 was repaired at G5 level by explicitly marking same-host-only
  PostgreSQL backup as a release-blocking gap (the original G5 acceptance
  condition allowed "external target OR explicit release-blocking gap"). However
  no concrete off-server backup target exists yet — `company/infrastructure/
  backup-recovery-plan.md` records none. Off-server encrypted backup + restore
  rehearsal remain release blockers (OQ-001).
- G5-BLK-001 security-group evidence, monitoring evidence, and TLS endpoint
  evidence are design-complete but evidence capture is expected at release.

## Reviewer Conflict Resolution

G5-ISS-002 received conflicting verdicts: UI Designer marked it Closed (the
metric-tile DTO and zero-state exist for rendering), while Business Analyst
marked it Block (BR-017 status groupings are not represented). As gate owner,
Architecture upholds the Business Analyst decision: the no-downgrade rule
protects BR-017 business groupings, and render-time DTO presence does not
satisfy business-rule fidelity. G5-ISS-002 is recorded open as G5-RE-001.

## Required Repair Sequence

1. Architecture extends the report metric contract for BR-017 status groupings
   and defines the `groups[]` row DTO (G5-RE-001).
2. Architecture + Infrastructure Ops reconcile the production target asset with
   the company server inventory and name the registered asset (G5-RE-002).
3. Business Analyst and UI Designer re-verify the report contract.
4. Infrastructure Ops re-verifies the named production asset and capacity.
5. Architecture produces a repair note addendum.
6. Required reviewers perform a focused G5 re-review on G5-RE-001 and G5-RE-002.

## Final Decision

G5 Architecture Design is `Gate Blocked`.

Two P1 blockers (G5-RE-001, G5-RE-002) remain open. The 14 previously blocked
items (10 P0 + 4 P1 from the assigned domains above) are confirmed closed in the
active architecture documents. The project may not proceed to G6 MDA Modeling
until G5-RE-001 and G5-RE-002 are repaired and the affected reviewers approve.
