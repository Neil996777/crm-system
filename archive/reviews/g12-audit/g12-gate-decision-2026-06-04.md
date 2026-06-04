# G12 Gate Decision — CRM System: **PASS (Audit cleared for release)**

## Document Control

- Project: CRM System
- Gate: G12 — Audit → Release/Rework
- Owner: Audit (Claude), independent of execution (author ≠ auditor)
- Decision date: 2026-06-04
- Decision: **GATE PASSED.** Independent audit clears the implementation for release.
  Whether to perform the production go-live is the release owner's (user's) call; the
  audit gate no longer blocks it.

## Audit trail (six rounds, author ≠ auditor throughout)

G12 was not a single pass — it was a six-round adversarial audit→rework cycle in which every
finding was registered on disk, kicked back to Codex with no-downgrade + no-over-reach
constraints, remediated, and independently re-verified before the next step.

| Round | Method | Outcome |
|---|---|---|
| 1 | 5 parallel author≠auditor passes | REWORK — BLK-G12-001..008 (audit events never delivered; reporting ingest unauthenticated; security-group evidence fabricated/contradicted; +5) |
| 2 | RE-AUDIT (5 passes + independent compile/test + Claude read-only live Volcengine query) | REWORK — 6/8 closed; BLK-G12-001/006 partial; +009/010 |
| 3 | FINAL re-audit (4 passes + compile/test: 11/11, 203 tests, 0 skips) | REWORK — 001/006/009/010 fixed; surfaced BLK-G12-011 (lead non-transactional outbox) |
| 4 | spot re-audit (3 passes + build/test) | micro-REWORK — 011 fixed; +012/013/014 (lead disqualify event / post-commit audit window / e2e delivery test) |
| 5 | **whole-repo systematic sweep (7 dimensions, NOT scoped to recent changes)** | REWORK — found release-blocking defects in never-audited code: BLK-G12-015..025 (2 IDOR BLOCKERs, 7 MAJOR, 2 batch) |
| 6 | full whole-repo RE-SWEEP (6 dimensions) verifying #5 | REWORK (consolidated tail) — all 11 of #5 FIXED, no over-reach, no downgrade; drained residual BLK-G12-026..030 |
| **PASS** | spot re-audit of BLK-G12-026..030 (3 agents) + independent build/test with Docker live | **all 5 FIXED; 11/11 modules build/vet/test green on real Postgres; frontend green** |

The round-5 systematic sweep was pivotal: the prior four rounds had only audited "what Codex just
changed," so code never touched by a fix was never inspected. The whole-repo sweep caught genuine
IDOR data-exposure holes (commercial by-id contract/quote reads, import-export import-run reads)
and an unreliable identity-authz audit path that the narrow rounds structurally could not have
found.

## Final verification (the passing round)

**All BLK-G12-026..030 FIXED** (`archive/reviews/g12-audit/g12-resweep-2026-06-04.md` opened them;
this round closed them):

- BLK-G12-026 (MAJOR) — identity-authz denial-path audit (login-failed, access-denied,
  unauthenticated, invalid-session, inactive-user, authz-version-stale, user-admin-denied) is now
  durable: each appends in a committed transaction and returns `503 DEPENDENCY_UNAVAILABLE` on
  persist failure instead of `log.Printf`-swallowing; the BLK-G12-017 dispatcher delivers them
  with retry retention. Fail-first `TestDenialAuditOutboxFailureIsSurfaced` /
  `TestPermissionDeniedAuditOutboxFailureIsSurfaced` ran and passed. Zero swallowed audit appends
  remain in the service. The 401/403 client contract is intact (503 only when audit truly cannot
  persist).
- BLK-G12-027 (MINOR) — dead `identity-authz/internal/authz/audit_client.go` deleted (−103);
  no production reference to `NewAuditClient`/`AuditClient`/`AppendOperationLog`; no
  `EVT-USER-ADMIN-CHANGED` remains repo-wide; cleanup guard asserts absence.
- BLK-G12-028 (MINOR) — all 8 catalog events (LEAD-CONVERTED, QUOTE-ACCEPTED, PAYMENT-RECORDED,
  OPPORTUNITY-WON, OPPORTUNITY-LOST, RECORD-ARCHIVED, IMPORT-RUN, EXPORT-RUN) now have a test
  pinning the derived catalog `eventId`; pure test additions (0 deletions), no production mapping
  changed.
- BLK-G12-029 (MINOR) — `account` contacts sub-resource (`contact_query.go`) returns safe
  `404 NOT_FOUND` for an existing-but-unreadable account, consistent with the Denial Contract;
  `TEST-DENIAL-CONTACTS-001` proves no existence leak.
- BLK-G12-030 (LOW) — (a) work S2S verifier aligned to `ServiceTokenClaims`/`ErrServiceAuthFailed`
  in-place; (b) concurrent lead-conversion duplicate-key now returns 200-with-existing-row (no
  duplicate), with race tests; (c) account duplicate-warning made single-transaction with an
  atomic-rollback test.

**Constraints honored:** no over-reach (the single rework-6 commit `14d0613` touched only
identity-authz/account/opportunity/work + tests; `shared/` byte-unchanged; no shared-library
refactor; `shared/contracts` preserved; the only go.mod change is a benign indirect→direct
testcontainers promotion in work). No downgrade (catalog tests are pure additions; 3 retired tests
remain absent; 0 skips; no test weakened).

**Independent build/test (this round, Docker live):** all 11 Go modules `go build` + `go vet` +
`go test ./... -count=1` pass against real `postgres:16-alpine` testcontainers (DB-backed suites
actually executed, not skipped); frontend `tsc --noEmit` + `vite build` pass; go.mod/go.sum
unchanged; no regressions.

## Release-evidence posture (carried, independently verified)

- Production HTTPS endpoint `https://118.196.44.193` (Let's Encrypt IP cert), HTTP→HTTPS redirect,
  renewal timer — BLK-G11-001 Resolved.
- Volcengine security group: CRM ENI bound only to dedicated `sg-366ptx1bxp9ts1e710babmc8y`
  (public TCP 22/80/443 only; no 8080/5432/8088/8443/3389) — **independently verified by Claude's
  own read-only live API query** (`docs/release/evidence/volcengine-security-group-verified-readonly-2026-06-03.json`),
  not by the contradicted earlier files. BLK-G12-003/007/BLK-G11-002 Resolved.
- Encrypted off-server backup + restore rehearsal + daily backup timer — TASK-040 Resolved.

## Decision

**G12 GATE PASSED.** The CRM implementation has passed independent audit. All P0/P1 acceptance
items are satisfied with no downgrade; all registered blockers BLK-G12-001..030 are Resolved with
cited artifacts; the build and tests are independently green on real infrastructure; the
production release-evidence (HTTPS, dedicated security group, off-server backup) is independently
verified.

The audit gate no longer blocks release. The production go-live decision is the release owner's.

- Audit decision: this file.
- Final re-audit detail: `archive/reviews/g12-audit/g12-resweep-2026-06-04.md` (residual tail) +
  this round's spot re-audit of BLK-G12-026..030.
- Blocker register: `planning/blockers.md` (all Resolved).
- Gate sync: `planning/gate-status.md` (G11/G12 Gate Passed).
