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
