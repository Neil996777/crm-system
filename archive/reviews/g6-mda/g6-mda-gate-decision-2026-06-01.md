# G6 (MDA Modeling) — Gate Decision

## Document Control

- Project: CRM System
- Gate: G6 — Architecture Design → MDA Modeling
- Owner: Domain Modeling + Architecture (Claude, planning platform)
- Reviewers: Product Manager, Business Analyst, UX Designer, UI Designer,
  Security Compliance, QA Test Design
- Date: 2026-06-01
- Decision: **Gate Passed**
- Archive note: Gate evidence only. Not design authority.

## Scope of the Gate

The G6 MDA package consists of five artifacts under `modeling/`:

- `CIM.md` — Computation Independent Model (47 active business concepts + 24 processes; CIM-016 retired)
- `PIM.md` — Platform Independent Model (domain objects, 10 active state machines + 1 retired, invariants, behaviors)
- `PSM.md` — Platform Specific Model (reflects the accepted G5 architecture; SVC-001..010, contracts, data ownership, flows, deployment)
- `traceability-matrix.md` — consolidated end-to-end ACC → CIM → PIM → PSM chain
- `test-model.md` — QA test design (acceptance/state/permission/invariant/edge/abuse coverage)

## How the package was produced and verified

Each artifact was authored by its responsible role (Domain Modeling for CIM/PIM/PSM/
traceability; QA Test Design for the test model) and put through independent
multi-agent audits (author ≠ reviewer) covering: tier-altitude/boundary discipline,
citation reality (no fabricated references), cross-tier consistency, no-invention,
and no-downgrade. A whole-package cross-tier audit confirmed the CIM→PIM→PSM→matrix→
test chain is consistent and that all 23 P0/P1 acceptances (ACC-001..023) trace
end-to-end.

## Formal Scope Change by User (2026-06-01)

During G6 the owner revised four committed P0 rules via the No-Downgrade-permitted
Formal Scope Change mechanism (decision-log.md DEC-017..020):

- DEC-017 — Won = related contract Signed (not full payment); `Contract Signed` and
  `Payment In Progress` opportunity stages removed.
- DEC-018 — exactly one quote per opportunity.
- DEC-019 — payment tracking retained but decoupled from Won (post-sale follow-up).
- DEC-020 — Opportunity `Status` field removed (Pipeline Stage is the sole lifecycle dimension).

The change was cascaded and re-audited across the G3/G4 baseline → G5 architecture
(consistency reconciliation only; service decomposition unchanged) → the full G6 MDA
→ UX/UI and security authority docs. Affected IDs were retired in place (CIM-016,
PIM-SM-003, PIM-INV-011, three test cases), not renumbered. This resolved the three
former G6 blockers (BLK-001/002/003 — see `planning/blockers.md` Resolution Log). No
P0/P1 capability was dropped (payment tracking retained, only decoupled from Won).

## Reviewer Sign-Off

| Reviewer | Verdict | Notes |
|---|---|---|
| Product Manager | SIGN-OFF | ACC 23/23 covered; no-downgrade; DEC-017..020 reflected; no out-of-scope creep. |
| Business Analyst | SIGN-OFF-WITH-CONCERNS | BR/BP/EDGE faithfully modeled; non-blocking traceability/clarity notes addressed. |
| UX Designer | SIGN-OFF-WITH-CONCERNS | All UX flows/states/reminders supported; stale UX authority docs reconciled to DEC-017..020. |
| UI Designer | SIGN-OFF-WITH-CONCERNS | UI-001..017 supported; stale `ui-spec` Won/status wording reconciled. |
| Security Compliance | SIGN-OFF-WITH-CONCERNS | Authz/abuse/audit/retention controls intact and not weakened; stale security-doc references reconciled. |
| QA Test Design | SIGN-OFF-WITH-CONCERNS | ACC 23/23, EDGE 37/37, ABUSE 22/22, state/invariant coverage; no Pending test remains. |

All raised concerns were non-blocking and have been reconciled (the scope-change
cascade was completed across the previously-missed `project-charter.md`,
`privacy-requirements.md` PRIV-005, `permission-matrix.md` PM-018, `audit-log-spec.md`,
`security-requirements.md`, `data-design.md`, and the full `docs/ux-ui/` set). A final
whole-repository sweep confirms no live old-rule wording remains; residual `status`
occurrences are other-entity statuses (lead/quote/contract/payment/task/customer) and
generic history/operation-log enumerations, which are correct.

## Decision

**Gate Passed.** The G6 MDA package is complete, internally and cross-tier consistent,
fully traceable (ACC-001..023), audited (author ≠ reviewer), reconciled to the
2026-06-01 Formal Scope Change, and signed off by all six reviewer roles. Proceed to
G7 (Task Planning — Domain Modeling + QA Test Design).

Carried-forward release blockers (off-server backup + restore rehearsal, HTTPS/TLS,
security-group, monitoring) remain release-time evidence for G11/G12, not G6 blockers.
