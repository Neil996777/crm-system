# UI Security Audit

## Document Control

- Project: CRM System
- Phase: G4 UI Security Intake
- Owner Agent: Security Compliance
- Reviewed Artifacts:
  - `docs/ux-ui/ui-spec.md`
  - `docs/ux-ui/component-spec.md`
  - `docs/ux-ui/responsive-spec.md`
- Reference Artifacts:
  - `docs/product/acceptance-matrix.md`
  - `docs/business/business-rules.md`
  - `docs/business/role-permission-scenarios.md`
  - `docs/ux-ui/screen-state-spec.md`
- Status: Approved with required changes

## Decision

Security Compliance approves the UI design direction with required changes.

The UI artifacts do not violate the no-downgrade rule and do not violate the
G8 implementation boundary. The UI also correctly states that navigation or
action hiding does not replace authorization enforcement.

The findings below must be resolved before formal Security Design can rely on
the UI artifacts as complete G4 inputs.

## Findings

### SEC-UI-001: Administrator Role And Status Changes Need Pre-Save Confirmation

- Severity: P0
- Source:
  - `docs/ux-ui/ui-spec.md`
  - `docs/ux-ui/component-spec.md`
- Related Acceptance:
  - ACC-001
  - ACC-002
  - ACC-022
- Status: Required change

Issue:
- `UI-017` defines Admin User/Role Management but only names save
  confirmation.
- `CMP-021` states that saving role changes creates success feedback and an
  operation log event.
- The UI does not yet require an explicit pre-save confirmation for user
  status changes, role changes, Administrator grants, Administrator removal, or
  account disablement.

Required Change:
- Add a role/status-change confirmation state to `UI-017`.
- Require confirmation content to show target user, old role, new role, old
  status, new status, access impact, and audit/log notice.
- State that UI confirmation does not replace backend authorization.
- Add a blocked state for governance-risk operations such as disabling or
  downgrading the last Administrator.

### SEC-UI-002: Export Needs Data Exfiltration Protection UI

- Severity: P1
- Source:
  - `docs/ux-ui/ui-spec.md`
  - `docs/ux-ui/component-spec.md`
- Related Acceptance:
  - ACC-020
  - ACC-022
- Status: Required change

Issue:
- `UI-013` covers import/export, but export confirmation is not explicit.
- `CMP-009` lists import run as a confirmation-modal case but does not list
  export run.
- CSV export is a bulk data-outflow workflow and must make scope and audit
  impact visible before execution.

Required Change:
- Add export run to `CMP-009`.
- Add an export confirmation state to `UI-013`.
- Confirmation must show object, filter conditions, authorization scope,
  estimated record count, archived inclusion or exclusion, and audit notice.
- Export result UI must not display unnecessary sensitive sample data.

### SEC-UI-003: Sensitive Data Redaction Rules Are Not Explicit Enough

- Severity: P1
- Source:
  - `docs/ux-ui/ui-spec.md`
  - `docs/ux-ui/component-spec.md`
- Related Acceptance:
  - ACC-014
  - ACC-020
  - ACC-022
- Status: Required change

Issue:
- History timeline, operation logs, and import row errors may display customer,
  contact, amount, payment, and before/after values.
- The UI artifacts say those values can appear, but do not yet define
  redaction or safe-summary rules for error messages, row errors, log tables,
  and permission-denied states.

Required Change:
- Add UI data-display rules for permission-denied messages, toasts, error
  summaries, import row errors, and log summaries.
- Default to safe summaries that do not echo restricted or sensitive raw
  values.
- State that event details may show before/after values only according to
  Security data classification and role authorization.
- Import row errors should show row number, field, and validation rule first,
  and must not default to full contact, amount, or payment value echoing.

### SEC-UI-004: Login Error Message May Reveal Account State

- Severity: P1
- Source:
  - `docs/ux-ui/ui-spec.md`
- Related Acceptance:
  - ACC-001
  - ACC-002
- Status: Required change

Issue:
- `UI-001` allows an invalid-credentials or disabled-account message on the
  unauthenticated sign-in screen.
- This may let an unauthenticated actor distinguish an existing disabled
  account from invalid credentials.

Required Change:
- Use one generic unauthenticated sign-in failure message for invalid
  credentials, disabled accounts, and unavailable accounts.
- Disabled-account detail may be visible only to authorized Administrators in
  user management.

## Passed Checks

- UI visibility is not treated as authorization enforcement.
- Permission-denied panel avoids revealing restricted record names or existence.
- Responsive behavior does not remove P0/P1 capability.
- G8 implementation boundary remains explicit.

## Required Follow-Up

1. UI Designer revises `docs/ux-ui/ui-spec.md` and
   `docs/ux-ui/component-spec.md` for SEC-UI-001 through SEC-UI-004.
2. Security Compliance performs UI re-audit.
3. Formal Security Design starts only after the UI re-audit has no unresolved
   P0/P1 required change.
