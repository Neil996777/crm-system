# CRM Deployment — Current-State vs CI/CD Standard Gap Analysis

Status: Planning input only. **No deployment, no application-code change.**
Date: 2026-06-06
Author (planning): Claude (CRM planning roles).
Purpose: Self-contained, evidence-verified gap list feeding the **G5 re-entry**
(deployment + pipeline design) of the CRM CI/CD migration.
Authority: `../../../standards/cicd-and-release-standard.md`,
`../../../company/policy-changes/2026-06-06-cicd-and-release-pipeline.md`.
Companions: `delivery/cicd-migration-brief.md`, `delivery/cicd-migration-plan.md`.

All evidence below was read directly from the listed files (not from chat
history). File:line citations are the as-of-2026-06-06 working tree.

---

## 1. Current deployment method (as built)

The current production deployment is an **on-host source build**. The operator
(`crm-deploy`) runs, on `srv-volcengine-sh-01`:

| Step | Source | What it does |
|---|---|---|
| 1 | `go-live-runbook.md:39-41` | `git fetch` + `git checkout <commit>` — pulls source onto the prod host |
| 2 | `go-live-runbook.md:43-46` | `backup.sh` + `offsite-copy.sh` — DB backup before migrate (**compliant; keep**) |
| 3 | `go-live-runbook.md:48-49` | `cd frontend && npm ci && npm run build` — **builds the SPA on the host** |
| 4 | `go-live-runbook.md:51-52` | `docker compose -f docker-compose.prod.yml up -d --build` — **builds 10 Go images on the host** |
| 5 | `go-live-runbook.md:54-56` | `migrate.sh up` — applies migrations |
| 6-7 | `go-live-runbook.md:58-63` | reload Nginx; confirm backup/cert timers |
| Rollback | `go-live-runbook.md:87-92` | `git checkout <prev>` + `up -d --build` — **rebuilds on host to roll back** |

`docker-compose.prod.yml` carries a `build:` block for every one of the 10 Go
services (`:78-242`), each tagged `${CRM_IMAGE_TAG:-latest}` — a **moving tag**,
not a digest. Postgres is `postgres:16-alpine` (pulled, fine).

Net: **source goes onto the host → host builds → host runs `latest` → host
rebuilds to roll back.** This is exactly the pattern the new standard prohibits.

## 2. Gap matrix (current state → standard clause → required change)

| # | Standard clause | Requirement | Current state (evidence) | Verdict | Required change |
|---|---|---|---|---|---|
| G-1 | §1.1 | Prod host MUST NOT build source | `npm run build` (`runbook:49`), `up -d --build` (`runbook:52,91`) | ❌ FAIL | Move all build off-host (CI/workstation) |
| G-2 | §1.2 | CI builds + tests + produces immutable image | No CI pipeline exists | ❌ FAIL | Stand up off-host pipeline (`templates/cicd-pipeline.md`) |
| G-3 | §1.3 | CD = pull/load specified tag/digest only; no host `git checkout`-to-build | `git checkout` to build (`runbook:40,90`) | ❌ FAIL | Rewrite runbook to pull/load digest (`templates/deployment-runbook.md`) |
| G-4 | §1.4 | Production identity = digest; no moving tag | `${CRM_IMAGE_TAG:-latest}` (`prod.yml:81..242`) | ❌ FAIL | Image-only compose pinned to digest; record digest→commit |
| G-5 | §1.5 | No long-term host dependence on source pulls / build cache | Whole flow depends on `git fetch/checkout` | ❌ FAIL | Host runs registry/loaded image only; no source tree to run |
| G-6 | §4 | Release records test results, digest→commit, deploy transcript, health check, rollback point | ACC-017 evidence template exists (`docs/release/acc-017-evidence-template.md`) but **no image digest→commit**; rollback is "rebuild from commit", no digest rollback point | ⚠️ PARTIAL | Extend evidence with per-service digest→commit + previous-good digest as rollback point |
| G-7 | §3 | Use approved registry OR approved export/load channel | Neither exists | ❌ FAIL | Decide registry vs `docker save/load` (M1 ADR) |
| G-8 | §1.2 / pre-deploy backup | DB backup before migrate | `backup.sh` + offsite (`runbook:43-46`) | ✅ PASS | Keep as-is in the new runbook |

## 3. Practical constraints that shape the redesign (not gaps, but binding)

- **Architecture mismatch.** Build host is likely Apple Silicon (arm64); server
  is Ubuntu 24.04 **linux/amd64**, and the **server docker has no `buildx`**.
  Images MUST be built `--platform linux/amd64` off-host. (`cicd-migration-brief.md:75-81`)
- **Host was cleared 2026-06-06.** `/opt/crm-system` (incl. postgres volume),
  CRM `.env`, CRM Nginx vhost, and the 80/443 listener are gone; TLS cert was
  valid to 2026-06-09. The redeploy is a **clean-host first deployment**, not an
  in-place migration → **no break-glass record required** (§6 not invoked).
  (`cicd-migration-brief.md:26-32,97-101`)
- **Compose internals to preserve** when removing `build:`: postgres + init/
  migration mounts, `crm_internal` internal network, `crm_edge`, `read_only`/
  `cap_drop`/`no-new-privileges`, loopback-only gateway `127.0.0.1:8080`,
  healthchecks, log mounts. (`cicd-migration-brief.md:81-85`)

## 4. Open decisions to resolve at G5 (carry into ADR / decision-log)

- **D1** Registry form vs export/load (`docker save`→`scp`→`docker load`).
  Export/load = zero new infra, fully §3-compliant; registry = better for
  multi-host/future. → M1 ADR.
- **D2** Frontend: containerize as an nginx image. This is confirmed in
  `delivery/cicd-migration-acceptance.md` D2 and formalized by ADR-CICD-001.
- **D3** Release content. **Confirmed (2026-06-12):** release content commit is
  **`66d2531`**, representing the audited backend G12 result plus the completed
  and gate-cleared UI/UX follow-on, with zh-CN preserved. This supersedes earlier
  non-specific HEAD planning language. Also decide seed data for the empty DB.
- **D4** 80/443 security-group re-open scope (keep co-location constraint: do not
  take host ingress ownership beyond the CRM `server_name`).

## 5. Scope guardrails (restated)

- **No application-code change**; no business-logic / API / data-model / security-
  fix edits. This is build+deploy mechanism only.
- **No downgrade** of any P0/P1 acceptance item or prior G12 fix (IDOR, durable
  audit, optimistic concurrency, idempotency, etc.).
- A green pipeline does **not** substitute for QA (G10) / Integration (G11) /
  Audit (G12); §4 evidence is in addition to existing acceptance evidence.

## 6. Recommended next gate move

Re-enter at **G5** with the M1 registry/channel ADR + image-only deployment
design, then G7/G8 task plan for M1–M6, then Codex executes G9–G11, then G12
verifies digest→commit and that no on-host build was used. Update
`planning/gate-status.md` on the re-entry.

**Sequencing vs UI/UX completion (per D3).** The CI/CD **mechanism** design
(M1–M4) is independent of UI content and can proceed. The **final build + deploy**
(M5/M6) targets commit `66d2531`.

> This document changes nothing on disk beyond itself. It is the diagnostic
> input; the redesign is performed at G5, not here.
