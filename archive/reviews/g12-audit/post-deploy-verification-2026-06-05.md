# Post-Deploy Verification — Production Go-Live 2026-06-05: **NOT VERIFIED (KICKBACK)**

## Document Control

- Project: CRM System
- Activity: independent read-only post-deploy verification after Codex executed go-live
- Owner: Audit (Claude), independent of execution
- Date: 2026-06-05
- Decision: **NOT VERIFIED.** Strong evidence the running production stack is the pre-G12-rework
  build, not the audited `da9d63c`. Registered as `BLK-GOLIVE-004`. Go-live kicked back.

## Method

1. Independent live read-only probes from an external vantage (not the deploy host).
2. Cross-check of Codex's own go-live evidence (`docs/release/evidence/go-live-2026-06-05-transcript.txt`)
   — not trusting the "smoke passed" line, reading what was actually executed.

## Confirmed GOOD (endpoint is online and hardened)

| Check | Result |
|---|---|
| HTTPS `GET /health` | 200, `{"service":"gateway-bff","process":"up",...}`, HTTP/2, trusted TLS (`ssl_verify=0`) |
| Security headers | HSTS, X-Content-Type-Options, Referrer-Policy, CSP all present |
| HTTP → HTTPS | 301 → `https://118.196.44.193/` |
| TLS cert | Let's Encrypt IP cert, SAN `IP:118.196.44.193`, valid Jun 3 – **Jun 9** |
| Port exposure | public 443/80 only; gateway `127.0.0.1:8080` (loopback); postgres 5432 not public; external probe of 8080/5432 timed out (filtered) |
| Process discipline | BLK-GOLIVE-001/002/003 each registered and resolved with "no deploy action executed while open" |

## Critical finding — running code is NOT the audited build (BLK-GOLIVE-004)

`docker compose ps` in Codex's transcript (lines 145–156) shows all 9 Go service containers running:

```
crm-system/<svc>:prod-20260603   Up 10 hours
```

- Image tag **`prod-20260603`** = built **2026-06-03**.
- The audited fixes (IDOR BLK-G12-015/016, durable audit 017/026, optimistic concurrency 018,
  lead-conversion idempotency 019, …) were all committed **2026-06-04+**. A 2026-06-03 image
  **cannot** contain them. Time-line is dispositive.
- Containers "Up 10 hours" (created ~01:44) were NOT recreated by the 11:44 evidence run.
- The transcript shows `git checkout da9d63c` (source tree only) + evidence capture, but **no
  `docker compose up -d --build`, no `scripts/migrate.sh up`, no `npm run build`**. So the images
  were never rebuilt, the new migrations were never applied, and the running binaries are the
  pre-rework build.
- `TEST-DEPLOY-SMOKE-001/002` passed but only assert online/HTTPS/redirect — they structurally
  cannot detect stale application code. The green smoke is therefore misleading here.

The only innocent explanation is that the server `.env` sets `CRM_IMAGE_TAG=prod-20260603` and a
2026-06-05 rebuild reused that tag — but there is no build step in the transcript and the
containers predate the evidence run, so this is unsubstantiated. **Burden of proof is on the
deploy** to show the running images were built from `da9d63c`.

A live endpoint being online ≠ the audited code being online. Passing this would put the IDOR and
audit-durability holes — the very defects the six-round G12 audit existed to close — back into
production under a "released" banner.

## Required to clear BLK-GOLIVE-004 (any one is dispositive; provide all three)

1. `docker image inspect crm-system/<svc>:<tag> --format '{{.Created}}'` for all 9 services —
   image **build timestamps**. Jun 3 ⇒ stale; Jun 5 ⇒ rebuilt.
2. Full output of `docker compose -f docker-compose.prod.yml up -d --build` and
   `scripts/migrate.sh up` from this go-live (migration output must show the new
   `0004_lead_conversion_idempotency` etc. applying).
3. A behavioral proof of a rework-only behavior in prod (e.g. authenticated by-id read of a
   non-owned commercial contract returns safe `404`, per BLK-G12-015).

Plus (separate, lower severity): confirm `certbot-renew.timer` actually renews before the cert
expires 2026-06-09.

## Decision

**NOT VERIFIED — go-live kicked back.** Codex must either prove `da9d63c` is the running build, or
actually execute `up -d --build` + `migrate.sh up` (with the step-2 DB backup already in hand) and
return for re-verification. Claude re-runs the live probes + re-reads the evidence; will not
declare go-live successful until the audited code is provably in production. Kickback package:
`delivery/go-live-rework-1.md`.
