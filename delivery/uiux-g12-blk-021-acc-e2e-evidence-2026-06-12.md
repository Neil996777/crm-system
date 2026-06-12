# BLK-UIUX-G12-021 ACC E2E Coverage Evidence

Date: 2026-06-12
Owner: Codex (QA Execution / Frontend execution)
Scope: test coverage only

## Coverage Added

| ACC | Test ID | File | Asserted sub-scenario |
|---|---|---|---|
| ACC-003 | `TEST-LEAD-TRANSFER-001/002` | `frontend/e2e/leads.spec.ts` | Sales Manager transfers a lead owner through the existing `/api/leads/{id}/owner-transfer` endpoint; the owner/version persist; record-local history contains `EVT-OWNER-CHANGED`; Sales receives `403 PERMISSION_DENIED`. |
| ACC-012 | `TEST-ACTIVITY-CONTEXT-001/002/003` | `frontend/e2e/work.spec.ts` | Existing task UI creates and persists tasks bound to real Lead, Contract, and Payment contexts; task detail and filtered task query show the related type/id; the related lead/contract/payment page can open the referenced record context. |
| ACC-019 | `TEST-DUPLICATE-CONTACT-001/002` | `frontend/e2e/duplicate.spec.ts` | Contact email and phone duplicate warnings appear from the existing customer-detail add-contact flow; proceeding creates new contact records; original and newly-created contacts keep distinct IDs, proving no silent merge/overwrite. |

## Constraint Notes

- Test-only code changes plus this evidence and blocker status update.
- No backend, `shared`, or root `api` implementation changes.
- No new color tokens, enum values, role values, `test.skip`, `test.only`, or `test.slow`.
- ACC-012 uses the existing task surface because current Lead/Contract/Payment detail pages do not expose the Opportunity-only `ActivityNoteTaskPanel`; no product UI panel was added.
- For Payment context, the existing frontend payment page is contract-scoped, so the task binds `relatedType=Payment` to the payment contract id that the current `回款` detail view opens; the setup still creates a real actual payment and records its `paymentId`.

## Verification

- `npm run build` in `frontend`: PASS.
- Targeted new coverage run:
  `npm run test:e2e -- e2e/leads.spec.ts e2e/work.spec.ts e2e/duplicate.spec.ts -g "TEST-LEAD-TRANSFER|TEST-ACTIVITY-CONTEXT|TEST-DUPLICATE-CONTACT"`: PASS, 3/3.
- Full e2e one-run smoke:
  `npm run test:e2e`: PASS, 61/61, workers:2, retries:1.
- `rg -n "test\.skip|test\.only|\.slow\(" frontend/e2e`: no matches.

Claude/release-owner still owns the requested repeated full-suite audit; Codex does not self-pass G12.
