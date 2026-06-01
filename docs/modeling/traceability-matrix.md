# Traceability Matrix

Consolidated end-to-end traceability — the fourth MDA artifact (G6). One row per
product acceptance item (ACC-001..023, all P0 and P1). Each cited ID is a real ID
that exists in its source document; no ID is invented. The CIM/PIM/PSM/Contract/
Service/ARCH-ACC chain reuses the already-audited PSM Traceability table
(`docs/modeling/PSM.md`) verbatim and stays consistent with it. Left-column
references (PRD/BR/UX/UI/Security/Architecture) were verified against their source
docs. Task/Test/Audit are forward placeholders produced at later gates.

## Document Control

- Project: CRM System
- Date: 2026-06-01
- Role: Domain Modeling
- Gate: G6 (MDA Modeling)
- Scope note: This matrix consolidates the verified end-to-end chain per
  acceptance item. Left-column refs (PRD, BR, UX, UI, Security, Architecture) are
  cited only where a covering ID actually exists in the source doc; where a
  dimension genuinely does not apply to an ACC the cell is `—` with a parenthetical
  reason (never a fabricated ID). The CIM/PIM/PSM/ARCH-ACC and Service/Contract
  columns reuse the audited `PSM.md` "PSM Traceability" chain. Task ID = `pending
  (G8)`, Test ID = `pending (G7)`, Audit ID = `pending (G12)` are forward
  placeholders. The three former open P0 PM/BA decisions (BLK-001/002/003) were
  RESOLVED 2026-06-01 by a Formal Scope Change by User (decision-log.md
  DEC-017..020): Won = related contract Signed (not full payment); exactly one quote
  per opportunity; payment retained but decoupled from Won; Opportunity Status
  removed. The ⚠ blocker flags are accordingly removed; affected CIM IDs retired in
  place (CIM-016) are no longer cited.

## ID source legend (verified prefixes)

- PRD-* — `docs/product/prd.md` (reused from `acceptance-matrix.md` Source column).
- BR-* — `docs/business/business-rules.md` (BR-001..021; inverted from its
  Acceptance IDs column).
- UX-* — `docs/ux-ui/ux-flows.md` (UX-001..011; inverted from its Acceptance IDs
  column).
- UI-* — `docs/ux-ui/ui-spec.md` Screen Index (UI-001..017; inverted from its
  Acceptance IDs column).
- SEC-* — `docs/security/security-requirements.md` (SEC-001..018; inverted from
  its Acceptance IDs column, including range-expressed rows SEC-003 "ACC-003 to
  ACC-015" and SEC-017 "ACC-001 to ACC-023").
- Architecture ID — `docs/architecture/service-architecture-adr.md`
  (ADR-ARCH-001..005) and `docs/architecture/risk-register.md` (ARCH-RISK-001..014),
  cited per ACC via the audited `PSM.md` "Architecture Acceptance" Source column.
- CIM-* — `docs/modeling/CIM.md`; PIM-* — `docs/modeling/PIM.md`; PSM-* /
  ARCH-ACC-* — `docs/modeling/PSM.md` (reused from the audited PSM Traceability +
  Architecture Acceptance tables).

## Matrix

| PRD ID | Acceptance ID | Business Rule ID | UX ID | UI ID | Security ID | Architecture ID | CIM ID | PIM ID | PSM ID | Architecture Acceptance ID | Task ID | Test ID | Audit ID |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| PRD-001, NFR-002, DEC-005 | ACC-001 | BR-001 | UX-001 | UI-001, UI-017 | SEC-001, SEC-002, SEC-008, SEC-016, SEC-017 | ADR-ARCH-005; ARCH-RISK-002, ARCH-RISK-009 | CIM-001, CIM-002, CIM-003; CIM-PROC-001, CIM-PROC-024 | PIM-001, PIM-002 | PSM-001 | ARCH-ACC-002, ARCH-ACC-008 | pending (G8) | pending (G7) | pending (G12) |
| PRD-002, NFR-002, DEC-010 | ACC-002 | BR-001, BR-002, BR-016 | UX-001, UX-011 | UI-001, UI-002, UI-016, UI-017 | SEC-002, SEC-003, SEC-004, SEC-005, SEC-006, SEC-007, SEC-014, SEC-015, SEC-016, SEC-017 | ADR-ARCH-002, ADR-ARCH-003, ADR-ARCH-005; ARCH-RISK-001, ARCH-RISK-002, ARCH-RISK-011 | CIM-001, CIM-002, CIM-003, CIM-004, CIM-005, CIM-006, CIM-007, CIM-037, CIM-039; CIM-PROC-002, CIM-PROC-014, CIM-PROC-020 | PIM-001, PIM-002, PIM-003 | PSM-001 | ARCH-ACC-001, ARCH-ACC-002, ARCH-ACC-009 | pending (G8) | pending (G7) | pending (G12) |
| PRD-003 | ACC-003 | BR-003, BR-004 | UX-002 | UI-004, UI-005 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-002, ARCH-RISK-012 | CIM-007, CIM-008, CIM-009, CIM-010; CIM-PROC-003 | PIM-004, PIM-003 | PSM-002 | ARCH-ACC-002, ARCH-ACC-010 | pending (G8) | pending (G7) | pending (G12) |
| PRD-004 | ACC-004 | BR-004 | UX-002 | UI-005 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-003, ARCH-RISK-008 | CIM-008, CIM-009; CIM-PROC-004 | PIM-004 | PSM-002 | ARCH-ACC-003, ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G12) |
| PRD-005 | ACC-005 | BR-002, BR-003 | UX-003 | UI-004, UI-006 | SEC-003, SEC-006, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-012, ARCH-RISK-013 | CIM-011, CIM-012; CIM-PROC-006 | PIM-005 | PSM-003 | ARCH-ACC-010, ARCH-ACC-011 | pending (G8) | pending (G7) | pending (G12) |
| PRD-006 | ACC-006 | BR-003 | UX-003 | UI-006 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-012 | CIM-013; CIM-PROC-006 | PIM-006 | PSM-003 | ARCH-ACC-010 | pending (G8) | pending (G7) | pending (G12) |
| PRD-007 | ACC-007 | BR-003, BR-005 | UX-004 | UI-007 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-008, ARCH-RISK-012 | CIM-014, CIM-029; CIM-PROC-007 | PIM-007, PIM-012 | PSM-004 | ARCH-ACC-007, ARCH-ACC-010 | pending (G8) | pending (G7) | pending (G12) |
| PRD-008 | ACC-008 | BR-005, BR-020 | UX-004 | UI-007 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-003, ARCH-RISK-008 | CIM-014, CIM-015; CIM-PROC-007 | PIM-007 | PSM-004 | ARCH-ACC-003, ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G12) |
| PRD-009, DEC-018 | ACC-009 | BR-003, BR-006 | UX-005 | UI-008 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-003, ARCH-RISK-012 | CIM-019, CIM-020, CIM-029; CIM-PROC-008 | PIM-008, PIM-012 | PSM-005 | ARCH-ACC-003, ARCH-ACC-010 | pending (G8) | pending (G7) | pending (G12) |
| PRD-010, DEC-006, DEC-007, DEC-016 | ACC-010 | BR-003, BR-007 | UX-005 | UI-009 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-012 | CIM-021, CIM-022, CIM-023, CIM-024, CIM-025, CIM-029; CIM-PROC-009 | PIM-009, PIM-012 | PSM-006 | ARCH-ACC-010 | pending (G8) | pending (G7) | pending (G12) |
| PRD-011, DEC-013, DEC-014, DEC-019 | ACC-011 | BR-003, BR-008 | UX-006 | UI-010 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-008, ARCH-RISK-003 | CIM-026, CIM-027, CIM-028, CIM-029; CIM-PROC-010 | PIM-010, PIM-011, PIM-012 | PSM-007 | ARCH-ACC-007, ARCH-ACC-003 | pending (G8) | pending (G7) | pending (G12) |
| PRD-012 | ACC-012 | BR-003, BR-009 | UX-007 | UI-002, UI-011 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-003 | CIM-030, CIM-031, CIM-032; CIM-PROC-012 | PIM-013, PIM-014, PIM-015 | PSM-008 | ARCH-ACC-003 | pending (G8) | pending (G7) | pending (G12) |
| PRD-013, DEC-012, DEC-017 | ACC-013 | BR-005, BR-020, BR-008 (payment history preserved on closure — not a Won gate, DEC-019) | UX-004, UX-006 | UI-007, UI-010 | SEC-003, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-008 | CIM-017, CIM-018; CIM-PROC-011 | PIM-007, PIM-011 | PSM-004, PSM-007 | ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G12) |
| PRD-014, NFR-004 | ACC-014 | BR-010, BR-016, BR-020 | UX-010, UX-011 | UI-004, UI-007, UI-014, UI-016 | SEC-003, SEC-007, SEC-009, SEC-014, SEC-015, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-003, ARCH-RISK-013 | CIM-035, CIM-038, CIM-048, CIM-049; CIM-PROC-017, CIM-PROC-020 | PIM-018, PIM-019 | PSM-009 | ARCH-ACC-003 | pending (G8) | pending (G7) | pending (G12) |
| PRD-015 | ACC-015 | BR-016 | UX-011 | UI-002, UI-003, UI-016 | SEC-003, SEC-004, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-002 | CIM-038, CIM-047; CIM-PROC-023, CIM-PROC-020 | PIM-026, PIM-020 | PSM-011 | ARCH-ACC-002 | pending (G8) | pending (G7) | pending (G12) |
| PRD-016, NFR-001, DEC-008 | ACC-016 | BR-002, BR-015 | — (persistence guarantee; no user-facing UX flow) | UI-004 (Entity Detail pattern, range ACC-003 to ACC-016) | SEC-006, SEC-015, SEC-017, SEC-018 | ADR-ARCH-001, ADR-ARCH-002, ADR-ARCH-004; ARCH-RISK-001, ARCH-RISK-004 | CIM-045; CIM-PROC-021 | PIM-001 to PIM-028 (persisted) | PSM-013 (+ PSM-001..012) | ARCH-ACC-004, ARCH-ACC-015 | pending (G8) | pending (G7) | pending (G12) |
| PRD-017, NFR-003, DEC-004 | ACC-017 | BR-015 | — (deployment/ops; no user-facing UX flow) | — (deployment/ops; no product screen) | SEC-017, SEC-018 | ADR-ARCH-001, ADR-ARCH-004, ADR-ARCH-005; ARCH-RISK-004, ARCH-RISK-009 | CIM-046; CIM-PROC-022 | (deployment; PIM-OPEN-004) | PSM-014 | ARCH-ACC-004, ARCH-ACC-008, ARCH-ACC-013, ARCH-ACC-014 | pending (G8) | pending (G7) | pending (G12) |
| PRD-018 | ACC-018 | BR-014, BR-016, BR-017 | UX-008 | UI-012, UI-015 | SEC-013, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-006 | CIM-004, CIM-043; CIM-PROC-015 | PIM-024 | PSM-010 | ARCH-ACC-006 | pending (G8) | pending (G7) | pending (G12) |
| PRD-019 | ACC-019 | BR-011, BR-019 | — (no dedicated duplicate-warning UX flow; surfaced inline in UX-002/UX-003 lead/account create) | — (no UI screen rows ACC-019; warning surfaced inline in UI-005/UI-006 create) | SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-008 | CIM-040; CIM-PROC-005 | PIM-021 | PSM-002, PSM-003 | ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G12) |
| PRD-020, DEC-015 | ACC-020 | BR-012, BR-018 | UX-009 | UI-013 | SEC-011, SEC-012, SEC-014, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-005, ARCH-RISK-014 | CIM-041, CIM-042; CIM-PROC-016 | PIM-022, PIM-023 | PSM-012 | ARCH-ACC-005, ARCH-ACC-012 | pending (G8) | pending (G7) | pending (G12) |
| PRD-021, DEC-015 | ACC-021 | BR-007, BR-008, BR-009, BR-013, BR-016, BR-021 | UX-007, UX-011 | UI-002, UI-009, UI-010, UI-011, UI-016 | SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-008 | CIM-023, CIM-032, CIM-033, CIM-034; CIM-PROC-013 | PIM-015, PIM-016, PIM-017 | PSM-008 | ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G12) |
| PRD-022, NFR-005 | ACC-022 | BR-010 | UX-010 | UI-014 | SEC-007, SEC-008, SEC-010, SEC-012, SEC-014, SEC-015, SEC-016, SEC-017 | ADR-ARCH-002, ADR-ARCH-003, ADR-ARCH-005; ARCH-RISK-003, ARCH-RISK-011 | CIM-003, CIM-036, CIM-048, CIM-049; CIM-PROC-018, CIM-PROC-024 | PIM-019, PIM-001 | PSM-009 | ARCH-ACC-003, ARCH-ACC-009 | pending (G8) | pending (G7) | pending (G12) |
| PRD-023, NFR-007 | ACC-023 | BR-014, BR-016, BR-017 | UX-008, UX-010, UX-011 | UI-012, UI-015, UI-016 | SEC-013, SEC-014, SEC-017 | ADR-ARCH-002, ADR-ARCH-003; ARCH-RISK-006 | CIM-038, CIM-044; CIM-PROC-019, CIM-PROC-020 | PIM-025, PIM-020 | PSM-010 | ARCH-ACC-006 | pending (G8) | pending (G7) | pending (G12) |

## Coverage & Gaps

- Row count: 23 acceptance rows (ACC-001..023). All P0 (ACC-001..017) and all P1
  (ACC-018..023) items appear; no row was removed and no P0/P1 item was downgraded,
  merged, or satisfied by a non-persistent/mock path. This is consistent with the
  audited PSM Traceability coverage note.
- Cells set to `—` (with reason; no fabricated ID used):
  - ACC-016 UX ID — persistence guarantee; there is no user-facing UX flow for
    "data survives refresh/re-login/restart" (UX-flows enumerate user journeys, not
    the persistence NFR).
  - ACC-017 UX ID and UI ID — deployment/operations acceptance (HTTPS endpoint,
    backup/restore, security group, monitoring); no product screen or UX journey;
    realized at PSM-014 + ARCH-ACC-004/008/013/014, owner includes infrastructure-ops.
  - ACC-019 UX ID and UI ID — no UX-flows row and no UI Screen-Index row carries
    ACC-019; the duplicate warning is an inline behavior surfaced during lead/account/
    contact create (UX-002/UX-003, UI-005/UI-006) rather than a standalone screen/flow.
    Cited the real inline-host screens parenthetically but did not invent an
    ACC-019-specific UX/UI ID.
- Formerly blocker-flagged rows — all RESOLVED 2026-06-01 by Formal Scope Change by
  User (decision-log.md DEC-017..020); ⚠ flags removed:
  - ACC-007 — BLK-001 resolved by DEC-020 (Opportunity Status field removed;
    Pipeline Stage is the sole lifecycle dimension; retired CIM-016 no longer cited).
  - ACC-009 — BLK-003 resolved by DEC-018 (exactly one quote per opportunity; no
    second quote to accept).
  - ACC-011 — BLK-002 resolved by DEC-017 + DEC-019 (Won = contract Signed; payment
    retained but decoupled, so multi-plan "fully paid" no longer gates Won).
  - ACC-013 — BLK-002 resolved by DEC-017 (Won = contract Signed, not full payment).
  - ACC-004 — its former cross-reference note to BLK-003 is removed (no second-quote
    decision remains).
- Real left-column coverage gaps surfaced (genuine source gaps, not fabrications):
  - ACC-019 has no dedicated UX-flows or UI Screen-Index ID. Duplicate warning is
    modeled only as inline create-time behavior. If G6/G7 wants an explicit
    UX/UI surface ID for the duplicate-warning interaction, UX/UI design must add it;
    MDA correctly declined to invent one.
  - ACC-016 and ACC-017 have no UX surface, and ACC-017 has no UI surface, by nature
    (persistence/deployment). These are covered at architecture/PSM altitude
    (ADR-ARCH-001/002/004/005, ARCH-ACC-004/008/013/014/015, PSM-013/PSM-014), not by
    a product screen — expected, not a defect.
- Architecture-deferred items (not gate blockers): BLK-A01 (overdue-evaluation
  trigger) was resolved in PSM "Resolved Mechanisms" (on-read evaluation) and is not
  flagged here. Carried-forward release-evidence blockers (off-server backup +
  restore rehearsal, HTTPS/TLS, security group, monitoring) are modeled at PSM via
  ARCH-ACC-004/008/013/014/015 and proven at G11/audited at G12; they touch ACC-016
  and ACC-017.

## Rules

- Every P0/P1 acceptance item appears in this matrix (23/23).
- Gaps are marked: `—` cells carry an explicit reason. The former open P0 decisions
  (BLK-001/002/003) were resolved 2026-06-01 by Formal Scope Change (DEC-017..020),
  so no ⚠ blocker flags remain.
- No row was removed to hide incomplete work.
