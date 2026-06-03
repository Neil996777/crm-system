# G12 Independent FINAL RE-AUDIT — Decision [3rd KICKBACK: single item]

## Document Control

- Project: CRM System
- Gate: G12 — Audit → Release/Rework (final focused re-audit after rework #2)
- Owner: Audit (Claude), independent of execution
- Predecessors: `g12-audit-decision-2026-06-03.md` (1st audit, REWORK),
  `g12-reaudit-2026-06-03.md` (re-audit, 2nd REWORK)
- Method: Four parallel author≠auditor passes (BLK-G12-001 tests; BLK-G12-006 residuals;
  BLK-G12-009/010; regression) + an independent compile/test run. Read-only.
- Date: 2026-06-03
- Decision: **REWORK (3rd, single item) — Gate NOT yet Passed.** All four kicked-back items
  verified FIXED; one new consistency defect (BLK-G12-011) elevated to a required fix per the
  release owner's decision; G12 passes after it is closed and spot-verified.

## Rework #2 items — all verified FIXED

| ID | Verdict | Basis |
|---|---|---|
| BLK-G12-001 | **FIXED** | commercial/work/account each now have real testcontainer dispatcher tests for delivery-to-audit-history + failure-retry + UID dedup (4/4 services; opportunity already had them). Consumer-side idempotency proven in audit-history. |
| BLK-G12-006 | **FIXED** | Restore counts reconciled (10/10 matching `\dn`/`\du`); distinct per-operator SSH key fingerprints; `sshd -T` proves `passwordauthentication no` + `kbdinteractive no` (`permitrootlogin without-password`, the OpenSSH synonym of `prohibit-password`); negative probe re-run from genuinely external `srv-aliyun-bj-01` (47.95.119.211). |
| BLK-G12-009 | **FIXED** | Checker now requires the four mutating-call RequestIds (http=200) under `--apply`; a fail-first test confirms the old 20:29 read-only file FAILS as `--apply` and a proper apply fixture PASSES; read-only exports accepted only as `--verification`. |
| BLK-G12-010 | **FIXED** | `X-Correlation-Id` propagated on reporting delivery for all four producers; Administrator all-scope report test `TEST-BASIC-REPORT-006` added. |

## Independent compile/test run

11/11 Go modules (10 services + shared/contracts) build + vet + test exit 0; **203 test
functions, 0 skips**, integration tests against real `postgres:16-alpine` testcontainers
(Docker available, `GOPROXY=off`, `-count=1`). Frontend `tsc --noEmit` + `npm run build`
green. Playwright e2e not run (needs live stack) — honestly reported as unverified.
Regression: 3 retired tests still absent; no new fakes; no test weakening; previously-fixed
items (reporting S2S, transactional outbox in the four services, append-only audit, no
cross-service DB, post-scope-change model) all intact.

## New finding — BLK-G12-011 (lead non-transactional outbox)

The regression auditor surfaced that the **`lead` service alone** still uses a
fire-and-forget outbox: `_ = h.outbox.Append(...)` is called after `repo.Create` commits and
its error is discarded (`lead_command.go:115/153/207`, `lead_qualify.go:78`,
`lead_convert.go:85`, `duplicate_check.go:43`). This is the SAME class of defect that
BLK-G12-004 required fixed for opportunity/commercial/account/work. If the Append fails, the
lead mutation commits with no outbox row, so `EVT-LEAD-QUALIFIED/DISQUALIFIED/CONVERTED` and
lead `EVT-OWNER-CHANGED` can be lost with no row to retry — an AUD-IMM-002 / ACC-014/022 gap.

It is **pre-existing** (TASK-032), not introduced by rework #2, and Claude had marked the lead
consistency item **Optional** in rework #2. Because the same bar was enforced on four other
services, the release owner elected (2026-06-03) to **elevate it to a required 3rd-kickback
fix** rather than accept it as a residual. Consistency + no-downgrade support this.

## Decision

**Gate NOT yet Passed — 3rd REWORK (single item).** G11 remains `Gate Blocked`. Codex applies
the same transactional-outbox pattern to the lead service (`delivery/G12-rework-3.md`), adds a
fail-first rollback test, and returns; Claude does a minimal spot re-audit of the lead change
(plus a quick build/test confirmation) and, if clean, **passes G12** with no other open items.
Everything else is verified release-ready: the codebase builds and tests green, the security
posture is independently verified, and all prior findings are closed. No-downgrade applies.
