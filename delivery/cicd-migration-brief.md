# CRM CI/CD Migration — Investigation & Execution Brief

For: CRM project agents (architecture / infrastructure-ops / task-planner /
backend-/frontend-engineer / qa-execution / integration-owner / audit).
Author: Claude (planning), 2026-06-06.
Companion: `delivery/cicd-migration-plan.md` (task package M1–M6 + gate path).
Authority: `../../../standards/cicd-and-release-standard.md`,
`../../../company/policy-changes/2026-06-06-cicd-and-release-pipeline.md`.

This brief is self-contained on purpose: do not rely on chat history. Investigate
from the files referenced here.

---

## 1. Background — what happened and why this exists

- The company adopted its **first CI/CD & release-pipeline standard** on
  2026-06-06 (`standards/cicd-and-release-standard.md`). Core rule:
  **production hosts must not build application source**; build off-host (CI),
  ship an **immutable image**, deploy by **pulling/loading a digest** that is
  **traceable to a commit**, and record release evidence.
- CRM previously violated this: it built **on the production host**
  (`docker-compose.prod.yml` used `build:` for all 10 Go services; the runbook ran
  `npm ci && npm run build` and `docker compose up -d --build` on the host). This
  caused disk pressure (root disk hit 87%) and a stale-build release incident.
- On **2026-06-06 the runtime host `srv-volcengine-sh-01` was cleared**: the CRM
  stack, images, `/opt/crm-system` (incl. the postgres data volume), and CRM
  Nginx config were removed; CRM data was **discarded by user decision**. The host
  now runs only OS + vendor cloud agents + base infra (Docker/Nginx/SSH/certbot).
- Therefore this is **not an in-place migration of a running system**. It is
  **CRM's first standard-compliant (re)deployment from a clean host**, using the
  already-audited application code. The DB starts **empty** (fresh migrations).

## 2. The standard you must comply with (summary; read the full file)

1. Off-host build only; no source build / `--build` / compiler on the prod host.
2. CI builds + runs tests + produces images; CD only pulls/loads a tag/digest.
3. Every deployed image is traceable to its source commit; production identity is
   a **digest**, not a moving tag.
4. Release must record: test results, image **digest→commit**, deploy transcript,
   health check, and a named rollback point.
5. No long-term host dependence on GitHub source pulls or build cache.
6. Pipeline does NOT bypass QA (G10) / Integration (G11) / Audit (G12).
7. On-host build is allowed only as break-glass (approval + risk record + expiry
   ≤30d + migration plan + full evidence). **Not needed here** (clean host).

## 3. Problem points / current gaps to resolve

- P1. `docker-compose.prod.yml` is build-based (`build:` on every Go service) →
  must become **image-only** (references tag/digest, no `build:`).
- P2. `deploy/ops/go-live-runbook.md` builds on the host (`npm run build`,
  `up -d --build`) → must become **pull/load-by-digest** (use
  `../../../templates/deployment-runbook.md`).
- P3. **No image registry or artifact channel exists** → must decide one (M1).
- P4. The frontend is currently built on the host into `frontend/dist` for Nginx
  → must be built **off-host** and shipped (or containerized). Decide.
- P5. The host was wiped → **secrets/`.env`, security-group 80/443, and TLS cert
  must be re-provisioned**; the old `/opt/crm-system/.env` is gone.
- P6. The **exact release commit to deploy** must be confirmed (the audited HEAD,
  including the later zh-CN localization work — see `planning/gate-status.md`).

## 4. What to do (task package — investigate → decide → execute through gates)

| # | Task | Output |
|---|---|---|
| M1 | Choose image registry **or** save/load channel | ADR in `docs/architecture/adr/` (registry vs export-load; auth; retention; provenance) |
| M2 | Off-host CI: build all 10 Go images + frontend, run tests, tag by commit, push/export digest-pinned | pipeline def (`../../../templates/cicd-pipeline.md`) + built images |
| M3 | Image-only production compose | `docker-compose.prod.images.yml` (no `build:`; `image:` by tag/digest; keep postgres, internal net, security opts, loopback gateway, healthchecks) |
| M4 | Digest-pinned deploy runbook | rewritten runbook (no host build) |
| M5 | Release evidence | test results, per-service digest→commit, deploy transcript, health check, rollback point |
| M6 | Commit-traceability verification | proof running images = audited commit (digest, not moving tag) |

## 5. What you must understand / verify up front (prerequisites & known facts)

**Three technical gotchas (will break the deploy if missed):**
1. **CPU architecture.** Build host (likely Apple Silicon, arm64) ≠ server
   (`srv-volcengine-sh-01` = Ubuntu 24.04, **linux/amd64**). Images MUST be built
   for `linux/amd64` (e.g. `docker buildx build --platform linux/amd64`).
   Note: the **server's docker has no `buildx`** — build on the workstation/CI,
   not the host.
2. **Image-only compose.** Keep everything in the current prod compose EXCEPT the
   `build:` blocks → replace with `image:`. Preserve: postgres + its init/migration
   mounts, `crm_internal` internal network, `crm_edge`, `read_only`/`cap_drop`/
   `no-new-privileges`, the loopback-only gateway (`127.0.0.1:8080`), healthchecks,
   log mounts.
3. **Frontend.** Decide: build `dist` off-host and `scp` to
   `/opt/crm-system/current/frontend/dist` for host Nginx, OR containerize the
   frontend as an nginx image. Either is fine; on-host `npm build` is not.

**Infrastructure facts (source of truth = `company/infrastructure/` registers; do
not copy secrets into docs):**
- Host: `srv-volcengine-sh-01`, public IP `118.196.44.193`, internal `172.31.8.67`,
  4 vCPU / 8 GiB / 40 GiB, root disk currently ~25% (cleared).
- SSH: see `ssh-access-register.md` (key under `.secrets/ssh-keys/`, local bind
  `-b 192.168.0.104`; root login for now). Never print key contents.
- Ports: 22 public; **80/443 currently have no listener** (CRM Nginx removed) —
  security-group re-open + Nginx vhost must be restored for ingress. Gateway is
  loopback `127.0.0.1:8080` via Nginx reverse proxy only.
- TLS: Let's Encrypt cert for `118.196.44.193` was valid to **2026-06-09**;
  certbot retained. Confirm/renew on redeploy.
- Vendor cloud agents (`openclaw-gateway`, `cloud-monitor-agent`, `elkeid-agent`)
  must remain untouched.
- `.env` / secrets: the prior file was deleted with `/opt/crm-system`. Required
  vars are referenced in the old runbook/compose (POSTGRES_*, per-service
  `*_DATABASE_URL`, `SERVICE_TOKEN_SECRET`, session/JWT secret). Re-provision on
  the host; never commit.

**Release content (corrected 2026-06-06 — see gap-analysis D3):**
- The release content is **NOT "the audited HEAD as-is."** The dashboard/overview/
  reports UI is committed **P1** (ACC-018/ACC-023, CAP-009) and its specified
  UX/UI design realization is **prior-intended work that was not completed** (a
  separate follow-on change: `delivery/uiux-completion-charter.md`). The
  deployable commit is **"audited backend + the gate-cleared UI/UX completion"**,
  not the current commit; shipping as-is would deploy an incomplete P1 UI
  (no-downgrade violation). The later zh-CN localization (live 2026-06-05) must be
  preserved.
- This CI/CD migration itself changes **no application code** and must **not**
  weaken any G12 fix; it only changes how the artifact is built and deployed.
- Confirm the exact release commit only **after** the UI/UX completion is
  gate-cleared (`delivery/cicd-current-state-gap-analysis.md` D3).
- DB starts empty → run all migrations fresh; decide if any seed data is needed.

## 6. Open decisions the project must make (record in decision-log / ADR)

- D1. Registry form vs save/load export form (M1). Save/load = zero new infra,
  fully compliant; registry = better for multi-host/future. Pick and ADR it.
- D2. Frontend: scp `dist` vs containerize as nginx image (P4).
- D3. Exact release commit (= audited backend + gate-cleared UI/UX completion,
  NOT as-is HEAD; see gap-analysis D3 + `uiux-completion-charter.md`) + whether
  any seed data is loaded into the empty DB.
- D4. 80/443 security-group re-open scope (keep co-location constraint: do not
  take host ingress ownership beyond the CRM `server_name`).

## 7. Hard constraints

- **No application code change**, no business-logic rework; this is build+deploy
  mechanism only. No downgrade/revert of any G12 fix or P0/P1 acceptance.
- **Off-host build**; the production host only pulls/loads + runs.
- **Digest-pinned + commit-traceable**; no `latest`-only production identity.
- **Capture §4 release evidence**; a green pipeline does not replace QA(G10)/
  Integration(G11)/Audit(G12).
- Follow the gate path: re-enter at **G5** (deployment + pipeline design incl. M1
  ADR) → **G7/G8** (tasks for M1–M6) → **G9–G11** (Codex executes) → **G12**
  (independent audit verifies digest→commit, §4 evidence, and that no on-host
  build was used). Update `planning/gate-status.md` on every handoff.
- Infrastructure changes (security group, ports, TLS, deploy) are
  `infrastructure-ops` execution with before/after evidence in
  `company/infrastructure/infrastructure-change-log.md`.

## 8. Definition of done

- Production runs **off-host-built, digest-pinned** images for all 10 services +
  frontend, each traceable to the audited commit.
- `docker-compose.prod.yml` (or the image-only variant) has **no `build:` keys**;
  the runbook has **no host build steps**.
- §4 release evidence recorded; ingress (80/443 + TLS) healthy; health checks pass.
- No P0/P1 acceptance item or prior G12 fix weakened.
- G12 deployment-compliance audit passed; `planning/gate-status.md` updated.
