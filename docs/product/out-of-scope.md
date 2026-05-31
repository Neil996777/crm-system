# Out-of-Committed-Scope Candidates (Pending Real-Client Decision)

Items in this file are **not** part of committed P0/P1/P2 scope. They are
candidate capabilities that were bounded out of the current committed CRM loop
during a planning phase that has **no real client (甲方)** yet.

Important framing (corrected 2026-05-31):

- None of these has been formally accepted out of scope by a real
  client/sponsor. Each is **left to a real-client decision**, not decided here.
- These are **not downgraded tasks**. None was ever a committed P0/P1/P2
  requirement (verified against acceptance items ACC-001…ACC-023). Removing a
  committed P0/P1 item from scope is a different action and may happen **only**
  through `Formal Scope Change by User`, recorded separately.
- Priority (P0/P1/P2) governs the **execution order of committed scope only**.
  It is never a mechanism to defer, weaken, or drop committed work, and these
  candidates are not assigned a priority until a real client commits them.

| ID | Item | Reason (why it is outside the current committed loop) | Status |
|---|---|---|---|
| OOS-001 | Contract approval workflow | Committed contract management is record/status/amount/linkage/attachment-or-notes only; an approval workflow is a separate governance capability not in committed scope. | Pending real-client decision |
| OOS-002 | Electronic signature | Requires external integration, legal workflow decisions, and provider selection. | Pending real-client decision |
| OOS-003 | Contract template generation | Requires template governance and document generation rules beyond record-based contract management. | Pending real-client decision |
| OOS-004 | Quote approval and discount approval workflow | Approval policies are not defined; the committed CRM loop runs on quote accept/reject without an approval workflow. | Pending real-client decision |
| OOS-005 | Invoice management | Payment tracking is committed P0 scope (ACC-011); invoice lifecycle is a separate finance process. | Pending real-client decision |
| OOS-006 | Email and calendar synchronization | Not required for the committed CRM loop; depends on external email/calendar integration. | Pending real-client decision |
| OOS-007 | Advanced analytics, forecasting, and sales performance reporting | Basic manager overview is committed P1 (ACC-018) and basic reports are committed P1 (ACC-023); advanced analytics is a separate capability outside committed scope. | Pending real-client decision |
| OOS-008 | External integrations with Feishu, DingTalk, WeCom, ERP, or finance systems | Integration targets and contracts are not defined; not part of the committed CRM loop. | Pending real-client decision |
| OOS-009 | AI sales summaries, next-step suggestions, and risk hints | Not required for committed CRM production validity; depends on undefined AI capability and data rules. | Pending real-client decision |
| OOS-010 | Dedicated mobile app | The committed delivery channel is the responsive web client; a separate native mobile app is not in committed scope. | Pending real-client decision |
| OOS-011 | Multi-tenant SaaS organization management | Single-team collaboration is the confirmed working assumption (DEC-009, OQ-015); multi-tenant SaaS is not in committed scope. | Pending real-client decision |
| OOS-012 | Complex contract version management | Committed contract management records core contract information and history; full version management is a separate capability not in committed scope. | Pending real-client decision |
| OOS-013 | Multi-currency, tax, and discount automation | The committed amount model is single-currency (DEC-013, OQ-023); multi-currency, tax, and discount automation require finance/business rules and are not in committed scope. | Pending real-client decision |
| OOS-014 | Automatic lead assignment rules | Manual assignment is committed P0 (ACC-003); automatic assignment rules are a separate automation capability not in committed scope. | Pending real-client decision |
| OOS-015 | XLSX import/export | CSV is the committed P1 import/export format (DEC-015, OQ-011); XLSX is not in committed scope. | Pending real-client decision |
| OOS-016 | Email, SMS, or chat reminder delivery | In-app reminders are the committed P1 reminder channel (DEC-015, OQ-012); email, SMS, or chat delivery is not in committed scope. | Pending real-client decision |
| OOS-017 | Contract attachment upload as a required P0 path | Contract notes are committed P0 (DEC-016, OQ-009); contract attachment upload is not a committed P0 path. | Pending real-client decision |
