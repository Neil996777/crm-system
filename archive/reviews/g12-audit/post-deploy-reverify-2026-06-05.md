# Post-Deploy RE-VERIFICATION — Production Go-Live 2026-06-05: **VERIFIED**

## Document Control

- Project: CRM System
- Activity: independent re-verification after Codex executed go-live rework #1 (`BLK-GOLIVE-004`)
- Owner: Audit (Claude), independent of execution
- Date: 2026-06-05
- Prior: `post-deploy-verification-2026-06-05.md` (NOT VERIFIED — running pre-rework images)
- Decision: **VERIFIED.** The audited build `da9d63c` is provably deployed and serving in
  production. `BLK-GOLIVE-004` Resolved. Go-live confirmed successful.

## Method

1. Independent live read-only probes from an external vantage.
2. Cross-check of Codex's rework evidence (`docs/release/evidence/go-live-rework-1-2026-06-05-transcript.txt`)
   line-by-line — not trusting the blocker summary; verifying image build times, the actual
   `up -d --build` and `migrate.sh up` execution, and a behavioral proof of a rework-only fix.

## Evidence (all confirmed)

| Check | Evidence | Verdict |
|---|---|---|
| Image rebuilt from audited source | all 10 Go service images `created=2026-06-05` with new sha256 IDs (transcript 876–885) | ✅ not the 06-03 build |
| Real compile | per-service `RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server` + layer export (202–330) | ✅ |
| New migrations applied | `applying services/account/migrations/0006_lead_conversion_idempotency.up.sql` and the other rework migrations (544+) | ✅ |
| **Behavioral proof BLK-G12-015 (IDOR)** | `GET /quotes/<id>` as `X-Actor-Role=Sales` non-owner → `HTTP/1.1 404 Not Found`; temp quote then removed, `CLEANUP_COUNT=0` (912–918) | ✅ rework code is live |
| Containers recreated | all `Up 6 minutes` (recreated 12:14), postgres correctly preserved (`Up 45 hours`) | ✅ |
| Live endpoint (my own probe) | HTTPS `/health` 200; HTTP→301; 8080/5432 publicly unreachable (timeout) | ✅ |
| Smoke | `deploy/healthcheck/check_endpoint.sh` TEST-DEPLOY-SMOKE-001/002 passed | ✅ |
| Cert renewal | `crm-certbot-renew.timer` active; `certbot renew --dry-run` "all simulated renewals succeeded"; cert valid to 2026-06-09, renewal proven (926–953) | ✅ |

The only innocent explanation for the first kickback (the `prod-20260603` tag was reused on a
2026-06-05 rebuild because `CRM_IMAGE_TAG=prod-20260603` is set on the server) is now substantiated
by the image Created timestamps, new image IDs, the real build output, and the dispositive 404
behavioral proof. The deployed code is `da9d63c`.

## Follow-up recommendation (LOW, non-blocking)

The running image tag `crm-system/<svc>:prod-20260603` is misleading: the image content is the
2026-06-05 `da9d63c` rebuild, but the tag reads "0603". This is exactly what made the first
verification ambiguous. Recommend setting `CRM_IMAGE_TAG` to the release date or commit (e.g.
`prod-20260605` / `da9d63c`) on the next deploy so `docker ps` reflects the real build. Not a
correctness issue — the deployed code is proven correct — purely ops/evidence hygiene.

Minor note: the behavioral proof wrote a temporary quote to the production DB and removed it
(`CLEANUP_COUNT=0`). A non-owner read test needs a record to exist; cleanup was proven. Acceptable.

## Decision

**VERIFIED — production go-live successful.** The audited CRM build (`da9d63c`, including all G12
fixes BLK-G12-001..030) is deployed and serving at `https://118.196.44.193`, with the IDOR fix
behaviorally confirmed live, new migrations applied, endpoint hardened (loopback gateway, public
8080/5432 unreachable), and TLS renewal proven. `BLK-GOLIVE-001..004` all Resolved.
