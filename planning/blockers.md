# Blockers — CRM System

This register tracks open blockers and pre-gate decisions surfaced during planning.
It is a companion to `planning/gate-status.md` (the gate sync point). On a kickback,
set the affected gate to `Gate Blocked` in gate-status, record it here, and return
to Claude. See `../../../company/collaboration-model.md`.

Status values: `Open` / `In Review` / `Resolved` / `Formal Scope Change by User`.
No-downgrade rule applies: a blocker touching a P0/P1 acceptance may only move to
`Resolved` (with the deciding source recorded) or `Formal Scope Change by User`.

## Open Modeling Decisions — must close before G8

These were surfaced by the G6 MDA modeling (CIM/PIM) and its independent multi-agent
audits. Each is a genuine upstream-source gap that MDA correctly declined to invent
(per the no-invention discipline).

> **Update 2026-06-01: ALL THREE RESOLVED by a Formal Scope Change by User**
> (decision-log.md DEC-017..020). BLK-001 → DEC-020 (Opportunity Status field
> removed); BLK-002 → DEC-017/019 (Won = contract Signed; payment decoupled, so no
> multi-plan "fully paid" Won-gate); BLK-003 → DEC-018 (one quote per opportunity;
> no second quote to accept). Cascade applied across baseline → architecture → MDA
> and re-audited. See Resolution Log below.

| ID | Blocker | Owner | Blocks | Touches | Status | Opened |
|---|---|---|---|---|---|---|
| BLK-001 | Opportunity Status enumerated value set distinct from Pipeline Stage is not enumerated in the accepted sources (CIM-016 commits Status as a persisted dimension; PRD-007/ACC-007 require status as a persisted field but enumerate no values). Decide whether a distinct Status enumeration exists or whether Status is realized through Stage. | Product Manager + Business Analyst | G8 handoff (and G7 test design determinism) | ACC-007 (P0) | Resolved — Formal Scope Change by User (DEC-020) | 2026-06-01 |
| BLK-002 | Multi-plan contract "fully paid" aggregation is undefined. Overpayment ceiling is fixed at contract level (EDGE-019) and Won requires contract-level full payment (BR-008, DEC-012), but the accepted sources do not define how multiple Payment Plans under one contract aggregate into "contract fully paid." Confirm the aggregation rule. | Product Manager + Business Analyst | G8 handoff (and G7 deterministic Won/overpayment test design) | ACC-011, ACC-013 (P0) | Resolved — Formal Scope Change by User (DEC-017/019) | 2026-06-01 |
| BLK-003 | Second-quote-accept observable outcome is unspecified. The "at most one Accepted quote per opportunity" invariant is fixed, but EDGE-012 does not specify the observable business result of accepting a second quote (reject the second accept vs auto-demote the current Accepted quote). A product decision is required before deterministic test design. | Product Manager | G7 test design / G8 handoff | ACC-009 (P0) | Resolved — Formal Scope Change by User (DEC-018) | 2026-06-01 |

### Source pointers
- BLK-001: `modeling/CIM.md` CIM-016; `modeling/PIM.md` PIM-SM-003, PIM-OPEN-003; `docs/product/acceptance-matrix.md` ACC-007; PRD-007.
- BLK-002: `modeling/PIM.md` PIM-SM-006, PIM-INV-007/023/025, PIM-OPEN-005; `docs/business/business-rules.md` BR-008; `docs/product/decision-log.md` DEC-012; `docs/business/edge-cases.md` EDGE-019.
- BLK-003: `modeling/CIM.md` CIM-PROC-008 Open/Blocked; `modeling/PIM.md` PIM-SM-004, PIM-OPEN-001; `docs/business/edge-cases.md` EDGE-012; ACC-009.

## Architecture / PSM-deferred (not PM/BA blockers; tracked for G7/PSM)

These are correctly deferred to PSM/Architecture at PIM altitude and are recorded for
the PSM artifact and G7 test design; they are not pre-G8 PM/BA decisions.

| ID | Item | Owner | Resolves at | Refs |
|---|---|---|---|---|
| BLK-A01 | Overdue-evaluation trigger (on-read vs scheduled) — needed for deterministic overdue test design; modeled at PIM as a Business-Date guard, mechanism deferred. | Architecture (in PSM) | PSM / G7 | PIM-OPEN-002, CIM-034, BR-021 |

## Open G11 Release Evidence Blockers

These are active execution blockers surfaced during Codex G9-G11 execution. Per the
no-invention discipline, Codex may not invent production endpoint, TLS, security-group,
monitoring, or runtime evidence.

| ID | Blocker | Owner | Blocks | Touches | Status | Opened |
|---|---|---|---|---|---|---|
| BLK-G11-001 | TASK-039 cannot be completed because production release requires a valid HTTPS endpoint with TLS certificate, but the deployment source says the domain is "Not specified yet" and absence of a valid HTTPS endpoint blocks production release. G11 also requires real runtime evidence on `srv-volcengine-sh-01` for endpoint, TLS status, security-group inbound rules, opened ports, health URL, deployment timestamp, operator, and smoke result. | Infrastructure Ops + Integration Owner | G11 handoff to G12; TASK-039 | ACC-017 (P0), ARCH-ACC-008/013/014/015, TEST-DEPLOY-SMOKE-001/002 | Resolved — `https://118.196.44.193` approved and provisioned with a Let's Encrypt IP certificate on 2026-06-03 | 2026-06-03 |
| BLK-G11-002 | TASK-039 has server-side HTTPS/TLS, health, redirect, service health, monitoring threshold, renewal, operator evidence, and Volcengine security-group API evidence. The API evidence binds instance `i-yemoz0an7kk36d2c9bp6` to security group `sg-1pm4k7f37z8xs643rg0fvk85e` through ENI `eni-13e8tbocd8f0g79iu5jer8idt` and confirms CRM gateway `8080` and PostgreSQL `5432` are not allowed from `0.0.0.0/0`. User approved releasing previous deployments; Codex stopped and removed host-network Hermes, so host-level `8642` no longer listens and CRM smoke still passes. The user removed old/non-CRM Volcengine security-group rules for TCP `8088`, TCP `8443`, and TCP `3389`; API post-cleanup verification confirms only public TCP `22`, `80`, and `443` remain. | Infrastructure Ops + Security Compliance | G11 handoff to G12; TASK-039 | ACC-017 (P0), ARCH-ACC-013, ARCH-ACC-014, TEST-DEPLOY-SMOKE-001/002 | Resolved | 2026-06-03 |

## Carried-forward Release Blockers (not gate blockers)

Recorded in `planning/gate-status.md` G12 row; repeated here for visibility. These were
release-time evidence items, not modeling blockers: encrypted off-server backup copy +
restore rehearsal, HTTPS/TLS endpoint, security-group, and monitoring evidence. As of
2026-06-03, Codex has recorded runtime evidence for these items under TASK-039 and
TASK-040; G12 still performs independent audit before any release decision. Refs:
OQ-001, RISK-002, `docs/architecture/deployment-notes.md`.

## G12 Audit Findings — Rework before release (2026-06-03)

Surfaced by the first G12 independent audit (`archive/reviews/g12-audit/g12-audit-decision-2026-06-03.md`),
five parallel author≠auditor passes; two dimensions FAILED; first rework package
`delivery/G12-rework.md`. Codex remediated and returned. Claude then ran a focused
multi-agent **G12 RE-AUDIT** (`archive/reviews/g12-audit/g12-reaudit-2026-06-03.md`),
including an **independent compile/test run** (11/11 Go modules build/vet/test green
against real PostgreSQL testcontainers; frontend tsc+build green) and a **Claude-run
read-only live Volcengine query**. Result: **6 of 8 findings genuinely closed; BLK-G12-001
and BLK-G12-006 remain Open (partially fixed)**. Second rework package:
`delivery/G12-rework-2.md`. No-downgrade applies. The Status column below reflects the
RE-AUDIT verdict, which supersedes Codex's self-reported "Resolved".

| ID | Severity | Blocker | Owner | Touches | Status (re-audit) | Opened |
|---|---|---|---|---|---|---|
| BLK-G12-001 | BLOCKER | `opportunity/commercial/work/account` wrote `EVT-*` events only to local `outbox_events` with no relay to audit-history. **Second rework fixed** — account, commercial, and work now each have real PostgreSQL testcontainer dispatcher coverage for successful S2S audit-history delivery, failed-delivery retry retention, and duplicate event UID idempotency; opportunity already had the same coverage. | backend-engineer | ACC-014, ACC-022 (P0); AUD-IMM-002 | Resolved | 2026-06-03 |
| BLK-G12-002 | BLOCKER | `reporting` `POST /internal/projections` had no S2S verification. **Re-audit: FIXED** — signed-token middleware enforced (aud/intent/≤5min, fail-closed); negative+positive tests present. | backend-engineer | SEC-SVC-BLK-002 | Resolved | 2026-06-03 |
| BLK-G12-003 | BLOCKER | Security-group cleanup was unsubstantiated and contradicted the only raw API export. **Re-audit: RESOLVED via independent live verification** — Claude ran a read-only Volcengine query (RequestIds `2026-06-03T22:08+08:00`) confirming the CRM ENI is bound only to dedicated SG `sg-366ptx1bxp9ts1e710babmc8y` (public TCP 22/80/443 only; no 8080/5432; no 8088/8443/3389). NOTE: the earlier `...dedicated-raw-...json` (20:29) recorded only Describe calls (no mutating RequestIds) and is NOT acceptable provenance; the END STATE is nonetheless independently confirmed correct. | infrastructure-ops + security-compliance | ACC-017 (P0); ARCH-ACC-013/014 | Resolved | 2026-06-03 |
| BLK-G12-004 | MAJOR | Outbox append not in the mutation's transaction and error discarded. **Re-audit: FIXED** — single-transaction `inTransaction` helper across all four services; append error returned (rollback + 503); testcontainer rollback test `TEST-HISTORY-TX-001`. | backend-engineer | AUD-IMM-002 | Resolved | 2026-06-03 |
| BLK-G12-005 | MAJOR | `reporting` projection ingest never invoked; authz over empty data. **Re-audit: FIXED** — four source domains deliver to reporting via signed dispatchers; consumer upserts read-model; end-to-end + authz tests present. | backend-engineer | PM-043..045 (P1) | Resolved | 2026-06-03 |
| BLK-G12-006 | MAJOR | Release evidence not independently verifiable. **Second rework fixed** — restore counts now match captured catalog output (roles/schemas/servicePermissionRoles all `10`), `crm-deploy` and `crm-ops` use distinct Ed25519 key fingerprints, `sshd -T` proves password auth is disabled and keyboard-interactive is disabled, and `8080`/`5432` negative probes were rerun from external edge `srv-aliyun-bj-01`. | infrastructure-ops | ACC-017 (P0); operator-access review conditions | Resolved | 2026-06-03 |
| BLK-G12-007 | MAJOR | CRM shared the `Default` security group. **Re-audit: RESOLVED via live verification** — CRM ENI bound only to dedicated SG `sg-366ptx...`; Default SG no longer governs the CRM instance. Same evidence as BLK-G12-003. | infrastructure-ops | ACC-017; deployment-notes Network Exposure | Resolved | 2026-06-03 |
| BLK-G12-008 | MINOR | Test-traceability + hardening gaps. **Re-audit: FIXED** — 7 P0 tests now tagged 1:1; `work` S2S 5-min cap enforced; `AUTHZ_VERSION_STALE` enforced; root screenshot removed; all regression checks PASS. | backend-engineer | traceability; minor security hardening | Resolved | 2026-06-03 |
| BLK-G12-009 | MAJOR | Evidence-integrity gap: `scripts/test_security_group_evidence.py` validated only JSON shape and accepted a read-only snapshot as "success". **Second rework fixed** — `--apply` remediation evidence now requires `CreateSecurityGroup`, `AuthorizeSecurityGroupIngress`, `RevokeSecurityGroupIngress`, and `ModifyNetworkInterfaceAttributes` RequestIds with HTTP 200; read-only exports are accepted only with `--verification`. | infrastructure-ops | evidence integrity | Resolved | 2026-06-03 |
| BLK-G12-010 | MINOR | Reporting projection delivery omits `X-Correlation-Id` (set on the audit-append call but not the reporting call), and no Administrator all-scope report test exists. Spec STB-003 traceability gap. | backend-engineer | STB-003; PM-043 | **Open** | 2026-06-03 |

> **BLK-G11-002 RE-RESOLVED (2026-06-03):** the original "Resolved" claim rested on a
> hand-authored file and was correctly reopened by BLK-G12-003. It is now genuinely
> resolved together with BLK-G12-003/007 by Claude's independent read-only live API
> verification (`docs/release/evidence/volcengine-security-group-verified-readonly-2026-06-03.json`),
> not by the contradicted/snapshot-only earlier files.

## Resolution Log

| Date | ID | Resolution | Deciding source |
|---|---|---|---|
| 2026-06-01 | BLK-001 | Opportunity `Status` field removed; Pipeline Stage is the sole lifecycle dimension. The "what are Status's values" gap is dissolved, not answered. | DEC-020 (Formal Scope Change by User) |
| 2026-06-01 | BLK-002 | Won decoupled from full payment (Won = contract Signed); multi-plan "fully paid" is no longer a Won precondition, so the aggregation gate is moot. Overpayment ceiling remains contract-level. | DEC-017 + DEC-019 (Formal Scope Change by User) |
| 2026-06-01 | BLK-003 | Each opportunity has exactly one quote; there is no second quote to accept, so the reject-vs-auto-demote question no longer exists. | DEC-018 (Formal Scope Change by User) |
| 2026-06-01 | BLK-A01 | Overdue-evaluation trigger resolved at PSM as on-read evaluation against the supplied `businessDate` (Asia/Shanghai); deterministic for G7 test design. | PSM "Resolved Mechanisms"; api-spec.md Reminder Query; FLOW-005 |
| 2026-06-03 | BLK-G11-001 | The user approved `https://118.196.44.193` as the production HTTPS endpoint. Codex installed Certbot 5.6.0, issued a Let's Encrypt IP certificate with SAN `IP Address:118.196.44.193`, enabled Nginx 443, verified HTTP→HTTPS redirect and server-side deploy smoke, and configured renewal timer/dry-run. | `docs/release/acc-017-evidence-template.md`; TASK-039 server evidence |
| 2026-06-03 | BLK-G11-002 | Volcengine API evidence confirms CRM `8080` and PostgreSQL `5432` are not publicly allowed; host-network Hermes `8642` was stopped/removed; old/non-CRM security-group rules TCP `8088`, TCP `8443`, and TCP `3389` were removed by the user and verified by API. | `docs/release/evidence/volcengine-security-group-post-cleanup-2026-06-03.json`; `docs/release/evidence/old-deployment-release-2026-06-03.json`; TASK-039 closure evidence |
| 2026-06-03 | TASK-040 release evidence | Encrypted PostgreSQL backup `crm-postgres-20260603T104620Z.sql.gz.enc` was produced on `srv-volcengine-sh-01`, copied to `srv-aliyun-bj-01`, verified by checksum, and restored in rehearsal run `20260603T104837Z`. `crm-backup.timer` is enabled and active for daily 02:00 backups. | `docs/release/acc-017-backup-evidence-template.md`; `docs/release/evidence/backup-restore-rehearsal-2026-06-03.json`; TASK-040 closure evidence |
| 2026-06-03 | BLK-G12-001 | Implemented outbox dispatchers for opportunity, account, commercial, and work services. Dispatchers read unpublished local `outbox_events`, send S2S-signed `audit.append` calls to audit-history with the outbox row ID as `eventUid`, leave failed deliveries unpublished for retry, and mark success with `published_at`. audit-history accepts producer event UID and is idempotent on duplicate UID. Compose/prod compose now configure `AUDIT_HISTORY_SERVICE_URL` for the four producers. | `services/*/internal/event/outbox.go`; `services/*/cmd/server/main.go`; `services/opportunity/internal/event/dispatcher_test.go`; `services/audit-history/internal/handler/append_test.go`; `go test ./...` in opportunity/account/commercial/work/audit-history |
| 2026-06-03 | BLK-G12-004 | Moved append-bearing mutations in opportunity, account, commercial, and work onto one local DB transaction with the outbox insert; append errors are no longer discarded. Added fail-first rollback coverage for opportunity create (`TEST-HISTORY-TX-001`) and removed dropped outbox append patterns from the four audited service handlers. | `services/opportunity/internal/handler/opportunity_command_test.go`; `services/{opportunity,account,commercial,work}/internal/handler`; `services/{opportunity,account,commercial,work}/internal/repo`; `go test ./...` in opportunity/account/commercial/work |
| 2026-06-03 | BLK-G12-002 | Added signed S2S verification to reporting `POST /internal/projections` with audience `reporting`, intent `reporting.projection_ingest`, max 5-minute token lifetime, and `SERVICE_AUTH_FAILED` rejection for missing/expired/wrong-audience/wrong-intent/wrong-signature calls. Negative cases assert no projection mutation; valid token succeeds. | `services/reporting/internal/authz/service_token.go`; `services/reporting/internal/handler/overview_query.go`; `services/reporting/internal/handler/overview_query_test.go`; `go test ./internal/handler -run TestProjectionIngestRequiresS2SToken`; `go test ./...` in services/reporting |
| 2026-06-03 | BLK-G12-005 | Wired source domain outbox events into reporting projection ingest. lead dispatches local outbox events to reporting; account/opportunity/commercial dispatchers now deliver reporting projections before marking an event published. reporting S2S projection ingest populates `record_projections`, Manager report queries read the aggregate, and Sales remains denied. | `services/lead/internal/event/outbox.go`; `services/{account,opportunity,commercial}/internal/event/outbox.go`; `services/reporting/internal/handler/overview_query_test.go`; `docker-compose.yml`; `docker-compose.prod.yml`; `go test ./...` in lead/account/opportunity/commercial/reporting |
| 2026-06-03 | BLK-G12-003 / BLK-G12-007 rework attempt blocked | Added fail-first evidence validation (`TEST-DEPLOY-SG-001`) requiring raw Volcengine API proof that CRM is bound to a dedicated non-Default security group with public TCP ingress limited to 22/80/443 and no public 8080/5432/8088/8443/3389. OpenAPI access succeeded, but remediation stopped at `CreateSecurityGroup` because the current operator account lacks `vpc:CreateSecurityGroup`; G11/G12 remain Gate Blocked until the required VPC IAM actions are granted and the script is rerun. | `scripts/test_security_group_evidence.py`; `scripts/volcengine_security_group_rework.py`; `docs/release/evidence/volcengine-security-group-rework-blocked-2026-06-03.md` |
| 2026-06-03 | BLK-G12-003 / BLK-G12-007 | After VPC IAM permission was granted, reran the OpenAPI remediation. CRM ENI `eni-13e8tbocd8f0g79iu5jer8idt` is now bound only to dedicated non-Default security group `sg-366ptx1bxp9ts1e710babmc8y` (`crm-system-prod-public`). The dedicated group's final raw ingress allows public TCP `22`, `80`, and `443` only. The final raw default-group export has no public TCP `8088`, `8443`, or `3389`, and the final evidence has no public TCP `8080` or `5432`. `TEST-DEPLOY-SG-001` passed. | `docs/release/evidence/volcengine-security-group-dedicated-raw-2026-06-03.json`; `docs/release/evidence/volcengine-security-group-rework-transcript-2026-06-03.txt`; `scripts/test_security_group_evidence.py` |
| 2026-06-03 | BLK-G12-006 | Second rework closed the residual evidence gaps: restore JSON and transcript now match the captured catalog counts (`roles=10`, `schemas=10`, `servicePermissionRoles=10`); `crm-deploy` uses `SHA256:ZGLqXBHGgqy29ZUMFysRjaw579Z3yx1980pIFWBb/b4` and `crm-ops` uses `SHA256:PHl9ZXjKKPzI5oiWrll9Jj60X04+5S7/TMpV1q3AYQA`; `sshd -T` shows `passwordauthentication no`, `kbdinteractiveauthentication no`, and root key-only login; `8080`/`5432` external negative probes ran from `srv-aliyun-bj-01` and timed out as expected. `TEST-RELEASE-EVIDENCE-001` passed. | `docs/release/evidence/external-negative-probes-2026-06-03.txt`; `docs/release/evidence/backup-restore-rehearsal-2026-06-03.json`; `docs/release/evidence/backup-restore-transcript-2026-06-03.txt`; `docs/release/evidence/operator-access-transcript-2026-06-03.txt`; `scripts/test_release_evidence_transcripts.sh` |
| 2026-06-03 | BLK-G12-009 | Hardened `scripts/test_security_group_evidence.py`: `--verification` accepts read-only Describe exports only as state verification, while `--apply` requires mutating API provenance with HTTP 200 RequestIds for `CreateSecurityGroup`, `AuthorizeSecurityGroupIngress`, `RevokeSecurityGroupIngress`, and `ModifyNetworkInterfaceAttributes`. Fail-first test proved the old 20:29 read-only file fails as `--apply`; Claude's read-only verification file passes as `--verification`; an apply fixture with all mutating RequestIds passes. | `scripts/test_security_group_evidence.py`; `scripts/test_security_group_evidence_checker.sh`; `python3 scripts/test_security_group_evidence.py --evidence docs/release/evidence/volcengine-security-group-dedicated-raw-2026-06-03.json --apply` fails as expected; `bash scripts/test_security_group_evidence_checker.sh` |
| 2026-06-03 | BLK-G12-008 | Added explicit P0 test IDs to existing permission, denied-mutation, owner-transfer, and open-work cascade tests; added fail-first `TEST-SVC-TOKEN-LIFETIME-001` for work S2S token lifetime and enforced a 5-minute cap; added fail-first `AUTHZ_VERSION_STALE` coverage and now returns that error when session authz version is stale; removed the stray root screenshot. | `services/identity-authz/internal/handler/{auth.go,auth_test.go,permission_test.go,user_admin.go}`; `services/work/internal/authz/service_token.go`; `services/work/internal/authz/service_token_test.go`; `services/{lead,opportunity,work}/internal/handler/*test.go`; `go test ./...` in identity-authz/lead/work/opportunity |
| 2026-06-03 | BLK-G12-003 / BLK-G12-007 (RE-AUDIT, supersedes the two rows above) | Claude G12 re-audit found the Codex `...dedicated-raw-...json` (20:29) recorded only Describe calls — no `CreateSecurityGroup`/`Authorize`/`Revoke`/`ModifyNetworkInterfaceAttributes` RequestIds — so it cannot serve as remediation provenance. Claude then ran an INDEPENDENT read-only live Volcengine query (Describe only) and confirmed the END STATE is genuinely correct: ENI `eni-13e8...` bound only to dedicated SG `sg-366ptx...` (`crm-system-prod-public`), public TCP `22/80/443` only, no public `8080/5432/8088/8443/3389`. Real API RequestIds `2026-06-03T22:08+08:00`, http=200. Underlying remediation accepted on the strength of this independent verification; provenance-capture and checker hardening tracked as BLK-G12-009. | `docs/release/evidence/volcengine-security-group-verified-readonly-2026-06-03.json` (Claude independent live read) |
| 2026-06-03 | BLK-G12-001 (second rework) | Added account, commercial, and work dispatcher integration tests mirroring opportunity: each uses real PostgreSQL testcontainers, appends a local outbox event, verifies `DispatchOnce` sends an S2S-signed audit-history append with stable `eventUid` and correlation ID, leaves failed deliveries unpublished for retry, and preserves the same event UID across duplicate delivery attempts. Fixed `work.outbox_events` migration to include `published_at`, matching the dispatcher contract. | `services/{account,commercial,work}/internal/event/dispatcher_test.go`; `services/work/migrations/0002_work.up.sql`; `go test ./internal/event -run TestOutboxDispatcherDeliversRetriesAndDedupesAuditHistoryEvents` in account/commercial/work; `go test ./...` in account/commercial/work |
| 2026-06-03 | BLK-G11-002 (re-resolved) | Re-resolved together with BLK-G12-003/007 by Claude's independent read-only live API verification, replacing the earlier contradicted/snapshot-only files as the deciding source. | `docs/release/evidence/volcengine-security-group-verified-readonly-2026-06-03.json` |
