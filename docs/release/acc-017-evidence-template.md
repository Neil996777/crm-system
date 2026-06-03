# ACC-017 Release Evidence

| Field | Evidence |
|---|---|
| Endpoint | `https://118.196.44.193` |
| HTTP redirect endpoint | `http://118.196.44.193` -> `https://118.196.44.193` |
| TLS certificate status | Valid server-side on 2026-06-03; certificate SAN is `IP Address:118.196.44.193` |
| TLS issuer / validity | Let's Encrypt `YE1`; notBefore `Jun 3 05:46:05 2026 GMT`; notAfter `Jun 9 21:46:04 2026 GMT` |
| TLS renewal | `crm-certbot-renew.timer` enabled; dry-run succeeded on 2026-06-03; TLS expiry check passes with 48h threshold |
| Security headers | Validated server-side on 2026-06-03: HSTS, X-Content-Type-Options, Referrer-Policy, CSP |
| Runtime host | `srv-volcengine-sh-01` / `118.196.44.193` |
| Deployment path | `/opt/crm-system/current` |
| PostgreSQL data path | `/opt/crm-system/volumes/postgres` |
| Log path | `/opt/crm-system/logs` |
| Open public ports | Host-level observed ports: `22` SSH, `80/443` Nginx CRM ingress, pre-existing non-CRM `8642` Hermes. Cloud security-group export is still pending. |
| Restricted ports | Host-level Docker/iptables evidence shows CRM gateway bound to `127.0.0.1:8080`; PostgreSQL has no host port mapping and is exposed only inside Compose. Cloud security-group evidence is still pending. |
| Health URL | `https://118.196.44.193/health` |
| Deployment timestamp | 2026-06-03 15:20 CST |
| Operator | Initial deployment by root; named `crm-deploy` and `crm-ops` accounts created with SSH keys and sudo boundaries on 2026-06-03 |
| Smoke result | Server-side `TEST-DEPLOY-SMOKE-001/002` passed on 2026-06-03 |
| Monitoring thresholds | `deploy/monitoring/thresholds.md` |
| Backup / restore evidence | TASK-040 |

## Smoke Test Commands

```bash
bash deploy/healthcheck/check_endpoint.sh
```

## G11 Evidence Notes

Production closure requires real evidence from `srv-volcengine-sh-01`: endpoint,
TLS certificate details, HTTP redirect, security headers, security-group inbound
rules, opened ports, service health, monitoring thresholds, deployment timestamp,
and operator identity.

## Remaining G11 Blocker

This evidence is not sufficient to mark TASK-039 Done until the Volcengine cloud
security-group inbound rules are exported or otherwise recorded. The host-level
state is consistent with restricted CRM exposure, but the provider security-group
record is still required for ACC-017 closure.

## Infrastructure Ops Read-Only Check

Infrastructure Ops read-only verification on 2026-06-03 confirmed:

- CRM Compose project `current` is running with all 11 containers healthy.
- Nginx exposes CRM through `80/443`; `80` redirects to HTTPS.
- `https://118.196.44.193/health` returns `HTTP/2 200` from the server.
- The CRM gateway is host-published only as `127.0.0.1:8080->8080/tcp`.
- PostgreSQL has no host port mapping; it is container/Compose internal only.
- Host-level evidence does not show CRM `8080` or PostgreSQL `5432` listening on
  public interfaces.
- Pre-existing public `8642` is Hermes, not CRM, but still requires owner and
  Security Compliance review because it is exposed on the same production
  candidate host.

Remaining external evidence needed: Volcengine security-group inbound rules for
instance `i-yemoz0an7kk36d2c9bp6`, including protocol, port, source CIDR, policy,
and timestamp.

## TASK-039 Closeout Attempt

Closeout verification on 2026-06-03 16:07-16:08 CST confirmed:

- Runtime instance hostname is `iv-yemoz0an7kk36d2c9bp6`, matching instance ID
  `i-yemoz0an7kk36d2c9bp6`.
- Server-side metadata probes did not return cloud security-group rules.
- Local Codex environment has no Volcengine CLI or Volcengine credential
  configuration available for read-only API export.
- Server-side `scripts/test_deploy_smoke.sh` passed for
  `https://118.196.44.193`.
- `crm-certbot-renew.timer` is active and enabled.
- TLS certificate remains issued by Let's Encrypt `YE1` with SAN
  `IP Address:118.196.44.193`.
- `https://118.196.44.193/health` returns `HTTP/2 200` with HSTS,
  X-Content-Type-Options, Referrer-Policy, and CSP headers.
- Current host listeners remain consistent with restricted CRM exposure:
  public `80/443` for Nginx, public `22` for SSH, public pre-existing non-CRM
  `8642` for Hermes, CRM gateway on `127.0.0.1:8080`, and no PostgreSQL host
  port mapping.
- Workspace infrastructure registers still list Hermes `8642` owner as `TBD`
  and require Security Compliance review.

TASK-039 cannot be marked Done from this closeout attempt. Required closure
evidence remains:

- Volcengine console/API export for the inbound security-group rules bound to
  instance `i-yemoz0an7kk36d2c9bp6`, proving CRM `8080` and PostgreSQL `5432`
  are not publicly allowed.
- Owner and Security Compliance disposition for the pre-existing public Hermes
  `8642` exposure on the same production candidate host.
