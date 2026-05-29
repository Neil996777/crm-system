# G3 Review Decision

## Gate

- Gate: G3
- Transition: PRD -> Acceptance Matrix
- Owner: Product Manager
- Required Reviewers: Business Analyst, QA TDD
- Date: 2026-05-26
- Decision: Gate Blocked

## Pass Condition

Every P0/P1 requirement has a verifiable acceptance item.

## Review Inputs

| Document | Review Result |
|---|---|
| `docs/product/prd.md` | Reviewed; P0/P1 requirements are mapped, but several need rule refinement for testability. |
| `docs/product/acceptance-matrix.md` | Reviewed; structurally complete, but several P0/P1 rows are not yet verifiably testable. |
| `docs/product/open-questions.md` | Reviewed; contains P0/P1-affecting unresolved questions. |

## Reviewer Decisions

| Reviewer | Decision | Summary |
|---|---|---|
| Business Analyst | Blocked | P0/P1 acceptance items depend on unresolved business rules, permissions, lifecycle transitions, and financial rules. |
| QA TDD | Blocked | P0/P1 requirements map to acceptance IDs, but several acceptance rows cannot yet produce reliable positive, negative, permission, and state-transition tests. |

## Coverage Result

| Check | Result | Notes |
|---|---|---|
| Every P0 requirement has an acceptance item | Pass | PRD-001 through PRD-017 map to ACC-001 through ACC-017. |
| Every P1 requirement has an acceptance item | Pass | PRD-018 through PRD-023 map to ACC-018 through ACC-023. |
| Acceptance rows include required fields | Pass | Matrix has role, precondition, trigger, expected result, failure behavior, completion standard, verification method, evidence, owner, status, related links, and blocker fields. |
| P0/P1 acceptance items are verifiable | Blocked | Multiple rows depend on unresolved business/security rules. |
| G3 pass allowed | No | G3 cannot pass while required reviewer blockers remain open. |

## Blockers

| ID | Severity | Area | Affected Acceptance IDs | Description | Required Resolution |
|---|---|---|---|---|---|
| G3-BLOCK-001 | P0 Blocker | Business rules | ACC-003, ACC-004, ACC-007, ACC-008, ACC-010, ACC-011, ACC-013 | Required fields, qualification rules, stage transitions, contract rules, payment rules, and won/lost rules remain open. | Resolve or explicitly block the related open questions; update acceptance rows with concrete rules. |
| G3-BLOCK-002 | P0 Blocker | Permissions | ACC-002, ACC-014, ACC-018, ACC-022, ACC-023 | Role/action/resource/condition rules are not concrete enough to test. | Add minimum permission matrix for Administrator, Sales Manager, and Sales across P0/P1 entities and actions. |
| G3-BLOCK-003 | P0 Blocker | Lifecycle/state transitions | ACC-004, ACC-008, ACC-009, ACC-010, ACC-011, ACC-012, ACC-013 | Lead, opportunity, quote, contract, payment, and task states are candidates, not stable transition rules. | Define allowed transitions, forbidden transitions, required fields, actor permissions, and history/audit events. |
| G3-BLOCK-004 | P0 Blocker | Quote/contract/payment money rules | ACC-009, ACC-010, ACC-011, ACC-013 | Multiple quotes, expired quotes, amount mismatch, overpayment, partial payment, overdue logic, and won criteria are unresolved. | Define v1 money model and closure rules. |
| G3-BLOCK-005 | P1 Blocker | P1 feature testability | ACC-019, ACC-020, ACC-021, ACC-022, ACC-023 | Duplicate detection, import/export, reminders, audit query, and report definitions lack concrete test inputs/outputs. | Define minimum v1 behavior for each P1 item or record a formal sponsor scope change. |
| G3-BLOCK-006 | P1 Issue | Acceptance granularity | ACC-003, ACC-012, ACC-015 | Some rows combine broad CRUD/search/filter/assignment or activity/note/task behavior. | Keep parent ACC rows, but add scenario-level sub-criteria or child rows for happy, negative, permission, and edge paths. |
| G3-BLOCK-007 | P1 Issue | Audit terminology | ACC-014, ACC-022 | Record-local collaboration history and admin/global operation logs overlap. | Define ACC-014 as record-local business history and ACC-022 as admin/global audit query with shared event IDs. |

## Required Next Work

1. Product Manager, Business Analyst, Security Compliance, and QA TDD refine P0/P1 rules before another G3 review.
2. Define permission matrix at product/business level before security deep-dive.
3. Define lifecycle/state-transition tables for Lead, Opportunity, Quote, Contract, Payment, and Task.
4. Define quote, contract, payment, and opportunity closure financial rules.
5. Define concrete P1 behavior for duplicates, import/export, reminders, audit query, and reports.
6. Update `docs/product/acceptance-matrix.md` so affected P0/P1 rows become verifiably testable.
7. Do not write implementation code.

## Gate Outcome

G3 is blocked.

The project remains in PRD / Acceptance Matrix refinement. Implementation is
not authorized and remains blocked until G8 passes.
