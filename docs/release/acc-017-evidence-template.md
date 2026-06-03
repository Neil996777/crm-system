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
| Open public ports | Host-level observed ports after old deployment release: `22` SSH and `80/443` Nginx CRM ingress. Volcengine security group post-cleanup verification shows public TCP `22`, `80`, and `443` only. |
| Restricted ports | Host-level Docker/iptables evidence shows CRM gateway bound to `127.0.0.1:8080`; PostgreSQL has no host port mapping and is exposed only inside Compose. Volcengine API evidence confirms CRM `8080` and PostgreSQL `5432` are not allowed from `0.0.0.0/0`. |
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

## G11 Closure

TASK-039 deployment/security-group closure evidence is complete as of
2026-06-03. The Volcengine provider security-group record confirms CRM `8080`
and PostgreSQL `5432` are not publicly allowed. Host-level Hermes `8642` has
been released, and the old/non-CRM public TCP `8088`, `8443`, and `3389` rules
have been removed from the Volcengine security group.

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

This intermediate evidence was superseded by the Volcengine API export and
post-cleanup verification recorded below.

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
- At that time, host listeners still included pre-existing non-CRM Hermes
  `8642`, and cloud security-group evidence had not yet been exported.

This intermediate closeout attempt did not close TASK-039 by itself. The later
Volcengine API export, old deployment release, and post-cleanup security-group
verification below provide the closure evidence.

## Volcengine Security Group API Export

Volcengine API export on 2026-06-03 identified the bound primary ENI and security
group:

- Instance: `i-yemoz0an7kk36d2c9bp6`
- Public IP: `118.196.44.193`
- Private IP: `172.31.8.67`
- Network interface: `eni-13e8tbocd8f0g79iu5jer8idt`
- Security group: `sg-1pm4k7f37z8xs643rg0fvk85e` (`Default`)
- Evidence files:
  - `docs/release/evidence/old-deployment-release-2026-06-03.json`
  - `docs/release/evidence/volcengine-ecs-describe-instance-2026-06-03.json`
  - `docs/release/evidence/volcengine-security-group-post-cleanup-2026-06-03.json`
  - `docs/release/evidence/volcengine-security-group-summary-2026-06-03.json`
  - `docs/release/evidence/volcengine-security-group-raw-2026-06-03.json`
  - `docs/release/evidence/volcengine-security-group-dedicated-raw-2026-06-03.json`
  - `docs/release/evidence/volcengine-security-group-rework-transcript-2026-06-03.txt`

Inbound security-group evidence:

- Publicly allowed for CRM ingress: TCP `80`, TCP `443`.
- Publicly allowed for administration: TCP `22`.
- Post-G12 rework CRM dedicated security group:
  `sg-366ptx1bxp9ts1e710babmc8y` (`crm-system-prod-public`).
- CRM ENI `eni-13e8tbocd8f0g79iu5jer8idt` is bound only to
  `sg-366ptx1bxp9ts1e710babmc8y`, not the shared `Default` group.
- Dedicated security-group public TCP rules: `22`, `80`, and `443`.
- CRM gateway TCP `8080` is not allowed from `0.0.0.0/0`.
- PostgreSQL TCP `5432` is not allowed from `0.0.0.0/0`.
- Old/non-CRM public TCP `8088`, `8443`, and `3389` are absent from the
  final raw security-group export.

API note: the initial IAM policy lacked `ecs:DescribeInstances`, so the first
ECS read returned `403`. After `ECSReadOnlyAccess` was added,
`ecs:DescribeInstances` returned HTTP `200` for `i-yemoz0an7kk36d2c9bp6` and
confirmed the same ENI and security group already identified through VPC API
evidence. `vpc:DescribeSecurityGroupAttributes` returned the inbound rules.

## Old Deployment Release

The user approved releasing previous deployments after the project transfer.
Codex released host-level Hermes exposure on 2026-06-03:

- Container `hermes` (`hermes-agent`, host network, restart `unless-stopped`) had
  been listening on `0.0.0.0:8642`.
- Codex set the container restart policy to `no`, stopped it, and removed the
  container.
- Post-release host listeners no longer include `8642`.
- CRM remained healthy, with all 11 CRM containers running and server-side
  `scripts/test_deploy_smoke.sh` passing.

Cloud security-group cleanup was completed after the user removed public TCP
`8088`, `8443`, and `3389` in Volcengine. API post-cleanup verification confirms
these old/non-CRM rules are gone and only public TCP `22`, `80`, and `443`
remain.

G12 rework supersedes the earlier hand-authored post-cleanup summary. Codex used
the Volcengine OpenAPI to create dedicated security group
`sg-366ptx1bxp9ts1e710babmc8y`, bind the CRM primary ENI only to that group, and
export the raw final `DescribeNetworkInterfaces`, `DescribeSecurityGroups`, and
`DescribeSecurityGroupAttributes` responses. `TEST-DEPLOY-SG-001` passed against
that raw evidence.
