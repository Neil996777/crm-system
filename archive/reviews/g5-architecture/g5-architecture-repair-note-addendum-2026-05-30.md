# G5 Architecture Repair Note — Addendum

## Document Control

- Project: CRM System
- Date: 2026-05-30
- Role: Architecture
- Purpose: Record the architecture repair for the two blockers left open by the
  G5 re-review decision (`g5-architecture-re-review-decision-2026-05-30.md`).
- Archive Note: Repair evidence only. Active design authority remains under
  `docs/architecture/` and current upstream design documents.

## Scope

No implementation code was written. No P0/P1 acceptance item was downgraded,
deleted, merged away, weakened, or accepted as partial work.

User decision incorporated: both target servers are personal servers with no
prior allocation. The project is a learning/experiment build held to
production-grade discipline. The user directed Architecture to allocate the two
servers.

## Server Allocation Decision (Architecture + Infrastructure Ops)

| Role | Registered Asset | Provider / Region | Spec | Rationale |
|---|---|---|---|---|
| Production runtime host | `srv-volcengine-sh-01` | Volcengine / Shanghai | 4 vCPU / 8 GiB / 40 GiB | Only asset sized for Docker Compose + multiple Go services + PostgreSQL; the inventory's designated prod candidate. |
| Off-server backup target | `srv-aliyun-bj-01` | Alibaba Cloud / Beijing | 2 vCPU / 2 GiB / 40 GiB | Different provider + region = DR isolation; already runs the Alibaba Cloud backup client; sized for backup storage only. |

This allocation reconciles the production target with
`company/infrastructure/server-inventory.md` and simultaneously gives the
previously-undefined off-server backup requirement a concrete named target.

## Blocker Closure Mapping

| Re-review Item | Priority | Repair Summary | Active Design Reference |
|---|---|---|---|
| G5-RE-001 | P1 | Report contract now returns a `breakdowns` object with the mandatory BR-014/BR-017 groupings (leadsByStatus, opportunitiesByStage, quotesByStatus, contractsByStatus, paymentsByStatus) and a defined `GroupRow` DTO (key, label, count, amount; payments also dueAmount/paidAmount). Zero-state and authz-before-aggregation preserved. | `api-spec.md` (Report Metrics), `service-acceptance-map.md` (ACC-018, ACC-023) |
| G5-RE-002 | P1 | Production runtime host named as `srv-volcengine-sh-01` and off-server backup target named as `srv-aliyun-bj-01`, both reconciled with the company server inventory. Co-location constraints (port-80 reverse proxy, disk 61%/headroom, shared-memory upgrade trigger) recorded. Pre-release endpoint IP corrected to the runtime host. | `deployment-notes.md`, `service-architecture-adr.md` (ADR-ARCH-001, ADR-ARCH-004), `service-acceptance-map.md` (ACC-017), `open-questions.md` (OQ-001) |

## Effect On Recorded Release Gaps

- The off-server backup release gap (formerly G5-BLK-002) now has a concrete
  named target (`srv-aliyun-bj-01`). It remains a release blocker until
  Infrastructure Ops records encrypted off-server copy evidence and restore
  rehearsal — the design gap is closed, the evidence is still pending at release.

## Remaining Gate Status

Architecture repair is ready for a focused G5 re-review on G5-RE-001 (Business
Analyst, UI Designer) and G5-RE-002 (Infrastructure Ops). G5 is not passed until
those reviewers approve without open P0/P1 blockers.
