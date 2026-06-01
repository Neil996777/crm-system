# Task Dependencies — Build Order DAG

Prerequisites are listed per task (schema field 9). This file is the consolidated
DAG and the recommended build order. Edges point prerequisite → dependent. The
build respects service boundaries: a slice's frontend depends on its backend, and
both depend on the foundation (auth + permission + audit-history + the owning
service scaffold).

## Foundation roots (no prerequisites except TASK-001)

```
TASK-001 (monorepo + compose + PostgreSQL + migrate + shared/contracts)
  ├─> TASK-002 (identity-authz: auth/session)
  │     └─> TASK-003 (identity-authz: permission/scope/S2S/last-admin)
  │           ├─> TASK-004 (audit-history: append-only history/oplog/outbox)
  │           └─> TASK-005 (gateway-bff: routing/correlation/safe-error)
  │                 └─> TASK-006 (frontend: shell + sign-in + work overview)
```

TASK-002 depends on TASK-001. TASK-003 depends on TASK-002. TASK-004 depends on
TASK-001 + TASK-003. TASK-005 depends on TASK-002 + TASK-003. TASK-006 depends on
TASK-002 + TASK-005.

## Capability slices (each depends on TASK-003 + TASK-004 unless noted)

Leads:
```
TASK-007 (lead CRUD/assign)  ──> TASK-009 (lead UI)
TASK-007 ──> TASK-008 (lead qualify/convert)  ──> TASK-009
TASK-009 also depends on TASK-006 (frontend shell); TASK-009 builds the reusable
  list/detail pattern that TASK-030 later generalizes (TASK-030 depends on the slices,
  NOT on TASK-009 — see Retrieval block)
TASK-008 also depends on TASK-010 (account create/link) + TASK-013 (opportunity create)  [FLOW-002 conversion]
```

Account/Contact:
```
TASK-010 (account CRUD) ──> TASK-011 (contacts) ──> TASK-012 (account/contact UI)
TASK-010 ──> TASK-012
```

Opportunity:
```
TASK-013 (opp create, needs TASK-010 for customer link)
  └─> TASK-014 (stage transitions)
        └─> TASK-015 (close Won/Lost; also needs TASK-019 contract Signed)
TASK-013/014/015 ──> TASK-016 (opportunity UI)
```

Commercial (quote → contract → payment):
```
TASK-017 (quote, needs TASK-013 opp link) ──> TASK-021 (quote UI)
TASK-017 ──> TASK-018 (contract create from Accepted quote)
                └─> TASK-019 (contract lifecycle / Signed)  ──> TASK-022 (contract UI)
TASK-018 ──> TASK-020 (payments/overpayment) ──> TASK-023 (payment UI)
TASK-019 ──> TASK-015 (Won requires Signed contract)
```

Work / reminders:
```
TASK-024 (activities/notes/tasks; refs TASK-007/010/013/018/020) ──> TASK-025 (work UI)
TASK-019 + TASK-020 + TASK-024 ──> TASK-026 (reminders on-read + Reminder Center)
```

History / audit / admin:
```
TASK-004 + slices emitting events (TASK-007/014/017/019/020/024) ──> TASK-027 (history timeline UI)
TASK-004 + TASK-003 + event-producing slices ──> TASK-028 (admin operation-log query + UI)
TASK-003 + TASK-006 ──> TASK-029 (admin user/role UI)
```

Retrieval / cross-cutting:
```
TASK-005 + owning Query APIs (TASK-007/010/011/013/017/018/020/024) ──> TASK-030 (list/detail/search/filter)
TASK-007 + TASK-010 + TASK-011 ──> TASK-031 (duplicate warning)
TASK-019 + TASK-020 + TASK-024 + TASK-030 ──> TASK-032 (archive lifecycle)
```

Reporting:
```
event sources (TASK-007/013/017/018/020/024) ──> TASK-033 (team overview + projection infra)
TASK-033 ──> TASK-034 (basic reports)
```

Import / export:
```
TASK-003 + TASK-004 + target services (TASK-007/010/011/013/017/018/020/024) ──> TASK-035 (CSV import)
TASK-035 + TASK-030 ──> TASK-036 (CSV export)
```

Persistence / classification (cross-cutting, late):
```
all persistence-bearing slices (TASK-007..036) ──> TASK-037 (persistence verification)
TASK-004 + data-owning slices ──> TASK-038 (classification & retention)
```

## Deployment / release evidence (Phase 2, last)

```
TASK-001..038 (complete deployable system) ──> TASK-039 (deploy/HTTPS/security/monitoring)
                                                  └─> TASK-040 (off-server backup + restore rehearsal)
```

## Critical path (longest dependency chain)

```
TASK-001 → TASK-002 → TASK-003 → TASK-004
        → TASK-013 (needs TASK-010) → TASK-017 → TASK-018 → TASK-019 → TASK-015
        → TASK-016 / TASK-026 / TASK-032 → TASK-037 / TASK-038
        → TASK-039 → TASK-040
```

Commercial (quote→contract→payment) plus the opportunity-close-Won-requires-Signed-
contract edge (TASK-019 → TASK-015) is the longest functional chain; deployment and
backup/restore (TASK-039, TASK-040) are terminal.

## Parallelization guidance

After the foundation (TASK-001..005) is done, these slice clusters can proceed in
parallel (each cluster internally ordered as above):

- Cluster A: Leads (TASK-007 → TASK-008 → TASK-009)
- Cluster B: Account/Contact (TASK-010 → TASK-011 → TASK-012)
- Cluster C: Opportunity + Commercial (TASK-013 → TASK-017 → TASK-018 → TASK-019 →
  TASK-020 → TASK-015; UIs TASK-016/021/022/023)
- Cluster D: Work + Reminders (TASK-024 → TASK-025 → TASK-026), starts once its
  referenced record services exist.

Note the cross-cluster edges: TASK-008 (lead convert) needs B (account) + C
(opportunity); TASK-015 (close Won) needs C contract Signed (TASK-019). Schedule B
and the opportunity/contract part of C before completing TASK-008/TASK-015.

History UI (TASK-027), admin log (TASK-028), retrieval (TASK-030), duplicate
(TASK-031), archive (TASK-032), reporting (TASK-033/034), import/export
(TASK-035/036) integrate across clusters and run after their source slices.
TASK-037/038 are cross-cutting verification near the end; TASK-039/040 are the
final deployment + release-evidence stage.
