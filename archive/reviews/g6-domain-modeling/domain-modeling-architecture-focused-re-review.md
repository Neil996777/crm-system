# Domain Modeling Architecture Focused Re-Review

## Document Control

- Project: CRM System
- Phase: G6 Domain Modeling Focused Re-Review
- Reviewer Agent: Domain Modeling
- Review Type: Focused re-review of Architecture repair
- Date: 2026-05-27
- Status: Passed

## Scope

This focused re-review checks only the two P0/P1 blockers raised in
`archive/reviews/g6-domain-modeling/domain-modeling-architecture-intake-review.md`
and the Architecture-owned document status labels affected by that repair.

Reviewed architecture inputs:

- `docs/architecture/architecture.md`
- `docs/architecture/module-boundaries.md`
- `docs/architecture/api-spec.md`
- `docs/architecture/data-design.md`
- `docs/architecture/integration-design.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/frontend-backend-contract.md`
- `docs/architecture/deployment-notes.md`
- `docs/architecture/risk-register.md`

No `modeling/` files were created or edited. No implementation code was
written.

## Review Conclusion

Passed.

Architecture has repaired the two original G6 handoff blockers without
downgrading, weakening, deleting, or merging away P0/P1 behavior. The repaired
architecture is now sufficient for Domain Modeling to express the affected
behavior in PSM without inventing API/action, permission, transaction, audit,
data, or frontend/backend contract semantics.

## Original Blocker Closure

| Original Finding | Priority | Re-Review Result | Closure Rationale |
|---|---|---|---|
| Invalid Lead restore was not architecturally mapped. | P0 | Closed | Architecture now defines `POST /leads/{id}/restore-pending`, Administrator/Sales Manager authorization, Sales denial, Invalid-only source state, owner-present guard, no conversion side effect, conversion blocked until restore commits, conversion-once enforcement, row-lock transaction behavior, and `EVT-STATUS-CHANGED` with reason code `LEAD_RESTORED_TO_PENDING_QUALIFICATION`. |
| Duplicate Warning contract was under-specified. | P1 | Closed | Architecture now defines `/duplicate-warnings/check`, `DUPLICATE_WARNING_REQUIRED`, safe warning payloads, warning tokens, `duplicateWarningAcknowledgement`, proceed-after-warning resubmission, create/edit mutation interception, authorization scope, restricted-match masking, duplicate normalization/index requirements, and no automatic merge/overwrite/link behavior. |

## Focus Area 1: P0 Invalid Lead Restore

Result: Passed.

The repaired architecture can be represented by PSM without guessing:

- API/action: `docs/architecture/api-spec.md` defines
  `POST /leads/{id}/restore-pending` as a dedicated business action.
- Permission: `docs/architecture/api-spec.md` and
  `docs/architecture/authz-architecture.md` restrict restore to Administrator
  and Sales Manager, with Sales denied.
- State rule: the source state must be Invalid, the target state is only
  Pending Qualification, the action does not convert, and Invalid leads remain
  blocked from conversion until restore succeeds.
- Conversion invariant: conversion remains one-time only and converted leads
  are denied repeat conversion.
- Transaction and locking: `docs/architecture/data-design.md` requires lead row
  locking and a durable transaction for qualification, restore, and conversion.
- History/audit: restore writes `EVT-STATUS-CHANGED` with reason code
  `LEAD_RESTORED_TO_PENDING_QUALIFICATION`.
- Frontend/backend contract: `docs/architecture/frontend-backend-contract.md`
  maps the restore UI/API behavior, success state, safe denial/blocking
  behavior, and committed-state requirement before conversion.
- PSM requirements: architecture documents explicitly require PSM mapping for
  the endpoint, DTO/action, authorization, transaction, event, and UI state.

No P0 blocker remains for Invalid Lead restore.

## Focus Area 2: P1 Duplicate Warning

Result: Passed.

The repaired architecture can be represented by PSM without guessing:

- API/action: `docs/architecture/api-spec.md` defines
  `POST /duplicate-warnings/check` for lead, company/customer, and contact
  create/edit preflight.
- Error/warning code: `DUPLICATE_WARNING_REQUIRED` is defined as the mutation
  interception response when acknowledgement is missing, invalid, stale, or
  no longer matches the intended payload.
- Safe payload: the warning DTO includes safe match categories, optional
  authorized match labels, restricted-match signals, request ID, and a
  server-issued warning token.
- Warning token: tokens are short-lived and bound to actor, operation, resource
  type, optional resource ID, payload hash, and duplicate rule set.
- Acknowledgement: create/edit mutations accept
  `duplicateWarningAcknowledgement` with `acknowledged=true` and the
  server-issued `warningToken`.
- Proceed-after-warning: acknowledged requests are revalidated before mutation;
  stale or changed payloads return a fresh warning without mutation.
- Mutation integration: lead, company/customer, and contact create/edit flows
  must run the same duplicate-warning check before durable mutation.
- Authorization scope: actor must be authorized for the target resource type
  and current record when editing.
- Unauthorized match masking: unauthorized matches may only produce
  `restrictedMatchSignal=true`; names, contact methods, owner names, values,
  counts, and links are not exposed.
- Normalization/indexing: `docs/architecture/data-design.md` requires
  normalized company, phone, and email fields plus non-unique active-record
  lookup indexes for warning support.
- No automatic merge/overwrite: architecture explicitly forbids automatic
  merge, overwrite, ownership transfer, contact linking, or conversion.
- Frontend/backend contract: frontend must handle both preflight and mutation
  interception, must not invent hidden details, and must display fresh warning
  state instead of assuming success.
- PSM requirements: architecture documents explicitly require PSM mapping for
  DTOs, warning token acknowledgement, safe display rules, repository lookup,
  normalized fields/indexes, and state transitions.

No P1 blocker remains for Duplicate Warning.

## Document Control Status Review

Result: Passed.

All reviewed architecture documents now use a non-draft status:

`Architecture Design Ready for G6 Focused Re-Review`

This status is acceptable for the focused re-review handoff and no longer
misleads Domain Modeling into treating the architecture package as an
unrepaired draft.

## New P0/P1 Blockers

None found in the focused re-review scope.

Open architecture risks remain tracked for MDA, QA, integration, and audit, but
the reviewed risks do not block the start of MDA Modeling because they are
explicitly mapped to downstream PSM/test/evidence work rather than left as
missing architecture behavior.

## Recommendation

Proceed to MDA Modeling.

Domain Modeling may now create CIM, PIM, PSM, domain model, state machines,
domain events, traceability matrix, and test model in the normal G6/G7 sequence,
while preserving all P0/P1 acceptance items and the repaired architecture
constraints.
