# Abuse Cases

## Document Control

- Project: CRM System
- Phase: G4 Security Design
- Owner Agent: Security Compliance
- Status: Accepted as Architecture Input

## Abuse Case Matrix

| ID | Priority | Scenario | Attack Path | Required Protection | Verification | Acceptance IDs |
|---|---|---|---|---|---|---|
| ABUSE-001 | P0 | Unauthenticated CRM data access | Directly open protected route or call protected API without session | Reject access; no CRM data exposed; safe sign-in or denied state; access failure logged where applicable | E2E and API negative tests | ACC-001, ACC-002 |
| ABUSE-002 | P0 | IDOR record access | Sales changes URL or API id to a non-owned lead, customer, opportunity, quote, contract, payment, task, or history record | Backend verifies ownership/assignment and related-parent scope; deny without existence or value leakage | API negative tests for each P0 resource group | ACC-002, ACC-003 to ACC-015 |
| ABUSE-003 | P0 | Unauthorized mutation | Sales submits edit, owner transfer, archive, or close request for a non-owned/non-assigned record | Backend denies; no data mutation; no misleading success response | API mutation tests and persistence checks | ACC-002, ACC-014, ACC-016 |
| ABUSE-004 | P0 | Privilege escalation to Administrator | Sales Manager or Sales attempts user/role API or changes request payload role fields | Backend denies user/role management; role assignment fields ignored or rejected unless Administrator action is authorized | API negative tests and operation-log checks | ACC-001, ACC-002, ACC-022 |
| ABUSE-005 | P0 | Last Administrator removal | Administrator disables or downgrades the only active Administrator | Operation blocked before save; no account state mutation; blocked operation logged | UI, API, and audit-log tests | ACC-001, ACC-002, ACC-022 |
| ABUSE-006 | P0 | Login enumeration | Attacker compares login messages for invalid credentials, disabled accounts, or unavailable accounts | Unified unauthenticated failure message; account details visible only to authorized Administrator in user management | E2E login negative tests | ACC-001 |
| ABUSE-007 | P1 | CSV formula injection | Import or export includes cells beginning with formula-like content that spreadsheet tools may execute | Import validation flags dangerous content or safely escapes it; export safely encodes dangerous cell prefixes | Import/export security tests | ACC-020 |
| ABUSE-008 | P1 | Import row authorization bypass | Sales Manager imports rows referencing records outside team scope or relationships they cannot manage | Row-level validation checks permissions and related references before mutation; invalid rows fail safely | Import integration tests with mixed valid/invalid rows | ACC-020 |
| ABUSE-009 | P1 | Import partial failure corruption | CSV contains valid and invalid rows that could overwrite or corrupt existing records | Valid rows handled according to import policy; invalid rows reported; existing valid records not corrupted; no silent overwrite | Import persistence tests | ACC-020, ACC-016 |
| ABUSE-010 | P1 | Export data leakage | Sales or unauthorized actor requests export; Sales Manager changes filters to include unauthorized or archived data | Backend denies Sales export; Sales Manager export limited to team records; archived inclusion requires explicit authorized filter; operation logged | Export negative and audit tests | ACC-020, ACC-022 |
| ABUSE-011 | P1 | Export confirmation bypass | User calls export API directly without UI confirmation | Backend requires authorized export request semantics and logs the export; UI confirmation remains required for user workflow | API and E2E tests | ACC-020, ACC-022 |
| ABUSE-012 | P1 | Global operation-log leakage | Sales Manager or Sales opens global log route/API or guesses log ids | Backend denies; no event names, target record names, before/after values, or user lists exposed | API and E2E negative tests | ACC-022 |
| ABUSE-013 | P0 | Record-local history leakage | Sales opens history for non-owned record or child event from another record | Backend checks related record permission for history query and event detail | History permission tests | ACC-014 |
| ABUSE-014 | P1 | Sensitive before/after value leakage | Generic error, toast, import row error, log summary, or permission-denied state includes contact, payment, contract note, or restricted diff value | UI and API return safe summaries by default; details require authorization and data classification checks | UI state tests and API response assertions | ACC-002, ACC-014, ACC-020, ACC-022 |
| ABUSE-015 | P1 | Report inference | Sales or unauthorized actor accesses reports, or report aggregate includes unauthorized records before filtering | Reports enforce authorization before aggregation; Sales denied manager/admin reports | Report authorization and aggregate tests | ACC-018, ACC-023 |
| ABUSE-016 | P0/P1 | Archived record bypass | User accesses archived records through active list, search, reminder, report, or direct id outside explicit archive context | Archived records excluded from active/default views; explicit archived filter and normal authorization required | E2E and API tests for active/default and archived filters | ACC-014, ACC-015, ACC-021, ACC-023 |
| ABUSE-017 | P0 | Hard-delete attempt | Actor calls a delete endpoint or UI route for core CRM data | No normal hard delete capability; backend rejects if attempted; data remains persisted | API negative and persistence tests | ACC-002, ACC-016 |
| ABUSE-018 | P0 | Business-rule bypass via direct API | Actor attempts forbidden stage transition, Won without a Signed contract, expired quote contract link, overpayment, or post-close stage edit | Backend validates business rule and authorization together; no mutation on failure | API business-rule negative tests | ACC-008 to ACC-013 |
| ABUSE-019 | P0/P1 | Acting on behalf of Sales user | Administrator or Sales Manager submits payload claiming another actor identity for business change | Backend records authenticated actor as actor; optional owner/assignee is separate from actor identity | Audit event tests | ACC-014, ACC-022 |
| ABUSE-020 | P1 | Duplicate-warning enumeration | User probes company/contact/lead values to reveal unauthorized record details | Duplicate warnings do not reveal restricted matched record names or contact details; no automatic merge | UI/API duplicate-warning tests | ACC-019 |
| ABUSE-021 | P0/P1 | Reminder leakage | Unauthorized reminder list reveals due/overdue tasks, contracts, or payment existence | Reminder queries apply record permissions and hide unauthorized reminders | Reminder visibility tests | ACC-021 |
| ABUSE-022 | P0/P1 | Stale session or stale role | User continues using old session after role/status change or disablement | Protected API requests re-evaluate active account and role; disabled users denied | Session/role-change integration tests | ACC-001, ACC-002 |

## Release-Blocking Security Expectations

- Open P0 abuse-case failures block release.
- Open P1 abuse-case failures block the committed release unless formally
  removed by sponsor scope change.
- Abuse-case verification must include backend/API tests, not only UI tests.
- Any implementation that satisfies a core path with mock, static-only,
  in-memory-only, or non-persistent behavior cannot satisfy P0/P1 acceptance.
