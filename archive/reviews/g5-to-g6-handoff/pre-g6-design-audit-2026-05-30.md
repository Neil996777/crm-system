# Pre-G6 Design Audit (G5 → G6 Handoff)

## Document Control

- Project: CRM System
- Date: 2026-05-30
- Role: Audit (independent reverse verification)
- Trigger: Before entering G6, audit the full design set, because G6 builds the
  MDA package (CIM/PIM/PSM, service mapping, state machines, domain events,
  traceability) from ALL design files as authority.
- Method: Audit delegated across three dimensions; findings consolidated here.
- Archive Note: Audit evidence only. Not design authority.

## Scope

All active design files under `docs/product/`, `docs/business/`, `docs/ux-ui/`,
`docs/security/`, and `docs/architecture/`, with focus on completeness,
internal consistency, and sufficiency as MDA/PSM input. No implementation code
written. No P0/P1 item downgraded, deleted, merged, or weakened.

## Dimension Results

| Dimension | Initial Verdict | After Repair |
|---|---|---|
| Acceptance & traceability completeness | Ready for G6 | Ready for G6 |
| Cross-document consistency | Not Ready (5 contradictions) | Ready for G6 (repaired) |
| Architecture sufficiency as MDA/PSM input | Ready for G6 | Ready for G6 |

### Dimension 1 — Acceptance & Traceability Completeness: Ready

All 23 P0/P1 items (ACC-001…ACC-023) reverse-trace cleanly: acceptance →
business capability (CAP-001…CAP-012) → service (SVC-001…SVC-010) → exactly one
Service Owner Agent → contract family → data ownership / forbidden access. No
P0/P1 item dropped, merged, or downgraded (17 P0 + 6 P1 intact). All eight
cross-capability flows (lead→opportunity→quote→contract→payment→won, owner
transfer, archive, import/export, reminders, reporting, audit, backup/restore)
have named owning services. gateway-bff correctly owns no domain data.

### Dimension 2 — Cross-Document Consistency: Not Ready → Repaired

The audit found the server-allocation edit was applied to `deployment-notes.md`,
`service-architecture-adr.md`, `open-questions.md`, and `service-acceptance-map.md`
but NOT to `architecture.md`, which still carried stale "one Alibaba Cloud ECS /
single host" runtime wording in 5 places (Overview, Constraints, Strategy table,
Topology diagram, ADR-ARCH-001 summary). Because G6 treats every design file as
authority, these would have corrupted the MDA deployment/host model.

Repair (2026-05-30): all 5 `architecture.md` references updated to runtime host
`srv-volcengine-sh-01` (Volcengine) + off-server backup target `srv-aliyun-bj-01`
(Alibaba); the topology diagram now places the backup target outside the runtime
subgraph. Verified by search: no stale runtime-host wording remains, and
`47.95.119.211` appears only as the Alibaba backup-target asset IP.

Report `breakdowns`/`GroupRow` contract was audited consistent across
`api-spec.md`, `service-acceptance-map.md` (ACC-018/ACC-023), and BR-014/BR-017;
no stale groups-only definition remains authoritative.

### Dimension 3 — Architecture Sufficiency as MDA/PSM Input: Ready

All seven PSM input dimensions are present with file references: service
list/bounded-context boundaries, aggregate/data ownership + forbidden access,
domain events (producers + failure behavior), API/event/error/permission
contracts, cross-service reliability (idempotency/retry/timeout/correlation ID),
and service-to-service trust boundaries. State machines are partial only in that
no single consolidated transition table exists — but every transition rule,
terminal state, and guard G6 needs is an accepted decision, so G6 models rather
than invents. No area forces G6 to assume architecture authority.

## Non-Blocking Notes For G6 (model, do not invent)

1. Build consolidated per-aggregate state-transition tables (lead, opportunity,
   quote, contract, payment, archive, owner transfer, close-lost) from the
   accepted distributed rules.
2. Tabulate event producer→consumer routing in `domain-events.md` from the
   accepted flow matrices.
3. Pin the deferred-by-design choices with accepted constraints: sync/async per
   flow, concrete outbox vs replay, import/export temp-file cleanup job/path.
4. Optional: tighten the payments `GroupRow` `dueAmount`/`paidAmount` from prose
   into an explicit schema variant in PSM (prose is currently unambiguous; not a
   contradiction).

## Audit Decision

After repair of the 5 `architecture.md` contradictions, the design set is
internally consistent, complete in traceability, and sufficient as MDA/PSM
input. The G5 → G6 handoff is `Ready for G6`.

Carried-forward release blockers (off-server backup copy + restore rehearsal,
HTTPS/TLS endpoint evidence, security-group/monitoring evidence) are correctly
deferred to release and do not block G6. Company-layer infrastructure register
reconciliation remains a deferred workspace-level item.
