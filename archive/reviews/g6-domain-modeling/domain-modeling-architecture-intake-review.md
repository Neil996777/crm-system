# Domain Modeling Architecture Intake Review

## Document Control

- Project: CRM System
- Phase: G6 Domain Modeling Intake
- Reviewer Agent: Domain Modeling
- Review Type: Receiving-agent architecture handoff audit
- Date: 2026-05-27
- Status: Blocked

## Scope

This review audits whether the architecture artifacts can be represented by MDA
PSM without guessing or weakening accepted P0/P1 requirements.

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

Upstream spot-check inputs:

- `docs/product/prd.md`
- `docs/product/acceptance-matrix.md`
- `docs/business/*.md`
- `docs/ux-ui/*.md`
- `docs/security/*.md`

No `modeling/` files were created or edited.

## Audit Conclusion

Blocked.

Architecture is broadly sufficient for the main CRM loop, Go backend,
PostgreSQL persistence, role-based authorization, audit/history, import/export,
reports, reminders, deployment, backup, and frontend/backend separation.

However, two P0/P1 handoff gaps require Architecture repair before MDA Modeling
starts. Starting PSM now would force Domain Modeling to invent API/action/event
mapping for committed behavior, which would violate the no-downgrade and
receiving-agent handoff rules.

## Blocking Findings

### P0 Blocker: Invalid Lead Restore Is Not Architecturally Mapped

- Upstream requirement:
  - `docs/product/acceptance-matrix.md:19` requires ACC-004 to enforce lead
    transition rules and says an Invalid lead cannot convert unless restored to
    Pending Qualification by Administrator or Sales Manager.
  - `docs/product/prd.md:342` states the same restore rule.
  - `docs/business/business-rules.md:115` to `docs/business/business-rules.md:117`
    makes the restore prerequisite part of accepted business rules.
  - `docs/business/edge-cases.md:19` treats this as a P0 edge case.
- Architecture gap:
  - `docs/architecture/api-spec.md:145` only defines `/leads/{id}/qualify` for
    marking Valid or Invalid.
  - `docs/architecture/api-spec.md:146` defines conversion but not restore.
  - `docs/architecture/module-boundaries.md:57` covers owner requirement,
    Unassigned denial, and conversion-once, but not Invalid-to-Pending restore.
  - `docs/architecture/data-design.md:101` groups lead qualification/conversion
    in a transaction but does not define restore action, guard, actor, audit
    event, or PSM mapping expectation.
- Why this blocks MDA:
  - PSM would need to guess whether restore is implemented as generic PATCH,
    a dedicated action endpoint, a qualification action variant, or a status
    transition endpoint.
  - PSM would also need to guess the audit event mapping for the restore.
  - Because ACC-004 is P0, this cannot be left to implementation interpretation.
- Required Architecture repair:
  - Define the Invalid-to-Pending Qualification restore action in API, authz,
    data/transaction, audit/event, and frontend/backend contract terms.
  - Preserve the rule that only Administrator or Sales Manager may restore an
    Invalid lead, and conversion remains blocked until restore succeeds.

### P1 Blocker: Duplicate Warning Contract Is Under-Specified

- Upstream requirement:
  - `docs/product/acceptance-matrix.md:34` makes ACC-019 a P1 committed
    capability for company, contact, and lead duplicate warnings.
  - `docs/product/prd.md:430` requires exact company name, contact phone/email,
    and lead company/contact match warnings.
  - `docs/business/business-rules.md:312` to
    `docs/business/business-rules.md:327` defines normalization and
    non-blocking behavior.
  - `docs/security/permission-matrix.md:81` requires safe warning behavior that
    does not expose unauthorized matched record details.
  - `docs/security/abuse-cases.md:33` requires tests against duplicate-warning
    enumeration.
  - `docs/ux-ui/interaction-spec.md:39` and
    `docs/ux-ui/interaction-spec.md:57` require warning before final save and
    allow proceed-after-warning.
- Architecture support present:
  - `docs/architecture/architecture.md:89` mentions duplicate checks inside
    domain policies.
  - `docs/architecture/architecture.md:133` maps ACC-019 to a duplicate-check
    service with safe warning payloads.
  - `docs/architecture/data-design.md:90` mentions normalized search fields.
  - `docs/architecture/risk-register.md:31` tracks duplicate warning leakage.
- Architecture gap:
  - `docs/architecture/api-spec.md` has no endpoint, request/response, mutation
    integration rule, or error/warning payload contract for duplicate warnings.
  - `docs/architecture/frontend-backend-contract.md` does not map duplicate
    warning UI state to API behavior.
  - `docs/architecture/data-design.md` does not require normalized duplicate
    columns/indexes separately from search indexes or define duplicate-scope
    lookup behavior.
- Why this blocks MDA:
  - PSM would need to guess whether duplicate checking is a separate endpoint,
    embedded in create/update validation responses, or both.
  - PSM would need to invent warning DTO shape, safe match signal semantics,
    proceed-after-warning semantics, and authorized-scope behavior.
  - Because ACC-019 is P1 committed release scope, this cannot be deferred as an
    implementation detail.
- Required Architecture repair:
  - Define duplicate warning API behavior and DTOs for lead, company/customer,
    and contact create/edit flows.
  - Define safe payload semantics: warning signal without unauthorized matched
    record names, contact methods, or restricted details.
  - Define normalization/index requirements for company name, contact phone,
    contact email, and lead company/contact fields.
  - Define proceed-after-warning behavior without automatic merge or overwrite.

## Non-Blocking Improvements

### P2: Architecture Document Status Labels Are Stale

`PROJECT_CONTEXT.md` says Architecture Design is complete and ready for Domain
Modeling intake, while the architecture documents still say "Architecture
Design Draft for G6 Review" in document control. This does not weaken a P0/P1
behavior, but Architecture should align status labels after repair.

### P2: Sensitive Confirmation Mapping Can Be Clearer

UX/UI and Security require confirmations for terminal or high-impact actions.
Architecture covers export confirmation and backend validation, but the
frontend/backend contract could more explicitly map confirmation flows for Won,
Lost, archive, contract termination, import, export, and role/status changes.
This is not blocking because backend authorization and business validation are
already defined, and confirmation UI behavior is covered by UX/UI.

### P2: Money Representation Is Deferred To PSM

`docs/architecture/frontend-backend-contract.md:40` leaves decimal string vs
minor-unit money representation for PSM. This is acceptable if PSM resolves it
before task planning and maps it consistently to Go DTOs, PostgreSQL columns,
report sums, and frontend display.

## MDA Modeling Pre-Conditions

Before MDA Modeling starts:

1. Architecture must repair the P0 Invalid Lead restore mapping.
2. Architecture must repair the P1 duplicate warning API/contract/data mapping.
3. Architecture should update document control status labels after repair.
4. Domain Modeling should re-review only the repaired sections plus any
   affected traceability references.
5. `modeling/` must remain untouched until the re-review passes.

## Recommendation

Do not enter MDA Modeling yet.

Return the two blocking findings to Architecture. After Architecture repairs
the gaps, Domain Modeling should perform a focused re-review. If no new P0/P1
gaps are introduced, MDA Modeling can proceed to CIM/PIM/PSM creation.
