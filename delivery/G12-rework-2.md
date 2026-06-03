# G12 Rework Package #2 — CRM System (Claude → Codex, 2nd KICKBACK)

## Document Control

- Source: G12 RE-AUDIT, 2026-06-03 (`archive/reviews/g12-audit/g12-reaudit-2026-06-03.md`)
- Decision: **REWORK (2nd)** — G11 stays `Gate Blocked`. Do NOT release.
- Status going in: 6 of 8 first-round findings closed and independently verified;
  the codebase builds and unit/integration tests pass (Claude ran them);
  the security-group network state is verified correct. **Only the items below remain.**
- Executor: Codex; then return to Claude for a final focused G12 re-audit. Do not self-pass.

## Rules (unchanged)

TDD fail-first; no weakening tests; real persistence/real captured evidence; no-downgrade.
On completion of each item: update `delivery/tasks.md` status, fill
`modeling/traceability-matrix.md`, set the matching `planning/blockers.md` row to `Resolved`
with the deciding artifact, and commit.

## What is ALREADY closed (do not redo)

BLK-G12-002, 004, 005, 008 — FIXED and verified. BLK-G12-003, 007 (and BLK-G11-002) —
security-group state independently verified by Claude's read-only live API query
(`docs/release/evidence/volcengine-security-group-verified-readonly-2026-06-03.json`); the
network is correct. Do not re-touch the security group itself.

---

## 🔴 BLK-G12-001 (still Open — partial) — Per-service audit-delivery tests

**Where:** `services/commercial`, `services/work`, `services/account`.

**State:** The outbox dispatcher + transactional delivery are implemented and wired for all
four producers, and `opportunity` has full test proof (`internal/event/dispatcher_test.go`).
But commercial/work/account have **no test** that the event actually reaches audit-history —
their command tests only assert a row was written to the local `outbox_events` table.

**Required fix (TDD):** For each of commercial, work, account, add an integration test
(real Postgres testcontainer) modeled on `services/opportunity/internal/event/dispatcher_test.go`:
1. perform a sensitive mutation, run `DispatchOnce`, assert the corresponding `EVT-*` event is
   delivered to audit-history `/internal/events/append` with the correct signed S2S headers;
2. failure path — audit append fails → row stays unpublished and is retried (no loss);
3. duplicate delivery is de-duplicated by event UID at the consumer.

**Acceptance:** commercial/work/account each have the three tests green; the
delivery-per-service matrix in the re-audit (currently "opportunity only") becomes 4/4.
(Optional consistency: bring `lead`'s remaining non-transactional `_ = h.outbox.Append`
to the same transactional pattern, since lead was cited as the durable reference.)

---

## 🟠 BLK-G12-006 (still Open — partial) — Release-evidence residuals

The transcripts are now real captures — keep them. Fix only these four residuals:

1. **Restore-count contradiction.** `backup-restore-rehearsal-2026-06-03.json` says
   `schemas: 9` / `servicePermissionRoles: 9`, but the attached `\dn` / `\du` output lists
   **10**. Reconcile: regenerate the JSON from the actual catalog output so counts match the
   transcript exactly (or explain the discrepancy in the evidence).
2. **Shared SSH key.** `crm-deploy` and `crm-ops` share one SSH key fingerprint
   (`SHA256:fr4e...`). Issue a **separate key per operator** so per-operator separation is real;
   re-capture `operator-access-transcript` showing distinct fingerprints.
3. **No `sshd -T` proof.** The effective sshd config still shows `passwordauthentication yes` /
   `permitrootlogin without-password`, overridden only by an unverified drop-in. Capture
   `sshd -T | grep -E 'passwordauthentication|permitrootlogin|kbdinteractive'` proving the
   hardening is the EFFECTIVE config.
4. **External-edge negative probe.** The 8080/5432 negative probe ran from LAN
   `192.168.0.107`. Re-run it from genuinely **outside** the host's network edge (e.g. a
   different network / the off-server host) and capture refused/timed-out for `8080` and `5432`.

**Acceptance:** all four residual captures on disk and referenced from the ACC-017 evidence;
`test_release_evidence_transcripts.sh` updated to assert the reconciled counts + distinct
fingerprints + `sshd -T` effective values.

---

## 🟠 BLK-G12-009 (new — MAJOR) — Harden the evidence checker

**Where:** `scripts/test_security_group_evidence.py` (and the evidence convention).

**What's wrong:** the checker validates only JSON shape, so a **read-only Describe snapshot**
passed as if it proved a mutation (this is exactly how the 20:29 file slipped through).

**Required fix:** a file that claims a remediation (`--apply`) must be REQUIRED to contain the
mutating-call RequestIds (`CreateSecurityGroup`, `AuthorizeSecurityGroupIngress`,
`RevokeSecurityGroupIngress`, `ModifyNetworkInterfaceAttributes`) with http=200; a read-only
export must be labeled verification-only and must NOT be accepted as remediation proof. Make
the checker fail-first on a read-only snapshot presented as remediation.

**Acceptance:** a test feeding the old 20:29 read-only file to the checker FAILS; a proper
apply-with-mutating-RequestIds export PASSES; Claude's verification-only file is accepted as
verification, not as remediation.

---

## 🟡 BLK-G12-010 (new — MINOR) — Correlation-id + admin-scope test

1. Reporting projection delivery omits `X-Correlation-Id` (it is set on the audit-append call
   but not the reporting call in `services/*/internal/event/outbox.go`). Add it so the
   projection-ingest path is traceable per STB-003.
2. Add a reporting test for **Administrator all-scope** report query (currently only
   Manager team-scope and Sales-denied are covered).

**Acceptance:** correlation id present on reporting delivery; admin-all-scope report test green.

---

## Definition of Done for rework #2

- BLK-G12-001 delivery/failure/dedup tests green for commercial, work, account (4/4).
- BLK-G12-006 four residuals captured + checker updated.
- BLK-G12-009 checker hardened with the failing/passing tests.
- BLK-G12-010 correlation id + admin-scope test.
- tasks.md / traceability-matrix.md / blockers.md updated to real artifacts; commits made.
- Return to Claude for the FINAL focused G12 re-audit (BLK-G12-001/006/009/010 + a spot
  regression). Do not self-pass G12.
