# UI/UX G12 Functional Completeness Continuation: Focus Heading Race Evidence

Date: 2026-06-11

Status: Codex continuation return for Claude re-verification. Codex does not
self-resolve `BLK-UIUX-G12-017`.

## Scope

This continuation keeps the functional-completeness package intact:

- A row action menus remain as previously returned.
- B row click / record navigation remains as previously returned.
- C formerly-dead buttons remain wired or removed as previously returned.

Claude's functional click-testing passed A/B/C, but re-kicked this continuation
because `TEST-UIUX-A7-001` still found a product focus race:

- `WorkOverview.tsx` moved focus to `[data-focus-heading]` with one
  `requestAnimationFrame`.
- The new focus-stage clickable rows use `tabIndex=0`.
- Under occasional timing/load, the single rAF could run before the heading was
  queryable, so focus did not land on the heading.

## Product Fix

Changed file:

- `frontend/src/pages/WorkOverview.tsx`

Implementation:

- Added a cancellable focus-heading request:
  - `focusHeadingFrameRef`
  - `focusHeadingRequestRef`
  - `focusHeadingWaitMs = 1200`
- Replaced the single-rAF `focusStageHeading()` with a rAF retry loop that:
  - waits until `focusRootRef.current?.querySelector('[data-focus-heading]')`
    returns a heading;
  - focuses only that heading;
  - stops after the bounded wait if the focus view has disappeared;
  - ignores stale requests through the request id.
- Cancelled stale heading-focus attempts when:
  - entering a new card focus view;
  - switching to another focus card;
  - exiting focus view;
  - unmounting the dashboard.
- Kept the existing return-focus behavior: on exit, focus still returns to the
  original `[data-dashboard-card="..."]`.

Expected behavior:

- Entering focus view: focus reliably lands on the focus heading.
- Switching focus cards: focus reliably lands on the new focus heading after the
  switch transition.
- Returning from focus view: focus returns to the original dashboard card.
- Clickable rows in the focus stage do not receive initial focus.

## Constraint Check

- Backend diff: none.
- `shared` diff: none.
- root API diff: none.
- Endpoint usage: unchanged.
- CSS/color: no CSS changed by this continuation; no new color token or literal.
- zh-CN: no UI copy changed.
- Enum values: unchanged.
- Role values/gates: unchanged.
- Playwright workers: unchanged at `workers: 2`.
- No e2e skip/retry/slow/only changes.

## Verification

Commands run by Codex on 2026-06-11:

| Command | Result |
|---|---|
| `npx tsc --noEmit` | PASS |
| `npm run build` | PASS (`index-DNPx4_mx.css`, `index-LB-hONTJ.js`) |
| `npm run test:e2e` | PASS, 54/54, 0 failed, 0 skipped, workers: 2, 31.9s |
| `git diff --check` | PASS |
| `git diff --name-only -- services shared api packages/shared apps/api` | empty |
| `git status --short -- services shared api packages/shared apps/api` | empty |

Claude acceptance remains pending: run at least 8 consecutive full
`npm run test:e2e` passes at workers:2, without rerun/skip/retry. Codex does not
self-resolve `BLK-UIUX-G12-017`.
