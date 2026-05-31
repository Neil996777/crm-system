# Process Gap Register (Project Layer)

- Project: CRM System
- Date: 2026-05-30
- Scope: Project-layer resolution of workflow gaps found during the
  2026-05-30 workflow review.
- Authority note: This register records **project-level strengthenings** only.
  Per `AGENTS.md`, project rules may only strengthen workspace rules, never
  weaken them. The matching **company-layer** fixes (editing
  `company/operating-model.md`, `workflows/`, `agents/`) are intentionally
  deferred to a separate workspace-level discussion and are NOT applied here.

## Background

A 2026-05-30 review of the delivery workflow against this project's execution
found that the 2026-05-29 default-microservice-governance policy was folded into
`company/operating-model.md` but the company **gate reviewer matrix and new
governance roles were not fully propagated**. This left two governance-critical
roles with responsibilities but no formal gate seat, plus a gate-semantics
ambiguity and no defined reset procedure. Until the company layer is updated,
this project closes the gaps locally by strengthening its own G5+ review rules.

## Gaps And Project-Layer Resolution

### GAP-PROC-001 — Infrastructure Ops has no formal gate reviewer seat (Severity: Medium)

- Finding: `company/operating-model.md` Delivery Gates table never lists
  `infrastructure-ops` as a gate owner or required reviewer (G1–G12), yet the
  microservice governance policy requires G5 to define deployment and
  observability, and the Release Rule requires backup / monitoring / security
  group / TLS evidence — all Infrastructure Ops domain. In the 2026-05-29 G5
  review, Infrastructure Ops raised two P0 blockers (G5-BLK-001, G5-BLK-002).
- Project resolution: Infrastructure Ops is a **Required Reviewer** for this
  project at **G5, G8, G11, and G12**. This is recorded as a project
  strengthening. (Already partially reflected as `GAP-G8-016` and in
  `AGENTS.md` / `STANDARD-APPLICATION-INDEX.md`; this register makes the gate
  positions explicit.)
- Company follow-up (deferred): add `infrastructure-ops` to the gate reviewer
  matrix at the company layer.

### GAP-PROC-002 — `Primary Flow Owner Agent` has no gate landing (Severity: Low/Medium)

- Finding: The governance policy introduces a `Primary Flow Owner Agent` for
  cross-capability flows, but no gate review item checks it. CRM's core value
  chain (lead → opportunity → quote → contract → payment → won) is exactly such
  a cross-service flow.
- Project resolution: The project G5 (and later G6/G7/G8) review **must verify
  that each CRM cross-capability flow names exactly one `Primary Flow Owner
  Agent`** in the architecture / service documents.
- Company follow-up (deferred): add a `Primary Flow Owner Agent` check to the
  gate pass conditions at the company layer.

### GAP-PROC-003 — G5/G6 gate semantics are ambiguous (Severity: Low)

- Finding: The G5 transition label previously read "Business/UX/UI/Security
  Design -> Architecture Design" (entry), but the G5 pass condition and Phase 6
  require architecture **outputs** to already exist, and G6 again re-states
  architecture acceptance. Where architecture is formally accepted is ambiguous.
- Update (2026-05-30): The company layer split the old single G4 into G4a–G4d
  and relabeled the G5 transition to `Design Closure -> Architecture Design`.
  The label ambiguity below is partly superseded; the binding interpretation
  still holds.
- Project interpretation (binding for this project):
  - **G5 = Architecture Acceptance.** G5 passes only when the architecture
    outputs (service strategy, service list, owner agents, contracts, data
    ownership, deployment, observability, reliability) exist and required
    reviewers report no open P0/P1 blocker.
  - **G6 = MDA/PSM fidelity to the accepted architecture.** G6 does not
    re-litigate accepted architecture; it verifies PSM faithfully represents it.
- Company follow-up (deferred): relabel G5 as `Architecture Design ->
  Architecture Accepted` at the company layer.

### GAP-PROC-004 — No defined reset / scope-rollback procedure (Severity: Low)

- Finding: The workflow is forward-only (G1→G12) with no standard procedure for
  invalidating already-passed artifacts and rolling back to an upstream gate.
  The 2026-05-29 reset was handled ad hoc via `STANDARD-APPLICATION-INDEX.md`
  and `STANDARD-APPLICATION-REVIEW.md`.
- Project resolution: For this project, the canonical reset record is:
  `PROJECT_CONTEXT.md` (Reset Boundary), `STANDARD-APPLICATION-INDEX.md`
  (retained inputs + required rebuild), and `STANDARD-APPLICATION-REVIEW.md`
  (per-agent gap closure). Any future reset must update these three.
- Company follow-up (deferred): promote this pattern into a reusable
  `workflows/` reset / scope-change procedure at the company layer.

## Project G5 Required Reviewer Set (Effective)

Gate owner: Architecture. Required reviewers for this project's G5 re-review:

| Reviewer | Source |
|---|---|
| Product Manager | company G5 |
| Business Analyst | company G5 |
| UX Designer | company G5 |
| UI Designer | company G5 |
| Security Compliance | company G5 |
| Infrastructure Ops | project strengthening (GAP-PROC-001) |

Additional G5 check (GAP-PROC-002): every CRM cross-capability flow names one
`Primary Flow Owner Agent`.
