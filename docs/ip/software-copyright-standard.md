# Software Copyright Registration Standard

## Document Control

- Project: CRM System
- Registration Target: Computer Software Copyright Registration in China
- Owner: Project Sponsor / IP Owner
- Status: Draft Preparation Standard
- Date: 2026-05-29

## Purpose

This document defines how the CRM project will prepare software copyright
registration materials after the software reaches a stable releasable version.
It does not replace the official application form, agency review, or legal
advice.

## Registration Timing

Prepare the formal software copyright package when all are true:

- software name and version are stable,
- executable or deployable version exists,
- core source code is complete enough to represent the claimed software,
- user manual or operation manual exists,
- ownership and contributor records are clear,
- no secrets or private data are included in submitted source/document samples.

## Required Registration Inputs

| Material | Project Source | Standard |
|---|---|---|
| Software full name | PRD / release decision | Must be consistent across application form, manual, source header, and package. |
| Software short name | Sponsor decision | Optional but must be consistent if used. |
| Version number | release tag | Must match source/document samples and application form. |
| Development completion date | release evidence | Must be supportable by Git history and release records. |
| First publication date | release/publication evidence | Use only if the software was actually published. |
| Copyright owner identity | sponsor/company identity records | Do not store identity documents in this public repo. |
| Ownership basis | employment, commission, assignment, or self-development record | Must be documented privately if not purely individual self-development. |
| Source program sample | source repository at release tag | Remove secrets, keys, credentials, and customer data. |
| Documentation sample | user manual, operation manual, or design/user documentation | Must match software version and actual features. |

## Official Material Rules To Preserve

Based on current official/service guidance:

- The application normally includes the registration application form, software
  identification material, and related proof documents.
- Software identification material includes program and document identification
  material.
- Program and document identification material is generally the first and last
  30 consecutive pages of source program and one qualifying document. If the
  entire program or document is under 60 pages, submit the full material.
- Unless a special rule applies, source program pages should contain at least
  50 lines per page, and document pages should contain at least 30 lines per
  page.
- Materials should use A4 paper, Chinese where required, consistent version
  naming, and required signature or seal rules.

Always verify the current China Copyright Protection Center / local service
requirements before filing.

## Project Package Structure

Formal packages must be prepared outside the public GitHub repository:

```text
private-ip-package/
  software-copyright/
    crm-system-vX.Y.Z/
      01-application-form/
      02-owner-identity/
      03-ownership-proof/
      04-source-identification/
      05-document-identification/
      06-release-evidence/
      07-submission-record/
```

## Public Repository Boundary

Allowed in public repo:

- process standard,
- source code after intentional open-source decision,
- public documentation,
- public release tags,
- non-sensitive evidence references.

Not allowed in public repo:

- ID card, business license scans, seals, signatures,
- account credentials, API keys, database URLs, private certificates,
- customer or sales data,
- unpublished invention disclosure details,
- draft patent claims,
- formal signed registration application forms.

## Source Code Sample Preparation Standard

When preparing source identification material:

- use the release tag intended for registration,
- include only code actually owned or properly licensed,
- exclude dependency folders and generated vendor code unless required and
  legally owned/allowed,
- replace or omit secrets and environment-specific private values,
- keep filenames, version headers, and package name consistent,
- record the exact commit hash used for extraction.

## Documentation Sample Preparation Standard

The preferred document sample for this CRM is a user or operation manual,
because it can align with final screens and workflows without exposing internal
patent-sensitive implementation details.

The document sample should include:

- software overview,
- roles and permissions,
- installation or access method,
- login and account behavior,
- lead/customer/contact workflows,
- opportunity/quote/contract/payment workflows,
- activity, task, reminder, report, import/export, and audit-log behavior,
- version and software name matching the application form.

## Readiness Checklist

| Item | Required Evidence | Status |
|---|---|---|
| Software name confirmed | sponsor decision | Pending |
| Version number confirmed | release tag | Pending |
| Copyright owner confirmed | private owner record | Pending |
| Ownership chain clear | employment/commission/assignment/self-development proof | Pending |
| Source code complete | release commit | Pending implementation |
| User/operation manual complete | final manual | Pending |
| Source sample extracted | first/last 30 pages or full source if under threshold | Pending |
| Document sample extracted | first/last 30 pages or full document if under threshold | Pending |
| Secrets and private data removed | review checklist | Pending |
| Official form generated | official registration system | Pending |
| Submission receipt archived | private package | Pending |

## References

- Beijing Government Service page summarizing official software copyright
  registration material rules:
  https://banshi.beijing.gov.cn/pubtask/task/1/110000000000/3e283672-76be-4c8c-98e8-0bebe9bd06bf.html?locationCode=110000000000

