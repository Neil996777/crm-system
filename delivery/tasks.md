# CRM System — G7 Task Plan (tasks.md)

Acceptance-driven, end-to-end, directly executable by the implementation agent
(Codex) at G9–G11. Every task has all 17 schema fields. Read
`delivery/delivery-plan.md` first for stack, repo layout, and global conventions
(common envelope, S2S token, outbox, optimistic concurrency, no-downgrade, TDD).
Conventions that apply to every task are stated once in the delivery plan and not
repeated per task; each task cites only the IDs specific to it.

Status legend: `Not Started` (Codex updates to `In Progress` / `Done` /
`Blocked` during G9–G11). No P0/P1 task may be marked Done via mock/stub/TODO/
static/non-persistent behavior (SVC-ACC-011, DEC-008).

Total tasks: 40 (TASK-001..TASK-040). Foundation TASK-001..006; capability slices
TASK-007..038; deployment/release evidence TASK-039..040.

---

## Phase 0 — Foundation / Platform

### TASK-001 — Monorepo, Docker Compose, PostgreSQL, migrate scaffold, shared contracts

1. **Task ID:** TASK-001
2. **Status:** Done
3. **Objective:** A running monorepo skeleton: Docker Compose brings up PostgreSQL
   and an empty-but-healthy container per service (SVC-001..010), each service a
   Go module with a `/health` endpoint and a `golang-migrate` scaffold against its
   own schema/DB-user; a `shared/contracts` package holds DTO/event/error/permission
   constants only.
4. **Business capability:** CAP-011 Persistence and production operation (primary ACC-016).
5. **Acceptance item:** ACC-016 (persistence platform foundation; full ACC-016 proven in TASK-037).
6. **Reference docs:** PSM-014, PSM Service Mapping (SVC-001..010), `data-design.md`
   Data Ownership Map (schemas + `crm_*_user` DB users), `module-boundaries.md`
   (Shared Package Boundary, Forbidden dependency types), `deployment-notes.md`
   (Container Deployment, paths `/opt/crm-system/...`), DEC-022, ADR-ARCH-001.
7. **File changes:**
   - `docker-compose.yml` (postgres + 10 service containers + internal network +
     per-service env + healthchecks; PostgreSQL NOT publicly exposed).
   - `services/<svc>/go.mod`, `services/<svc>/cmd/server/main.go` (health endpoint),
     `services/<svc>/migrations/0001_init_schema.up.sql` + `.down.sql` (create the
     service schema + role grants for `crm_<svc>_user` only) for SVC-002..010.
   - `services/gateway-bff/go.mod`, `services/gateway-bff/cmd/server/main.go` (no DB).
   - `shared/contracts/` (`errors.go` error codes, `envelope.go` request/response,
     `permission_actions.go`, `correlation.go`, event-schema structs) — NO domain
     methods/repos/business rules.
   - `Makefile` or `scripts/` for `migrate up/down` per service.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** None.
10. **Definition of Done:** `docker compose up` starts PostgreSQL + all 10 services;
    each service `/health` returns process-up + owned-DB connectivity; each schema
    exists and is reachable ONLY by its own `crm_<svc>_user`; `shared/contracts`
    compiles and contains no domain logic; gateway-bff has no DB credentials.
11. **Acceptance method:** Compose smoke (all health endpoints green); a static
    check that no service module imports another service's internal package and no
    service connects to a schema other than its own (ARCH-ACC-001, SVC-ACC-006/007).
12. **Automated tests:** Integration `TEST-PERSISTENCE-005` (scaffold-level: real
    PostgreSQL container per testcontainers reachable; no in-memory substitute);
    health-endpoint integration test per service. Type: Integration.
13. **Manual verification:** Run `docker compose up`; `curl` each `/health`; connect
    as `crm_lead_user` and confirm it cannot read `identity_authz` tables.
14. **Traceability:** CIM-045 → PIM-BEH-033 → PSM-014 / Data Ownership Map →
    CONTRACT-019 (S2S scaffolding) → ACC-016 → TEST-PERSISTENCE-005.
15. **TDD:** Write the health-endpoint + cross-schema-isolation integration tests
    first (fail with no services), then scaffold until green.
16. **No-downgrade items:** Real PostgreSQL container (no SQLite/in-memory); real
    per-service DB users with no cross-schema grants; no stubbed health endpoint.
17. **Blocker:** None.

### TASK-002 — identity-authz: authentication + session

1. **Task ID:** TASK-002
2. **Status:** Done
3. **Objective:** A user signs in with valid credentials and receives a persisted
   session bound to their single assigned role; invalid credentials, disabled
   accounts, and unauthenticated access are denied with one unified failure message.
4. **Business capability:** CAP-001 Identity and role access (primary ACC-001).
5. **Acceptance item:** ACC-001.
6. **Reference docs:** CIM-001/002, PIM-001/002, PIM-SM-011 (Active/Disabled),
   PIM-INV-045/047, PIM-BEH-001, PSM-001, CONTRACT-001/002, `api-spec.md` "Permission
   Check"/Service API Summary (identity-authz), `authz-architecture.md` (opaque
   `HttpOnly`/`Secure`/`SameSite=Lax` session cookie), PM-001/002, SEC-001,
   AUTH-001..006, ABUSE-001/006/022, DEC-005, TEST-AUTH-LOGIN.
7. **File changes:**
   - `services/identity-authz/internal/domain/user.go`, `role.go`, `session.go`.
   - `services/identity-authz/internal/handler/auth.go` (sign in, sign out,
     current user, session check).
   - `services/identity-authz/internal/repo/user_repo.go`, `session_repo.go`.
   - `services/identity-authz/migrations/0002_users_sessions.up.sql` (+ down):
     users, roles, sessions tables; seed first Administrator.
   - `services/identity-authz/internal/event/outbox.go` (UserSignedIn/Out,
     UserAccessDenied → SVC-008).
   - tests: `services/identity-authz/internal/handler/auth_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-001.
10. **Definition of Done:** Valid login persists an authenticated session (opaque
    cookie) carrying the role; session survives refresh and re-login; sign-out
    revokes server-side and clears cookie; invalid creds / disabled / unauthenticated
    all return one generic failure; login success/failure and access-denied emit
    operation-log events via outbox.
11. **Acceptance method:** ACC-001 scenario — sign in reaches role-scoped CRM;
    invalid/disabled/unauthenticated rejected without exposing account state.
12. **Automated tests:** `TEST-AUTH-LOGIN-001..006` (valid login binds role; invalid
    creds unified message; disabled denied; unauthenticated route/API denied; session
    persists across refresh/re-login; stale session re-evaluated after disable/role
    change), `TEST-ABUSE-UNAUTH-001`, `TEST-ABUSE-ENUM-001`. Type: Integration + E2E.
13. **Manual verification:** Sign in as seeded admin → reach workspace; sign in with
    wrong password → generic error; disable a user then attempt login → same generic
    error; refresh → still signed in.
14. **Traceability:** CIM-001/CIM-PROC-001 → PIM-001/PIM-SM-011/PIM-BEH-001 →
    PSM-001 → CONTRACT-001/002 → ACC-001 → TEST-AUTH-LOGIN-001..006.
15. **TDD:** Write TEST-AUTH-LOGIN + abuse tests first (fail), then implement.
16. **No-downgrade items:** Real session persistence in PostgreSQL; real disabled-user
    check server-side; unified failure message (no account-state disclosure).
17. **Blocker:** None.

### TASK-003 — identity-authz: permission decisions, three-role scope, S2S token, user-admin invariants

1. **Task ID:** TASK-003
2. **Status:** Done
3. **Objective:** Every protected action is authorized server-side by a central
   permission check enforcing the three-role scope (admin=all, manager=team,
   sales=owned/assigned); service-to-service calls require a signed token; the
   last-active-Administrator invariant is enforced.
4. **Business capability:** CAP-001 Identity and role access (primary ACC-002).
5. **Acceptance item:** ACC-002.
6. **Reference docs:** CIM-003/004/005/006/007, CIM-PROC-002/024, PIM-002/003,
   PIM-SM-011, PIM-INV-046, PIM-BEH-002/003, PSM-001, CONTRACT-001/019,
   `api-spec.md` "Permission Check", `authz-architecture.md` (Permission/Denial +
   S2S: Bearer service-token, `X-Service-Id`, `X-Intent`, audience, 5-min lifetime,
   `SERVICE_AUTH_FAILED`), PM-003..PM-009/PM-029, SEC-002..006/008, ABUSE-003/004/005,
   STB-003, SVC-ACC-006/007/008, ARCH-ACC-002/009, TEST-AUTHZ-SCOPE/TEST-USER-ADMIN.
7. **File changes:**
   - `services/identity-authz/internal/handler/permission.go`
     (`POST /internal/permissions/check`: actor/action/resource/scope → allow + scope).
   - `services/identity-authz/internal/handler/user_admin.go` (create user, change
     role/status; reject last-active-admin removal before save).
   - `services/identity-authz/internal/domain/permission_policy.go`,
     `last_admin_guard.go`.
   - `services/identity-authz/internal/authz/service_token.go` (sign + verify S2S token).
   - `shared/contracts/permission_actions.go` (action constants).
   - migrations `0003_permission_policy.up.sql`.
   - tests: `permission_test.go`, `user_admin_test.go`, `service_token_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-002.
10. **Definition of Done:** Permission check returns allow/deny + scope per
    PM-001..029; denied actions expose/mutate nothing; S2S calls without a valid
    signed token return `SERVICE_AUTH_FAILED`; disabling/downgrading the last active
    Administrator is blocked before save with an operation-log event; no hard delete
    path exists for any role.
11. **Acceptance method:** ACC-002 — admin governs all, manager team, sales owned-only;
    unauthorized actions denied without data exposure/mutation; Sales cannot view
    global logs; no hard delete.
12. **Automated tests:** `TEST-AUTHZ-SCOPE-001..006`, `TEST-USER-ADMIN-001..004`,
    `TEST-INV-LASTADMIN-001`, `TEST-PERM-USERADMIN-001..003`, `TEST-ABUSE-PRIVESC-001`,
    `TEST-ABUSE-MUTATE-001`, `TEST-ABUSE-S2S-001`, `TEST-INV-NODELETE-001` (route
    unavailable). Note: `TEST-AUTHZ-SCOPE-001/002/003` subsume the positive three-role
    CRUD-scope cases (admin=all / manager=team / sales=owned-assigned), i.e. they cover the
    `TEST-PERM-CRUD-ADMIN-001 / -MGR-001 / -SALES-001` positive mappings.
    Type: Integration (backend/API negative), E2E where user-visible.
13. **Manual verification:** As Sales call another user's lead by id → denied, no
    leakage; attempt to disable the only admin → blocked state; call an internal
    endpoint without S2S token → `SERVICE_AUTH_FAILED`.
14. **Traceability:** CIM-PROC-002/024 → PIM-002/003/PIM-INV-046/PIM-BEH-002/003 →
    PSM-001 → CONTRACT-001/019 → ACC-002 → TEST-AUTHZ-SCOPE-001..006 / TEST-INV-LASTADMIN-001.
15. **TDD:** Write scope/last-admin/S2S negative tests first (fail), then implement.
16. **No-downgrade items:** Real server-side permission check (not UI-only,
    ARCH-ACC-002); real S2S signed-token verification (not network trust); real
    last-admin guard; no hard-delete endpoint.
17. **Blocker:** None.

### TASK-004 — audit-history: append-only record history + admin operation log + outbox sink

1. **Task ID:** TASK-004
2. **Status:** Done
3. **Objective:** A trusted internal contract appends record-local history events
   and admin operation-log events to an append-only, tamper-evident store; history
   is queryable by record permission and operation logs by Administrator only.
4. **Business capability:** CAP-008 Collaboration history and operation audit
   (primary ACC-014; also ACC-022).
5. **Acceptance item:** ACC-014 (record-local history; ACC-022 admin oplog query proven in TASK-028).
6. **Reference docs:** CIM-035/036, CIM-PROC-017/018, PIM-018/019/027, PIM-BEH-028/029,
   PSM-009/013, CONTRACT-013/014, `api-spec.md` Service API Summary (audit-history)
   + Event Contract Requirements, `audit-log-spec.md` (Common Event Schema, EVT-*
   catalog, Query Requirements, AUD-IMM-001..005), `data-design.md` History And
   Audit Data (`eventId`,`prevHash`,`eventHash`,…), PM-024/025/040..042, PRIV-010/011/016,
   SEC-009/010, ABUSE-013, ABUSE-019, AUTHZ-009, ARCH-ACC-003, FLOW-006,
   TEST-HISTORY/TEST-OPLOG/TEST-ABUSE-ACTAS.
7. **File changes:**
   - `services/audit-history/internal/handler/append.go` (trusted internal append),
     `history_query.go`, `oplog_query.go`.
   - `services/audit-history/internal/domain/event.go` (hash chain `prevHash`/`eventHash`).
   - `services/audit-history/internal/repo/event_repo.go` (append-only; no update/delete).
   - migrations `0002_history_oplog.up.sql` (append-only tables, no UPDATE/DELETE grant).
   - `shared/contracts/events.go` (Common Event Schema struct).
   - tests: `append_test.go`, `history_query_test.go`, `oplog_query_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-001, TASK-003.
10. **Definition of Done:** Append accepts only valid S2S calls; events are
    append-only (no update/delete path, AUD-IMM-002) with a verified hash chain;
    record-history query enforces related-record permission (PM-024/025); oplog query
    is Administrator-only (PM-040..042); events carry required schema fields and safe
    summaries.
11. **Acceptance method:** ACC-014 — record-local history shows actor/event/resource/
    timestamp/before-after by record permission; non-owned history denied; history not
    editable via normal CRM actions.
12. **Automated tests:** `TEST-HISTORY-001..004`, `TEST-OPLOG-005` (not editable),
    `TEST-RETENTION-001` (append-only/no shorten), `TEST-ABUSE-ACTAS-001` (act-on-behalf-of:
    recorded actor is the authenticated principal; payload-claimed actor ignored; owner/assignee
    are separate fields from actor — ABUSE-019, AUTHZ-009), G12 rework
    `TEST-HISTORY-DISPATCH-RETRY-001`, `TEST-HISTORY-IDEMPOTENT-001`,
    `TEST-HISTORY-TX-001`; G12 second rework per-service dispatcher tests in
    account, commercial, and work covering successful audit-history delivery,
    failed-delivery retry retention, and duplicate event UID idempotency; G12 third
    rework lead transactional outbox rollback coverage (`TEST-HISTORY-TX-001`) proves
    lead create rolls back when the local outbox append fails; G12 fourth micro-rework
    lead dispatcher tests prove lead audit-history S2S delivery, retry retention, and
    distinct `EVT-LEAD-QUALIFIED` / `EVT-LEAD-DISQUALIFIED` event IDs. Type:
    Integration.
13. **Manual verification:** Append an event via S2S; query as record owner → visible;
    query as non-owner → denied; attempt to edit an event → unavailable.
14. **Traceability:** CIM-035/036 → PIM-018/019/PIM-BEH-028/029 → PSM-009 →
    CONTRACT-013/014 → ACC-014 → TEST-HISTORY-001..004.
15. **TDD:** Write history-permission + non-editable + hash-chain tests first (fail).
    G12 rework fail-first evidence: opportunity dispatcher test initially failed on missing
    `DispatchOnce`/`DispatchConfig`; audit-history duplicate UID test initially failed before
    idempotent append support; outbox transaction rollback test initially returned 201 before
    transaction coupling. G12 second rework fail-first evidence: work dispatcher test
    initially failed because `work.outbox_events` lacked `published_at`; after the
    migration fix, account/commercial/work dispatcher tests and full service test suites
    passed. G12 third rework fail-first evidence: lead create rollback test initially
    returned 201 because `_ = h.outbox.Append` discarded a forced outbox failure; after
    transactional coupling it returned an error and left no lead row, and `go test ./...`
    passed in services/lead. G12 fourth micro-rework fail-first evidence: lead dispatcher
    test did not compile before `AuditHistoryServiceURL`/audit mapping existed; after routing
    lead audit delivery through the transactional outbox and removing the post-commit
    `audit.AppendRecordHistory` call, `go test ./internal/event -run TestLeadOutboxDispatcher`
    and `go test ./... -count=1` passed in services/lead.
16. **No-downgrade items:** Real append-only persistence (no in-memory log); real
    hash chain; real record-permission gate on history query (ABUSE-013); actor identity
    is the authenticated principal, never a client-supplied field (ABUSE-019, AUTHZ-009).
17. **Blocker:** None.

### TASK-005 — gateway-bff: routing, correlation propagation, safe error normalization

1. **Task ID:** TASK-005
2. **Status:** Done
3. **Objective:** A single external API edge that authenticates request context,
   routes commands/queries to owning services, propagates correlation IDs, and
   normalizes safe error responses — without owning data or deciding business state.
4. **Business capability:** CAP-007 Core CRM navigation and record retrieval
   (primary ACC-015 spine).
5. **Acceptance item:** ACC-015 (retrieval aggregation spine; full list/detail/
   search/filter proven in TASK-030).
6. **Reference docs:** PIM-026, PSM-011, CONTRACT-001 (authz), `frontend-backend-contract.md`
   (BFF aggregation, invalid-filter error), `module-boundaries.md` Gateway Boundary,
   `api-spec.md` Common Request/Response + Error Codes, ARCH-ACC-002, TEST-NAV-RETRIEVE.
7. **File changes:**
   - `services/gateway-bff/internal/handler/router.go`, `proxy.go`, `aggregate.go`.
   - `services/gateway-bff/internal/authz/context.go` (calls SVC-002 session/permission).
   - `services/gateway-bff/internal/middleware/correlation.go`, `safe_error.go`.
   - tests: `router_test.go`, `safe_error_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-002, TASK-003.
10. **Definition of Done:** Gateway authenticates via SVC-002, routes to owning
    services with a propagated `correlationId`, returns only fields the target
    authorized, and normalizes errors to the safe envelope; gateway has no DB and
    decides no business state.
11. **Acceptance method:** ACC-015 spine — authorized list/detail requests routed and
    aggregated; invalid filter returns validation feedback; unauthorized records hidden.
12. **Automated tests:** `TEST-NAV-RETRIEVE-001` (list+detail routing), `-004`
    (invalid filter feedback), `-005` (unauthorized hidden). Type: Integration + E2E.
13. **Manual verification:** Hit a list endpoint through gateway → correlationId in
    response and logs; send an invalid filter → safe validation error.
14. **Traceability:** PIM-026 → PSM-011 → CONTRACT-001 → ACC-015 → TEST-NAV-RETRIEVE-001/004/005.
15. **TDD:** Write routing + safe-error + correlation tests first (fail).
16. **No-downgrade items:** Gateway never writes a DB or decides business state
    (module-boundaries.md); no leakage of unauthorized fields.
17. **Blocker:** None.

### TASK-006 — frontend: app shell, auth flow, Sign In + Work Overview screens

1. **Task ID:** TASK-006
2. **Status:** Done
3. **Objective:** The React+TS app shell with role-aware navigation, a working sign-in
   screen, and a Work Overview landing that talks only to gateway-bff.
4. **Business capability:** CAP-001 Identity and role access (primary ACC-001).
5. **Acceptance item:** ACC-001 (frontend half; backend in TASK-002).
6. **Reference docs:** `ui-spec.md` UI-001 (Sign In), UI-002 (Work Overview),
   Navigation + Visibility rules, Data Display Safety; `frontend-backend-contract.md`;
   CONTRACT-001, ACC-001/002/015/021, TEST-AUTH-LOGIN.
7. **File changes:**
   - `frontend/src/app/Shell.tsx`, `frontend/src/app/Nav.tsx` (role-aware sections).
   - `frontend/src/pages/SignIn.tsx`, `frontend/src/pages/WorkOverview.tsx`.
   - `frontend/src/api/client.ts` (gateway-bff only), `frontend/src/api/auth.ts`.
   - `frontend/src/auth/SessionProvider.tsx`.
   - tests: `e2e/auth.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-002, TASK-005.
10. **Definition of Done:** Sign-in routes to role-scoped workspace; one generic
    error on failure; nav shows only role-allowed sections (UI hiding does NOT
    replace backend authz); all calls go through gateway-bff.
11. **Acceptance method:** ACC-001 E2E — sign in as each role, reach Work Overview,
    see only allowed sections.
12. **Automated tests:** E2E `TEST-AUTH-LOGIN-001/005` (UI path), nav-visibility E2E.
    Type: E2E.
13. **Manual verification:** Sign in as Sales, Manager, Admin → confirm distinct nav;
    sign out → return to Sign In.
14. **Traceability:** CIM-001 → PIM-BEH-001 → PSM-011 → CONTRACT-001 → ACC-001 →
    TEST-AUTH-LOGIN-001/005.
15. **TDD:** Write the auth E2E spec first (fails until UI+API wired).
16. **No-downgrade items:** Frontend never the authorization authority; no hardcoded
    user list; session from real backend.
17. **Blocker:** None.

---

## Phase 1 — Capability Vertical Slices

### TASK-007 — lead-service: create/edit/assign/transfer with required fields and owner rules

1. **Task ID:** TASK-007
2. **Status:** Done
3. **Objective:** Create, view, edit, search/filter, and assign/transfer leads with
   required fields persisted; Unassigned leads supported; owner required before
   Pending Qualification; owner-change history preserved.
4. **Business capability:** CAP-002 Lead intake and qualification (primary ACC-003).
5. **Acceptance item:** ACC-003.
6. **Reference docs:** CIM-008/009/010, CIM-PROC-003, PIM-004/003, PIM-SM-001 (create/
   assign rows), PIM-SM-008, PIM-INV-001/004/005, PIM-BEH-004/005, PSM-002,
   CONTRACT-003/004/020, `api-spec.md` Service API Summary (lead) + "Owner Transfer"
   + "Editable Record Concurrency", PM-010..015, PIM-SM-008, EDGE-003, EDGE-024,
   ABUSE-002/003, TEST-LEAD-CREATE/TEST-LEAD-ASSIGN/TEST-OWNER-TRANSFER.
7. **File changes:**
   - `services/lead/internal/domain/lead.go` (aggregate, required-field guards),
     `ownership.go`.
   - `services/lead/internal/handler/lead_command.go` (create/update/assign/transfer),
     `lead_query.go` (list/detail/search).
   - `services/lead/internal/repo/lead_repo.go`.
   - `services/lead/internal/event/outbox.go` (LeadCreated, LeadOwnerChanged).
   - migrations `0002_leads.up.sql` (leads, owner, status, source; `version`,`updated_at`).
   - tests: `lead_command_test.go`, `lead_query_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-003, TASK-004.
10. **Definition of Done:** Lead persists with lead/company name + source + status;
    owner required before Pending Qualification; missing required fields blocked;
    only manager/admin assign/transfer (Sales denied); owner change emits history;
    `expectedVersion` enforced.
11. **Acceptance method:** ACC-003 — create/view/edit/search/filter/assign for allowed
    roles, denied for unauthorized, Unassigned before assignment, owner-change history.
12. **Automated tests:** `TEST-LEAD-CREATE-001/002/003`, `TEST-LEAD-ASSIGN-001/002`,
    `TEST-OWNER-TRANSFER-001` (manager/admin assign allowed — PIM-SM-008, PM-014),
    `TEST-OWNER-TRANSFER-002` (manager/admin transfer allowed — PM-014),
    `TEST-OWNER-TRANSFER-003` (Sales assign/transfer denied; out-of-scope target rejected — PM-015, PIM-SM-008),
    `TEST-AUTHZ-SCOPE-004` (Sales non-owned denied), `TEST-ABUSE-IDOR-001`;
    G12 third rework `TEST-HISTORY-TX-001` proves lead create mutation rolls back if
    its history/outbox insert fails. Type:
    Integration + E2E.
13. **Manual verification:** Create a lead missing source → blocked; create Unassigned
    → allowed; Manager assigns owner → history shows owner change; Sales tries assign
    → denied.
14. **Traceability:** CIM-008/CIM-PROC-003 → PIM-004/PIM-SM-001/PIM-SM-008/PIM-BEH-004/005 →
    PSM-002 → CONTRACT-003/004 → ACC-003 → TEST-LEAD-CREATE-001..003 / TEST-LEAD-ASSIGN-001/002 /
    TEST-OWNER-TRANSFER-001..003; EDGE-024 open-work cascade = TEST-OWNER-TRANSFER-004 (owned by TASK-024).
15. **TDD:** Write create/assign + IDOR negative tests first (fail).
16. **No-downgrade items:** Real DB persistence; real owner-required guard; real
    record-local owner-change history event; real scope check.
17. **Blocker:** None.

### TASK-008 — lead-service: qualification (Valid/Invalid/restore) + conversion with history

1. **Task ID:** TASK-008
2. **Status:** Done
3. **Objective:** Qualify a Pending lead Valid or Invalid (with reason), restore an
   Invalid lead (admin/manager only), and convert a Valid lead to an opportunity while
   preserving original lead history and preventing re-conversion.
4. **Business capability:** CAP-002 Lead intake and qualification (primary ACC-004).
5. **Acceptance item:** ACC-004.
6. **Reference docs:** CIM-PROC-004, PIM-004, PIM-SM-001 (qualify/restore/convert rows),
   PIM-INV-001/002/003/005, PIM-BEH-006, PSM-002, CONTRACT-003/004, `api-spec.md`
   "Lead Conversion" (idempotencyKey; failure cases), FLOW-002 (lead→opportunity),
   PM-013/014, EDGE-004/005/006, ABUSE-018, TEST-LEAD-QUALIFY.
7. **File changes:**
   - `services/lead/internal/domain/qualification.go`, `conversion.go` (conversion-once guard).
   - `services/lead/internal/handler/lead_qualify.go`, `lead_convert.go`.
   - `services/lead/internal/event/outbox.go` (LeadQualified, LeadConverted).
   - calls account-service + opportunity-service via S2S (create/link account+contact,
     create opportunity) per FLOW-002.
   - tests: `qualification_test.go`, `conversion_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-007, TASK-010 (account create/link), TASK-013 (opportunity create).
10. **Definition of Done:** Valid/Invalid/Converted persisted with reason/downstream
    link; Unassigned lead cannot qualify/convert; Invalid cannot convert until restored
    by admin/manager; converted lead cannot reconvert; conversion preserves lead history;
    convert is idempotent by key.
11. **Acceptance method:** ACC-004 — transition rules pass for allowed actors; forbidden
    transitions rejected without mutation; conversion preserves history.
12. **Automated tests:** `TEST-LEAD-QUALIFY-001..007`, `TEST-ABUSE-BRBYPASS-001` (subset);
    G12 third rework regression in services/lead covers the same transactional outbox
    helper used by qualification and conversion mutations. G12 fourth micro-rework adds
    distinct disqualify event coverage (`LeadDisqualified` → `EVT-LEAD-DISQUALIFIED`) and
    dispatcher audit retry coverage for qualification events.
    Type: Integration + E2E.
13. **Manual verification:** Mark Valid → convert → opportunity created, lead history
    intact; mark Invalid then convert → rejected; restore as Sales → denied.
14. **Traceability:** CIM-PROC-004 → PIM-SM-001/PIM-INV-002/003/PIM-BEH-006 → PSM-002
    → CONTRACT-003/004 + FLOW-002 → ACC-004 → TEST-LEAD-QUALIFY-001..007.
15. **TDD:** Write qualify/convert positive + all reject tests first (fail).
16. **No-downgrade items:** Real conversion-once guard in aggregate; real cross-service
    create via S2S (not a stub); preserved history is a real event.
17. **Blocker:** None.

### TASK-009 — frontend: Lead List + Lead Detail & Qualification screens

1. **Task ID:** TASK-009
2. **Status:** Done
3. **Objective:** UI to create/list/search leads and run qualification/conversion,
   reflecting Unassigned/denied/invalid-reason/converted states.
4. **Business capability:** CAP-002 Lead intake and qualification (primary ACC-003).
5. **Acceptance item:** ACC-003 (lead UI; ACC-004 qualification UI covered by same screen).
6. **Reference docs:** `ui-spec.md` UI-003 (List Pattern), UI-005 (Lead Detail And
   Qualification), `frontend-backend-contract.md`, CONTRACT-003, ACC-003/004,
   TEST-LEAD-CREATE/QUALIFY (E2E legs).
7. **File changes:**
   - `frontend/src/pages/leads/LeadList.tsx`, `LeadDetail.tsx`.
   - `frontend/src/api/leads.ts`.
   - `frontend/src/components/QualificationActions.tsx`, `ConvertLeadDialog.tsx`.
   - tests: `e2e/leads.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-007, TASK-008, and the frontend shell (TASK-006). The reusable
   list/detail pattern is built here first; TASK-030 generalizes it later across all entities
   (TASK-030 depends on the slices, not the reverse).
10. **Definition of Done:** Lead list with search/filter; create form with required
    fields + validation; qualification actions disabled on Unassigned for Sales;
    invalid-reason required; converted shows read-only conversion; stale-edit conflict
    surfaced.
11. **Acceptance method:** ACC-003/004 E2E through gateway.
12. **Automated tests:** E2E `TEST-LEAD-CREATE-002` (validation), `TEST-LEAD-QUALIFY-003`
    (convert), `TEST-LEAD-QUALIFY-004` (Unassigned denied UI). Type: E2E.
13. **Manual verification:** Create lead, qualify Valid, convert; observe history.
14. **Traceability:** CIM-008 → PIM-BEH-004/006 → PSM-002/PSM-011 → CONTRACT-003 →
    ACC-003/004 → TEST-LEAD-CREATE-002 / TEST-LEAD-QUALIFY-003/004.
15. **TDD:** Write lead E2E spec first (fail).
16. **No-downgrade items:** No static lead list; data from real API; UI disablement
    backed by real backend denial.
17. **Blocker:** None.

### TASK-010 — account-service: company/customer CRUD with required fields, no hard delete

1. **Task ID:** TASK-010
2. **Status:** Done
3. **Objective:** Create/view/edit/search/filter ToB companies/customers with required
   fields persisted; no hard delete; unauthorized access denied.
4. **Business capability:** CAP-003 Account and contact management (primary ACC-005).
5. **Acceptance item:** ACC-005.
6. **Reference docs:** CIM-011/012, CIM-PROC-006, PIM-005, PIM-SM-010 (no hard delete),
   PIM-BEH-007, PSM-003, CONTRACT-005/006/020, `api-spec.md` Service API Summary
   (account/contact), PM-008/009/016/017/029, EDGE (required fields), ABUSE-002/017,
   TEST-CUSTOMER-CRUD/TEST-INV-NODELETE.
7. **File changes:**
   - `services/account/internal/domain/account.go`.
   - `services/account/internal/handler/account_command.go`, `account_query.go`.
   - `services/account/internal/repo/account_repo.go`.
   - `services/account/internal/event/outbox.go` (AccountCreated, OwnerChanged).
   - migrations `0002_accounts.up.sql` (accounts; `version`,`updated_at`; no DELETE grant).
   - tests: `account_command_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-003, TASK-004.
10. **Definition of Done:** Account persists with company name + customer status +
    owner; missing required fields blocked; no hard-delete route; unauthorized denied;
    `expectedVersion` enforced; history emitted on mutation.
11. **Acceptance method:** ACC-005 scenarios for authorized/unauthorized.
12. **Automated tests:** `TEST-CUSTOMER-CRUD-001..004`, `TEST-INV-NODELETE-001`. Type:
    Integration + E2E.
13. **Manual verification:** Create without status → blocked; edit persists across
    restart; attempt delete → unavailable.
14. **Traceability:** CIM-011 → PIM-005/PIM-BEH-007 → PSM-003 → CONTRACT-005/006 →
    ACC-005 → TEST-CUSTOMER-CRUD-001..004.
15. **TDD:** Write CRUD + no-delete tests first (fail).
16. **No-downgrade items:** Real persistence; real no-hard-delete; real scope check.
17. **Blocker:** None.

### TASK-011 — account-service: multiple contacts under a company with required link

1. **Task ID:** TASK-011
2. **Status:** Done
3. **Objective:** Create/link multiple contacts under a company/customer, each
   requiring a related company and at least one contact method or role note.
4. **Business capability:** CAP-003 Account and contact management (primary ACC-006).
5. **Acceptance item:** ACC-006.
6. **Reference docs:** CIM-013, PIM-006, PIM-BEH-008, PSM-003, CONTRACT-005/006,
   `api-spec.md` Service API Summary (account/contact), PM-016/017, EDGE-007,
   TEST-CONTACT-LINK.
7. **File changes:**
   - `services/account/internal/domain/contact.go`.
   - `services/account/internal/handler/contact_command.go`, `contact_query.go`.
   - `services/account/internal/repo/contact_repo.go`.
   - migrations `0003_contacts.up.sql` (contacts FK to account; contact method/role note).
   - `services/account/internal/event/outbox.go` (ContactCreated).
   - tests: `contact_command_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-010.
10. **Definition of Done:** Contact persists with name + related company + ≥1 contact
    method or role note; save without company blocked; multiple contacts visible in the
    company context; unauthorized denied.
11. **Acceptance method:** ACC-006 — multiple contacts per company, visible in context.
12. **Automated tests:** `TEST-CONTACT-LINK-001..004`. Type: Integration + E2E.
13. **Manual verification:** Add 2 contacts to a company → both listed; add contact with
    no method/note → blocked.
14. **Traceability:** CIM-013 → PIM-006/PIM-BEH-008 → PSM-003 → CONTRACT-005 → ACC-006
    → TEST-CONTACT-LINK-001..004.
15. **TDD:** Write contact-link tests first (fail).
16. **No-downgrade items:** Real FK link to company; real persistence.
17. **Blocker:** None.

### TASK-012 — frontend: Customer/Contact List + Detail screens

1. **Task ID:** TASK-012
2. **Status:** Done
3. **Objective:** UI to manage companies and their contacts, with related opportunities/
   contracts/payments/history sections.
4. **Business capability:** CAP-003 Account and contact management (primary ACC-005).
5. **Acceptance item:** ACC-005 (and ACC-006 contacts via the same Customer/Contact Detail).
6. **Reference docs:** `ui-spec.md` UI-006 (Customer/Contact Detail), UI-003 (List),
   CONTRACT-005, ACC-005/006, TEST-CUSTOMER-CRUD/TEST-CONTACT-LINK (E2E).
7. **File changes:**
   - `frontend/src/pages/accounts/AccountList.tsx`, `AccountDetail.tsx`.
   - `frontend/src/components/ContactTable.tsx`, `AddContactDialog.tsx`.
   - `frontend/src/api/accounts.ts`.
   - tests: `e2e/accounts.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-010, TASK-011.
10. **Definition of Done:** Account list/detail; contacts table with add; required-field
    validation; unauthorized states; related sections render from real APIs.
11. **Acceptance method:** ACC-005/006 E2E.
12. **Automated tests:** E2E `TEST-CUSTOMER-CRUD-002`, `TEST-CONTACT-LINK-003`. Type: E2E.
13. **Manual verification:** Create company + 2 contacts; view in context.
14. **Traceability:** CIM-011/013 → PIM-BEH-007/008 → PSM-003/PSM-011 → CONTRACT-005 →
    ACC-005/006 → TEST-CUSTOMER-CRUD-002 / TEST-CONTACT-LINK-003.
15. **TDD:** Write accounts E2E spec first (fail).
16. **No-downgrade items:** No static data; real API; UI matches backend rules.
17. **Blocker:** None.

### TASK-013 — opportunity-service: create/edit with Pipeline Stage (no Status field)

1. **Task ID:** TASK-013
2. **Status:** Done
3. **Objective:** Create/view/edit/search/filter opportunities with required links/
   fields and a Pipeline Stage (the sole lifecycle dimension; no Status field); no hard
   delete; owner/role visibility enforced.
4. **Business capability:** CAP-004 Opportunity pipeline (primary ACC-007).
5. **Acceptance item:** ACC-007.
6. **Reference docs:** CIM-014/015 (DEC-020 Status removed), CIM-PROC-007, PIM-007/012,
   PIM-SM-002 (create row), PIM-BEH-009, PSM-004, CONTRACT-007/008/020, `api-spec.md`
   Service API Summary (opportunity) + "Editable Record Concurrency", PM-018/019,
   DEC-020, TEST-OPP-CREATE.
7. **File changes:**
   - `services/opportunity/internal/domain/opportunity.go` (Pipeline Stage enum: New
     Opportunity, Needs Confirmed, Quote, Contract Negotiation, Won, Lost — NO Status).
   - `services/opportunity/internal/handler/opportunity_command.go`, `opportunity_query.go`.
   - `services/opportunity/internal/repo/opportunity_repo.go`.
   - migrations `0002_opportunities.up.sql` (stage column only — no status column;
     `version`,`updated_at`).
   - `services/opportunity/internal/event/outbox.go` (OpportunityCreated).
   - tests: `opportunity_command_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-003, TASK-004, TASK-010 (customer link).
10. **Definition of Done:** Opportunity persists with related customer + owner + stage +
    expected amount + expected close date; NO separate status field exists; missing
    links/fields blocked; unauthorized edits denied; no hard delete; `expectedVersion`.
11. **Acceptance method:** ACC-007 — create/view/edit/search/filter with persisted data
    and role/owner visibility; confirm no Status field.
12. **Automated tests:** `TEST-OPP-CREATE-001..004` (incl. -003 plain create, no Status;
    retired TEST-OPP-STATUS-ENUM-001 NOT implemented). Type: Integration + E2E.
13. **Manual verification:** Create opportunity; confirm form/detail has Stage but no
    Status; missing amount → blocked.
14. **Traceability:** CIM-014/CIM-PROC-007 → PIM-007/PIM-SM-002/PIM-BEH-009 → PSM-004 →
    CONTRACT-007 → ACC-007 → TEST-OPP-CREATE-001..004.
15. **TDD:** Write create + missing-field + no-Status tests first (fail).
16. **No-downgrade items:** Real persistence; schema has no Status field (DEC-020 honored,
    not a hidden/unused column); real scope check.
17. **Blocker:** None.

### TASK-014 — opportunity-service: pipeline stage transitions with history

1. **Task ID:** TASK-014
2. **Status:** Done
3. **Objective:** Move an opportunity through allowed forward stage transitions,
   rejecting forbidden transitions and rollbacks without mutation, emitting a history
   event per change.
4. **Business capability:** CAP-004 Opportunity pipeline (primary ACC-008).
5. **Acceptance item:** ACC-008 (stage transitions; closure in TASK-015).
6. **Reference docs:** CIM-015, PIM-007, PIM-SM-002 (transition matrix), PIM-INV-006/010,
   PIM-BEH-010, PSM-004, CONTRACT-007/008, `api-spec.md` Service API Summary (change
   stage), PM-018, EDGE-008, ABUSE-018, TEST-OPP-STAGE.
7. **File changes:**
   - `services/opportunity/internal/domain/stage_machine.go` (allowed transitions).
   - `services/opportunity/internal/handler/opportunity_stage.go`.
   - `services/opportunity/internal/event/outbox.go` (OpportunityStageChanged).
   - tests: `stage_machine_test.go`, `opportunity_stage_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-013.
10. **Definition of Done:** Allowed forward transitions persist + emit history; forbidden
    transitions and arbitrary rollback rejected (`INVALID_TRANSITION`) with no mutation.
11. **Acceptance method:** ACC-008 transition table enforced for allowed/forbidden +
    history events.
12. **Automated tests:** `TEST-OPP-STAGE-001/002/003`, `TEST-HISTORY-002` (stage history),
    `TEST-ABUSE-BRBYPASS-001` (subset). Type: Integration + E2E (Unit for stage_machine).
13. **Manual verification:** Advance New→Needs Confirmed→Quote; attempt skip/rollback →
    rejected; history shows each change.
14. **Traceability:** CIM-015 → PIM-SM-002/PIM-INV-006/PIM-BEH-010 → PSM-004 →
    CONTRACT-007/008 → ACC-008 → TEST-OPP-STAGE-001/002/003.
15. **TDD:** Write allowed + forbidden + rollback tests first (fail).
16. **No-downgrade items:** Real transition guard; real history event; rejected
    transitions cause no DB change.
17. **Blocker:** None.

### TASK-015 — opportunity-service: close Won (contract Signed) / Lost (reason), terminal lock

1. **Task ID:** TASK-015
2. **Status:** Done
3. **Objective:** Close an opportunity Won only when its related contract is Signed
   (verified with commercial-service), or Lost with a recorded reason; Won/Lost are
   terminal and non-reopenable; closure preserves related history; post-close notes/
   tasks still allowed.
4. **Business capability:** CAP-004 Opportunity pipeline (primary ACC-013; also ACC-008 closure).
5. **Acceptance item:** ACC-013.
6. **Reference docs:** CIM-017/018, CIM-PROC-011, PIM-007/011, PIM-SM-009 (close),
   PIM-INV-007/009/035/036/037/038/039, PIM-BEH-011, PSM-004/007, CONTRACT-007/008/009/010,
   `api-spec.md` "Close Opportunity Won" (verify Signed contract via commercial-service),
   "Close Opportunity Lost", FLOW-004 (contract signing → Won), DEC-017/019, EDGE-009/010/011/036,
   ABUSE-018, TEST-OPP-CLOSE / TEST-INV-WONAFTERPAY / TEST-INV-TERMINAL.
7. **File changes:**
   - `services/opportunity/internal/domain/closure.go` (terminal lock; Won requires
     Signed contract; Lost requires reason).
   - `services/opportunity/internal/handler/close_won.go`, `close_lost.go`.
   - `services/opportunity/internal/authz/commercial_client.go` (S2S: query contract
     Signed status — NO cross-service DB).
   - `services/opportunity/internal/event/outbox.go` (OpportunityClosedWon/Lost).
   - tests: `closure_test.go`, `close_won_test.go`, `close_lost_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-014, TASK-019 (contract status: Signed).
10. **Definition of Done:** Won persists only after commercial-service confirms a Signed
    contract (else `EARLY_WON_BLOCKED`); full payment NOT required; Lost requires reason
    (`LOST_REASON_REQUIRED`); Won/Lost terminal — reopen/rollback/re-close rejected
    (`TERMINAL_RECORD_READ_ONLY`); related history preserved; post-close notes/tasks via
    work-service still allowed.
11. **Acceptance method:** ACC-013 — Signed-contract Won, lost-reason Lost, forbidden
    reopen/early-won rejected.
12. **Automated tests:** `TEST-OPP-CLOSE-001..006`, `TEST-INV-WONAFTERPAY-001`,
    `TEST-INV-TERMINAL-001`. Type: Integration + E2E.
13. **Manual verification:** Try Won with no Signed contract → blocked; sign contract →
    Won succeeds; reopen → rejected; add a note after close → allowed.
14. **Traceability:** CIM-017/CIM-PROC-011 → PIM-SM-009/PIM-INV-035/037/PIM-BEH-011 →
    PSM-004 → CONTRACT-007/008 + FLOW-004 → ACC-013 → TEST-OPP-CLOSE-001/002/005.
15. **TDD:** Write Won-without-Signed reject, Lost-without-reason reject, reopen reject
    first (fail).
16. **No-downgrade items:** Real S2S contract-Signed verification (not assumed/stubbed);
    real terminal lock; payment is NOT a Won gate (DEC-019) but overpayment still blocked
    elsewhere; real history preservation.
17. **Blocker:** None.

### TASK-016 — frontend: Opportunity Detail (stage stepper, close Won/Lost)

1. **Task ID:** TASK-016
2. **Status:** Done
3. **Objective:** Opportunity UI with a stage stepper, blocked-transition alerts, and
   Won/Lost closure confirmations reflecting the Signed-contract and lost-reason rules.
4. **Business capability:** CAP-004 Opportunity pipeline (primary ACC-008).
5. **Acceptance item:** ACC-008 (and ACC-007/ACC-013 via same screen).
6. **Reference docs:** `ui-spec.md` UI-007 (Opportunity Detail; no Status field, Won
   blocked until contract Signed), CONTRACT-007, ACC-007/008/013/014, TEST-OPP-STAGE/CLOSE (E2E).
7. **File changes:**
   - `frontend/src/pages/opportunities/OpportunityList.tsx`, `OpportunityDetail.tsx`.
   - `frontend/src/components/StageStepper.tsx`, `CloseOpportunityDialog.tsx`.
   - `frontend/src/api/opportunities.ts`.
   - tests: `e2e/opportunities.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-013, TASK-014, TASK-015.
10. **Definition of Done:** Stage stepper shows only valid next steps; forbidden
    transition alert; Won disabled/blocked until contract Signed; Lost requires reason;
    terminal Won/Lost read-only except notes/tasks; no Status control.
11. **Acceptance method:** ACC-008 E2E.
12. **Automated tests:** E2E `TEST-OPP-STAGE-002` (blocked), `TEST-OPP-CLOSE-002`
    (Won blocked), `TEST-OPP-CLOSE-003` (Lost). Type: E2E.
13. **Manual verification:** Advance stages; attempt Won early → blocked message; close Lost.
14. **Traceability:** CIM-015 → PIM-BEH-010/011 → PSM-004/PSM-011 → CONTRACT-007 →
    ACC-008 → TEST-OPP-STAGE-002 / TEST-OPP-CLOSE-002/003.
15. **TDD:** Write opportunities E2E spec first (fail).
16. **No-downgrade items:** No static stage list; UI blocking backed by real backend
    guards; no Status field shown.
17. **Blocker:** None.

### TASK-017 — commercial-service: quote lifecycle, exactly one quote per opportunity

1. **Task ID:** TASK-017
2. **Status:** Done
3. **Objective:** Create/send/accept/reject/expire a quote with required fields; enforce
   exactly one quote per opportunity; an expired quote cannot link to a new contract.
4. **Business capability:** CAP-005 Commercial execution (primary ACC-009).
5. **Acceptance item:** ACC-009.
6. **Reference docs:** CIM-019/020, CIM-PROC-008, PIM-008/012, PIM-SM-004, PIM-INV-012/013/014/015,
   PIM-BEH-012/013, PSM-005, CONTRACT-009/010/020, `api-spec.md` Service API Summary
   (commercial: create/update quote, change quote status), FLOW-003, DEC-018, PM-020,
   EDGE-012/013, TEST-QUOTE-LIFECYCLE / TEST-QUOTE-ACCEPT / TEST-INV-ONEACCEPT.
7. **File changes:**
   - `services/commercial/internal/domain/quote.go` (one-quote-per-opportunity guard;
     status machine Draft/Sent/Accepted/Rejected/Expired).
   - `services/commercial/internal/handler/quote_command.go`, `quote_query.go`.
   - `services/commercial/internal/repo/quote_repo.go` (unique constraint on opportunityId).
   - migrations `0002_quotes.up.sql` (quotes; UNIQUE(opportunity_id); `version`,`updated_at`).
   - `services/commercial/internal/event/outbox.go` (QuoteAccepted).
   - tests: `quote_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-003, TASK-004, TASK-013 (opportunity link).
10. **Definition of Done:** Quote persists with opportunity + customer + amount + status
    + validity end date + owner; a second quote on the same opportunity is rejected
    (DB UNIQUE + domain guard); missing amount/status/validity blocked; expired quote
    cannot link to a contract; quote acceptance emits history + oplog.
11. **Acceptance method:** ACC-009 — quote transitions pass; exactly one quote per
    opportunity; contract linkage requires the Accepted quote.
12. **Automated tests:** `TEST-QUOTE-LIFECYCLE-001..003`, `TEST-QUOTE-ACCEPT-001/002`,
    `TEST-INV-ONEACCEPT-001`; G12 systematic rework `TEST-AUTHZ-SCOPE-005` covers
    quote single-record read authorization and no data leak on denied by-id reads.
    Type: Integration + E2E (Unit for one-quote guard).
13. **Manual verification:** Create quote; try a second on same opportunity → rejected;
    accept; expire → cannot link contract.
14. **Traceability:** CIM-019/CIM-PROC-008 → PIM-008/PIM-SM-004/PIM-INV-012/PIM-BEH-013 →
    PSM-005 → CONTRACT-009/010 → ACC-009 → TEST-QUOTE-ACCEPT-001 / TEST-INV-ONEACCEPT-001.
15. **TDD:** Write one-quote-per-opportunity reject + lifecycle tests first (fail).
    G12 systematic rework fail-first evidence: commercial quote by-id read test first
    returned 200 with non-owned quote details to another Sales user; after owner-scope
    check it returns safe `NOT_FOUND` with no record data, and `go test ./... -count=1`
    passed in services/commercial.
16. **No-downgrade items:** Real DB UNIQUE + domain guard (no second quote); real
    persistence; real acceptance history event. Do NOT implement any multi-quote or
    second-accept path (DEC-018, retired TEST-QUOTE-ACCEPT-003).
17. **Blocker:** None.

### TASK-018 — commercial-service: contract create from Accepted quote (note, expected signed date, amount-diff reason)

1. **Task ID:** TASK-018
2. **Status:** Done
3. **Objective:** Create a Pending Signature contract from an Accepted, non-Expired quote
   with a required note and expected signed date; record an amount-difference reason when
   contract amount differs from the quote.
4. **Business capability:** CAP-005 Commercial execution (primary ACC-010).
5. **Acceptance item:** ACC-010 (create half; lifecycle in TASK-019).
6. **Reference docs:** CIM-021/023/024/025, CIM-PROC-009, PIM-009/012, PIM-SM-005 (create
   row), PIM-INV-016/018/019/020, PIM-BEH-014/016, PSM-006, CONTRACT-009/010/020,
   `api-spec.md` Service API Summary (create/update contract), DEC-006/007/014/016,
   EDGE-013/014/017, TEST-CONTRACT-CREATE / TEST-CONTRACT-AMOUNT-DIFF / TEST-INV-CONTRACTQUOTE.
7. **File changes:**
   - `services/commercial/internal/domain/contract.go` (create guards; quote-link rule;
     amount-diff reason).
   - `services/commercial/internal/handler/contract_command.go`, `contract_query.go`.
   - `services/commercial/internal/repo/contract_repo.go`.
   - migrations `0003_contracts.up.sql` (contracts FK quote/opportunity/customer;
     `version`,`updated_at`).
   - tests: `contract_create_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-017.
10. **Definition of Done:** Pending Signature contract persists with customer +
    opportunity + Accepted quote + amount + status + note + expected signed date
    (no signed date required yet); missing note/link/amount/expected-signed-date rejected;
    expired/non-Accepted quote link rejected; amount differing from quote requires reason.
11. **Acceptance method:** ACC-010 create — Pending Signature valid without signed date but
    requires expected signed date; notes P0-required; amount-diff reason persisted.
12. **Automated tests:** `TEST-CONTRACT-CREATE-001/002/003`, `TEST-CONTRACT-AMOUNT-DIFF-001`,
    `TEST-INV-CONTRACTQUOTE-001`, `TEST-INV-NOAPPROVAL-001`; G12 systematic rework
    `TEST-AUTHZ-SCOPE-005` covers contract single-record read authorization and no data
    leak on denied by-id reads. Type: Integration + Unit.
13. **Manual verification:** Create contract from Accepted quote → ok; from expired quote →
    rejected; without note → blocked; differing amount without reason → blocked.
14. **Traceability:** CIM-021/CIM-PROC-009 → PIM-009/PIM-SM-005/PIM-INV-016/018/PIM-BEH-014/016 →
    PSM-006 → CONTRACT-009/010 → ACC-010 → TEST-CONTRACT-CREATE-001..003.
15. **TDD:** Write create + reject (no note / expired quote) + amount-diff tests first (fail).
    G12 systematic rework fail-first evidence: commercial contract by-id read test first
    returned 200 with non-owned contract details to another Sales user; after owner-scope
    check it returns safe `NOT_FOUND` with no record data.
16. **No-downgrade items:** Real Accepted-quote link check; real required-note guard; no
    approval/e-sign/template (DEC-007); real persistence.
17. **Blocker:** None.

### TASK-019 — commercial-service: contract lifecycle (sign/activate/complete/terminate) with date guards

1. **Task ID:** TASK-019
2. **Status:** Done
3. **Objective:** Transition a contract Pending Signature → Signed → Active → Completed,
   or Terminate (pre/post signature), requiring a signed/effective date for signed states;
   contract Signed is the event that enables Opportunity Won.
4. **Business capability:** CAP-005 Commercial execution (primary ACC-010).
5. **Acceptance item:** ACC-010 (lifecycle).
6. **Reference docs:** CIM-022/024, PIM-009, PIM-SM-005 (sign/activate/complete/terminate),
   PIM-INV-017/021, PIM-BEH-015, PSM-006, CONTRACT-009/010, `api-spec.md` change contract
   status, FLOW-004 (ContractStatusChanged(Signed)), EDGE-015/016, PM-021, ABUSE-018,
   TEST-CONTRACT-LIFECYCLE / TEST-INV-CONTRACTDATE.
7. **File changes:**
   - `services/commercial/internal/domain/contract_lifecycle.go` (date guards).
   - `services/commercial/internal/handler/contract_status.go`.
   - `services/commercial/internal/event/outbox.go` (ContractStatusChanged incl. Signed →
     consumed by opportunity-service for Won verification).
   - tests: `contract_lifecycle_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-018.
10. **Definition of Done:** Sign/Activate/Complete and post-signature Terminate require a
    signed/effective date (rejected if missing); pre-signature Terminate allowed without
    it; each change emits history + oplog and a `ContractStatusChanged` event; Signed
    status is queryable by opportunity-service via S2S/event projection.
11. **Acceptance method:** ACC-010 lifecycle transition rules pass; signed states without
    signed date rejected.
12. **Automated tests:** `TEST-CONTRACT-LIFECYCLE-001/002/003`, `TEST-INV-CONTRACTDATE-001`;
    G12 systematic rework `TEST-AUTHZ-SCOPE-005` covers contract by-id read scope.
    Type: Integration.
13. **Manual verification:** Sign without effective date → rejected; sign with date → ok;
    confirm opportunity can now close Won.
14. **Traceability:** CIM-022 → PIM-SM-005/PIM-INV-017/PIM-BEH-015 → PSM-006 →
    CONTRACT-009/010 + FLOW-004 → ACC-010 → TEST-CONTRACT-LIFECYCLE-001/002.
15. **TDD:** Write date-guard reject + sign-success tests first (fail).
16. **No-downgrade items:** Real date guards; real status event for Won verification (not
    assumed); real persistence.
17. **Blocker:** None.

### TASK-020 — commercial-service: payment plans + actual payments, status, overpayment block (post-sale)

1. **Task ID:** TASK-020
2. **Status:** Done
3. **Objective:** Create payment plans and record actual payments updating Unpaid/
   Partially Paid/Paid; reject zero/negative and contract-level overpayment; single
   currency; payment tracking is post-sale and does NOT gate Won.
4. **Business capability:** CAP-005 Commercial execution (primary ACC-011).
5. **Acceptance item:** ACC-011.
6. **Reference docs:** CIM-026/027/028/029, CIM-PROC-010, PIM-010/011/012, PIM-SM-006,
   PIM-INV-022/023/024/025, PIM-BEH-017/018, PSM-007, CONTRACT-009/010, `api-spec.md`
   "Record Payment" (`paymentStatus`,`remainingAmount`), DEC-013/014/019, EDGE-018/019,
   PM-022, ABUSE-018, TEST-PAYMENT-RECORD / TEST-PAYMENT-GUARD / TEST-INV-OVERPAY /
   TEST-INV-PAYAMOUNT / TEST-INV-CURRENCY.
7. **File changes:**
   - `services/commercial/internal/domain/payment.go` (plan + actual payment; cumulative-
     vs-contract-remaining ceiling; the Contract aggregate owns plans+payments so the
     overpayment invariant is atomic).
   - `services/commercial/internal/handler/payment_command.go`, `payment_query.go`.
   - migrations `0004_payments.up.sql` (payment_plans, actual_payments; idempotency key).
   - `services/commercial/internal/event/outbox.go` (PaymentRecorded).
   - tests: `payment_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-018.
10. **Definition of Done:** Plan persists Unpaid; partial → Partially Paid, full → Paid;
    zero/negative rejected (`INVALID_AMOUNT`); cumulative beyond contract remaining rejected
    (`OVERPAYMENT_BLOCKED`) across plans; single currency enforced; record-payment idempotent;
    payment does NOT affect Opportunity Won (decoupled, DEC-019).
11. **Acceptance method:** ACC-011 — payment transition rules; single-currency; overpayment
    blocked.
12. **Automated tests:** `TEST-PAYMENT-RECORD-001..003`, `TEST-PAYMENT-GUARD-001..004`,
    `TEST-INV-OVERPAY-001`, `TEST-INV-PAYAMOUNT-001`, `TEST-INV-CURRENCY-001`. Type:
    Integration + Unit. (Do NOT implement retired TEST-PAYMENT-FULLPAID-AGG-001.)
13. **Manual verification:** Record partial then full → Paid; record beyond remaining →
    blocked; record negative → blocked.
14. **Traceability:** CIM-026/CIM-PROC-010 → PIM-010/011/PIM-SM-006/PIM-INV-023/PIM-BEH-018 →
    PSM-007 → CONTRACT-009/010 → ACC-011 → TEST-PAYMENT-GUARD-003 / TEST-INV-OVERPAY-001.
15. **TDD:** Write zero/negative/overpayment reject tests first (fail).
16. **No-downgrade items:** Real overpayment ceiling at contract level; real persistence;
    payment is decoupled from Won (no full-payment-gates-Won logic).
17. **Blocker:** None.

### TASK-021 — frontend: Quote Detail screen

1. **Task ID:** TASK-021
2. **Status:** Done
3. **Objective:** Quote UI with status/amount/validity, accept/reject actions, expired
   warning, and contract-link indicator; reflects exactly-one-quote rule.
4. **Business capability:** CAP-005 Commercial execution (primary ACC-009).
5. **Acceptance item:** ACC-009 (UI).
6. **Reference docs:** `ui-spec.md` UI-008 (Quote Detail), CONTRACT-009, ACC-009,
   TEST-QUOTE-LIFECYCLE/ACCEPT (E2E).
7. **File changes:**
   - `frontend/src/pages/quotes/QuoteList.tsx`, `QuoteDetail.tsx`.
   - `frontend/src/api/quotes.ts`.
   - tests: `e2e/quotes.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-017.
10. **Definition of Done:** Quote create/detail; accept/reject; expired warning; expired
    quote blocks contract link in UI; no UI path to add a second quote.
11. **Acceptance method:** ACC-009 E2E.
12. **Automated tests:** E2E `TEST-QUOTE-ACCEPT-001`, `TEST-QUOTE-LIFECYCLE-002`. Type: E2E.
13. **Manual verification:** Create + accept a quote; observe single-quote constraint.
14. **Traceability:** CIM-019 → PIM-BEH-012/013 → PSM-005/PSM-011 → CONTRACT-009 →
    ACC-009 → TEST-QUOTE-ACCEPT-001.
15. **TDD:** Write quotes E2E spec first (fail).
16. **No-downgrade items:** No static quote data; real API; no second-quote UI affordance.
17. **Blocker:** None.

### TASK-022 — frontend: Contract Detail screen

1. **Task ID:** TASK-022
2. **Status:** Done
3. **Objective:** Contract UI with status/dates/note/amount-diff reason, status actions,
   and pending-signature reminder warning.
4. **Business capability:** CAP-005 Commercial execution (primary ACC-010).
5. **Acceptance item:** ACC-010 (UI).
6. **Reference docs:** `ui-spec.md` UI-009 (Contract Detail), CONTRACT-009, ACC-010/021,
   TEST-CONTRACT-CREATE/LIFECYCLE (E2E).
7. **File changes:**
   - `frontend/src/pages/contracts/ContractList.tsx`, `ContractDetail.tsx`.
   - `frontend/src/api/contracts.ts`.
   - tests: `e2e/contracts.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-018, TASK-019.
10. **Definition of Done:** Create requires note + expected signed date; signed states
    require signed/effective date; amount mismatch requires reason; pending-signature
    past expected signed date shows reminder warning.
11. **Acceptance method:** ACC-010 E2E.
12. **Automated tests:** E2E `TEST-CONTRACT-CREATE-002`, `TEST-CONTRACT-LIFECYCLE-002`. Type: E2E.
13. **Manual verification:** Create + sign a contract; observe date guards.
14. **Traceability:** CIM-021 → PIM-BEH-014/015 → PSM-006/PSM-011 → CONTRACT-009 →
    ACC-010 → TEST-CONTRACT-CREATE-002 / TEST-CONTRACT-LIFECYCLE-002.
15. **TDD:** Write contracts E2E spec first (fail).
16. **No-downgrade items:** No static contract data; UI guards backed by backend.
17. **Blocker:** None.

### TASK-023 — frontend: Payment Detail screen

1. **Task ID:** TASK-023
2. **Status:** Done
3. **Objective:** Payment UI for plan + actual payment with remaining amount, overdue
   signal, and zero/negative/overpayment validation; payment shown as post-sale tracking.
4. **Business capability:** CAP-005 Commercial execution (primary ACC-011).
5. **Acceptance item:** ACC-011 (UI).
6. **Reference docs:** `ui-spec.md` UI-010 (Payment Detail; post-sale, does not gate Won),
   CONTRACT-009, ACC-011/013/021, TEST-PAYMENT-RECORD/GUARD (E2E).
7. **File changes:**
   - `frontend/src/pages/payments/PaymentList.tsx`, `PaymentDetail.tsx`.
   - `frontend/src/api/payments.ts`.
   - tests: `e2e/payments.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-020.
10. **Definition of Done:** Record payment updates status + remaining; zero/negative/
    overpayment blocked with safe messages; overdue signal shown; no UI implying payment
    gates Won.
11. **Acceptance method:** ACC-011 E2E.
12. **Automated tests:** E2E `TEST-PAYMENT-RECORD-002`, `TEST-PAYMENT-GUARD-003`. Type: E2E.
13. **Manual verification:** Record partial/full; attempt overpayment → blocked.
14. **Traceability:** CIM-026/027 → PIM-BEH-017/018 → PSM-007/PSM-011 → CONTRACT-009 →
    ACC-011 → TEST-PAYMENT-RECORD-002 / TEST-PAYMENT-GUARD-003.
15. **TDD:** Write payments E2E spec first (fail).
16. **No-downgrade items:** No static payment data; UI validation backed by backend.
17. **Blocker:** None.

### TASK-024 — work-service: activities, notes, tasks against related records

1. **Task ID:** TASK-024
2. **Status:** Done
3. **Objective:** Create activities, notes, and tasks linked to any CRM record with
   required fields; manage task status (Open/Completed/Cancelled/Overdue); completed/
   cancelled tasks are not active reminders.
4. **Business capability:** CAP-006 Work activity and reminders (primary ACC-012).
5. **Acceptance item:** ACC-012.
6. **Reference docs:** CIM-030/031/032, CIM-PROC-012, PIM-013/014/015, PIM-SM-007,
   PIM-INV-028/029/030/031/033, PIM-BEH-020/021, PSM-008, CONTRACT-011/012, `api-spec.md`
   Service API Summary (work) + "Owner Transfer", PM-023, EDGE-023/024, TEST-ACTIVITY-NOTE /
   TEST-TASK-LIFECYCLE / TEST-OWNER-TRANSFER-004.
7. **File changes:**
   - `services/work/internal/domain/activity.go`, `note.go`, `task.go` (status machine).
   - `services/work/internal/handler/work_command.go`, `work_query.go`.
   - `services/work/internal/repo/work_repo.go`.
   - migrations `0002_work.up.sql` (activities, notes, tasks; related-record ref).
   - `services/work/internal/event/outbox.go` (WorkItemCreated, TaskStatusChanged).
   - tests: `work_command_test.go`, `task_lifecycle_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-003, TASK-004 (related-record services exist for refs:
   TASK-007/010/013/018/020).
10. **Definition of Done:** Activity/note/task persists with related entity + actor/owner +
    timestamp + content/title + due date (where applicable) + status; missing related
    record or required fields blocked; task complete/cancel; overdue evaluated on-read by
    business date; completed/cancelled not active reminders; on a parent record owner change,
    open tasks/follow-ups transfer to the new owner unless explicitly reassigned (EDGE-024,
    PIM-INV-030/033); unauthorized create/view denied.
11. **Acceptance method:** ACC-012 — activity/note/task scenarios against lead/customer/
    opportunity/contract/payment, incl. invalid due date and permission-denied.
12. **Automated tests:** `TEST-ACTIVITY-NOTE-001..003`, `TEST-TASK-LIFECYCLE-001..004`,
    `TEST-INV-TASKREMINDER-001`, `TEST-OWNER-TRANSFER-004` (EDGE-024 open-work
    cascade: on parent record owner change, open tasks/follow-ups transfer to the new owner
    unless explicitly reassigned — PIM-INV-030/033; this is the test home for the cascade).
    Type: Integration + E2E.
13. **Manual verification:** Add note to a lead; create task with due date; complete it →
    no longer a reminder.
14. **Traceability:** CIM-030/032/CIM-PROC-012 → PIM-013/015/PIM-SM-007/PIM-BEH-020/021 →
    PSM-008 → CONTRACT-011/012 → ACC-012 → TEST-ACTIVITY-NOTE-001 / TEST-TASK-LIFECYCLE-001..004.
15. **TDD:** Write activity/note/task create + missing-link reject + completed-not-reminder
    tests first (fail).
16. **No-downgrade items:** Real related-record link; real persistence; real history events.
17. **Blocker:** None.

### TASK-025 — frontend: Activities/Notes/Tasks UI (embedded in record detail + standalone list)

1. **Task ID:** TASK-025
2. **Status:** Done
3. **Objective:** UI to add/view activities, notes, and tasks within a record's detail and
   as a standalone task list with status actions.
4. **Business capability:** CAP-006 Work activity and reminders (primary ACC-012).
5. **Acceptance item:** ACC-012 (UI).
6. **Reference docs:** `ui-spec.md` UI-004 (Entity Detail related sections), UI-002 (active
   work), CONTRACT-011, ACC-012, TEST-ACTIVITY-NOTE / TEST-TASK-LIFECYCLE (E2E).
7. **File changes:**
   - `frontend/src/components/ActivityNoteTaskPanel.tsx`, `TaskList.tsx`.
   - `frontend/src/api/work.ts`.
   - tests: `e2e/work.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-024.
10. **Definition of Done:** Add activity/note/task in record detail; task status actions;
    validation on missing fields; permission-denied state; data from real API.
11. **Acceptance method:** ACC-012 E2E.
12. **Automated tests:** E2E `TEST-ACTIVITY-NOTE-002`, `TEST-TASK-LIFECYCLE-002`. Type: E2E.
13. **Manual verification:** Add task in opportunity detail; complete it.
14. **Traceability:** CIM-030/032 → PIM-BEH-020/021 → PSM-008/PSM-011 → CONTRACT-011 →
    ACC-012 → TEST-ACTIVITY-NOTE-002 / TEST-TASK-LIFECYCLE-002.
15. **TDD:** Write work E2E spec first (fail).
16. **No-downgrade items:** No static work data; real API.
17. **Blocker:** None.

### TASK-026 — Reminders: on-read due/overdue evaluation + Reminder Center

1. **Task ID:** TASK-026
2. **Status:** Done
3. **Objective:** An in-app reminder query computes due/overdue tasks, pending-signature
   contracts past expected signed date, and due/overdue payments on-read against a supplied
   business date, hiding inactive and unauthorized items; the Reminder Center renders them.
4. **Business capability:** CAP-006 Work activity and reminders (primary ACC-021; also CAP-012).
5. **Acceptance item:** ACC-021.
6. **Reference docs:** CIM-033/034, CIM-PROC-013, PIM-015/016/017, PIM-SM-005/006/007,
   PIM-INV-026/028/031, PIM-BEH-019/022, PSM-008, CONTRACT-011/010, `api-spec.md` "Reminder
   Query" (`businessDate`, `Asia/Shanghai`), FLOW-005, PSM Resolved Mechanisms (on-read),
   DEC-015, PM-046/047, EDGE-020/021/022/023/037, ABUSE-021, TEST-REMINDER / TEST-PAYMENT-OVERDUE.
7. **File changes:**
   - `services/work/internal/handler/reminder_query.go` (aggregates own tasks + commercial-
     service eligibility via S2S; on-read overdue).
   - `services/work/internal/domain/reminder.go` (derivation rules; business-date guard).
   - `services/commercial/internal/handler/reminder_eligibility.go` (contract pending-
     signature past expected signed date; payment due/overdue — derives a transient
     PaymentOverdue reminder signal on read against the supplied business date; this is NOT
     an outbox/audit event and writes nothing).
   - `frontend/src/pages/reminders/ReminderCenter.tsx`, `frontend/src/api/reminders.ts`.
   - tests: `reminder_query_test.go`, `reminder_eligibility_test.go`, `e2e/reminders.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-019, TASK-020, TASK-024.
10. **Definition of Done:** Reminder query returns authorized due/overdue tasks, pending-
    signature contracts past expected signed date, due/overdue payments for the supplied
    business date; completed/cancelled tasks, signed/terminated/fully-paid contracts,
    archived and unauthorized items excluded; UI groups by type.
11. **Acceptance method:** ACC-021 — reminder scenarios for due/overdue tasks, contract
    pending-signature, payments, inactive suppression, permission filtering.
12. **Automated tests:** `TEST-REMINDER-001..005`, `TEST-REMINDER-BOUNDARY-001`,
    `TEST-PAYMENT-OVERDUE-001`, `TEST-TASK-LIFECYCLE-003`. Type: Integration + E2E.
13. **Manual verification:** Set a task due yesterday → appears overdue for that business
    date; complete it → disappears; non-owned reminder hidden.
14. **Traceability:** CIM-033/CIM-PROC-013 → PIM-016/PIM-BEH-022 → PSM-008 →
    CONTRACT-011/010 + FLOW-005 → ACC-021 → TEST-REMINDER-001..005.
15. **TDD:** Write reminder positive + inactive-suppression + permission-hidden + boundary
    tests first (fail).
16. **No-downgrade items:** Real on-read evaluation against persisted data; real permission
    filter (ABUSE-021); no static reminder list.
17. **Blocker:** None.

### TASK-027 — frontend: Record-local history timeline (in record detail)

1. **Task ID:** TASK-027
2. **Status:** Done
3. **Objective:** Show the permitted record-local history timeline inside each record's
   detail, with safe before/after values per role and classification.
4. **Business capability:** CAP-008 Collaboration history and operation audit (primary ACC-014).
5. **Acceptance item:** ACC-014 (UI/query path; append covered in TASK-004 and per-slice events).
6. **Reference docs:** CIM-035, CIM-PROC-017, PIM-018, PSM-009, CONTRACT-013, `api-spec.md`
   record history query, `ui-spec.md` UI-014 (History) + Security display rules, PM-024/025,
   PRIV-010, SEC-009, ABUSE-013, EDGE-029, TEST-HISTORY.
7. **File changes:**
   - `services/gateway-bff/internal/handler/history_aggregate.go` (route to SVC-008 history query).
   - `frontend/src/components/HistoryTimeline.tsx`, `frontend/src/api/history.ts`.
   - tests: `e2e/history.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-004, TASK-005, and slices that emit history (TASK-007/014/017/019/020/024).
10. **Definition of Done:** History timeline shows actor/event/resource/timestamp/safe
    before-after by record permission; non-owned history denied; not editable in UI.
11. **Acceptance method:** ACC-014 — record-local history visible per record permission.
12. **Automated tests:** E2E `TEST-HISTORY-001`, `TEST-HISTORY-003` (non-owned denied),
    `TEST-HISTORY-004` (not editable). Type: Integration + E2E.
13. **Manual verification:** Change a lead owner → history shows the change; view as
    non-owner → denied.
14. **Traceability:** CIM-035/CIM-PROC-017 → PIM-018/PIM-BEH-028 → PSM-009 → CONTRACT-013 →
    ACC-014 → TEST-HISTORY-001/003/004.
15. **TDD:** Write history E2E + non-owned-denied tests first (fail).
16. **No-downgrade items:** Real permission gate on history query; safe summaries; no editing.
17. **Blocker:** None.

### TASK-028 — audit-history + frontend: Admin global operation-log query

1. **Task ID:** TASK-028
2. **Status:** Done
3. **Objective:** Administrator-only global operation-log query covering required event
   classes, with a UI table; non-admins denied; logs not editable.
4. **Business capability:** CAP-008 Collaboration history and operation audit (primary ACC-022).
5. **Acceptance item:** ACC-022.
6. **Reference docs:** CIM-036, CIM-PROC-018, PIM-019/001, PIM-BEH-029, PSM-009, CONTRACT-013/014/002,
   `api-spec.md` admin operation log query, `audit-log-spec.md` Event Catalog + Query Requirements,
   `ui-spec.md` UI-014 (Admin Logs), PM-040/041/042, PRIV-011/016, SEC-010, ABUSE-012, EDGE-029,
   TEST-OPLOG.
7. **File changes:**
   - `services/audit-history/internal/handler/oplog_query.go` (admin-only filters).
   - `services/gateway-bff/internal/handler/oplog_aggregate.go`.
   - `frontend/src/pages/admin/OperationLogs.tsx`, `frontend/src/api/oplog.ts`.
   - tests: `oplog_query_test.go`, `e2e/oplog.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-004, TASK-003 (admin authz), and event-producing slices.
10. **Definition of Done:** Admin sees global events (login/access failures, role/status
    changes, last-admin-blocked, owner changes, stage/status changes, quote acceptance,
    contract changes, payments, archive, import, export) with id/actor/action/resource/
    timestamp/result/before-after; Manager and Sales denied; not editable.
11. **Acceptance method:** ACC-022 — admin global operation-log query content + denial for
    non-admins.
12. **Automated tests:** `TEST-OPLOG-001..005`. Type: Integration + Manual.
13. **Manual verification:** As admin, query logs after performing actions; as manager/sales
    → denied.
14. **Traceability:** CIM-036/CIM-PROC-018 → PIM-019/PIM-BEH-029 → PSM-009 → CONTRACT-013/014 →
    ACC-022 → TEST-OPLOG-001..005.
15. **TDD:** Write admin-allow + manager/sales-deny + not-editable tests first (fail).
16. **No-downgrade items:** Real admin-only gate; real append-only logs; no editing path.
17. **Blocker:** None.

### TASK-029 — frontend: Admin User/Role Management screen

1. **Task ID:** TASK-029
2. **Status:** Done
3. **Objective:** Admin UI to create users, change roles/status with confirmation, and show
   the last-Administrator-blocked state — backed by backend governance.
4. **Business capability:** CAP-001 Identity and role access (primary ACC-002; also ACC-001/022).
5. **Acceptance item:** ACC-002 (user/role admin surface; backend in TASK-003).
6. **Reference docs:** `ui-spec.md` UI-017 (Admin User/Role Management, confirmation + blocked
   state), CIM-PROC-024, PIM-SM-011, PIM-INV-046, CONTRACT-001/002, PM-003/006/007, ABUSE-004/005,
   TEST-USER-ADMIN / TEST-INV-LASTADMIN.
7. **File changes:**
   - `frontend/src/pages/admin/UserManagement.tsx`, `frontend/src/api/users.ts`.
   - `frontend/src/components/RoleStatusChangeDialog.tsx`.
   - tests: `e2e/user-admin.spec.ts`.
8. **Owner agent:** frontend-engineer
9. **Prerequisites:** TASK-003, TASK-006.
10. **Definition of Done:** Create user; change role/status with confirmation showing old/new
    + access impact + audit notice; last-admin disable/downgrade shown as blocked with no
    confirm; non-admin denied; UI never replaces backend authz.
11. **Acceptance method:** ACC-002 — three-role governance surface; last-admin protection.
12. **Automated tests:** E2E `TEST-USER-ADMIN-001`, `TEST-INV-LASTADMIN-001` (blocked state),
    `TEST-PERM-USERADMIN-002/003` (non-admin denied). Type: E2E.
13. **Manual verification:** Create a user; try to disable the only admin → blocked state.
14. **Traceability:** CIM-PROC-024 → PIM-SM-011/PIM-INV-046 → PSM-001 → CONTRACT-001/002 →
    ACC-002 → TEST-USER-ADMIN-001 / TEST-INV-LASTADMIN-001.
15. **TDD:** Write user-admin E2E + last-admin-blocked spec first (fail).
16. **No-downgrade items:** UI confirmation does not replace backend enforcement; real
    last-admin block.
17. **Blocker:** None.

### TASK-030 — gateway-bff + frontend: core list/detail/search/basic filter across all P0 entities

1. **Task ID:** TASK-030
2. **Status:** Done
3. **Objective:** Role-scoped list, detail, search, and basic filter across all P0 entities
   with empty-state, invalid-filter feedback, and permission hiding.
4. **Business capability:** CAP-007 Core CRM navigation and record retrieval (primary ACC-015).
5. **Acceptance item:** ACC-015.
6. **Reference docs:** CIM-047, CIM-PROC-023, PIM-026/020, PIM-BEH-030, PSM-011, CONTRACT-001
   + each owning Query API, `frontend-backend-contract.md`, `ui-spec.md` UI-003 (List) + UI-004
   (Detail), PM-008..033, EDGE-002/031, ABUSE-016, TEST-NAV-RETRIEVE.
7. **File changes:**
   - `services/gateway-bff/internal/handler/list_aggregate.go`, `search_filter.go`,
     `invalid_filter.go`.
   - `frontend/src/components/EntityList.tsx` (reusable list pattern), `EntityDetail.tsx`.
   - `frontend/src/api/retrieval.ts`.
   - tests: `search_filter_test.go`, `e2e/retrieval.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-005, and the owning Query APIs (TASK-007/010/011/013/017/018/020/024).
10. **Definition of Done:** List/detail/search/filter for leads, companies/customers,
    contacts, opportunities, quotes, contracts, payments, activities, tasks; empty state;
    invalid filter validation feedback; unauthorized records hidden; archived excluded from
    active views by default.
11. **Acceptance method:** ACC-015 — happy/empty/invalid-filter/permission-hidden/permission-
    denied across all P0 entities.
12. **Automated tests:** `TEST-NAV-RETRIEVE-001..006`, `TEST-ABUSE-ARCHIVED-001`. Type:
    E2E + Integration.
13. **Manual verification:** List each entity; apply a valid then invalid filter; confirm
    unauthorized records hidden.
14. **Traceability:** CIM-047/CIM-PROC-023 → PIM-026/PIM-BEH-030 → PSM-011 → CONTRACT-001 →
    ACC-015 → TEST-NAV-RETRIEVE-001..006.
15. **TDD:** Write list/detail/empty/invalid-filter/permission-hidden tests first (fail).
16. **No-downgrade items:** Real permission filtering at backend (not UI-only); real data;
    archived excluded from active views.
17. **Blocker:** None.

### TASK-031 — Duplicate warning on lead/company/contact create/edit

1. **Task ID:** TASK-031
2. **Status:** Done
3. **Objective:** Raise a non-blocking duplicate warning on normalized company name /
   contact phone-email / lead company-contact match; allow proceed-after-warning to create a
   new record only; no merge/overwrite; reveal no unauthorized matched-record detail.
4. **Business capability:** CAP-002 + CAP-003 (primary ACC-019).
5. **Acceptance item:** ACC-019.
6. **Reference docs:** CIM-040, CIM-PROC-005, PIM-021, PIM-BEH-025, PSM-002/003, CONTRACT-003/005,
   `api-spec.md` "Duplicate Warning" (normalize fields; `warningToken`; proceed token),
   FLOW-012, DEC (BR-011/019), PM-048, EDGE-025/033/034/035, ABUSE-020, TEST-DUPLICATE-WARN.
7. **File changes:**
   - `services/lead/internal/domain/duplicate.go`, `services/account/internal/domain/duplicate.go`
     (normalization rules).
   - `services/lead/internal/handler/duplicate_check.go`, `services/account/internal/handler/duplicate_check.go`.
   - `services/lead/internal/handler/lead_command.go` / account command: accept `proceedWarningToken`.
   - `services/lead/internal/event/outbox.go` / account: DuplicateWarningRaised.
   - `frontend/src/components/DuplicateWarning.tsx`.
   - tests: `duplicate_test.go` (both services), `e2e/duplicate.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-007, TASK-010, TASK-011.
10. **Definition of Done:** Warning raised on exact company name (case/space normalized),
    contact phone/email (normalized), lead company/contact match; proceed-after-warning with
    a single-use token creates a new record only (no merge/overwrite); unique data → no
    warning; matched-record details safe.
11. **Acceptance method:** ACC-019 — warning scenarios + proceed-after-warning + no-warning
    unique data.
12. **Automated tests:** `TEST-DUPLICATE-WARN-001..006`, `TEST-ABUSE-DUPENUM-001`;
    G12 third rework regression in services/lead covers the same transactional outbox
    helper used by lead duplicate-warning token + event persistence. Type:
    Integration + E2E.
13. **Manual verification:** Create a company; create another with same name differing case →
    warning; proceed → second record created, first untouched.
14. **Traceability:** CIM-040/CIM-PROC-005 → PIM-021/PIM-BEH-025 → PSM-002/003 →
    CONTRACT-003/005 + FLOW-012 → ACC-019 → TEST-DUPLICATE-WARN-001..006.
15. **TDD:** Write match + proceed-no-merge + unique-no-warning + enum-safety tests first (fail).
16. **No-downgrade items:** Real normalized matching against persisted records; no silent
    merge/overwrite; no unauthorized matched-detail leakage.
17. **Blocker:** None.

### TASK-032 — Archive lifecycle with active-obligation blocking (record-owning services + UI)

1. **Task ID:** TASK-032
2. **Status:** Done
3. **Objective:** Archive eligible records (admin/manager only) when no unresolved active
   downstream obligations exist; archived records leave active views/reminders/default
   reports but remain retrievable via explicit archived filter; no hard delete.
4. **Business capability:** CAP-012 Archive and lifecycle governance (primary ACC-002; also
   ACC-014/015/023).
5. **Acceptance item:** ACC-002 (archive governance; archive history feeds ACC-014, filtering ACC-015/023).
6. **Reference docs:** CIM-037/038/039, CIM-PROC-020, PIM-020, PIM-SM-010, PIM-INV-040/041/042/043/044,
   PIM-BEH-024, PSM-002/003/004/006 (eligible records), CONTRACT-005/007/009/011, `api-spec.md`
   "Archive Eligibility" (`canArchive`, obligations), FLOW-010, `ui-spec.md` UI-016 (Archive
   Confirmation), PM-026/027/028/029/030..033, EDGE-031/032, ABUSE-016/017, ARCH-ACC-011,
   TEST-ARCHIVE / TEST-INV-ARCHIVEBLOCK / TEST-INV-NODELETE.
7. **File changes:**
   - `services/<record-owning>/internal/handler/archive.go` (account, opportunity,
     commercial-contract, lead) + `archive_eligibility.go` (queries work-service/commercial-
     service obligations via S2S).
   - `services/<record-owning>/internal/domain/archive.go`.
   - `frontend/src/components/ArchiveConfirmation.tsx`, archived-filter in `EntityList.tsx`.
   - tests: `archive_test.go` per service, `e2e/archive.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-019, TASK-020, TASK-024, TASK-030.
10. **Definition of Done:** Archive blocked (`ARCHIVE_BLOCKED_ACTIVE_OBLIGATION`) when open
    tasks / pending-signature contracts / unpaid payments exist; allowed once resolved;
    Sales denied; archive emits history + oplog; archived excluded from active views/
    reminders/default reports, retrievable via explicit archived filter; no hard delete.
11. **Acceptance method:** ACC-002 archive governance — admin/manager allow, Sales deny,
    obligation-blocked, no hard delete.
12. **Automated tests:** `TEST-ARCHIVE-001..004`, `TEST-INV-ARCHIVEBLOCK-001`,
    `TEST-INV-NODELETE-001`, `TEST-ABUSE-ARCHIVED-001`. Type: Integration + E2E.
13. **Manual verification:** Try to archive a customer with an open task → blocked; resolve
    task → archive succeeds; confirm hidden from active list, visible under archived filter.
14. **Traceability:** CIM-037/CIM-PROC-020 → PIM-020/PIM-SM-010/PIM-INV-041/PIM-BEH-024 →
    PSM-002/003/004/006 → CONTRACT-005/007/009/011 + FLOW-010 → ACC-002 → TEST-ARCHIVE-003 /
    TEST-INV-ARCHIVEBLOCK-001.
15. **TDD:** Write obligation-blocked reject + Sales-denied + archived-filter tests first (fail).
16. **No-downgrade items:** Real obligation check across services via S2S (not assumed); real
    archive state; no hard delete.
17. **Blocker:** None.

### TASK-033 — reporting-service + frontend: Manager Team Overview

1. **Task ID:** TASK-033
2. **Status:** Done
3. **Objective:** A Sales Manager team overview aggregating team leads/opportunities/quotes/
   contracts/payments/tasks and pipeline status from read models; Sales denied; empty state.
4. **Business capability:** CAP-009 Team overview and reports (primary ACC-018).
5. **Acceptance item:** ACC-018.
6. **Reference docs:** CIM-043, CIM-PROC-015, PIM-024/020, PIM-BEH-031, PSM-010, CONTRACT-015/016,
   `api-spec.md` "Report Metrics" / overview, `data-design.md` Reporting Data (read model,
   authz-before-aggregate), `ui-spec.md` UI-012 (Manager Overview), PM-044/045, ARCH-ACC-006,
   EDGE-028, ABUSE-015, TEST-TEAM-OVERVIEW.
7. **File changes:**
   - `services/reporting/internal/handler/overview_query.go`.
   - `services/reporting/internal/repo/projection_repo.go`, `internal/projection/consumer.go`
     (consumes source domain events; rebuild — NO source DB access).
   - migrations `0002_reporting_projections.up.sql`.
   - `frontend/src/pages/reports/ManagerOverview.tsx`, `frontend/src/api/reports.ts`.
   - tests: `overview_query_test.go`, `e2e/overview.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-007, TASK-013, TASK-017, TASK-018, TASK-020, TASK-024 (event sources).
10. **Definition of Done:** Manager sees team records + pipeline status from read models;
    Sales denied (`PERMISSION_DENIED`); empty data → empty state; non-team records excluded;
    authorization applied before aggregation; archived excluded by default.
11. **Acceptance method:** ACC-018 — team overview from persisted authorized team records;
    Sales denied.
12. **Automated tests:** `TEST-TEAM-OVERVIEW-001..004`; G12 rework
    `TEST-REPORTING-S2S-001..005` for internal projection ingest authentication
    and no-mutation denial; `TEST-REPORTING-PROJECTION-INGEST-001/006` for
    producer dispatcher projection delivery and manager aggregate query; G12 second
    rework `TEST-REPORTING-CORRELATION-001` for projection ingest correlation-id
    propagation; G12 fourth micro-rework lead dispatcher test proves lead outbox→reporting
    S2S headers, `X-Correlation-Id`, and failed-delivery retry retention. Type: E2E + Integration.
13. **Manual verification:** As manager, open overview; as Sales → denied.
14. **Traceability:** CIM-043/CIM-PROC-015 → PIM-024/PIM-BEH-031 → PSM-010 → CONTRACT-015/016 →
    ACC-018 → TEST-TEAM-OVERVIEW-001..004.
15. **TDD:** Write manager-allow + Sales-deny + empty-state tests first (fail).
    G12 rework fail-first evidence: `TestProjectionIngestRequiresS2SToken`
    initially failed because `Config.ServiceID` / `ServiceTokenSecret` and S2S
    verification were absent, then passed after signed-token middleware. Projection
    dispatcher test initially failed on missing `ReportingServiceURL`, then passed
    after producer outbox dispatchers called reporting ingest. G12 second rework
    fail-first evidence: opportunity dispatcher test failed because reporting
    projection delivery omitted `X-Correlation-Id`, then passed after all producer
    reporting dispatchers set it. G12 fourth micro-rework fail-first evidence: lead
    dispatcher test initially failed before lead had the same audit/reporting dispatcher
    contract, then passed with real PostgreSQL testcontainers.
16. **No-downgrade items:** Read model from events/contracts, NOT source DB (ARCH-ACC-006);
    real authz-before-aggregate; no static numbers.
17. **Blocker:** None.

### TASK-034 — reporting-service + frontend: Basic sales reports

1. **Task ID:** TASK-034
2. **Status:** Done
3. **Objective:** Basic reports with counts/sums for leads by status, opportunities by stage,
   quotes/contracts/payments by status/amount over persisted authorized records; empty/zero
   state; Sales denied; archived excluded by default.
4. **Business capability:** CAP-009 Team overview and reports (primary ACC-023).
5. **Acceptance item:** ACC-023.
6. **Reference docs:** CIM-044, CIM-PROC-019, PIM-025/020, PIM-BEH-032, PSM-010, CONTRACT-015/016,
   `api-spec.md` "Report Metrics" (`metrics`, `breakdowns`, `GroupRow`), `ui-spec.md` UI-015
   (Reports), PM-043/044/045, ARCH-ACC-006, BR-014/017, EDGE-028/031, ABUSE-015, TEST-BASIC-REPORT.
7. **File changes:**
   - `services/reporting/internal/handler/report_query.go` (metrics + breakdowns).
   - `frontend/src/pages/reports/BasicReports.tsx`.
   - tests: `report_query_test.go`, `e2e/reports.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-033 (projection infrastructure).
10. **Definition of Done:** Reports return required breakdowns (leadsByStatus count-only;
    opportunitiesByStage summed amount; quotes/contracts/payments by status with amounts);
    numbers traceable to persisted authorized records; empty → zero/empty; Sales denied;
    unauthorized records excluded; archived excluded by default.
11. **Acceptance method:** ACC-023 — report numbers traceable to records and matching metric
    groupings; Sales denied.
12. **Automated tests:** `TEST-BASIC-REPORT-001..006`; G12 rework
    `TEST-REPORTING-PROJECTION-INGEST-006` covers real S2S-populated reporting
    read model used by report queries. Type: Integration + E2E.
13. **Manual verification:** Create records, open reports, confirm counts/sums match; as Sales
    → denied.
14. **Traceability:** CIM-044/CIM-PROC-019 → PIM-025/PIM-BEH-032 → PSM-010 → CONTRACT-015/016 →
    ACC-023 → TEST-BASIC-REPORT-001..006.
15. **TDD:** Write grouping-traceability + empty + Sales-deny + archived-excluded tests first (fail).
    G12 second rework added `TEST-BASIC-REPORT-006` for Administrator all-scope
    report query over active records across teams.
16. **No-downgrade items:** Numbers derived from real persisted records via read model; no
    hardcoded metrics; archived excluded.
17. **Blocker:** None.

### TASK-035 — import-export-service + frontend: CSV import with row-level errors

1. **Task ID:** TASK-035
2. **Status:** Done
3. **Objective:** Authorized CSV import that validates rows, imports valid rows via target
   domain Command APIs, reports row-level errors without corrupting existing records, and
   denies Sales; CSV formula-injection-safe.
4. **Business capability:** CAP-010 Data import and export (primary ACC-020).
5. **Acceptance item:** ACC-020 (import half; export in TASK-036).
6. **Reference docs:** CIM-041, CIM-PROC-016, PIM-022, PIM-BEH-026, PSM-012, CONTRACT-017/018,
   `api-spec.md` "Import And Export Runs" (validation, dangerous-cell handling, 24h temp,
   Operation Status Contract), `integration-design.md` Import/Export Scope, FLOW-008,
   `ui-spec.md` UI-013, PM-034/035/036, ARCH-ACC-005/012, EDGE-026, ABUSE-007/008/009,
   TEST-CSV-IMPORT / TEST-ABUSE-CSVINJECT / TEST-ABUSE-IMPORTAUTHZ.
7. **File changes:**
   - `services/import-export/internal/handler/import_run.go`, `run_status.go`.
   - `services/import-export/internal/domain/row_validation.go` (required fields, enum, date,
     amount, dangerous-cell), `csv_safe.go`.
   - `services/import-export/internal/repo/run_repo.go`; calls target domain Command APIs via
     S2S (NO direct domain DB writes).
   - `services/import-export/internal/cleanup/temp_cleanup.go` (24h scheduled cleanup).
   - migrations `0002_import_export_runs.up.sql`.
   - `frontend/src/pages/importexport/Import.tsx`, `frontend/src/api/importexport.ts`.
   - tests: `import_run_test.go`, `row_validation_test.go`, `e2e/import.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-003, TASK-004, target services (TASK-007/010/011/013/017/018/020/024).
10. **Definition of Done:** Valid rows imported with success summary via target Command APIs;
    invalid rows reported with row-level errors and no corruption/overwrite; unsupported format
    rejected before mutation; Sales denied; per-row authorization; dangerous CSV cells rejected/
    escaped; operation logged; temp files cleaned after 24h.
11. **Acceptance method:** ACC-020 import — valid import + row-level errors + no corruption +
    authorization.
12. **Automated tests:** `TEST-CSV-IMPORT-001..004`, `TEST-ABUSE-CSVINJECT-001`,
    `TEST-ABUSE-IMPORTAUTHZ-001`. Type: Integration + E2E.
13. **Manual verification:** Import a mixed CSV → valid rows in, invalid reported, existing
    untouched; as Sales → denied.
14. **Traceability:** CIM-041/CIM-PROC-016 → PIM-022/PIM-BEH-026 → PSM-012 → CONTRACT-017/018 +
    FLOW-008 → ACC-020 → TEST-CSV-IMPORT-001..004.
15. **TDD:** Write mixed-row + unsupported-format + Sales-denied + formula-injection tests first (fail).
16. **No-downgrade items:** Mutations only via target Command APIs (no direct domain DB,
    ARCH-ACC-005); real row-level isolation; real authorization; real formula-injection safety.
17. **Blocker:** None.

### TASK-036 — import-export-service + frontend: CSV export (authorized records, archived excluded, confirm + log)

1. **Task ID:** TASK-036
2. **Status:** Done
3. **Objective:** Authorized CSV export including only authorized records, excluding archived
   by default, with explicit confirmation and audit log; Sales denied; export safely encodes
   dangerous cells.
4. **Business capability:** CAP-010 Data import and export (primary ACC-020).
5. **Acceptance item:** ACC-020 (export).
6. **Reference docs:** CIM-042, CIM-PROC-016, PIM-023, PIM-BEH-027, PSM-012, CONTRACT-017/018,
   `api-spec.md` "Import And Export Runs" (export metadata, dangerous-cell escaping),
   `ui-spec.md` UI-013 (Export confirmation), PM-037/038/039, ARCH-ACC-012, EDGE-027,
   ABUSE-010/011, TEST-CSV-EXPORT / TEST-ABUSE-EXPORTLEAK / TEST-ABUSE-EXPORTCONFIRM.
7. **File changes:**
   - `services/import-export/internal/handler/export_run.go`, `export_metadata.go`.
   - `services/import-export/internal/domain/export_scope.go` (authz scope, archived default-excluded).
   - `frontend/src/pages/importexport/Export.tsx` (confirmation: object, filters, archived
     inclusion, scope, count, audit notice).
   - tests: `export_run_test.go`, `e2e/export.spec.ts`.
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** TASK-035 (run infrastructure), TASK-030 (query scope).
10. **Definition of Done:** Export includes only authorized records (manager=team), archived
    excluded unless explicit authorized filter; Sales denied; explicit confirmation required;
    operation logged with filters/archived-inclusion/count/result; dangerous cells safely
    encoded; temp export file cleaned after 24h.
11. **Acceptance method:** ACC-020 export — authorized records only, Sales denied.
12. **Automated tests:** `TEST-CSV-EXPORT-001/002`, `TEST-ABUSE-EXPORTLEAK-001`,
    `TEST-ABUSE-EXPORTCONFIRM-001`. Type: Integration + E2E.
13. **Manual verification:** As manager, export team records (no archived) with confirmation;
    as Sales → denied.
14. **Traceability:** CIM-042/CIM-PROC-016 → PIM-023/PIM-BEH-027 → PSM-012 → CONTRACT-017/018 →
    ACC-020 → TEST-CSV-EXPORT-001/002.
15. **TDD:** Write Sales-denied + archived-excluded + confirm-required + formula-safe tests first (fail).
16. **No-downgrade items:** Real authorization-before-export; archived excluded by default;
    real audit log; formula-injection safe.
17. **Blocker:** None.

### TASK-037 — Persistence verification across refresh / re-login / restart

1. **Task ID:** TASK-037
2. **Status:** Done
3. **Objective:** Prove all core CRM data persists across browser refresh, logout/login, and
   service restart; failed saves are surfaced and never appear as successful persistent changes.
4. **Business capability:** CAP-011 Persistence and production operation (primary ACC-016).
5. **Acceptance item:** ACC-016.
6. **Reference docs:** CIM-045, CIM-PROC-021, PIM-BEH-033, all PSM-001..012 owned data,
   PSM-013, CONTRACT-020/014, `data-design.md`, DEC-008, ARCH-ACC-015, EDGE-030, SEC-018,
   TEST-PERSISTENCE.
7. **File changes:**
   - `e2e/persistence.spec.ts` (create → refresh → re-login → assert; restart compose service
     → assert).
   - `services/<svc>/internal/repo/*_test.go` additions for restart-survival integration tests
     where missing.
   - `services/<svc>/internal/handler/*` ensure failed saves return error (no false success).
8. **Owner agent:** backend-engineer, frontend-engineer
9. **Prerequisites:** All data-owning slices (TASK-007..036 persistence-bearing tasks).
10. **Definition of Done:** Data created in each owning service survives refresh, logout/login,
    and a container restart; a forced save failure surfaces an error and creates no record; no
    P0 path uses mock/in-memory-only behavior.
11. **Acceptance method:** ACC-016 — persistence verified after refresh, re-login, restart;
    failed save surfaced.
12. **Automated tests:** `TEST-PERSISTENCE-001..005`. Type: Integration + E2E + Manual.
13. **Manual verification:** Create records; `docker compose restart <svc>`; confirm data
    present; simulate DB error → save fails visibly.
14. **Traceability:** CIM-045/CIM-PROC-021 → PIM-BEH-033 → PSM-013 (+ PSM-001..012) →
    CONTRACT-020/014 → ACC-016 → TEST-PERSISTENCE-001..005.
15. **TDD:** Write refresh/re-login/restart survival + failed-save tests first (fail).
16. **No-downgrade items:** Real PostgreSQL persistence; restart survival proven; no in-memory
    substitute satisfies ACC-016 (SVC-ACC-011).
17. **Blocker:** None.

### TASK-038 — Data classification & retention enforcement (cross-cutting)

1. **Task ID:** TASK-038
2. **Status:** Done
3. **Objective:** Carry the committed data classification on records/logs and enforce minimum
   retention anchored to lifecycle events (never shortened, never hard-deleted); mask sensitive
   before/after values per classification in denied/error/log states.
4. **Business capability:** CAP-008/CAP-011 cross-cutting (primary ACC-014; also ACC-016/022).
5. **Acceptance item:** ACC-014 (classification/retention reinforcing history/audit; ACC-022/016 supported).
6. **Reference docs:** CIM-048/049, PIM-027/028, PIM-INV-049/050/051/052, PIM-BEH-034, PSM-013,
   CONTRACT-013/014, `privacy-requirements.md` (PRIV-001..016 + Retention Policy durations),
   `data-design.md` (Backup/retention), COMP-013, ABUSE-014, SEC-014, TEST-RETENTION.
7. **File changes:**
   - `shared/contracts/classification.go` (classification tag enum).
   - `services/<svc>/internal/domain/*`: attach classification + retention metadata to records/events.
   - `services/audit-history/internal/domain/event.go`: enforce append-only retention, masked safe summaries.
   - `services/<svc>/internal/handler/*`: safe-summary masking on errors/denials.
   - migrations: add classification/retention columns where modeled.
   - tests: `retention_test.go`, `classification_mask_test.go`.
8. **Owner agent:** backend-engineer
9. **Prerequisites:** TASK-004, and data-owning slices.
10. **Definition of Done:** Each record/log carries a committed classification; retention
    metadata anchored to lifecycle events; no path shortens retention below the committed
    minimum or hard-deletes; sensitive before/after values masked per classification in
    error/denial/log states.
11. **Acceptance method:** ACC-014 reinforced — classification/retention metadata present and
    masking enforced in denied/error/log states.
12. **Automated tests:** `TEST-RETENTION-001` (retention ≥ minimum, append-only),
    `TEST-RETENTION-002` (classification masking), `TEST-ABUSE-... value leakage` covered by
    `TEST-RETENTION-002`. Type: Integration.
13. **Manual verification:** Inspect a record/log for classification + retention metadata;
    trigger a denial → no sensitive raw value leaked.
14. **Traceability:** CIM-048/049 → PIM-027/028/PIM-INV-050/PIM-BEH-034 → PSM-013 →
    CONTRACT-013/014 → ACC-014 → TEST-RETENTION-001/002.
15. **TDD:** Write retention-not-shortened + masking tests first (fail).
16. **No-downgrade items:** Real classification/retention metadata persisted; real masking
    (not cosmetic); no hard delete; retention never shortened below committed minimum.
17. **Blocker:** None.

---

## Phase 2 — Deployment + Release Evidence (G11/G12 evidence)

### TASK-039 — Deploy on runtime host: Docker Compose, reverse proxy, HTTPS/TLS, security group, health/monitoring

1. **Task ID:** TASK-039
2. **Status:** Done
3. **Objective:** Deploy the full stack on the committed runtime host via Docker Compose behind
   the existing reverse proxy with HTTPS/TLS, restricted network exposure, and health/monitoring
   evidence; CRM reachable and connected to persistent services.
4. **Business capability:** CAP-011 Persistence and production operation (primary ACC-017).
5. **Acceptance item:** ACC-017 (deploy/operate + HTTPS/TLS/security-group/monitoring evidence).
6. **Reference docs:** CIM-046, CIM-PROC-022, PIM-OPEN-004, PSM-014, `deployment-notes.md`
   (Deployment Target `srv-volcengine-sh-01`, co-location, Network Exposure, Production Endpoint
   Strategy, Health And Observability, Operator Access, paths `/opt/crm-system/...`),
   `authz-architecture.md` (Secure/HttpOnly/SameSite cookie), ARCH-ACC-008/013/014/015,
   NFR-003, DEC-004, TEST-DEPLOY-SMOKE.
7. **File changes:**
   - `docker-compose.prod.yml` (prod env, restart policies, volumes
     `/opt/crm-system/volumes/postgres`, logs `/opt/crm-system/logs`).
   - `deploy/nginx/crm.conf` (reverse-proxy server_name/subpath, HTTPS, redirect 80→443,
     security headers HSTS/X-Content-Type-Options/Referrer-Policy/CSP).
   - `deploy/healthcheck/` aggregated health probe; `deploy/monitoring/` notes for CPU/mem/disk/
     container/backup monitoring + thresholds.
   - `docs/release/acc-017-evidence-template.md` (endpoint, TLS status, security-group inbound
     rules, opened ports, health URL, timestamp, operator, smoke result).
   - `deploy/ops/operator-access.md` (named least-privilege deploy/ops user; root restricted to
     provisioning/emergency only; SSH key ownership + sudo boundary) — per `deployment-notes.md`
     "Operator Access".
8. **Owner agent:** backend-engineer (implementation) + infrastructure-ops (runtime environment + release evidence, per `deployment-notes.md` / `process/process-gap-register.md`)
9. **Prerequisites:** TASK-001..038 (a complete deployable system).
10. **Definition of Done:** Stack runs on `srv-volcengine-sh-01` via Compose; HTTPS endpoint
    with valid TLS (or release blocked until configured); HTTP→HTTPS redirect; secure cookies;
    security headers set; PostgreSQL/internal ports/backup dir NOT publicly exposed; per-service
    health green; monitoring + thresholds recorded; a named least-privilege deploy/ops user exists
    with root restricted to provisioning/emergency and SSH key ownership + sudo boundary documented
    (deployment-notes "Operator Access"); ACC-017 evidence template completed at G11.
11. **Acceptance method:** ACC-017 — deployment smoke (reachable, configured, persistent
    services connected); release-evidence verified at G11/G12 (HTTPS/TLS, security group,
    monitoring).
12. **Automated tests:** `TEST-DEPLOY-SMOKE-001/002` (reachable+configured; misconfig/unavailable
    dependency blocks readiness). Type: Manual + Integration.
13. **Manual verification:** Hit the HTTPS endpoint; verify TLS + security headers; confirm
    PostgreSQL port not reachable publicly; check each `/health`.
14. **Traceability:** CIM-046/CIM-PROC-022 → (PIM-OPEN-004 deployment) → PSM-014 →
    deployment-notes.md → ACC-017 → TEST-DEPLOY-SMOKE-001/002.
15. **TDD:** Write the deploy smoke checks first (fail until deployed); evidence template filled
    at G11.
16. **No-downgrade items:** Real HTTPS/TLS (no self-signed-only for production); real restricted
    exposure; real health/monitoring — not a screenshot of a local run. ARCH-ACC-008/013/014/015
    are `Release-evidence pending` and proven at G11, audited at G12.
17. **Blocker:** Reopened by G12 on 2026-06-03. BLK-G12-003/007 are resolved:
    Codex added fail-first `TEST-DEPLOY-SG-001`, created dedicated security group
    `sg-366ptx1bxp9ts1e710babmc8y`, moved CRM ENI `eni-13e8tbocd8f0g79iu5jer8idt`
    to that group only, exported raw Volcengine API evidence, and verified public
    TCP ingress is limited to `22`, `80`, and `443` with no public `8080`, `5432`,
    `8088`, `8443`, or `3389`. Evidence:
    `docs/release/evidence/volcengine-security-group-dedicated-raw-2026-06-03.json`;
    `docs/release/evidence/volcengine-security-group-rework-transcript-2026-06-03.txt`.
    BLK-G12-006 is also resolved: captured HTTPS/redirect/openssl/certbot,
    external-edge negative probes from `srv-aliyun-bj-01`, restore catalog counts,
    `sshd -T` effective hardening, and distinct `crm-deploy`/`crm-ops` SSH key
    fingerprints under `docs/release/evidence/`; `scripts/test_release_evidence_transcripts.sh`
    passed. BLK-G12-009 is resolved: `scripts/test_security_group_evidence.py`
    now treats read-only exports as verification-only and requires mutating
    `CreateSecurityGroup`/`AuthorizeSecurityGroupIngress`/`RevokeSecurityGroupIngress`/
    `ModifyNetworkInterfaceAttributes` RequestIds with HTTP 200 for `--apply`;
    `scripts/test_security_group_evidence_checker.sh` passed. TASK-039 is complete
    again, pending independent G12 re-audit.
    Previous G11 text retained for historical context: Resolved on 2026-06-03. The HTTPS endpoint blocker was resolved
    with approved endpoint `https://118.196.44.193` and a Let's Encrypt IP certificate.
    Volcengine security-group API evidence was exported for security group
    `sg-1pm4k7f37z8xs643rg0fvk85e` bound to instance `i-yemoz0an7kk36d2c9bp6` via ENI
    `eni-13e8tbocd8f0g79iu5jer8idt`; it confirms CRM gateway `8080` and PostgreSQL
    `5432` are not publicly allowed from `0.0.0.0/0`. The user approved releasing
    previous deployments; Codex stopped and removed the host-network Hermes container, so
    host-level `8642` no longer listens and CRM smoke still passes. The user removed
    old/non-CRM Volcengine security-group rules for TCP `8088`, TCP `8443`, and TCP
    `3389`; API post-cleanup verification confirms only public TCP `22`, `80`, and
    `443` remain. **Pre-G8
    condition:** Security Compliance must review the operator-access design (SSH access, key
    ownership, sudo boundary) before G8 implementation tasks are approved (deployment-notes
    "Operator Access").

### TASK-040 — Encrypted off-server backup + restore rehearsal (release evidence)

1. **Task ID:** TASK-040
2. **Status:** Done
3. **Objective:** Automated encrypted PostgreSQL backups with 7-day local retention, an encrypted
   off-server copy to the Beijing backup target, and a documented restore rehearsal proving
   recoverability — the carried release blocker for ACC-016/017.
4. **Business capability:** CAP-011 Persistence and production operation (primary ACC-017; also ACC-016).
5. **Acceptance item:** ACC-017 (backup + restore evidence; supports ACC-016 durability).
6. **Reference docs:** PSM-014, FLOW-009 (Backup & restore), Backup ownership row (PSM Data
   Ownership), `deployment-notes.md` (PostgreSQL Backup, Restore Requirement, Production Backup
   Release Rule, off-server target `srv-aliyun-bj-01`), ARCH-ACC-004, RISK-002/004,
   `privacy-requirements.md` (backups must not shorten retention).
7. **File changes:**
   - `deploy/backup/backup.sh` (timestamped, no-overwrite, encrypt, 7-day prune, log success/
     failure, path `/opt/crm-system/backups/postgres`).
   - `deploy/backup/offsite-copy.sh` (encrypted copy to `srv-aliyun-bj-01`).
   - `deploy/backup/restore-rehearsal.md` (procedure: restore into controlled target, verify
     users/roles/records/history/logs/DB permissions, record command/timestamp/file/checksum/
     encryption step/operator/result).
   - cron/compose entry for daily ~02:00 backup.
   - `docs/release/acc-017-backup-evidence-template.md`.
8. **Owner agent:** infrastructure-ops (+ backend-engineer for backup container wiring)
9. **Prerequisites:** TASK-039.
10. **Definition of Done:** Daily encrypted backup produced (no overwrite, 7-day prune, logged);
    encrypted off-server copy to `srv-aliyun-bj-01` succeeds; at least one restore rehearsal
    verifies users/roles/records/history/logs/DB permissions with documented command/timestamp/
    file/checksum/encryption/operator/result; backups do not shorten application retention.
11. **Acceptance method:** ACC-017 backup/restore evidence (ARCH-ACC-004), proven at G11 and
    audited at G12.
12. **Automated tests:** Backup-job integration check (file produced, encrypted, pruned);
    restore rehearsal is documented evidence (Manual/Integration). Type: Integration + Manual.
13. **Manual verification:** Run backup; verify encrypted file + off-server copy + checksum;
    perform restore into a controlled target and verify data.
14. **Traceability:** PSM-014/FLOW-009 → deployment-notes.md Backup/Restore → ARCH-ACC-004 →
    ACC-017/ACC-016 → TEST-DEPLOY-SMOKE (release-evidence extension).
15. **TDD:** Write the backup-artifact + restore-verification checks first; they FAIL until the
    backup job exists and a real encrypted artifact is produced (fail-first); rehearsal evidence
    recorded at G11.
16. **No-downgrade items:** Real encrypted off-server copy (same-host-only backup does NOT
    satisfy release); real restore rehearsal; encryption keys outside repo/service config.
    ARCH-ACC-004 is `Release-evidence pending` until proven at G11 and audited at G12.
17. **Blocker:** Resolved on 2026-06-03. Encrypted backup
    `crm-postgres-20260603T104620Z.sql.gz.enc` was produced on `srv-volcengine-sh-01`
    with checksum `f7381eaa9d246126cac93b304a147ea01721de404f52de76b340e3cfa9ba9d2a`,
    copied to `srv-aliyun-bj-01`, and verified with `sha256sum -c` returning `OK`.
    Restore rehearsal run `20260603T104837Z` restored the off-server encrypted backup
    into a controlled PostgreSQL target and passed verification for roles, `crm_system`,
    service schemas, `identity_authz.users`, audit-history tables, and service permission
    roles. G12 second rework reconciled the restore catalog evidence to the captured
    transcript counts: roles `10`, schemas `10`, and service permission roles `10`.
    `crm-backup.timer` is enabled and active with next run
    `2026-06-04 02:00:00 CST`. Evidence:
    `docs/release/acc-017-backup-evidence-template.md` and
    `docs/release/evidence/backup-restore-rehearsal-2026-06-03.json`.
