# CRM CI/CD Migration ACC Traceability

Status: G5/G8 handoff traceability for Claude audit.

## Deliverables

| ID | Deliverable | Purpose |
|---|---|---|
| DLV-CICD-ADR | `docs/architecture/adr/ADR-CICD-001-image-channel-and-frontend-runtime.md` | M1 registry/channel ADR, D1/D2 decisions, image set, tag/digest strategy. |
| DLV-CICD-ARCH | `docs/architecture/cicd-release-architecture.md` | M2/M3/M4/M5/M6 architecture design, compose diff, runbook design, Infrastructure Ops review. |
| DLV-CICD-TASKS | `delivery/cicd-migration-g8-task-package.md` | G8 task package for M1-M6. |
| DLV-CICD-EVIDENCE | `delivery/cicd-release-evidence-template.md` | Release evidence template for G10/G11. |
| DLV-CICD-DECLOG | `docs/product/decision-log.md` | D1/D2/D3 durable project decisions. |
| DLV-CICD-BRIEFPLAN | `delivery/cicd-migration-brief.md`, `delivery/cicd-migration-plan.md`, `delivery/cicd-current-state-gap-analysis.md` | Aligns old release-content wording to `66d2531`. |

## ACC-CICD Mapping

| ACC | Required proof | Delivery mapping |
|---|---|---|
| ACC-CICD-001 | Off-host build; no prod build/compile. | DLV-CICD-ADR rejects host build; DLV-CICD-ARCH M2/M3/M4; DLV-CICD-TASKS M2/M3/M4; DLV-CICD-EVIDENCE no-host-build audit. |
| ACC-CICD-002 | CI build + tests + images + export availability. | DLV-CICD-ARCH M2 CI stages and outputs; DLV-CICD-TASKS M2; DLV-CICD-EVIDENCE CI test results and image manifest. |
| ACC-CICD-003 | CD load/run only, specified digest/tag, no host build or checkout. | DLV-CICD-ADR export/load decision; DLV-CICD-ARCH M4; DLV-CICD-TASKS M4; DLV-CICD-EVIDENCE deploy transcript and no-host-build audit. |
| ACC-CICD-004 | Digest to commit traceability; no `latest` as release identity. | DLV-CICD-ADR tag/digest strategy; DLV-CICD-ARCH M6; DLV-CICD-TASKS M2/M6; DLV-CICD-EVIDENCE digest to commit table. |
| ACC-CICD-005 | Image-only Compose, no `build:`, frontend image. | DLV-CICD-ARCH M3 compose design and diff table; DLV-CICD-TASKS M3; DLV-CICD-EVIDENCE compose static check. |
| ACC-CICD-006 | Five release evidence items. | DLV-CICD-ARCH M5; DLV-CICD-TASKS M5; DLV-CICD-EVIDENCE sections 1-5 and G11 return notes. |
| ACC-CICD-007 | Production host not dependent on source checkout/build cache. | DLV-CICD-ARCH migration release artifact design; DLV-CICD-TASKS M3/M4/M6; DLV-CICD-EVIDENCE source-dependence audit. |
| ACC-CICD-008 | Zero application diff / no downgrade; release content commit fixed. | DLV-CICD-ADR release commit and scope; DLV-CICD-ARCH scope boundary; DLV-CICD-TASKS global constraints; DLV-CICD-DECLOG DEC-023..025; DLV-CICD-BRIEFPLAN wording alignment. |

## Reviewer Hooks

| Reviewer | What to inspect |
|---|---|
| Claude G8 handoff audit | ACC rows above, C1-C6 constraints, D1/D2/D3 consistency, no implementation before G8. |
| Infrastructure Ops | export/load channel, loopback ports, 80/443 co-location, disk retention, secret handling, backup/restore, monitoring. |
| QA Execution | CI suite includes backend tests + current e2e suite; no skip/only/slow weakening. |
| Integration Owner | G11 evidence has digest mapping, deploy transcript, health checks, rollback point. |
| Audit G12 | Running image digests and labels match commit `66d2531`; no host source build occurred. |
