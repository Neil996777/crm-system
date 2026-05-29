# User Journeys

## Document Control

- Project: CRM System
- Phase: G4 UX Design
- Owner Agent: UX Designer
- Source: `docs/product/prd.md`, `docs/product/acceptance-matrix.md`,
  `docs/business/*`
- Status: Accepted as Architecture Input

## UX Principles

- CRM users need a work-focused, repeatable operational experience.
- High-frequency sales work should prioritize fast scanning, predictable
  navigation, inline validation, and clear next actions.
- UX must not change product scope or business rules.
- Permission denial, validation failure, empty state, conflict, and recovery
  paths must be explicit for P0/P1 flows.
- Architecture reset on 2026-05-29: implementation is blocked until the restarted delivery flow passes G8.

## Journey Index

| ID | Priority | Role | Journey | Goal | Acceptance IDs | Status |
|---|---|---|---|---|---|---|
| JRN-001 | P0 | Sales | Daily assigned work | Find assigned work, follow up, and update records | ACC-002, ACC-003, ACC-012, ACC-015, ACC-021 | Accepted as Architecture Input |
| JRN-002 | P0 | Sales | Lead to opportunity | Qualify a lead and create customer/opportunity context | ACC-003, ACC-004, ACC-005, ACC-006, ACC-007 | Accepted as Architecture Input |
| JRN-003 | P0 | Sales | Opportunity to quote and contract | Move a deal through quote and contract records | ACC-007, ACC-008, ACC-009, ACC-010, ACC-014 | Accepted as Architecture Input |
| JRN-004 | P0 | Sales | Payment and closure | Track payment and close opportunity correctly | ACC-011, ACC-013, ACC-014, ACC-016 | Accepted as Architecture Input |
| JRN-005 | P0/P1 | Sales Manager | Team management | Review pipeline, assign work, and manage team risk | ACC-002, ACC-018, ACC-021, ACC-023 | Accepted as Architecture Input |
| JRN-006 | P0/P1 | Administrator | Governance and audit | Manage access and review operation evidence | ACC-001, ACC-002, ACC-014, ACC-022 | Accepted as Architecture Input |
| JRN-007 | P1 | Administrator, Sales Manager | CSV import/export | Bulk load or extract authorized CRM records | ACC-020 | Accepted as Architecture Input |
| JRN-008 | P1 | Sales Manager, Administrator | Reports | Review persisted sales metrics | ACC-018, ACC-023 | Accepted as Architecture Input |

## JRN-001: Sales Daily Assigned Work

Entry:
- Sales signs in and lands on assigned work overview.

User goals:
- See assigned leads, opportunities, tasks, reminders, and urgent payment or
  contract items.
- Continue work without opening unauthorized records.

Journey:
1. Sales signs in.
2. UX shows assigned active work and due/overdue reminders.
3. Sales filters or searches assigned records.
4. Sales opens a record detail.
5. Sales adds note, activity, or task.
6. Sales receives save feedback and remains in context.

Required states:
- Empty assigned work.
- Loading list/detail.
- Permission denied for unauthorized detail.
- Validation error for incomplete task.
- Success feedback after save.

UX handoff notes:
- UI should support dense scanning and quick return to the active work list.

## JRN-002: Sales Lead To Opportunity

Entry:
- Sales creates or opens an assigned lead.

User goals:
- Capture lead information, qualify it, and create customer/opportunity
  context.

Journey:
1. Sales opens lead create or assigned lead detail.
2. Sales enters required lead fields.
3. UX validates required fields before save or status transition.
4. Sales qualifies lead as Valid or Invalid.
5. If Valid, Sales creates or links customer and contact.
6. Sales creates opportunity from qualified need.
7. UX confirms conversion and links to created opportunity.

Required states:
- Unassigned lead cannot be edited or qualified by Sales.
- Invalid lead requires reason.
- Converted lead cannot be converted again.
- Missing company/contact data shows field-level errors.

UX handoff notes:
- Conversion should preserve context and provide a clear path to the created
  opportunity.

## JRN-003: Sales Opportunity To Quote And Contract

Entry:
- Sales opens an owned/assigned opportunity.

User goals:
- Advance pipeline, create quote, accept quote, and create contract.

Journey:
1. Sales opens opportunity detail.
2. UX shows current stage, required next data, and history.
3. Sales changes stage when required data exists.
4. Sales creates quote from opportunity.
5. Sales sends and marks quote Accepted.
6. Sales creates Pending Signature contract from Accepted quote.
7. UX requires expected signed date and contract note.
8. Sales later signs contract with signed/effective date.

Required states:
- Forbidden stage transition blocked with reason.
- Expired quote cannot be linked to new contract.
- Only one Accepted quote remains per opportunity.
- Contract amount mismatch requires reason.
- Pending Signature does not require signed/effective date.

UX handoff notes:
- Stage actions should expose missing requirement reasons, not silently fail.

## JRN-004: Sales Payment And Closure

Entry:
- Sales opens contract or payment area from an opportunity.

User goals:
- Record payment plan, actual payment, and close the opportunity correctly.

Journey:
1. Sales creates payment plan on contract.
2. Sales records actual payment.
3. UX shows unpaid, partially paid, paid, or overdue status.
4. Sales attempts Won closure only after full payment.
5. Sales closes as Lost with lost reason when applicable.
6. UX confirms terminal closure and preserves history.

Required states:
- Zero, negative, or overpayment amount blocked.
- Won blocked before full payment.
- Lost requires reason.
- Won/Lost cannot be reopened.

UX handoff notes:
- Closure actions need confirmation because terminal states cannot be reopened.

## JRN-005: Sales Manager Team Management

Entry:
- Sales Manager opens team overview.

User goals:
- See team pipeline, identify risk, transfer ownership, and inspect records.

Journey:
1. Sales Manager opens team overview.
2. UX summarizes leads, opportunities, quotes, contracts, payments, tasks, and
   reminders.
3. Sales Manager opens record detail.
4. Sales Manager assigns/transfers team work.
5. UX shows task/follow-up transfer effect.
6. Sales Manager archives eligible team record when needed.

Required states:
- Empty team data.
- Unauthorized manager-only views hidden from Sales.
- Archive blocked by active downstream obligations with related-record links.
- Transfer success feedback and history visibility.

UX handoff notes:
- Team overview must support risk scanning without becoming an analytics
  dashboard beyond P1 basic reports.

## JRN-006: Administrator Governance And Audit

Entry:
- Administrator opens governance or audit area.

User goals:
- Manage role access and review operation evidence.

Journey:
1. Administrator signs in.
2. UX exposes user/role governance entry.
3. Administrator opens user/role management.
4. Administrator views user list and user detail.
5. Administrator updates user status or assigned role.
6. Administrator reviews role capability summary.
7. Administrator opens global operation logs.
8. Administrator filters or searches events.
9. Administrator opens related record where authorized.

Required states:
- Sales and Sales Manager denied from user/role management.
- User list loading and empty states.
- User save validation and failure states.
- Sales and Sales Manager denied from global operation logs.
- Log edit unavailable.
- Empty log result.
- Access failure events visible to Administrator.

UX handoff notes:
- Global logs and record-local history must look related but remain clearly
  separate concepts.

## JRN-007: CSV Import/Export

Entry:
- Administrator or Sales Manager opens import/export.

User goals:
- Import or export authorized CRM records without corrupting data.

Journey:
1. User chooses import or export.
2. User selects supported object and CSV file or export scope.
3. UX shows validation progress.
4. UX shows success count and row-level failures.
5. User downloads/export authorized records.

Required states:
- Sales denied from import/export.
- Unsupported format rejected.
- Partial failure result with failed row list.
- Export excludes unauthorized records.

UX handoff notes:
- Long-running import should provide progress and final summary.

## JRN-008: Reports

Entry:
- Administrator or Sales Manager opens reports.

User goals:
- Review persisted counts and sums for committed CRM groupings.

Journey:
1. User opens reports.
2. UX shows report groups for leads, opportunities, quotes, contracts, and
   payments.
3. User applies authorized filters.
4. UX shows zero/empty state if no data.
5. User drills into related records where authorized.

Required states:
- Sales denied from reports.
- Unauthorized records excluded.
- Archived records excluded by default unless explicit archived filter is used.

UX handoff notes:
- Reports are basic operational summaries, not forecasting or advanced
  analytics.
