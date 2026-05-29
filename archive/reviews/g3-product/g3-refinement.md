# G3 Refinement

## Purpose

This document records the refinement work completed after the first G3 review
returned `Gate Blocked`.

## Date

2026-05-26

## Inputs

- `archive/reviews/g3-product/g3-review-decision.md`
- `docs/product/prd.md`
- `docs/product/acceptance-matrix.md`
- `docs/product/open-questions.md`
- `docs/product/decision-log.md`
- `docs/product/out-of-scope.md`

## Blocker Resolution Summary

| Blocker | Status | Resolution |
|---|---|---|
| G3-BLOCK-001 Business rules | Addressed for re-review | Required fields, qualification rules, stage transitions, contract rules, payment rules, and won/lost rules are now defined in `docs/product/prd.md`. |
| G3-BLOCK-002 Permissions | Addressed for re-review | Product-level minimum permission matrix is now defined in `docs/product/prd.md`. |
| G3-BLOCK-003 Lifecycle/state transitions | Addressed for re-review | Lead, opportunity, quote, contract, payment, and task transitions are now defined in `docs/product/prd.md`. |
| G3-BLOCK-004 Quote/contract/payment money rules | Addressed for re-review | Single-currency model, quote-contract amount difference rule, overpayment block, partial payment, overdue, and won criteria are now defined. |
| G3-BLOCK-005 P1 feature testability | Addressed for re-review | Duplicate warning rules, CSV import/export, in-app reminders, admin operation logs, and basic report metrics are now defined. |
| G3-BLOCK-006 Acceptance granularity | Addressed for re-review | Acceptance matrix rows were tightened and scenario-level coverage notes were added for broad items. |
| G3-BLOCK-007 Audit terminology | Addressed for re-review | PRD distinguishes P0 record-local business history from P1 admin/global operation logs and defines shared event IDs. |

## Documents Updated

| Document | Update |
|---|---|
| `docs/product/prd.md` | Added permission matrix, required fields, state transitions, financial rules, P1 minimum behavior, and history/audit distinction. |
| `docs/product/acceptance-matrix.md` | Rewrote P0/P1 rows with concrete success/failure behavior, completion standards, Ready status, and no G3 blockers. |
| `docs/product/open-questions.md` | Marked G3-blocking questions resolved for G3 and kept downstream security/architecture/launch questions open. |
| `docs/product/decision-log.md` | Added product decisions DEC-009 through DEC-016. |
| `docs/product/out-of-scope.md` | Added XLSX import/export, non-in-app reminder delivery, and required contract attachment upload exclusions. |

## Re-Review Fixes

| Finding | Resolution |
|---|---|
| Contract required fields conflicted with Pending Signature lifecycle. | Contract required fields are now status-specific; Pending Signature does not require signed/effective date, while Signed/Active/Completed and post-signature Terminated contracts do. |
| Lead owner requirement conflicted with Unassigned state. | Lead owner is now conditional; Unassigned leads can exist before assignment but cannot be qualified, edited by Sales, or converted. |
| Sales opportunity closure conflicted with archive/close wording. | Permission matrix now distinguishes sales closure as Won/Lost from archive actions. Sales can close owned opportunities as Won/Lost but cannot archive records. |
| Contract reminders were in PRD but missing from P1 acceptance. | P1 reminder behavior and ACC-021 now include contracts pending signature past expected signed date. |
| QA TDD re-review found that expected signed date was referenced by ACC-021 but not defined on Contract. | Contract required fields now define expected signed date for Pending Signature contracts, PRD rules identify it as the planned signature deadline, and ACC-010/ACC-021 verify it. |

## Re-Review Recommendation

Recommendation: move to G3 re-review by Business Analyst and QA TDD.

This document does not pass G3. Only the Product Manager gate owner with
required reviewer approval can record G3 as passed.

Implementation remains blocked until G8 passes.
