# Role Permission Scenarios

## Document Control

- Project: CRM System
- Phase: G4 Business Design
- Owner Agent: Business Analyst
- Source: `docs/product/prd.md`, `docs/product/acceptance-matrix.md`
- Status: Accepted as Architecture Input

## Permission Principles

- Permission rules must identify actor, action, resource, condition, and
  expected result.
- Sales Manager scope is all team records in v1.
- Sales scope is owned/assigned records and related child records only.
- Administrator and Sales Manager act as themselves; no silent acting on behalf
  of Sales users.
- Unauthorized actions must not expose or mutate data.

## Scenario Matrix

| ID | Priority | Role | Action | Resource | Condition | Expected Result | Acceptance IDs | Status |
|---|---|---|---|---|---|---|---|---|
| PERM-001 | P0 | Unauthenticated | Access | Core CRM data | No authenticated session | Denied; no CRM data exposed | ACC-001, ACC-002 | Accepted as Architecture Input |
| PERM-002 | P0 | Administrator | Manage | Users and roles | Authenticated as Administrator | Allowed | ACC-001, ACC-002 | Accepted as Architecture Input |
| PERM-003 | P0 | Sales Manager | Manage | Users and roles | Authenticated as Sales Manager | Denied | ACC-002 | Accepted as Architecture Input |
| PERM-004 | P0 | Sales | Create | Lead | Authenticated as Sales | Allowed; assigned to self by default unless later assigned by manager | ACC-003 | Accepted as Architecture Input |
| PERM-005 | P0 | Sales | Qualify | Unassigned lead | Lead has no owner | Denied | ACC-003, ACC-004 | Accepted as Architecture Input |
| PERM-006 | P0 | Sales Manager | Assign/transfer | Team lead | Lead is in team scope | Allowed; owner-change history recorded | ACC-003, ACC-014 | Accepted as Architecture Input |
| PERM-007 | P0 | Sales | View/edit | Non-owned lead | Lead is not owned/assigned | Denied; no data exposed | ACC-002, ACC-003 | Accepted as Architecture Input |
| PERM-008 | P0 | Sales | View/edit | Owned customer/contact | Related to owned/assigned record | Allowed | ACC-005, ACC-006 | Accepted as Architecture Input |
| PERM-009 | P0 | Sales | View/edit | Unrelated customer/contact | No ownership or assignment relation | Denied | ACC-002, ACC-005, ACC-006 | Accepted as Architecture Input |
| PERM-010 | P0 | Sales | Create/edit | Owned opportunity | Opportunity owned/assigned to user | Allowed | ACC-007, ACC-008 | Accepted as Architecture Input |
| PERM-011 | P0 | Sales | Close Won/Lost | Owned opportunity | Closure rules satisfied | Allowed; history recorded | ACC-013, ACC-014 | Accepted as Architecture Input |
| PERM-012 | P0 | Sales | Archive | Core CRM record | Any condition | Denied | ACC-002 | Accepted as Architecture Input |
| PERM-013 | P0 | Sales Manager | Archive | Eligible team record | Record is in team scope, eligible, and active downstream obligations are resolved or archived | Allowed; history/log event recorded | ACC-002, ACC-014, ACC-022 | Accepted as Architecture Input |
| PERM-014 | P0 | Sales | Create/edit | Quote | Related opportunity is owned/assigned | Allowed | ACC-009 | Accepted as Architecture Input |
| PERM-015 | P0 | Sales | Create/edit | Contract | Related opportunity is owned/assigned | Allowed | ACC-010 | Accepted as Architecture Input |
| PERM-016 | P0 | Sales | Record | Payment | Related contract is owned/assigned | Allowed when payment rules pass | ACC-011 | Accepted as Architecture Input |
| PERM-017 | P0 | Sales | View | Record-local history | Related record is owned/assigned | Allowed | ACC-014 | Accepted as Architecture Input |
| PERM-018 | P0 | Sales | View | Record-local history | Related record is not owned/assigned | Denied | ACC-014 | Accepted as Architecture Input |
| PERM-019 | P1 | Sales Manager | View | Team overview | Authenticated as Sales Manager | Allowed | ACC-018 | Accepted as Architecture Input |
| PERM-020 | P1 | Sales | View | Team overview | Authenticated as Sales | Denied | ACC-018 | Accepted as Architecture Input |
| PERM-021 | P1 | Administrator | View | Global operation logs | Authenticated as Administrator | Allowed | ACC-022 | Accepted as Architecture Input |
| PERM-022 | P1 | Sales Manager | View | Global operation logs | Authenticated as Sales Manager | Denied | ACC-022 | Accepted as Architecture Input |
| PERM-023 | P1 | Sales | View | Global operation logs | Authenticated as Sales | Denied | ACC-022 | Accepted as Architecture Input |
| PERM-024 | P1 | Administrator | Import/export | CRM records | Authenticated as Administrator | Allowed for governed records | ACC-020 | Accepted as Architecture Input |
| PERM-025 | P1 | Sales Manager | Import/export | CRM records | Team scope | Allowed for team records | ACC-020 | Accepted as Architecture Input |
| PERM-026 | P1 | Sales | Import/export | CRM records | Any condition | Denied | ACC-020 | Accepted as Architecture Input |
| PERM-027 | P1 | Administrator | View | Basic reports | Authenticated as Administrator | Allowed | ACC-023 | Accepted as Architecture Input |
| PERM-028 | P1 | Sales Manager | View | Basic reports | Team scope | Allowed for team reports | ACC-023 | Accepted as Architecture Input |
| PERM-029 | P1 | Sales | View | Manager/admin reports | Authenticated as Sales | Denied | ACC-023 | Accepted as Architecture Input |
| PERM-030 | P0/P1 | Administrator, Sales Manager | View | Archived records | Explicit archived filter or audit/history context | Allowed within governed/team scope | ACC-014, ACC-015, ACC-023 | Accepted as Architecture Input |
| PERM-031 | P0/P1 | Sales | View | Archived non-owned/non-assigned records | Any condition | Denied | ACC-002, ACC-014, ACC-015 | Accepted as Architecture Input |

## Required Permission Test Patterns

- Positive path: actor can perform allowed action in valid scope.
- Negative path: actor cannot perform disallowed action.
- Visibility path: unauthorized records are hidden or denied.
- Mutation path: denied action does not change data.
- Audit path: permission-sensitive actions create required history/log events.
