# UX/UI Service State Mapping

## Document Control

- Project: CRM System
- Phase: G5 Pre-Architecture Input Supplement
- Owner Agents: UX Designer, UI Designer
- Status: Ready for Architecture Intake
- Date: 2026-05-29
- Sources:
  - `docs/product/business-capability-map.md`
  - `docs/product/acceptance-matrix.md`
  - `docs/ux-ui/ux-flows.md`
  - `docs/ux-ui/user-journeys.md`
  - `docs/ux-ui/screen-flows.md`
  - `docs/ux-ui/interaction-spec.md`
  - `docs/ux-ui/screen-state-spec.md`
  - `docs/ux-ui/ui-spec.md`
  - `docs/ux-ui/component-spec.md`
  - `docs/ux-ui/responsive-spec.md`

## Purpose

This document maps UX/UI paths and states to business capabilities and
service-backed state needs. It does not define APIs, final service boundaries,
database schemas, or implementation architecture.

## Flow To Capability Mapping

| Screen Flow ID | UX/UI Flow | Business Capabilities | Acceptance IDs | Service-Backed State Needs |
|---|---|---|---|---|
| SF-001 | Sign in to role workspace | CAP-001 | ACC-001, ACC-002 | Auth loading, invalid credentials, disabled user, expired session, permission-filtered navigation. |
| SF-002 | Lead create and qualification | CAP-002, CAP-003, CAP-004, CAP-008 | ACC-003, ACC-004, ACC-005, ACC-006, ACC-007, ACC-014 | Save pending, validation errors, owner transfer, invalid restore block, conversion success/failure, duplicate warning. |
| SF-003 | Opportunity pipeline | CAP-004, CAP-005, CAP-008 | ACC-007, ACC-008, ACC-013, ACC-014 | Stage transition pending, forbidden transition, required data prompt, Won-before-payment block, lost reason required. |
| SF-004 | Quote and contract | CAP-005, CAP-008 | ACC-009, ACC-010, ACC-014 | Accepted quote conflict, expired quote block, missing expected signed date, amount mismatch reason, contract lifecycle status. |
| SF-005 | Payment and closure | CAP-004, CAP-005, CAP-008 | ACC-011, ACC-013, ACC-014 | Payment validation, overpayment block, partial/full payment status, early Won block, terminal close confirmation. |
| SF-006 | Tasks and reminders | CAP-006, CAP-007, CAP-012 | ACC-012, ACC-021 | Reminder list loading, inactive reminder exclusion, permission-denied related record, stale reminder refresh. |
| SF-007 | Team overview | CAP-001, CAP-009 | ACC-018, ACC-023 | Manager-only access, empty team state, authorized aggregate loading, unauthorized aggregate exclusion. |
| SF-008 | Import/export | CAP-010, CAP-008, CAP-011 | ACC-020, ACC-022, ACC-016 | CSV validation, long-running progress, row-level partial failure, export scope confirmation, unsupported format, permission denied. |
| SF-009 | History and logs | CAP-008 | ACC-014, ACC-022 | Record-local history loading, global log admin-only access, safe event detail, read-only state. |
| SF-010 | Archive | CAP-012, CAP-007, CAP-008 | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 | Archive confirmation, active downstream obligation block, active/default filtering refresh, archived filter visibility. |
| SF-011 | Entity list/detail/search/filter pattern | CAP-007 and target record capability | ACC-015 | List loading, empty state, invalid filter, permission-hidden rows, permission-denied detail, stale result refresh. |
| SF-012 | Administrator user and role management | CAP-001, CAP-008 | ACC-001, ACC-002, ACC-022 | Last-admin block, role/status confirmation, denied non-admin access, operation log evidence. |

## Common Service-Backed UI States

| State ID | State | Required UX/UI Behavior | Contract Support Needed From Architecture | Acceptance IDs |
|---|---|---|---|---|
| UXSS-001 | Loading | Show non-blocking progress for list/detail/save/import/export/report operations. | Operation status or request lifecycle semantics for slow operations. | ACC-015, ACC-020, ACC-023 |
| UXSS-002 | Empty | Show useful empty state without implying missing permission. | Distinguish authorized empty result from permission denial where needed. | ACC-015, ACC-018, ACC-023 |
| UXSS-003 | Validation error | Place field-level errors near source input and preserve user-entered values. | Field error contract with safe messages and machine-readable codes. | ACC-003 to ACC-013, ACC-020 |
| UXSS-004 | Permission denied | Show safe denial and return path; do not reveal restricted names, values, or existence. | Permission denial contract with safe reason category. | ACC-001, ACC-002, ACC-014, ACC-015, ACC-020, ACC-022 |
| UXSS-005 | Disabled action | Disable or hide actions that are not available in current role/state, while backend remains authoritative. | Capability/action availability hints may be provided but cannot be authoritative. | ACC-002 to ACC-013, ACC-021 |
| UXSS-006 | Blocked transition | Explain why action is blocked and what user can do next. | Business rule error code and safe user-facing message. | ACC-004, ACC-008 to ACC-013, ACC-021 |
| UXSS-007 | Conflict / stale data | Inform user data changed and offer refresh/retry path. | Conflict or version mismatch signal where Architecture adopts concurrency control. | ACC-003 to ACC-013, ACC-016 |
| UXSS-008 | Partial failure | Summarize successes and failures without corrupting existing records. | Row-level result contract, failure counts, safe row summaries. | ACC-020 |
| UXSS-009 | Read-only audit/history | History/logs are inspectable but not editable. | Read-only query contract and no mutation route for normal CRM actions. | ACC-014, ACC-022 |
| UXSS-010 | Long-running operation | Show progress/pending/completed/failed states for import/export or other long tasks. | Run status contract and result retrieval contract. | ACC-020 |
| UXSS-011 | Sensitive value display | Display or mask values according to role/scope and data classification. | Data classification, permission, masking, and safe summary contracts. | ACC-002, ACC-014, ACC-020, ACC-022, ACC-023 |
| UXSS-012 | Archived context | Distinguish active/default results from explicit archived views. | Active/default filter and explicit archived filter semantics. | ACC-014, ACC-015, ACC-021, ACC-023 |

## UI Component Contract Support Needs

| Component / Pattern | Needs From Service Contracts | Acceptance IDs |
|---|---|---|
| Role-scoped navigation | Current user, role, allowed top-level areas, and safe denial on protected route. | ACC-001, ACC-002 |
| Entity list/table | Pagination/filter/search result, empty state distinction, permission-filtered rows, archived filter semantics. | ACC-015 |
| Entity form | Field validation, required-field errors, save success/failure, safe related-record lookup. | ACC-003 to ACC-013 |
| Status action control | Current status, allowed transitions hint, authoritative backend transition response, blocked reason. | ACC-004, ACC-008 to ACC-013 |
| Duplicate warning banner | Warning reason category, safe duplicate summary, proceed-after-warning token or equivalent if Architecture chooses one. | ACC-019 |
| Import result table | Run status, successful count, failed count, row number, field, rule, safe summary. | ACC-020 |
| Export confirmation/result | Export scope, filters, archived inclusion, count, result availability, safe failure. | ACC-020 |
| Reminder list | Reminder type, due/overdue state, related record safe display, permission-filtered output. | ACC-021 |
| History timeline | Event ID, actor, action, timestamp, safe before/after values, related record permission. | ACC-014 |
| Admin operation log table | Admin-only query, event ID, actor, resource, result, timestamp, safe detail. | ACC-022 |
| Report metric card/table | Metric grouping, authorized aggregate, zero state, denied state for Sales. | ACC-018, ACC-023 |

## Handoff Notes

- UX/UI states require backend contracts but do not define those contracts.
- Hidden or disabled UI is never authorization enforcement.
- Architecture must define final contract shapes and error categories.
- QA must test UI states against acceptance IDs and permission outcomes.

