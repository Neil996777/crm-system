# Gate Status — CRM System

> ONBOARDING NOTE (2026-05-30): Added when the project was brought under the
> platform collaboration model. Statuses below reflect PROJECT_CONTEXT.md as of
> 2026-05-30 (G5 passed, G6 current). Confirm and maintain going forward.

This file is the single synchronization point between Claude (planning + audit)
and Codex (execution). Update it on every gate change, handoff, and kickback.
See `../../../../company/collaboration-model.md`.

Status values: `Gate Draft` / `Gate Review` / `Gate Blocked` / `Gate Passed`.

| Gate | Transition | Owner | Platform | Status | Reviewer sign-off | Date | Current blocker |
|---|---|---|---|---|---|---|---|
| G1 | Idea -> Requirement Discussion | Product Manager | Claude | Gate Passed | | | |
| G2 | Requirement Discussion -> PRD | Product Manager | Claude | Gate Passed | | | |
| G3 | PRD -> Acceptance Matrix | Product Manager | Claude | Gate Passed | | | |
| G4a | Acceptance Matrix -> Business Design | Business Analyst | Claude | Gate Passed | | | |
| G4b | Business Design -> UX Design | UX Designer | Claude | Gate Passed | | | |
| G4c | UX Design -> UI Design | UI Designer | Claude | Gate Passed | | | |
| G4d | UI Design -> Security Design (Design Closure) | Security Compliance | Claude | Gate Passed | PM, BA, UX, UI (retroactive closure) | 2026-05-30 | |
| G5 | Design Closure -> Architecture Design | Architecture | Claude | Gate Passed | all required reviewers (incl. Infrastructure Ops) | 2026-05-30 | |
| G6 | Architecture Design -> MDA Modeling | Domain Modeling + Architecture | Claude | Gate Passed | PM, BA, UX, UI, Security, QA Test Design (all signed off 2026-06-01) | 2026-06-01 | MDA package (CIM/PIM/PSM/Traceability/Test Model) complete, multi-agent audited, and signed off by all six reviewer roles. Formal Scope Change by User 2026-06-01 (DEC-017..020) applied across baseline → architecture → MDA → UX/UI → security and re-audited; BLK-001/002/003 RESOLVED. Decision: `archive/reviews/g6-mda/g6-mda-gate-decision-2026-06-01.md`. |
| G7 | MDA + Test Model -> Task Planning | Domain Modeling + QA Test Design | Claude | Gate Passed | PM, Architecture, Security, QA Test Design, Infrastructure Ops, Task Planner (all signed off 2026-06-01) | 2026-06-01 | Acceptance-driven delivery plan in `delivery/` (40 tasks, ACC 23/23, Codex-executable, multi-agent audited). Decision: `archive/reviews/g7-task-planning/g7-gate-decision-2026-06-01.md`. Pre-G8 entry condition: Security Compliance review of operator-access (deployment-notes), recorded on TASK-039. |
| G8 | Task Planning -> Implementation **[HANDOFF: Claude -> Codex]** | Task Planner | Claude -> Codex | Gate Passed | Task Planner, Infrastructure Ops, Security Compliance, Audit (signed off 2026-06-01) | 2026-06-01 | Self-contained execution handoff package on disk (`delivery/G8-handoff.md` + delivery plan). Pre-G8 operator-access Security review done (approved w/ conditions). HANDED OFF to Codex for G9–G11; Claude resumes at G12. Decision: `archive/reviews/g8-handoff/g8-gate-decision-2026-06-01.md`. |
| G9 | Implementation -> QA | Frontend / Backend Engineer | Codex | Gate Passed | TASK-001..038 implementation tasks completed with recorded task status, traceability, and commits | 2026-06-03 | |
| G10 | QA -> Integration | QA Execution | Codex | Gate Passed | Automated and E2E evidence recorded per task; release smoke and backup checks completed for TASK-039/040 | 2026-06-03 | |
| G11 | Integration -> Audit **[RETURN: Codex -> Claude]** | Integration Owner | Codex -> Claude | **Gate Passed** | Cleared by G12 independent audit (six rounds) | 2026-06-04 | All G12 blockers BLK-G12-001..030 Resolved; integration accepted into audit and audit passed. |
| G12 | Audit -> Release/Rework | Audit | Claude | **Gate Passed** | Audit (independent): 6-round audit→rework cycle (1st audit + re-audit + final re-audit + spot re-audit + systematic sweep + full RE-SWEEP + passing spot re-audit) | 2026-06-04 | **GATE PASSED.** Final spot re-audit (3 agents) confirmed BLK-G12-026..030 all FIXED; no over-reach, no downgrade; independent build/test GREEN with Docker live (11/11 Go modules build/vet/test on real postgres:16-alpine; frontend tsc+build). All BLK-G12-001..030 Resolved with cited artifacts; release-evidence (HTTPS `118.196.44.193`, dedicated SG `sg-366ptx...` verified by Claude read-only live API, off-server backup+restore) independently verified. Audit gate no longer blocks release; production go-live is the release owner's call. Decision: `archive/reviews/g12-audit/g12-gate-decision-2026-06-04.md`. |

## Handoff Log

| Date | From platform | To platform | Gate | Note |
|---|---|---|---|---|
| 2026-06-01 | Claude (planning) | Codex (execution) | G8 | Task planning complete and gate-passed; self-contained execution handoff package delivered (`delivery/G8-handoff.md`). Codex executes G9–G11; Claude resumes for independent G12 audit. |
| 2026-06-03 | Codex (execution/integration) | Claude (audit) | G11 | TASK-001..040 complete on disk. Runtime deployment, HTTPS/TLS, security-group cleanup, encrypted off-server backup, daily backup timer, and restore rehearsal evidence recorded. Claude resumes for mandatory G12 audit. |
| 2026-06-03 | Claude (audit) | Codex (execution) | G12 → G11 **KICKBACK** | G12 independent audit (5 parallel author≠auditor passes) found release-blocking defects: audit events not delivered to audit-history (BLK-G12-001); reporting internal write unauthenticated (BLK-G12-002); security-group closure unsubstantiated and contradicts raw API export (BLK-G12-003, reopens BLK-G11-002); + 4 MAJOR + 1 MINOR. Decision REWORK, no release. Package: `delivery/G12-rework.md`. Decision: `archive/reviews/g12-audit/g12-audit-decision-2026-06-03.md`. |
| 2026-06-03 | Codex (execution) | Claude (audit) | G11 → G12 (rework #1 return) | Codex remediated BLK-G12-001..008 with commits/evidence and returned for re-audit. |
| 2026-06-03 | Claude (audit) | Codex (execution) | G12 → G11 **2nd KICKBACK** | RE-AUDIT (5 parallel passes + independent compile/test run + Claude-run read-only live Volcengine query) closed 6/8: build/test GREEN; security-group state independently verified correct (BLK-G12-003/007/BLK-G11-002 Resolved). Codex's 20:29 SG "success" file was read-only-only (no mutating RequestIds) and rejected as provenance; END STATE confirmed via Claude's live read. Still Open: BLK-G12-001 (per-service audit-delivery tests), BLK-G12-006 (restore counts / per-operator SSH keys / sshd -T / external probe), + new BLK-G12-009 (evidence-checker hardening), BLK-G12-010 (corr-id/admin test). Decision 2nd REWORK. Package: `delivery/G12-rework-2.md`. Decision: `archive/reviews/g12-audit/g12-reaudit-2026-06-03.md`. |
| 2026-06-03 | Codex (execution) | Claude (audit) | G11 → G12 (rework #2 return) | Codex remediated BLK-G12-001, BLK-G12-006, BLK-G12-009, and BLK-G12-010 with commits `662f6e9`, `77f5f2a`, `789c72e`, and `45e4084`. G11/G12 remain Gate Blocked pending Claude final focused G12 re-audit; Codex does not self-pass G12. |
| 2026-06-03 | Claude (audit) | Codex (execution) | G12 → G11 **3rd KICKBACK (single item)** | FINAL re-audit (4 parallel passes + independent compile/test: 11/11 modules, 203 tests, 0 skips, real Postgres) confirmed BLK-G12-001/006/009/010 genuinely FIXED and surfaced BLK-G12-011 (`lead` non-transactional outbox — same class as BLK-G12-004). Release owner elected to fix it before pass. Package: `delivery/G12-rework-3.md`. Decision: `archive/reviews/g12-audit/g12-reaudit-final-2026-06-03.md`. G12 passes after a minimal spot re-audit. |
| 2026-06-03 | Codex (execution) | Claude (audit) | G11 → G12 (rework #3 return) | Codex made the lead outbox transactional (commit `b0a9855`) with fail-first rollback test `TEST-HISTORY-TX-001`; returned for spot re-audit. |
| 2026-06-04 | Claude (audit) | Codex (execution) | G12 → G11 **4th micro-KICKBACK** | Spot re-audit (3 parallel passes + independent build/test 11/11 green) confirmed BLK-G12-011 FIXED. 3 non-blocking MINOR follow-ups surfaced and elected for fix: BLK-G12-012 (disqualify distinct event id), BLK-G12-013 (lead audit-history still has a post-commit non-transactional window — make transactional), BLK-G12-014 (lead→reporting e2e delivery test). Package: `delivery/G12-rework-4.md`. G12 passes after minimal spot re-audit. |
| 2026-06-04 | Codex (execution) | Claude (audit) | G11 → G12 (rework #4 return) | Codex resolved BLK-G12-012/013/014 with lead distinct disqualify event, lead transactional outbox audit-history dispatch, removal of post-commit audit call, and real PostgreSQL lead dispatcher reporting/audit retry tests. Returned for minimal spot re-audit; Codex does not self-pass G12. |
| 2026-06-04 | Claude (audit) | Codex (execution) | G12 → G11 **5th KICKBACK (systematic sweep)** | Instead of passing on the narrow spot re-audit, ran a whole-repo systematic sweep (7 dimensions, not scoped to recent changes). Found release-blocking defects in never-audited code: commercial + import-export IDOR (BLK-G12-015/016, BLOCKER, real data exposure); identity-authz non-durable audit (017); work missing optimistic concurrency on ACC-012 (018); lead-conversion idempotency (019); 7 sync clients missing correlation-id (020); 8/24 EVT catalog events unemitted (021); account/commercial reporting delivery untested (022); api-spec close-won schema (023); + MINOR/LOW batches (024/025). Codex constraints recorded: no downgrade AND no over-reach. Package: `delivery/G12-rework-5.md`. Decision: `archive/reviews/g12-audit/g12-systematic-sweep-2026-06-04.md`. |
| 2026-06-04 | Codex (execution) | Claude (audit) | G11 → G12 (rework #5 return) | Codex remediated BLK-G12-015..025 (10 commits `3e7d0a8`..`5f34ae5`): commercial/import-export by-id IDOR fixed; identity-authz mutation audit made transactional + dispatcher wired; work optimistic concurrency; lead-conversion downstream idempotency (migration 0004); 7 sync clients correlation-id; 8 EVT catalog ids; account/commercial reporting dispatcher tests; close-won api-spec; denial/error uniformity; dead-code cleanup. `shared/contracts` RAISED not deleted. Returned for the full systematic re-sweep; did not self-pass G12. |
| 2026-06-04 | Claude (audit) | Codex (execution) | G12 → G11 **6th KICKBACK (consolidated tail)** | Full whole-repo RE-SWEEP (6 parallel author≠auditor agents, same systematic discipline) confirmed all 11 (BLK-G12-015..025) genuinely FIXED, **no over-reach** (`shared/` byte-unchanged, no new shared lib, contracts raised not deleted, 2 judgment calls surfaced; traceability 6/6 honest) and **no downgrade** (0 skips, 3 retired tests absent, only 403→404 strengthenings); 11/11 build/vet + frontend green. Drained a residual tail BLK-G12-026..030: identity-authz denial-path audit durability (MAJOR), dead audit-client carrying stale `EVT-USER-ADMIN-CHANGED` (MINOR), 8 unpinned EVT-id tests (MINOR), contacts sub-resource existence leak (MINOR), consistency/robustness batch (LOW). None release-blocking; release owner elected one consolidated micro-rework before pass. Package: `delivery/G12-rework-6.md`. Decision: `archive/reviews/g12-audit/g12-resweep-2026-06-04.md`. G12 passes after a minimal spot re-audit of these five. |
| 2026-06-04 | Codex (execution) | Claude (audit) | G11 → G12 (rework #6 return) | Codex resolved BLK-G12-026..030 in single commit `14d0613`: identity-authz denial-path audits made durable (committed-tx + 503 on persist failure, dispatcher retry), dead `authz/audit_client.go` deleted, 8 EVT catalog ids pinned in tests, account contacts return safe 404, work verifier shape aligned, concurrent lead-conversion duplicate-key returns 200-existing, duplicate-warning made single-transaction. Returned for spot re-audit; did not self-pass G12. |
| 2026-06-04 | Claude (audit) | — | G12 **GATE PASSED** | Spot re-audit (3 parallel author≠auditor agents) confirmed BLK-G12-026..030 all FIXED with file:line evidence; no over-reach (single commit touched only identity-authz/account/opportunity/work + tests; `shared/` byte-unchanged; only go.mod change a benign indirect→direct testcontainers promotion); no downgrade (catalog tests pure additions, 3 retired tests absent, 0 skips). Independent build/test with Docker live: 11/11 Go modules build/vet/test GREEN on real postgres:16-alpine (DB suites actually executed); frontend tsc+build GREEN; go.mod/go.sum unchanged; no regressions. All BLK-G12-001..030 Resolved. **G12 PASSED — audit clears the implementation for release.** Production go-live remains the release owner's decision. Decision: `archive/reviews/g12-audit/g12-gate-decision-2026-06-04.md`. |
| 2026-06-05 | Claude (audit) | Codex (execution) | Post-G12 **RELEASE EXECUTION HANDOFF** | Release owner authorized production go-live. Per the platform contract (deployment/integration is Codex's G9–G11 role, not the audit platform's), go-live execution is handed to Codex per `deploy/ops/go-live-runbook.md` against the audited commit `da9d63c`. A fresh deploy is required (G11 runtime predates the G12 rework; read-only probe 2026-06-05 found the endpoint not serving). Claude precondition checks done: local static deploy smoke PASSED. Codex executes (backup → frontend build → `docker compose up -d --build` → `migrate.sh up` → nginx reload), records ACC-017 release evidence, and returns to Claude for independent read-only post-deploy verification. Codex must not modify application code during go-live; any deploy-blocking issue is raised as a blocker, not worked around. |
| 2026-06-05 | Codex (execution) | Claude (audit) | Go-live return (evidence) | Codex resolved 3 pre-flight blockers (BLK-GOLIVE-001 SSH routing, 002 `.env` perms, 003 target not a git repo), checked out `da9d63c`, and recorded go-live evidence (`docs/release/evidence/go-live-2026-06-05-transcript.txt`, commit `e46b608`); smoke 001/002 reported passing. |
| 2026-06-05 | Claude (audit) | Codex (execution) | Post-deploy verification **NOT VERIFIED — KICKBACK** | Independent live probes confirm the endpoint is online and hardened (HTTPS 200, HTTP→301, 8080/5432 publicly unreachable). BUT cross-checking Codex's own transcript: `docker compose ps` shows all 9 services running image `crm-system/<svc>:prod-20260603` (built 2026-06-03) — the **pre-G12-rework build**; the transcript has NO `up -d --build`, NO `migrate.sh up`, NO frontend build. So the audited fixes (IDOR BLK-G12-015/016, durable audit 017/026, etc., all committed 2026-06-04+) are NOT in production and the new migrations were never applied; green smoke only checks online/HTTPS/redirect and cannot detect stale code. Registered `BLK-GOLIVE-004`. Codex must prove `da9d63c` is the running build or actually `up -d --build` + `migrate.sh up` and return. Package: `delivery/go-live-rework-1.md`. Decision: `archive/reviews/g12-audit/post-deploy-verification-2026-06-05.md`. (Also: TLS cert expires 2026-06-09 — confirm renewal.) |

## Notes

- G8 may not pass until the self-contained G8 execution handoff package exists on
  disk (see `company/collaboration-model.md`).
- On a kickback, set the affected gate back to `Gate Blocked`, record it here and
  in `planning/blockers.md`, and return to Claude.
- G6 MDA progress (2026-06-01): CIM and PIM authored by the Domain Modeling role
  and each passed independent multi-agent audit (author ≠ reviewer), including a
  dedicated tier-altitude (CIM/PIM/PSM boundary) pass. The audits surfaced three
  upstream-source gaps that MDA correctly declined to invent; they are registered
  in `planning/blockers.md` as BLK-001 (Opportunity Status enumeration, ACC-007),
  BLK-002 (multi-plan contract full-payment aggregation, ACC-011/013), and BLK-003
  (second-quote-accept observable outcome, ACC-009). Per the no-downgrade rule
  these P0-touching items must be resolved by PM/BA (or formally scope-changed by
  the user) before the G8 handoff. BLK-A01 (overdue-evaluation trigger) is
  PSM/Architecture-deferred for G7/PSM, not a PM/BA blocker.
- Formal Scope Change by User (2026-06-01): the owner revised four committed P0
  rules — DEC-017 (Won = related contract Signed, not full payment; `Contract
  Signed`/`Payment In Progress` opportunity stages removed), DEC-018 (exactly one
  quote per opportunity), DEC-019 (payment tracking retained but decoupled from
  Won), DEC-020 (Opportunity `Status` field removed). Recorded in
  `decision-log.md` (originals retained + annotated). Cascade applied and
  re-audited (author ≠ reviewer) across: G3/G4 baseline (prd, requirements,
  business-rules, acceptance-matrix, edge-cases, business-glossary,
  business-processes, user-scenarios, business-capability-map,
  service-governance-inputs, open-questions), G5 architecture (consistency
  reconciliation only — service decomposition unchanged: api-spec, architecture,
  integration-design, data-design, module-boundaries), and the full G6 MDA
  (CIM/PIM/PSM/Traceability/Test Model; affected IDs retired in place, not
  renumbered). This resolved BLK-001/002/003 (see `planning/blockers.md`
  Resolution Log). No P0/P1 capability dropped (payment tracking retained).
  Process tracker: `planning/scope-change-2026-06-01-TEMP.md`.
- G12 audit is mandatory and performed on Claude before any release decision.
- Carried-forward release blockers (not current gate blockers): encrypted
  off-server backup copy + restore rehearsal, HTTPS/TLS endpoint, security-group,
  and monitoring evidence.
- G1–G4c are the retained design baseline: produced/passed under the pre-split
  combined G4, retained through the 2026-05-29 reset, supplemented for
  service-boundary governance, and re-validated by the G4d retroactive closure
  check and the pre-G6 design audit.
- Evidence pointers:
  - G4d closure: `archive/reviews/g4-design-closure/g4d-design-closure-decision-2026-05-30.md`
  - G5 final decision: `archive/reviews/g5-architecture/g5-architecture-final-decision-2026-05-30.md`
  - Pre-G6 design audit: `archive/reviews/g5-to-g6-handoff/pre-g6-design-audit-2026-05-30.md`
  - G6 MDA gate decision (Gate Passed, six-role sign-off): `archive/reviews/g6-mda/g6-mda-gate-decision-2026-06-01.md`
  - G7 task-planning gate decision (Gate Passed, six-role sign-off): `archive/reviews/g7-task-planning/g7-gate-decision-2026-06-01.md`
  - G8 handoff gate decision (Gate Passed, Claude→Codex): `archive/reviews/g8-handoff/g8-gate-decision-2026-06-01.md`
  - Pre-G8 operator-access Security review: `archive/reviews/g8-handoff/security-operator-access-review-2026-06-01.md`
  - Pre-G6 re-verification (2026-05-31, Domain Modeling role): after the v1
    wording cleanup and the planning/process folder move, the design set was
    re-audited and confirmed `Ready for G6` — no semantic drift, no new
    inconsistency, MDA/PSM input sufficiency intact. Held in session record;
    not separately filed (user elected to proceed to G6).
- Project-layer strengthening (`process/process-gap-register.md`):
  Infrastructure Ops is a required reviewer at G5/G8/G11/G12; every CRM
  cross-capability flow must name a `Primary Flow Owner Agent`.
- DECIDED (2026-06-01) — repository layout, separating concerns:
  - `docs/` = design only (product, business, ux-ui, security, architecture).
  - `modeling/` = the MDA package (CIM/PIM/PSM/traceability-matrix/test-model),
    moved out of `docs/` on 2026-06-01 so modeling is not mixed with design.
  - `planning/` = gate governance only (`gate-status.md`, `blockers.md`).
  - `delivery/` = G7/G8 execution artifacts (`tasks.md`, `task-dependencies.md`,
    `delivery-plan.md`, `acceptance-task-map.md`) — to be created at G7/G8.
  - `process/` = process register.
  This keeps governance (gate-status, blockers) separate from execution artifacts.
  `STANDARD-APPLICATION-REVIEW.md` updated to list the task files under `delivery/`.
