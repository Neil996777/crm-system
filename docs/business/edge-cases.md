# Edge Cases

## Document Control

- Project: CRM System
- Phase: G4 Business Design
- Owner Agent: Business Analyst
- Source: `docs/product/prd.md`, `docs/product/acceptance-matrix.md`
- Status: Accepted as Architecture Input

## Edge Case Matrix

| ID | Priority | Scenario | Expected Behavior | Acceptance IDs | Status |
|---|---|---|---|---|---|
| EDGE-001 | P0 | Unauthenticated user opens CRM route | Access denied; no core CRM data exposed | ACC-001, ACC-002 | Accepted as Architecture Input |
| EDGE-002 | P0 | Sales tries to access non-owned/non-assigned record | Record hidden or denied; no data mutation | ACC-002 | Accepted as Architecture Input |
| EDGE-003 | P0 | Lead save missing lead name/company name, source, or status | Save blocked with validation feedback | ACC-003 | Accepted as Architecture Input |
| EDGE-004 | P0 | Sales tries to qualify Unassigned lead | Action denied | ACC-003, ACC-004 | Accepted as Architecture Input |
| EDGE-005 | P0 | Invalid lead is converted to opportunity | Conversion rejected unless restored to Pending Qualification by Administrator or Sales Manager | ACC-004 | Accepted as Architecture Input |
| EDGE-006 | P0 | Converted lead is converted again | Conversion rejected | ACC-004 | Accepted as Architecture Input |
| EDGE-007 | P0 | Contact save without related company/customer | Save blocked | ACC-006 | Accepted as Architecture Input |
| EDGE-008 | P0 | Opportunity moves through forbidden transition | Transition rejected; no data mutation | ACC-008 | Accepted as Architecture Input |
| EDGE-009 | P0 | Opportunity is marked Won without a Signed contract | Won closure rejected | ACC-013 | Accepted as Architecture Input — amended 2026-06-01 (DEC-017): Won precondition is a Signed contract, not full payment |
| EDGE-010 | P0 | Opportunity is marked Lost without lost reason | Lost closure rejected | ACC-013 | Accepted as Architecture Input |
| EDGE-011 | P0 | Won or Lost opportunity is reopened | Reopen rejected in the committed scope | ACC-008, ACC-013 | Accepted as Architecture Input |
| EDGE-012 | P0 | A second quote is created for an opportunity that already has one | Blocked — each opportunity has exactly one quote | ACC-009 | Accepted as Architecture Input — amended 2026-06-01 (DEC-018): one quote per opportunity; previous multiple-Accepted scenario no longer applies |
| EDGE-013 | P0 | Expired quote is linked to new contract | Link rejected | ACC-009, ACC-010 | Accepted as Architecture Input |
| EDGE-014 | P0 | Pending Signature contract lacks expected signed date | Save rejected | ACC-010, ACC-021 | Accepted as Architecture Input |
| EDGE-015 | P0 | Pending Signature contract lacks signed/effective date | Save allowed if other required fields exist | ACC-010 | Accepted as Architecture Input |
| EDGE-016 | P0 | Signed/Active/Completed contract lacks signed/effective date | Save or transition rejected | ACC-010 | Accepted as Architecture Input |
| EDGE-017 | P0 | Contract amount differs from accepted quote amount | Difference reason required | ACC-010 | Accepted as Architecture Input |
| EDGE-018 | P0 | Actual payment amount is zero or negative | Payment rejected | ACC-011 | Accepted as Architecture Input |
| EDGE-019 | P0 | Actual payment exceeds remaining contract amount | Overpayment rejected | ACC-011 | Accepted as Architecture Input |
| EDGE-020 | P0/P1 | Payment due date passes with unpaid amount | Payment becomes Overdue and reminder is available to authorized user | ACC-011, ACC-021 | Accepted as Architecture Input |
| EDGE-021 | P1 | Pending Signature contract passes expected signed date | Authorized user sees in-app reminder | ACC-021 | Accepted as Architecture Input |
| EDGE-022 | P1 | Signed, terminated, or fully paid contract appears in reminders | No active reminder shown | ACC-021 | Accepted as Architecture Input |
| EDGE-023 | P0/P1 | Completed or cancelled task reaches due date | No active reminder shown | ACC-012, ACC-021 | Accepted as Architecture Input |
| EDGE-024 | P0 | Parent record owner changes while open tasks exist | Open tasks and follow-ups transfer unless manually reassigned | ACC-002, ACC-012, ACC-014 | Accepted as Architecture Input |
| EDGE-025 | P1 | Duplicate company/contact/lead data is entered | Warning shown; save may continue; no automatic merge | ACC-019 | Accepted as Architecture Input |
| EDGE-026 | P1 | CSV import contains valid and invalid rows | Valid rows import; invalid rows reported; existing records not corrupted | ACC-020 | Accepted as Architecture Input |
| EDGE-027 | P1 | CSV export requested by Sales user | Export denied | ACC-020 | Accepted as Architecture Input |
| EDGE-028 | P1 | Report opened with no data | Empty or zero report state shown | ACC-023 | Accepted as Architecture Input |
| EDGE-029 | P0/P1 | User tries to edit record-local history or global operation logs | Edit rejected or unavailable | ACC-014, ACC-022 | Accepted as Architecture Input |
| EDGE-030 | P0 | User refreshes, logs out/in, or service restarts after save | Persisted data remains available according to permissions | ACC-016 | Accepted as Architecture Input |
| EDGE-031 | P0/P1 | Archived record would appear in active list, active reminder, or default report | Archived record is hidden from active/default views and available only through explicit archived filters or audit/history views | ACC-014, ACC-015, ACC-021, ACC-023 | Accepted as Architecture Input |
| EDGE-032 | P0/P1 | User tries to archive record with unresolved active downstream obligations | Archive is blocked or user must resolve/archive related active tasks, pending-signature contracts, or unpaid payment items first | ACC-002, ACC-014, ACC-021 | Accepted as Architecture Input |
| EDGE-033 | P1 | Duplicate company name differs only by case or leading/trailing spaces | Duplicate warning shown | ACC-019 | Accepted as Architecture Input |
| EDGE-034 | P1 | Duplicate contact phone differs only by spaces, hyphens, or parentheses | Duplicate warning shown | ACC-019 | Accepted as Architecture Input |
| EDGE-035 | P1 | Duplicate contact email differs only by case or leading/trailing spaces | Duplicate warning shown | ACC-019 | Accepted as Architecture Input |
| EDGE-036 | P0 | User tries to edit opportunity stage after Won or Lost | Edit rejected; only notes/tasks may be added through normal permissions | ACC-008, ACC-013, ACC-014 | Accepted as Architecture Input — amended 2026-06-01 (DEC-020): `status` removed |
| EDGE-037 | P1 | Reminder due date is evaluated around local date boundary | Workspace business date is used; deployment timezone detail remains Architecture input | ACC-021 | Accepted as Architecture Input |

## G4 Blocking Edge Cases

No G4-blocking business edge case is currently open.

Remaining downstream questions:
- Exact production deployment provider/domain/backup belongs to Architecture.
- Data classification and retention belongs to Security Compliance.
- Data migration or initial seed data belongs to launch planning.
