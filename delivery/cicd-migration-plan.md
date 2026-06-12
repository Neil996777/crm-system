# CRM CI/CD Migration Plan (Follow-On Project Change)

Status: Planned — not yet started. Planning artifact only; no implementation here.
Date: 2026-06-06
Owner (planning): Claude (CRM planning roles). Execution: Codex (G9–G11).
Authority: `../../../standards/cicd-and-release-standard.md`,
`../../../company/policy-changes/2026-06-06-cicd-and-release-pipeline.md`.

## Why

The company adopted a CI/CD & release-pipeline standard (2026-06-06) requiring
off-host build, digest-pinned deployment, commit-traceable images, and release
evidence. CRM previously built **on the production host** (`build:` for all 10 Go
services + on-server `npm run build`), which caused disk pressure and a
stale-build release incident. This change migrates CRM to the standard. It is a
**separate project change** from the governance change and runs through the CRM
project's normal gates.

## Current State (entry condition)

- The CRM runtime host (`srv-volcengine-sh-01`) was **cleared on 2026-06-06**;
  CRM is decommissioned and its data was discarded per user decision.
- Therefore there is **no running on-host-built CRM instance to grandfather** and
  **no break-glass record is required** — the redeploy will be CRM's first
  standard-compliant, off-host-build, digest-pinned deployment on a clean host.
- Release content (**confirmed 2026-06-12 — see acceptance C3/D3**): the release
  content commit is **`66d2531`**, which includes the audited backend G12 result
  plus the gate-cleared UI/UX completion, with zh-CN preserved. This migration
  must not weaken or revert any G12 fix (IDOR, durable audit, optimistic
  concurrency, idempotency, etc.); it changes *how* the artifact is built and
  deployed, not *what* the code does. The final build+deploy (M5/M6) targets
  `66d2531`.

## Scope

In scope (deployment mechanism only):
- Image registry selection and access (ADR).
- Off-host CI that builds all CRM service images + the frontend, runs the
  existing test suite, tags by commit, and pushes/exports digest-pinned images.
- Convert `docker-compose.prod.yml` from `build:` to **image-only** (references
  tag/digest; no `build:` keys; no on-server frontend build).
- Rewrite `deploy/ops/go-live-runbook.md` to pull/load a digest and run (use
  `templates/deployment-runbook.md`); remove on-server `npm/build/--build` steps.
- Release evidence capture (test results, digest→commit, deploy transcript,
  health check, rollback point).

Out of scope (must not change in this migration):
- CRM business/application logic, APIs, data model, security fixes, or any
  P0/P1 acceptance behavior. No downgrade of any prior G12 result.

## Task Package (to be detailed at the CRM change's G8)

| # | Task | Standard ref | Notes |
|---|---|---|---|
| M1 | Image registry selection ADR | §3, §9 | `docs/architecture/adr/` — registry, auth, retention, provenance; both "registry" and "export/load" forms are allowed by the standard. |
| M2 | Off-host CI pipeline (`templates/cicd-pipeline.md`) | §1.2, §3 | Build 10 Go service images + nginx frontend image for commit `66d2531`; run tests; tag `:66d2531`; export digest-pinned artifacts. |
| M3 | Image-only production compose | §1.1, §1.3 | Remove all `build:` keys; reference image tag/digest; remove on-server frontend build. |
| M4 | Digest-pinned deploy runbook (`templates/deployment-runbook.md`) | §1.3, §3 | Pull/load digest + run; no host build; backup before migrate. |
| M5 | Release evidence | §4 | Test results, per-service digest→commit, deploy transcript, post-deploy health check, named rollback point. |
| M6 | Commit-traceability verification | §1.4, §8 | Prove the running images map to the audited commit (digest, not moving tag). |

## Gate Path

This change affects deployment architecture, so it re-enters at **G5**
(Architecture / deployment + pipeline design, incl. M1 ADR) → **G7/G8** (task
plan for M1–M6 with CI/release-evidence tasks per Task Planner's new duty) →
**G9–G11** (Codex executes: CI build/push, image-only deploy, evidence) →
**G12** (independent audit verifies digest→commit, §4 evidence, and that no
on-host build was used). The no-downgrade, release, and evidence rules apply in
full; a green pipeline does not substitute for QA/Integration/Audit.

## Acceptance (this migration is done when)

- Production runs images built **off-host**, pulled/loaded by **digest**, each
  traceable to commit `66d2531`.
- `docker-compose.prod.yml` has **no `build:` keys**; the runbook has **no host
  build steps**.
- The host does not depend on GitHub source pulls or a build cache to run.
- §4 release evidence is recorded and G12-audited.
- No prior P0/P1 acceptance item or G12 fix was weakened.
