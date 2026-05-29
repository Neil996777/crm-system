# G2 Review Prep

## Gate

- Gate: G2
- Transition: Requirement Discussion -> PRD
- Owner: Product Manager
- Required Reviewers: Business Analyst, UX Designer, Security Compliance
- Status: Gate Review Prep
- Date: 2026-05-26

## G2 Pass Condition

Core users, scope, non-goals, risks, and open questions are documented.

G2 does not authorize implementation. Implementation remains blocked until G8
passes.

## Review Inputs

| Document | Purpose | Status |
|---|---|---|
| `docs/product/project-charter.md` | Business goal, target users, core loop, scope, constraints, risks. | Ready for G2 review |
| `docs/product/requirements.md` | P0/P1/P2 requirement draft and business loop coverage. | Ready for G2 review |
| `docs/product/open-questions.md` | Known unresolved product, business, permission, security, and launch questions. | Ready for G2 review |
| `docs/product/out-of-scope.md` | Explicit non-goals and deferred capabilities. | Ready for G2 review |
| `docs/product/decision-log.md` | Sponsor and product decisions from requirement discussion. | Ready for G2 review |
| `docs/product/acceptance-matrix.md` | Draft P0/P1 acceptance items for downstream G3 work. | Draft, sufficient for G2 review |

## Product Manager Review

| Check | Result | Notes |
|---|---|---|
| Business goal documented | Pass for G2 review | Production-ready ToB CRM for team sales. |
| Core users documented | Pass for G2 review | Administrator, Sales Manager, Sales. |
| v1 scope documented | Pass for G2 review | Full loop from lead through quote, contract, payment, win/loss, and history. |
| P0/P1/P2 boundaries documented | Pass for G2 review | P0/P1/P2 captured in requirements; P2 and exclusions captured in out-of-scope. |
| No-downgrade rule preserved | Pass for G2 review | P0/P1 governance stated in requirements and acceptance matrix. |
| Implementation authorization | Not authorized | No code before G8. |

## Business Analyst Review Prep

| Check | Result | Notes |
|---|---|---|
| Main business loop documented | Pass for G2 review | End-to-end ToB CRM loop is explicit. |
| Core business objects identified | Pass for G2 review | Lead, company/customer, contact, opportunity, quote, contract, payment, activity, task. |
| Business state candidates documented | Pass for G2 review | Candidate states exist; final lifecycle rules remain open for G3/G4. |
| Exceptions and edge questions captured | Pass for G2 review | Amount mismatch, overpayment, quote expiry, reopen, deletion/archive, owner transfer captured as open questions. |
| Business policy gaps hidden | No | Unconfirmed policies are recorded in open questions. |

## UX Designer Review Prep

| Check | Result | Notes |
|---|---|---|
| Core user roles are clear enough for UX discovery | Pass for G2 review | Three roles are documented. |
| Core user tasks are identifiable | Pass for G2 review | Lead, customer, opportunity, quote, contract, payment, activity, task, and history flows are in scope. |
| UX states finalized | Not yet | Loading, empty, error, success, permission-denied, conflict, and recovery states belong to G4 UX work. |
| UX blockers for G2 | None identified | G2 can proceed because UX gaps are documented as downstream work, not hidden assumptions. |

## Security Compliance Review Prep

| Check | Result | Notes |
|---|---|---|
| Security-sensitive scope identified | Pass for G2 review | Customer, contact, contract, payment, owner, and audit-sensitive data identified. |
| Role model documented | Pass for G2 review | Administrator, Sales Manager, Sales. |
| Permission details finalized | Not yet | Exact visibility and action rules are open questions for security and permission matrix work. |
| Audit-sensitive actions identified | Pass for G2 review | Ownership, stage, quote, contract, payment, and access changes are identified. |
| Security blockers for G2 | None identified | Open security questions are visible and can be resolved before G5. |

## Open Questions Handling

The open questions in `docs/product/open-questions.md` do not prevent G2 review
because G2 requires open questions to be documented, not fully resolved.

Any open question that affects a P0/P1 acceptance item must remain traceable and
must block later gate passage if unresolved when that gate requires the decision.

## Gate Prep Result

Recommendation: move to G2 review.

Decision outcome is recorded in
`archive/reviews/g2-product/g2-review-decision.md`.

## Next If G2 Passes

1. Create and complete `docs/product/prd.md`.
2. Refine `docs/product/acceptance-matrix.md` for G3.
3. Ensure every P0/P1 requirement maps to at least one verifiable acceptance item.
4. Keep all unresolved P0/P1-affecting questions visible as blockers or open questions.
5. Do not begin implementation.
