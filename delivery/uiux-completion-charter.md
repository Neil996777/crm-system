# CRM UI/UX Completion — Change Charter (Follow-On Project Change)

Status: **Chartered (planning).** Not started. Planning artifact only; no
implementation here. No deployment.
Date: 2026-06-06
Owner (planning): Claude (CRM PM / UX / UI / Architecture-impact roles).
Execution: Codex (G9–G11). Audit: Claude (G12).
Authority: project gates (`../../../workflows/software-delivery.md`),
`../../../company/operating-model.md` (no-downgrade),
`../../../standards/cross-cutting-requirements-standard.md`.
Companions: `delivery/cicd-migration-brief.md`,
`delivery/cicd-current-state-gap-analysis.md` (release-content premise, D3).

This charter exists because the UI/UX work was floating with **no plan and no
gate-status entry**. It registers the work, classifies it correctly, and defines
the gate path. It does not perform design or implementation.

---

## 1. What this change is (verified scope, not chat history)

The work in progress under `docs/ux-ui/` is a **completion of the specified
UX/UI design and its realization**, not a new feature:

- New `docs/ux-ui/design-system.md` (36 KB) — design system not previously
  realized.
- `docs/ux-ui/interaction-spec.md` (+936 lines) and `screen-state-spec.md`
  (+93 lines) — elaborated interaction / screen-state behavior.
- `docs/ux-ui/mockups/` — dashboard design iterations (v2–v7, variants a/b) +
  `_src` HTML.

These trace to the UX/UI design gates (**G4b UX / G4c UI**) and their specs, and
they back the committed dashboard/overview/reports surface.

## 2. Why it is "completing committed scope" (not new, not optional)

- The dashboard / team-overview / basic-reports capability is **committed P1**:
  `docs/product/acceptance-matrix.md` **ACC-018** (View team overview, P1) and
  **ACC-023** (View basic sales reports, P1), capability **CAP-009**
  (`docs/product/business-capability-map.md`), PRD-018 / PRD-023 / NFR-007.
- The **functional** acceptance for these was implemented and **G12-passed
  (2026-06-04)** — the pages exist and work: `frontend/src/pages/WorkOverview.tsx`,
  `frontend/src/pages/reports/ManagerOverview.tsx`,
  `frontend/src/pages/reports/BasicReports.tsx`, `frontend/src/api/reports.ts`.
- What was **not** completed is the **specified UX/UI design realization** — the
  design system and the dashboard/interaction design quality the ux-ui specs call
  for. Per release owner: this was **prior-intended work that should have been
  done before**, surfaced (not created) by the go-live exposure. It is therefore
  governed by **no-downgrade**: valid states are only `Done`, `Blocked`, or
  `Formal Scope Change by User`. It may not be silently dropped or shipped as a
  placeholder.

Contrast with the zh-CN localization (gate-status 2026-06-05): that was a
genuinely **NEW** requirement (UI language never specified). This is **not** new —
it realizes design that already backs committed capabilities.

## 3. Honest classification points to confirm (PM / UX / UI / Audit)

These are open and must be resolved as part of chartering — recorded here, not
assumed:

- **C1 — Visual-only vs behavior-changing.** If the elaborated
  interaction-spec / screen-state-spec only refine presentation, the path is
  lighter. If they materially change committed P0/P1 **interaction/state
  behavior**, the affected acceptance items must be reviewed/updated and the path
  is fuller. PM + UX to delimit the boundary.
- **C2 — Prior G12 relationship.** G12 audited the **functional** acceptance
  (ACC-018/023 etc.), which passed; the **design-realization** gap sat outside
  the strict acceptance items audited. This should be **recorded honestly**;
  whether the G12 record needs an annotation (not a re-open of passed functional
  items) is the release owner's / Audit's call. No silent rewrite of history.
- **C3 — Acceptance evidence for "design done."** Design realization needs a
  verifiable "done" definition (which screens, which states, accessibility /
  responsive targets) so it is not an open-ended polish loop. Tie to the
  cross-cutting standard (accessibility, target devices/responsive baseline are
  mandatory categories).

## 4. Proposed gate path (to confirm at charter sign-off)

Frontend/design change; **no new services, no backend/API/data-model change**
expected — Architecture confirms zero service/contract impact (a quick G5 impact
check, not a redesign).

1. **G4b / G4c (UX/UI design closure)** for the elaborated specs + design-system +
   chosen mockup direction → sign-off by PM, UX, UI (Security if any
   auth/permission-bearing surface changes; Accessibility per cross-cutting std).
2. **G7 / G8** — frontend implementation task plan (Codex-executable), tied to
   ACC-018/023 + the affected UI surface; explicit "design realized" acceptance.
3. **G9–G11** — Codex implements + QA (E2E kept green; zh-CN preserved).
4. **G12** — Claude independent audit (design realized to spec; no P0/P1
   functional downgrade; no security/zh-CN regression).

## 5. Relationship to the other in-flight change (CI/CD migration)

- This charter defines release **CONTENT**; the CI/CD migration
  (`delivery/cicd-migration-plan.md`) defines release **MECHANISM**.
- The CI/CD mechanism design (M1–M4) is independent and may proceed in parallel.
  The **final compliant build+deploy** (M5/M6) ships the commit that includes
  **this gate-cleared UI/UX completion** — not the current as-is HEAD
  (see `delivery/cicd-current-state-gap-analysis.md` D3).
- Recommended primary sequence: **UI/UX completion is the main line** (it
  determines what ships); CI/CD mechanism is built alongside; deploy last.

## 6. Scope guardrails

- **Frontend / design only.** No backend, API, data-model, or service change.
- **No downgrade** of any P0/P1 functional acceptance or any G12 security fix
  (IDOR, durable audit, optimistic concurrency, idempotency, etc.).
- **Preserve zh-CN localization** (phase 1 + 2, live as of 2026-06-05) — the
  redesigned UI must not regress to English or change enum/role compare-values.
- No deployment as part of this charter; go-live is a separate, gated step.

## 7. Registration

- Registered in `planning/gate-status.md` Handoff Log (2026-06-06).
- Also flagged: the **CI/CD migration** follow-on change is likewise not yet in
  the Handoff Log and should be registered for the same reason.
