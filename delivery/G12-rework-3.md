# G12 Rework Package #3 — CRM System (Claude → Codex, 3rd KICKBACK — single item)

## Document Control

- Source: G12 FINAL re-audit, 2026-06-03 (`archive/reviews/g12-audit/g12-reaudit-final-2026-06-03.md`)
- Decision: **REWORK (3rd, single item)** — G11 stays `Gate Blocked`. Do NOT release.
- Status going in: ALL other G12 findings (BLK-G12-001..010) are closed and independently
  verified; build + tests are green (11/11 modules, 203 tests, 0 skips, real Postgres);
  security posture verified. **This is the last open item.**
- Executor: Codex; then return to Claude for a minimal spot re-audit → G12 pass.

## The only item — BLK-G12-011: make the `lead` outbox transactional

**Where:** `services/lead/internal/handler/` — `lead_command.go` (lines ~115, 153, 207),
`lead_qualify.go` (~78), `lead_convert.go` (~85), `duplicate_check.go` (~43).

**What's wrong:** lead uses a fire-and-forget outbox — `_ = h.outbox.Append(...)` is called
AFTER `repo.Create` commits, and the error is discarded. If the Append fails, the lead
mutation is already committed with no outbox row, so the required `EVT-LEAD-QUALIFIED`,
`EVT-LEAD-DISQUALIFIED`, `EVT-LEAD-CONVERTED`, and lead `EVT-OWNER-CHANGED` events can be
permanently lost (no row for the dispatcher to retry). Violates AUD-IMM-002 "same durable
workflow"; touches ACC-014/022 (P0). This is the same defect BLK-G12-004 fixed for the other
four services.

**Required fix (TDD):** Apply the **same `inTransaction` + `txOutbox.Append` pattern** already
used in opportunity/commercial/account/work:
1. Perform the lead business write and the outbox enqueue in ONE DB transaction.
2. Do NOT discard the Append error — propagate it; on failure roll back the mutation and
   return an error (e.g. 503), so a success is never returned without the enqueued event.
3. Cover every cited call site (create, qualify/disqualify, convert, owner-change, and the
   duplicate-check create path).

**Acceptance / verify:**
- Add a fail-first rollback test (real Postgres testcontainer) modeled on
  `services/opportunity/internal/handler/opportunity_command_test.go` `TEST-HISTORY-TX-001`:
  force the outbox insert to fail → assert the lead mutation row is ABSENT and an error is
  returned (you can never observe a persisted lead change without its enqueued event).
- Confirm no remaining `_ = h.outbox.Append` in `services/lead/internal/handler/`.
- `go test ./... -count=1` green for `services/lead` (real Postgres).

## Definition of Done

- lead writes are transactional with their outbox enqueue; no discarded Append errors.
- Fail-first rollback test green; no `_ = ...Append` left in lead handlers.
- `delivery/tasks.md`, `modeling/traceability-matrix.md`, `planning/blockers.md` (BLK-G12-011
  → Resolved with the deciding test) updated; commit.
- Return to Claude for a minimal spot re-audit of the lead change + a build/test confirmation.
  Do not self-pass G12.
