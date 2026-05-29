# Patent Readiness Standard

## Document Control

- Project: CRM System
- Owner: Project Sponsor / IP Owner
- Status: Draft Preparation Standard
- Date: 2026-05-29
- Applies To: Invention discovery, confidentiality, technical disclosure,
  patentability review, and patent application preparation.

## Purpose

This document defines how the CRM project preserves patent options while the
software is designed and developed. It does not determine patentability and
does not replace a patent attorney or patent agency.

## Critical Confidentiality Rule

Potential inventions must be reviewed before public disclosure.

Do not commit unpublished invention disclosures, draft claims, special
algorithm details, architecture novelty arguments, benchmark results, or
customer-specific technical advantages to the public GitHub repository until
the sponsor has decided whether to file a patent application.

The current public repository should contain product/process/design materials
only. Patent-sensitive details must be stored in a private controlled package.

## Candidate Patent Areas For This CRM

The project should watch for genuine technical solutions, not ordinary business
rules. Possible candidate areas may include:

- permission and ownership transfer mechanisms with auditable consistency,
- quote-contract-payment state consistency algorithms,
- duplicate warning and safe merge-prevention mechanisms,
- evidence traceability automation across PRD, MDA, task, test, integration,
  and audit artifacts,
- secure import/export validation with row-level isolation and rollback
  evidence,
- CRM workflow risk detection or reminder prioritization if a technical method
  is invented later.

These are only watch areas. They are not patent claims and must not be treated
as patentable without professional review.

## Invention Disclosure Template

Create one private disclosure per candidate invention:

| Field | Required Content |
|---|---|
| Disclosure ID | Stable private ID, for example `INV-001` |
| Title | Technical title, not marketing name |
| Inventors | Contributors who made technical creative contributions |
| Owner / Assignee | Individual or company expected to own rights |
| Problem | Technical problem being solved |
| Existing approaches | Known alternatives and their limitations |
| Technical solution | System structure, data flow, algorithm, state machine, or protocol |
| Key novelty | What is technically different |
| Beneficial effect | Measurable technical effect, reliability, security, performance, consistency, or resource improvement |
| Implementation evidence | commit, prototype, experiment, design, or model reference |
| Public disclosure status | none / GitHub / demo / article / customer meeting / other |
| Confidentiality status | private / NDA / public |
| Filing recommendation | file / hold / abandon / consult agent |
| Review owner | sponsor, architect, patent agent, legal counsel |

## Private Package Structure

```text
private-ip-package/
  patents/
    INV-001/
      01-invention-disclosure.md
      02-prior-art-notes.md
      03-technical-diagrams/
      04-implementation-evidence/
      05-review-decision.md
      06-application-draft/
```

Do not place this package in the public repository.

## Patent Application File Awareness

For China patent filing, the official guidance distinguishes invention,
utility model, and design patent materials. For invention patent applications,
the application documents include request, abstract, claims, specification,
and drawings when needed. For utility model applications, request, abstract,
claims, specification, and drawings are required. For design patents, request,
images/photos, and a brief description are required.

The CRM project is software-centered, so the most likely route, if any, would
be invention-related technical solution review. A patent professional should
determine filing type and claim strategy.

## Public Disclosure Control

Before publishing any of the following, run patent review:

- new algorithm or technical method,
- architecture claimed as novel,
- workflow consistency mechanism with technical effect,
- security or audit mechanism with technical implementation details,
- benchmark or experimental proof of technical improvement,
- screenshots or diagrams that reveal a candidate invention,
- open-source code implementing the candidate invention.

## Readiness Checklist

| Item | Required Evidence | Status |
|---|---|---|
| Patent-sensitive areas identified | invention watch list | Pending |
| Private disclosure storage exists | private controlled location | Pending |
| Contributor/inventor record maintained | contribution log | Pending |
| Public disclosure review added to release checklist | release checklist item | Pending |
| Prior art search performed | search notes | Pending per candidate |
| Patentability review performed | agent/legal review | Pending per candidate |
| Filing decision recorded | file/hold/abandon decision | Pending per candidate |
| Application file prepared | official draft package | Pending per filed candidate |

## References

- CNIPA patent application matters:
  https://www.cnipa.gov.cn/art/2020/6/5/art_1517_92472.html

