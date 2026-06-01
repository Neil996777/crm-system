# Privacy Requirements

## Document Control

- Project: CRM System
- Phase: G4 Security Design
- Owner Agent: Security Compliance
- Status: Accepted as Architecture Input

## OQ-014 Resolution

OQ-014 is resolved for the committed release as follows:

- CRM data is classified into Internal, Confidential, Restricted, and Security
  Critical classes.
- Customer, contact, contract, payment, log, report, import, and export data
  require role and record-scope visibility controls.
- Core CRM records are retained for business continuity and auditability; no
  normal product workflow hard-deletes core CRM records.
- the committed retention expectations are explicit below and must be implemented or
  carried into Architecture risk tracking without weakening P0/P1 acceptance.

This is a product security policy for the project, not formal legal advice.
Formal legal or regulatory review can strengthen these requirements later but
cannot weaken P0/P1 security, audit, or data-integrity behavior without formal
sponsor scope change.

## Data Classification Levels

| Class | Meaning | Default Handling |
|---|---|---|
| Internal | Non-public CRM operational metadata with low sensitivity. | Authenticated access only; role/scope filtering applies. |
| Confidential | Business, customer, contact, opportunity, and sales-work data. | Role and record-scope authorization; safe errors and summaries. |
| Restricted | Contract, payment, export, import file content, sensitive before/after values, and detailed commercial values. | Least privilege; masking in generic feedback; audit on sensitive operations. |
| Security Critical | Authentication, role assignment, access failures, operation-log integrity, and authorization decisions. | Administrator-only where visible; strong auditability; no frontend-only enforcement. |

## Data Handling Matrix

| ID | Data Category | Classification | Visible To | Masking / Display Rule | Retention Expectation | Archive / Delete Boundary | Audit Requirement | Acceptance IDs |
|---|---|---|---|---|---|---|---|---|
| PRIV-001 | User account identity, role, and status | Security Critical | Administrator; current user may see own session role where UI requires it | Non-admin cannot see user lists or role-management details | Retain while account exists and for 7 years after deactivation in operation logs | User deactivation is allowed; hard delete is not a normal CRM workflow | Role/status changes logged globally | ACC-001, ACC-002, ACC-022 |
| PRIV-002 | Lead data, source, status, owner, need summary | Confidential | Administrator; Sales Manager; owning/assigned Sales | Denied states do not reveal lead name, owner, or existence | Retain active records; retain archived records for 7 years after archive or final linked opportunity closure, whichever is later | Archive allowed for Administrator/Sales Manager; no normal hard delete | Owner, qualification, conversion, archive history | ACC-003, ACC-004, ACC-014 |
| PRIV-003 | Company/customer data | Confidential | Administrator; Sales Manager; Sales with owned/assigned relation | Generic errors do not reveal unauthorized company names | Retain active records; retain archived records for 7 years after archive or final related contract/payment closure, whichever is later | Archive allowed; no normal hard delete | Create/edit/archive history where applicable | ACC-005, ACC-014 |
| PRIV-004 | Contact names, titles, phone, email, contact methods, notes | Confidential | Administrator; Sales Manager; Sales with related owned/assigned record | Contact methods are not shown in permission denied, generic errors, import row summaries, or unauthorized duplicate details | Retain with related company/customer; 7 years after archive or final related business closure, whichever is later | Archive with related record policy; no normal hard delete | Create/edit/archive history where applicable | ACC-006, ACC-014, ACC-019 |
| PRIV-005 | Opportunity stage, amount, expected close date, close reason | Confidential | Administrator; Sales Manager; owning/assigned Sales | Unauthorized reports and lists exclude both row and aggregate values | Retain active records; retain 7 years after Won/Lost or archive, whichever is later | Won/Lost terminal; archive allowed; no normal hard delete | Stage/closure history | ACC-007, ACC-008, ACC-013, ACC-014, ACC-023 |
| PRIV-006 | Quote amount, validity, accepted/rejected state | Restricted | Administrator; Sales Manager; Sales with related owned/assigned opportunity | Quote values not exposed in generic errors or unauthorized aggregates | Retain 7 years after related opportunity closure or archive, whichever is later | Archive allowed; no normal hard delete | Quote acceptance and status history | ACC-009, ACC-014, ACC-022 |
| PRIV-007 | Contract amount, dates, status, note, difference reason | Restricted | Administrator; Sales Manager; Sales with related owned/assigned opportunity/contract | Contract note and amount are hidden from unauthorized views and safe summaries | Retain 7 years after contract completion, termination, opportunity closure, or archive, whichever is later | Archive allowed only after obligation checks; no normal hard delete | Signature, termination, status, amount-difference history | ACC-010, ACC-014, ACC-022 |
| PRIV-008 | Payment plan, actual payment amount, due/payment dates, status | Restricted | Administrator; Sales Manager; Sales with related owned/assigned contract | Payment values are masked from generic feedback, import row summaries, and unauthorized reports | Retain 7 years after full payment, contract closure, or archive, whichever is later | Archive allowed only after obligation checks; no normal hard delete | Payment recorded, completed, overdue history and operation log | ACC-011, ACC-013, ACC-014, ACC-022 |
| PRIV-009 | Activity, note, task content and follow-up details | Confidential | Administrator; Sales Manager; Sales with related owned/assigned record | Task/note content not exposed in unauthorized reminders or safe errors | Retain with related record; 7 years after related record archive/closure, whichever is later | Archive allowed; no normal hard delete | Creation, completion, cancellation, owner-transfer history where relevant | ACC-012, ACC-014, ACC-021 |
| PRIV-010 | Record-local history before/after values | Restricted | Users authorized for the related record | Before/after details follow the related data class and are not exposed in generic summaries | Retain at least as long as the related record | Not editable through normal CRM actions; no normal hard delete | Event itself is the audit evidence | ACC-014 |
| PRIV-011 | Admin/global operation log | Security Critical | Administrator only | Summaries avoid sensitive raw values by default; detail allowed only to Administrator and still follows classification | Retain 7 years for sensitive business operations; retain at least 3 years for access/login failures | Append-only through normal CRM actions; no normal hard delete | Global event model required | ACC-022 |
| PRIV-012 | CSV import file content and row errors | Restricted | Importing Administrator or Sales Manager within scope | Row errors show row number, field, rule, and safe summary; no full sensitive raw value by default | Raw uploaded import file retained only for processing duration; import result metadata retained 1 year; imported records follow their own class | Failed rows do not corrupt existing records | Import run operation log with counts and result | ACC-020, ACC-022 |
| PRIV-013 | CSV export file content | Restricted | Exporting Administrator or Sales Manager within scope | Export confirmation shows scope, filters, archived inclusion, and count; no sensitive sample rows | Generated export file is not retained server-side beyond delivery unless Architecture explicitly designs secure temporary storage with short expiration; export metadata retained 7 years | Export file deletion follows storage expiration; exported business records remain retained in CRM | Export run operation log with scope and count | ACC-020, ACC-022 |
| PRIV-014 | Basic reports and team overview metrics | Confidential / Restricted when amounts included | Administrator; Sales Manager for team scope | Unauthorized records excluded from rows and aggregates; Sales denied manager/admin reports | Report snapshots are not retained unless Architecture adds a controlled report artifact; source records follow their own retention | Archived records excluded by default unless explicit authorized archived filter | Report access failures and sensitive queries logged where applicable | ACC-018, ACC-023 |
| PRIV-015 | Duplicate warning match signals | Confidential | Actor creating/editing authorized record | Warning may indicate a possible duplicate but must not expose unauthorized matched record details | No separate retention except normal record and operation data | No automatic merge or overwrite | No operation log required unless Architecture treats as suspicious enumeration | ACC-019 |
| PRIV-016 | Authentication failures and denied authorization attempts | Security Critical | Administrator in global operation logs | User-facing messages remain generic | Retain at least 3 years in operation logs | Append-only through normal CRM actions | Access/login failure operation events | ACC-001, ACC-002, ACC-022 |

## Retention Policy

| Data Class | Minimum Committed Retention |
|---|---|
| Core active CRM records | Retained while active and in use. |
| Archived core CRM records | Retained 7 years after archive or final related business closure, whichever is later. |
| Contract and payment records | Retained 7 years after contract completion/termination, full payment, opportunity closure, or archive, whichever is later. |
| Record-local history | Retained at least as long as the related CRM record. |
| Admin/global operation logs | Retained 7 years for business-sensitive operations; at least 3 years for login/access failures. |
| Import result metadata | Retained 1 year; imported records follow their own classification. |
| Raw import files | Retained only for processing duration unless Architecture defines secure temporary storage with short expiration. |
| Generated export files | Not retained server-side beyond delivery unless Architecture defines secure temporary storage with short expiration. |
| Backups | Architecture must define backup retention and restore controls without shortening application-level retention below these expectations. |

## Deletion And Archive Boundary

- Normal committed CRM workflows do not hard-delete core CRM records.
- Archive is the normal way to remove records from active work views.
- Archived records remain available through explicit archived filters,
  record-local history, operation logs, and audit/report evidence according to
  role and scope.
- Exceptional legal, privacy, or operational removal requests are outside
  normal CRM user workflows and require a controlled operational process,
  Administrator authorization, audit evidence, and Architecture-defined backup
  handling.
- Any future deletion or anonymization feature must be introduced as formal
  scope and must not break P0/P1 history, reporting, payment, contract, or audit
  integrity.

## Masking And Safe Summary Rules

- Permission-denied messages show safe denial text and a safe return path.
- Toasts, table errors, import row errors, and log summaries do not echo
  contact methods, contract notes, payment amounts, customer-sensitive text, or
  sensitive before/after values unless the actor is authorized to see the
  detailed record.
- Administrator operation-log detail may include restricted before/after values
  only when necessary for audit review.
- Report aggregates must be computed after authorization filters, not filtered
  only after aggregation.
