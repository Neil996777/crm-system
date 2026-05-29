# Evidence Register Template

## Purpose

Use this file as the project evidence ledger. Each entry must be reproducible,
traceable, and tied to a requirement, Gate, acceptance item, task, test, review,
or release decision.

## Evidence Rules

- Evidence must point to a file, commit, command result, screenshot, report, or
  signed/approved decision.
- Evidence must not contain secrets, private customer data, unpublished
  patent-sensitive details, or identity documents.
- Evidence for P0/P1 work must be sufficient for QA, integration, and audit to
  reproduce or inspect the result.

## Register

| Evidence ID | Date | Phase / Gate | Related IDs | Evidence Type | Location | Owner | Status | Notes |
|---|---|---|---|---|---|---|---|---|
| EV-001 | 2026-05-29 | Architecture Reset | PROJECT_CONTEXT | Git baseline | commit `0cddbc6` | Process Owner | Ready | Initial public repository baseline after discarded engineering artifacts were removed. |

## Evidence Types

Allowed values:

- Document
- Review Decision
- Repair Note
- Commit
- Command Result
- Test Report
- Screenshot
- Manual Verification
- Integration Evidence
- Audit Finding
- Release Decision
- Registration Material
- Patent Disclosure Material

## Closure Standard

An evidence entry may be marked Ready only when:

- location is stable,
- owner is named,
- related IDs are present,
- content can be inspected by a reviewer,
- no prohibited secret or private material is included.

