# CMMI-Oriented Process Standard

## Document Control

- Project: CRM System
- Owner: Project Sponsor / Process Owner
- Status: Draft Internal Standard
- Date: 2026-05-29
- Applies To: Product, business, UX/UI, security, architecture, MDA, task
  planning, implementation, QA, integration, audit, release, and maintenance.

## Purpose

This document defines the CRM project's internal process and evidence standard
for CMMI-style maturity readiness. It is not an appraisal result and must not
be represented as CMMI certification.

ISACA owns and governs CMMI. Formal appraisal work must be performed through
authorized CMMI channels. This project uses this document only to keep process
evidence complete, repeatable, and reviewable.

## Governance Rules

- Workspace no-downgrade and Gate rules remain mandatory.
- P0/P1 scope cannot be downgraded, deleted, merged away, weakened, or accepted
  as partial.
- Each Gate must have owner, required reviewers, pass condition, decision, and
  evidence.
- Every P0/P1 requirement must trace to acceptance, design/model, task, test,
  integration evidence, audit result, and release decision.
- No implementation starts before G8 passes.
- No release is valid with open P0/P1 blocker, missing evidence, or known mock,
  stub, TODO, static-only, or non-persistent core path.

## Required Process Assets

| Area | Required Project Assets | Evidence |
|---|---|---|
| Requirements | PRD, requirements, open questions, out-of-scope, decision log, acceptance matrix | requirement version, owner, status, priority, acceptance mapping |
| Project Planning | delivery workflow, Gate records, task plan after MDA | plan baseline, dependency map, blocker list, scope-change records |
| Monitoring And Control | status updates, blocker register, Gate decisions | issue age, owner, severity, closure evidence |
| Configuration Management | Git history, branch policy, release tags, artifact index | commit IDs, review records, versioned documents |
| Supplier / Tool Control | dependency inventory, third-party license register, deployment provider decisions | source, version, license, risk owner |
| Measurement | quality metrics, test coverage indicators, defect trends, cycle time | metric definition, collection date, trend, action taken |
| Process Quality Assurance | QA reports, review checklists, audit reports | reviewer, findings, severity, closure result |
| Verification | unit, integration, contract, E2E, manual verification | test ID, acceptance ID, environment, command/result |
| Validation | user-flow evidence, integration report, release smoke checks | role, starting state, steps, expected/actual result |
| Risk Management | risk register, security risks, compliance risks, patent/soft-copyright risks | mitigation owner, trigger, status |
| Decision Management | decision log, architecture decisions, tradeoff records | options, rationale, date, owner |

## Gate Evidence Standard

Each Gate record must include:

- Gate ID and transition.
- Owner and required reviewers.
- Input document list with version or commit reference.
- Findings grouped by severity.
- Explicit pass/block decision.
- P0/P1 blocker list.
- Required repair actions and owner.
- Re-review result after repair.
- Final evidence location.

## Traceability Standard

Each P0/P1 capability must maintain this chain:

```text
Requirement ID
  -> Acceptance ID
  -> Business / UX / UI / Security source
  -> Architecture decision
  -> MDA CIM/PIM/PSM model IDs
  -> Task ID
  -> Test ID / test file / manual verification
  -> Integration evidence
  -> Audit result
  -> Release decision
```

If any link is missing after its phase is expected to exist, the item is not
Done. If the missing link blocks a P0/P1 item, the item must be marked Blocked.

## Review And Audit Standard

- Reviewers must not review only their own work when a handoff reviewer exists.
- Intake review is performed by the receiving role before accepting upstream
  work.
- Review findings are archived outside active design directories.
- Design documents must contain current design truth only, not obsolete review
  noise.
- Repair notes must state what was repaired and which finding was closed.
- Old discarded architecture, MDA, task, or implementation artifacts are not
  current design authority after the 2026-05-29 architecture reset.

## Configuration Management Standard

- All project documents and source files must be version controlled unless they
  contain secrets, generated dependency folders, build outputs, or private
  personal identity material.
- Commit messages must describe the actual change.
- Release candidates must be tagged after QA, integration, and audit evidence
  is complete.
- Public repositories must not contain secrets, credentials, private customer
  data, unpublished patent-sensitive invention details, or personal identity
  documents used for registration.

## Measurement Standard

At minimum, track:

- P0/P1 acceptance coverage.
- Open blocker count by severity and age.
- Requirements changed after acceptance.
- Review findings opened/closed by Gate.
- Test coverage by acceptance ID.
- Defects by severity and root cause.
- Rework caused by unclear requirements, architecture gaps, or test gaps.

## CMMI Readiness Checklist

| Item | Required Evidence | Status |
|---|---|---|
| Process owner assigned | named owner and responsibilities | Pending |
| Gate evidence complete | Gate records from G1 onward | Partial |
| Acceptance traceability complete | acceptance-to-task-to-test chain after G8 | Pending new architecture/MDA/G8 |
| Configuration baseline exists | Git repository and initial baseline | Ready |
| Quality assurance process exists | QA plan and review standard after new MDA | Pending |
| Measurement definitions exist | metric names, collection method, owner | Pending |
| Risk and blocker management exists | active blocker/risk register | Pending new architecture |
| External appraisal scope defined | target CMMI domain, organization unit, appraisal provider | Not started |

## References

- ISACA CMMI Performance Solutions: https://www.isaca.org/enterprise/cmmi-performance-solutions

