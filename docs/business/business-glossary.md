# Business Glossary

## Document Control

- Project: CRM System
- Phase: G4 Business Design
- Owner Agent: Business Analyst
- Source: `docs/product/prd.md`, `docs/product/acceptance-matrix.md`
- Status: Accepted as Architecture Input

## Glossary

| Term | Definition | Source | Notes |
|---|---|---|---|
| Administrator | Role responsible for user/role governance, full governed CRM visibility, and global operation logs. | PRD-001, PRD-002 | P0 role |
| Sales Manager | Role responsible for team record visibility, assignment/transfer, team overview, and team reports. | PRD-001, PRD-002, PRD-018 | P0 role |
| Sales | Role responsible for owned/assigned leads, customers, opportunities, quotes, contracts, payments, activities, notes, and tasks. | PRD-001, PRD-002 | P0 role |
| Lead | Potential ToB sales target before qualification or conversion. | PRD-003, PRD-004 | May be Unassigned before assignment |
| Unassigned Lead | Lead without owner before Pending Qualification or later states. | PRD-003 | Cannot be qualified, edited by Sales, or converted |
| Qualified Lead | Lead marked Valid and eligible for downstream customer/opportunity work. | PRD-004 | Converts to opportunity context |
| Invalid Lead | Lead marked not suitable with invalid reason. | PRD-004 | Cannot convert unless restored |
| Company / Customer | ToB account or organization record used for contacts, opportunities, quotes, contracts, and payments. | PRD-005 | v1 uses company/customer concept together |
| Contact | Person or role linked under a company/customer. | PRD-006 | Requires contact name and contact method or role note |
| Opportunity | Sales deal record linked to customer, owner, expected amount, expected close date, stage, and status. | PRD-007 | Moves through v1 pipeline |
| Pipeline Stage | Business state describing opportunity progress. | PRD-008 | Includes terminal Won/Lost states |
| Quote | Commercial offer linked to opportunity and customer. | PRD-009 | Multiple quotes allowed; one Accepted quote per opportunity |
| Accepted Quote | Quote selected for contract linkage. | PRD-009 | Contract can link only to Accepted quote |
| Expired Quote | Quote whose validity end date has passed. | PRD-009 | Cannot link to new contract |
| Contract | Record-based contract object linked to customer, opportunity, and Accepted quote. | PRD-010 | Contract note is P0 required |
| Pending Signature Contract | Contract awaiting signature. | PRD-010, PRD-021 | Requires expected signed date, not signed/effective date |
| Expected Signed Date | Planned signature deadline for Pending Signature contract reminders. | PRD-010, PRD-021 | Not a substitute for signed/effective date |
| Signed/Effective Date | Date required once contract is Signed, Active, Completed, or terminated after signing. | PRD-010 | Required for signed lifecycle states |
| Payment Plan | Planned payment record linked to contract with due amount, due date, and status. | PRD-011 | Used for due/overdue tracking |
| Actual Payment | Recorded payment event linked to contract. | PRD-011 | Zero, negative, and overpayment are rejected |
| Partial Payment | Payment state where cumulative paid amount is greater than zero and less than contract amount. | PRD-011 | Supported in v1 |
| Overdue Payment | Payment state where due date has passed and unpaid amount remains. | PRD-011, PRD-021 | Creates authorized reminder |
| Activity | Business interaction record linked to CRM record. | PRD-012 | Preserves follow-up history |
| Note | Textual business note linked to CRM record. | PRD-012 | Visible by record permission |
| Task | Follow-up work item with owner, due date, status, and title. | PRD-012 | Due/overdue tasks create reminders |
| Reminder | In-app notification for due/overdue tasks, Pending Signature contracts past expected signed date, or due/overdue payments. | PRD-021 | In-app only for v1 |
| Owner | User responsible for a CRM record. | PRD-002 | Drives Sales visibility |
| Assigned Record | Record assigned to or owned by a user. | PRD-002 | Sales visibility depends on owner/assignment |
| Archive | Non-delete record action available to Administrator and Sales Manager for eligible records. | PRD-002 | Sales cannot archive |
| Archived Record | Record removed from active work views without hard deletion. | ACC-002, ACC-014, ACC-015 | Available through explicit archived filters and audit/history views |
| Hard Delete | Permanent deletion of core CRM records. | PRD-002 | Not allowed in v1 |
| Record-Local History | Business timeline visible from a related CRM record according to record permission. | PRD-014 | P0 |
| Admin / Global Operation Log | Administrator-only operational audit query across records and access-sensitive actions. | PRD-022 | P1 |
| Duplicate Warning | Non-blocking warning when company/contact/lead data matches configured duplicate rules. | PRD-019 | No automatic merge |
| CSV Import | v1 import format for authorized bulk data entry. | PRD-020 | Validates row by row |
| CSV Export | v1 export format for authorized CRM records. | PRD-020 | Sales cannot export |
| Basic Report | Counts and sums for committed CRM groupings using persisted authorized records. | PRD-023 | Not advanced analytics |
| Active Report | Default report view excluding archived records. | ACC-023 | Archived records require explicit filter |
| Opportunity Amount | Expected amount on an opportunity used for opportunity reporting. | PRD-007, PRD-023 | Platform-independent business metric |
| Quote Amount | Amount on a quote used for quote reporting. | PRD-009, PRD-023 | Summed by quote status |
| Contract Amount | Amount on a contract used for contract reporting. | PRD-010, PRD-023 | May differ from quote with reason |
| Due Amount | Amount expected by a payment plan. | PRD-011, PRD-023 | Used for payment plan reporting |
| Paid Amount | Actual amount recorded as payment. | PRD-011, PRD-023 | Used for actual payment reporting |
| Business Date | Workspace local date used for due/overdue reminder evaluation until Architecture defines timezone handling. | PRD-021 | Deployment timezone is Architecture input |
| Won | Terminal opportunity state after full payment is recorded. | PRD-013 | Cannot reopen in v1 |
| Lost | Terminal opportunity state requiring lost reason. | PRD-013 | Cannot reopen in v1 |
| No-Downgrade Rule | Governance rule that P0/P1 items cannot be downgraded, deleted, weakened, merged away, or accepted as partial work. | Project rules | Applies to all downstream artifacts |
