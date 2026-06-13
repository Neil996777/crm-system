# Decision Log

| ID | Date | Decision | Reason | Owner |
|---|---|---|---|---|
| DEC-001 | 2026-05-26 | The CRM is for ToB sales. | Sponsor confirmed the sales model during requirement discussion. | Product Manager |
| DEC-002 | 2026-05-26 | The CRM must support team collaboration, not only solo use. | Sponsor clarified team collaboration is required. | Product Manager |
| DEC-003 | 2026-05-26 | The committed release must cover the complete business loop. | Sponsor stated the committed release must cover the complete CRM business loop. | Product Manager |
| DEC-004 | 2026-05-26 | The project goal is full production launch. | Sponsor clarified the final project goal is production launch, not a demo. | Product Manager |
| DEC-005 | 2026-05-26 | The committed role model has three roles: Administrator, Sales Manager, and Sales. | Sponsor accepted the three-layer permission model as sufficient. | Product Manager / Security Compliance |
| DEC-006 | 2026-05-26 | Quote, contract, and payment management are included in committed P0. _[Amended 2026-06-01 by DEC-018 (one quote per opportunity) and DEC-019 (payment decoupled from Won); quote/contract/payment remain in committed P0 scope.]_ | Sponsor confirmed the complete loop must include quote, contract, and payment management. | Product Manager |
| DEC-007 | 2026-05-26 | Committed contract management is record-based and does not include approval workflow, electronic signature, or contract template generation. | Sponsor accepted the recommended boundary. | Product Manager |
| DEC-008 | 2026-05-26 | Core CRM paths must use persistent data and cannot be satisfied by mock, static-only, TODO, or non-persistent behavior. | Workspace no-downgrade rule and production-launch goal. | Product Manager / QA TDD / Audit |
| DEC-009 | 2026-05-26 | The committed scope is single team / single organization. | Keeps team collaboration clear without introducing multi-tenant SaaS complexity before architecture. | Product Manager |
| DEC-010 | 2026-05-26 | Sales Manager can view and manage all team records; Sales can view and manage owned/assigned records only. | Resolves G3 permission testability blocker. | Product Manager / Business Analyst |
| DEC-011 | 2026-05-26 | Core CRM records cannot be hard-deleted in the committed scope. | Preserves data integrity, history, and auditability. | Business Analyst / Security Compliance |
| DEC-012 | 2026-05-26 | Opportunity is Won only after full payment is recorded; Won and Lost are terminal in the committed scope. _[Superseded in part 2026-06-01 by DEC-017 (Won = contract Signed) and DEC-019 (payment decoupled from Won); Won/Lost remain terminal.]_ | Makes closure behavior testable for quote-contract-payment loop. | Product Manager / Business Analyst |
| DEC-013 | 2026-05-26 | The committed money model uses one currency and excludes tax, discount, and multi-currency automation from P0/P1. | Keeps quote, contract, and payment acceptance testable. | Product Manager / Business Analyst |
| DEC-014 | 2026-05-26 | Overpayment is blocked; contract amount may differ from accepted quote only with a recorded difference reason. | Defines core payment and amount-integrity behavior. | Business Analyst / QA TDD |
| DEC-015 | 2026-05-26 | P1 import/export is CSV only and P1 reminders are in-app only. | Defines minimum committed behavior for committed P1 items. | Product Manager / QA TDD |
| DEC-016 | 2026-05-26 | Contract notes are P0 required; contract attachment upload is not required for P0. | Keeps record-based contract management testable without requiring storage architecture before G5. | Product Manager |
| DEC-017 | 2026-06-01 | Opportunity is **Won when the related contract is Signed**; full payment is **not** a Won precondition. Won and Lost remain terminal (no reopen). The `Payment In Progress` pipeline stage is removed. Supersedes the payment-gate clause of DEC-012. | Aligns "won" with industry practice (deal won at signing); payment collection is post-sale follow-up, not a closing gate. Formal Scope Change by User. | Product Manager / Business Analyst |
| DEC-018 | 2026-06-01 | Each opportunity has **exactly one quote**; the system records the quote result, not multi-round negotiation. Amends DEC-006 and BR-006 (removes multiple-quotes and the one-Accepted-at-a-time constraint). | Opportunity↔quote is 1:1 in the committed flow; simplifies the model and removes the second-accept ambiguity. Formal Scope Change by User. | Product Manager / Business Analyst |
| DEC-019 | 2026-06-01 | Payment tracking (plans, actual payments, status, overdue reminders, reports) is **retained but decoupled from Won** — it is post-sale collection/visibility, not a closing gate. Overpayment still blocked; single currency unchanged. Amends DEC-012/DEC-006. | Sales/manager need collection follow-up and team receivables visibility; full AR/accounting remains finance's concern. Formal Scope Change by User. | Product Manager / Business Analyst |
| DEC-020 | 2026-06-01 | The separate Opportunity **`Status` field is removed**; Pipeline Stage (including terminal Won/Lost) is the sole opportunity lifecycle dimension. Amends PRD-007/REQ-007. | "Status" was never assigned values and is redundant with Stage. Formal Scope Change by User. | Product Manager / Business Analyst |
| DEC-021 | 2026-06-01 | **Web frontend stack = React + TypeScript** (talks only to gateway-bff per frontend-backend-contract). Complements the committed Go backend microservices (ADR-ARCH-001). | Stack was left open at G5; pinned at G7 so task artifacts can specify concrete, Codex-executable file changes for UI work. Decided by User. | User / Architecture |
| DEC-022 | 2026-06-01 | **Repo layout + tooling**: single monorepo — `services/<service>/` each an independent Go module + `frontend/` (React+TS); PostgreSQL via `golang-migrate`; Go `testing` + `testcontainers` for unit/integration; Playwright for E2E; Docker Compose orchestration (ADR-ARCH-001). | Concrete, conventional defaults so G8→Codex can build file-for-file without guessing. Decided by User. | User / Architecture |
| DEC-023 | 2026-06-12 | **CI/CD migration release content commit = `66d2531`.** This is the only application commit targeted by the release-mechanism migration. | Release owner confirmed D3 in `delivery/cicd-migration-acceptance.md`: audited backend G12 plus gate-cleared UI/UX completion. This supersedes earlier non-specific HEAD wording in CI/CD planning drafts. | Release Owner / Product Manager |
| DEC-024 | 2026-06-12 | **CI/CD image distribution channel = export/load** (`docker save` -> `scp` -> `docker load`) for this single-host CRM release. | Company infrastructure has no approved registry and the target is one production host. Export/load is allowed by the CI/CD standard, adds no new infrastructure, and is sufficient when paired with digest manifest, checksums, and G11/G12 evidence. See ADR-CICD-001. | Architecture / Infrastructure Ops |
| DEC-025 | 2026-06-12 | **Frontend release unit = nginx container image**, not loose `dist` files copied to the production host. | Keeps frontend and 10 Go services on one digest-recorded image path, preserves commit traceability, and avoids a second non-image deployment unit. See ADR-CICD-001. | Architecture / Infrastructure Ops |
| DEC-026 | 2026-06-13 | **Production DB/service secrets are generated ON the host at G11 deploy time (strong random, e.g. `openssl rand`) and stored only in `/opt/crm-system/secrets/prod.env` (0600).** Values are never written to the repo, image labels, deploy transcript, or release evidence (paths only) — C6. Each service password must be identical in its `CRM_DB_PASSWORD_<SVC>` and matching `<SVC>_DATABASE_URL`; no `*_dev_password` (the migration runner rejects dev values). | Release owner confirmed for the CRM CI/CD G11 deploy: host was wiped 2026-06-06 so secrets are provisioned fresh; on-host generation keeps secrets off every shared/tracked surface and out of Claude/Codex sessions. | Release Owner / Infrastructure Ops |

## Amendments — 2026-06-01 Formal Scope Change by User

This change set was raised by the owner while reviewing the three G6 open blockers
(BLK-001/002/003); it revises committed P0 rules via the No-Downgrade-permitted
Formal-Scope-Change mechanism. Original decisions are retained above and annotated;
the new decisions DEC-017..020 govern going forward.

| New | Supersedes / Amends | Old rule | New rule |
|---|---|---|---|
| DEC-017 | DEC-012 (payment-gate clause) | Won only after full payment recorded | Won when contract Signed; payment not a gate; `Payment In Progress` stage removed |
| DEC-018 | DEC-006, BR-006 | Multiple quotes per opportunity, one Accepted at a time | Exactly one quote per opportunity; record result not negotiation |
| DEC-019 | DEC-012, DEC-006 | Full payment was the Won precondition | Payment tracking retained but decoupled from Won (post-sale follow-up) |
| DEC-020 | PRD-007 / REQ-007 | Opportunity carries `stage` and `status` | `status` removed; Stage is the sole lifecycle dimension |

Resolves blockers: BLK-001 (by DEC-020), BLK-002 (by DEC-017/019 — Won no longer
depends on multi-plan payment aggregation), BLK-003 (by DEC-018 — no second quote
to accept). Downstream cascade tracked in `planning/scope-change-2026-06-01-TEMP.md`.
