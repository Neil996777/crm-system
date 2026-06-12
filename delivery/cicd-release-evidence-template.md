# CRM CI/CD Release Evidence Template

Status: Template for G10/G11 execution evidence. Fill during deployment. Do not
record secret values.

## Document Control

| Field | Evidence |
|---|---|
| Release content commit | `66d2531` |
| CI run ID / URL / local artifact path | `<fill>` |
| Release bundle path | `<fill>` |
| Target host | `srv-volcengine-sh-01` / `118.196.44.193` |
| Executor | `<fill>` |
| Deployment timestamp | `<fill>` |
| G10/G11 status | `<fill>` |
| G12 auditor | Claude, independent audit |

## 1. CI Test Results

| Suite | Command / job | Result | Artifact |
|---|---|---|---|
| Backend Go tests | `<fill>` | `<pass/fail>` | `<path>` |
| Frontend build | `npm run build` from `frontend/` | `<pass/fail>` | `<path>` |
| Playwright e2e | `npm run test:e2e` from `frontend/` | `<pass/fail>` | `<path>` |
| Static release checks | no host-build commands; no `latest`; no `build:` | `<pass/fail>` | `<path>` |

## 2. Image Digest To Commit

Application images must map to commit `66d2531`. Production identity is digest,
not tag.

| Runtime | Image tag | Image digest / ID | Tar SHA-256 | Source commit | Revision label verified |
|---|---|---|---|---|---|
| gateway-bff | `crm-system/gateway-bff:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| identity-authz | `crm-system/identity-authz:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| lead | `crm-system/lead:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| account | `crm-system/account:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| opportunity | `crm-system/opportunity:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| commercial | `crm-system/commercial:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| work | `crm-system/work:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| audit-history | `crm-system/audit-history:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| reporting | `crm-system/reporting:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| import-export | `crm-system/import-export:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| frontend-web | `crm-system/frontend-web:66d2531` | `<sha256>` | `<sha256>` | `66d2531` | `<yes/no>` |
| postgres | `postgres:16-alpine@sha256:<digest>` | `<sha256>` | n/a | n/a, upstream image | n/a |

Attach:

- `image-manifest.json`
- `image-manifest.sha256`
- `release-crm-system-66d2531.sha256`
- `docker image inspect` output after load
- `docker inspect` output for running containers

## 3. Deploy Transcript

Attach the full transcript, not a summary:

| Step | Transcript artifact |
|---|---|
| Bundle transfer | `<path>` |
| Checksum verification | `<path>` |
| `docker load` | `<path>` |
| Loaded image digest verification | `<path>` |
| Pre-migration backup | `<path>` |
| Image-only `docker compose up -d` | `<path>` |
| Migration from release artifact SQL | `<path>` |
| Nginx config test and reload | `<path>` |
| Health checks | `<path>` |

Transcript must not contain secret values.

## 4. Post-Deploy Health Checks

| Check | Expected | Actual | Evidence |
|---|---|---|---|
| Compose health | All CRM containers healthy | `<fill>` | `<path>` |
| HTTPS `/health` | HTTP 200 | `<fill>` | `<path>` |
| HTTP redirect | 301/308 to HTTPS | `<fill>` | `<path>` |
| Security headers | HSTS, X-Content-Type-Options, Referrer-Policy, CSP | `<fill>` | `<path>` |
| TLS certificate | Valid for endpoint | `<fill>` | `<path>` |
| Public gateway negative check | `8080` not public | `<fill>` | `<path>` |
| Public PostgreSQL negative check | `5432` not public | `<fill>` | `<path>` |
| Frontend loopback | `frontend-web` only loopback/internal | `<fill>` | `<path>` |
| Secret handling | secret path exists; values not printed | `<fill>` | `<path>` |

## 5. Rollback Point

| Field | Evidence |
|---|---|
| Previous-good image manifest | `<path or first clean deploy note>` |
| Previous-good release bundle checksum | `<sha256 or n/a for first clean deploy>` |
| Pre-deploy database backup | `<path>` |
| Backup checksum | `<sha256>` |
| Off-server copy evidence | `<path>` |
| Restore path / rehearsal reference | `<path>` |
| Rollback command transcript or dry-run plan | `<path>` |

For the first clean standard deployment after `/opt/crm-system` was removed, if
there is no previous CRM runtime, record that fact and the pre-deploy empty/seed
database backup state. After the first deployment, retain this bundle as the
previous-good rollback point for the next release.

## 6. No Host Build / No Source Dependence Audit

| Check | Result | Evidence |
|---|---|---|
| Production transcript contains no `npm run build` | `<pass/fail>` | `<path>` |
| Production transcript contains no `docker build` | `<pass/fail>` | `<path>` |
| Production transcript contains no `docker compose build` | `<pass/fail>` | `<path>` |
| Production transcript contains no `docker compose up --build` | `<pass/fail>` | `<path>` |
| Production transcript contains no `git checkout` for build | `<pass/fail>` | `<path>` |
| `docker-compose.prod.yml` contains no `build:` | `<pass/fail>` | `<path>` |
| Host runtime does not require `/opt/crm-system/current/.git` | `<pass/fail>` | `<path>` |
| Migrations read from release artifact, not Git checkout | `<pass/fail>` | `<path>` |

## 7. G11 Return Notes

| ACC | Evidence reference | Status for Claude G12 audit |
|---|---|---|
| ACC-CICD-001 | `<fill>` | `<fill>` |
| ACC-CICD-002 | `<fill>` | `<fill>` |
| ACC-CICD-003 | `<fill>` | `<fill>` |
| ACC-CICD-004 | `<fill>` | `<fill>` |
| ACC-CICD-005 | `<fill>` | `<fill>` |
| ACC-CICD-006 | `<fill>` | `<fill>` |
| ACC-CICD-007 | `<fill>` | `<fill>` |
| ACC-CICD-008 | `<fill>` | `<fill>` |

Codex does not self-pass G12.
