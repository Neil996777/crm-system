# Security Compliance Review — Operator Access (pre-G8 condition)

## Document Control

- Project: CRM System
- Role: Security Compliance (planning)
- Date: 2026-06-01
- Trigger: G7 gate review recorded a pre-G8 condition — Security Compliance must
  review the operator-access design (SSH access, key ownership, sudo boundary)
  before G8 implementation tasks are approved (`deployment-notes.md` "Operator
  Access"; recorded on `delivery/tasks.md` TASK-039 field 17).
- Decision: **Approved with conditions** (G8 implementation tasks may be approved).
- Archive note: Review evidence only. Not design authority.

## What was reviewed

The committed operator-access requirements:
- `docs/architecture/deployment-notes.md` "Operator Access" (lines ~229–237):
  - Long-term production operation must not rely on routine root SSH use.
  - Root may be used only for initial provisioning or emergency recovery.
  - A named deploy/ops user with least required privileges must exist before
    production release.
  - SSH access, key ownership, and sudo boundary must be reviewed by Security
    Compliance before G8 implementation tasks are approved.
- `delivery/tasks.md` TASK-039 — now tasks the operator-access artifact
  (`deploy/ops/operator-access.md`), with DoD requiring the named least-privilege
  deploy/ops user, root restricted to provisioning/emergency, and documented SSH
  key ownership + sudo boundary; owner `infrastructure-ops`.
- `docs/architecture/deployment-notes.md` "Network Exposure" — PostgreSQL/internal
  ports/backup dir/secrets/admin-debug endpoints forbidden from public exposure.

## Assessment

The operator-access design is adequate and consistent with the security baseline:
- Least-privilege deploy/ops account + root-only-for-provisioning/emergency is the
  correct posture and matches SEC-* least-privilege intent.
- The requirement is now an explicit task with a checkable DoD (TASK-039), not a
  loose reference — so it cannot be silently skipped at G9.
- It composes with the forbidden-public-exposure rule and the S2S signed-token
  service boundary (no reliance on host network trust).

No design defect found. The design does not, by itself, weaken any P0/P1 control.

## Conditions carried into G9/G11 (must be satisfied as implementation evidence)

1. The named deploy/ops user and its sudo boundary must be created and recorded in
   `deploy/ops/operator-access.md` during G9 (TASK-039), not deferred to launch.
2. SSH key ownership/rotation and the prohibition of routine root SSH must be
   evidenced at G11 alongside the other ACC-017 release evidence.
3. Secrets/config files and backup directory permissions must be confirmed
   non-public and non-world-readable (deployment-notes "Network Exposure").
4. G12 audit re-verifies operator-access evidence before any release decision.

## Decision

**Approved with conditions.** The operator-access design satisfies the pre-G8
condition; G8 implementation tasks may be approved. The four conditions above are
implementation/release evidence to be produced at G9/G11 and re-audited at G12.
