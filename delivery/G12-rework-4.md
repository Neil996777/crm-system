# G12 Rework Package #4 — CRM System (Claude → Codex, micro follow-ups)

## Document Control

- Source: G12 BLK-G12-011 spot re-audit, 2026-06-04 (`archive/reviews/g12-audit/g12-reaudit-final-2026-06-03.md`)
- Decision: 3 non-blocking MINOR follow-ups; release owner elected to fix them before G12 passes.
- Status going in: BLK-G12-001..011 all Resolved/verified; build+tests green (11/11 modules);
  security posture verified. These three are the only remaining items.
- Executor: Codex; then return to Claude for a minimal spot re-audit → G12 pass.

## Rules (unchanged)

TDD fail-first; no weakening tests; no-downgrade. On completion of each item: update
`delivery/tasks.md`, `modeling/traceability-matrix.md`, set the `planning/blockers.md` row to
`Resolved` with the deciding artifact, and commit.

---

## 🟡 BLK-G12-012 — Distinct event for lead disqualify

**Where:** `services/lead/internal/handler/lead_qualify.go` (valid vs invalid path),
`services/lead/internal/event/outbox.go` (event-type → `EVT-*` mapping), and the event
constants.

**What's wrong:** qualify (valid) and disqualify (invalid) both emit `event.LeadQualified`,
distinguished only by payload. So `EVT-LEAD-DISQUALIFIED` (a distinct catalog event in
`docs/security/audit-log-spec.md`) never appears with its own `event_id` in the audit log.

**Required fix:** emit a distinct event for the invalid/disqualify transition that maps to
`EVT-LEAD-DISQUALIFIED`; keep the valid transition mapping to `EVT-LEAD-QUALIFIED`. Preserve
the transactional-outbox behavior from BLK-G12-011 (same transaction, propagated error).

**Acceptance (TDD):** a test asserts the disqualify path produces an event resolving to
`EVT-LEAD-DISQUALIFIED` (with the invalid reason), and qualify still produces
`EVT-LEAD-QUALIFIED`.

---

## 🟡 BLK-G12-013 — Remove the non-transactional post-commit audit window in lead qualify

**Where:** `services/lead/internal/handler/lead_qualify.go` (~97-112), the post-commit
`audit.AppendRecordHistory` call (502 on failure).

**First determine** lead's audit-history delivery path:
- If lead's audit-history events are delivered by the transactional outbox dispatcher (like
  the other services after BLK-G12-001/011), then this post-commit `AppendRecordHistory` is
  REDUNDANT and reintroduces a non-transactional commit-then-audit-fail window — **remove it**.
- If the post-commit call is in fact lead's ONLY audit-history delivery path (i.e. lead's
  outbox dispatches only to reporting, not audit-history), then **route lead's audit-history
  event through the transactional outbox** (same dispatcher pattern as the other services), so
  no post-commit non-transactional audit window remains.

**Why it matters:** BLK-G12-011 made the lead *outbox enqueue* transactional, but if the real
audit-history delivery still rides a post-commit non-transactional call, the AUD-IMM-002
"same durable workflow" guarantee for lead audit events is still not fully met. Resolve which
path is authoritative and make it transactional + retryable.

**Acceptance (TDD):** no post-commit non-transactional audit call remains on the lead qualify
path; a test proves the lead qualify/disqualify audit event is delivered durably (in/after the
transaction via the outbox dispatcher, retryable on failure, not lost on a post-commit error).

---

## 🟡 BLK-G12-014 — Lead end-to-end outbox→reporting delivery test

**Where:** `services/lead/internal/event/` (add a dispatcher test).

**What's wrong:** lead's transactional enqueue and reporting's ingest are each tested, but no
test exercises the lead-side dispatcher poll→POST to reporting end-to-end (the other four
services have this).

**Required fix (TDD):** add a dispatcher test modeled on
`services/opportunity/internal/event/dispatcher_test.go` asserting lead's `DispatchOnce`
delivers the projection to reporting `/internal/projections` with the correct signed S2S
headers + `X-Correlation-Id`, and that delivery failure leaves the row unpublished for retry.

**Acceptance:** real-Postgres-testcontainer test green; lead delivery matrix matches the other
four services.

---

## Definition of Done

- BLK-G12-012/013/014 fixed with the named tests green; no `_ = ...Append` or post-commit
  non-transactional audit window reintroduced.
- tasks.md / traceability-matrix.md / blockers.md updated to real artifacts; commits made.
- Return to Claude for a minimal spot re-audit (these three + a quick build/test confirm).
  If clean, Claude passes G12. Do not self-pass G12.
