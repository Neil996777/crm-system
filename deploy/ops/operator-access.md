# Operator Access

Production deployment target: `srv-volcengine-sh-01` (`118.196.44.193`).

CRM operations must use named least-privilege accounts after initial provisioning:

| Account | Purpose | Required boundary |
|---|---|---|
| `crm-deploy` | Deploy CRM artifacts under `/opt/crm-system/current`, run Docker Compose for CRM, reload Nginx after a validated config test | Password login disabled; SSH key only; sudo limited to CRM deployment commands and `systemctl reload nginx` after `nginx -t` |
| `crm-ops` | Read logs, run health checks, inspect service status, execute approved backup/restore procedures | Password login disabled; SSH key only; sudo limited to read-only diagnostics plus approved backup commands |
| `root` | Initial provisioning and emergency recovery only | Not used for routine deployment; key ownership and use recorded in infrastructure registers |

Required filesystem ownership:

- `/opt/crm-system/current`: writable by `crm-deploy`, readable by `crm-ops`.
- `/opt/crm-system/volumes/postgres`: writable only by PostgreSQL container runtime/root.
- `/opt/crm-system/logs`: readable by `crm-ops`; service log directories created with least required write access.
- `/opt/crm-system/backups`: restricted to `crm-ops`/root and backup job identity.

SSH requirements:

- Disable password authentication for CRM operator users.
- Keep operator keys out of the repository.
- Record key fingerprints in the infrastructure SSH access register.
- Replace current root-only access with named operators before production closure.

G12 second rework operator key evidence (2026-06-03):

| Account | Local private key path | Public key fingerprint |
|---|---|---|
| `crm-deploy` | `/Users/neil/practice/software/.secrets/ssh-keys/crm-deploy-volcengine-sh-20260603` | `SHA256:ZGLqXBHGgqy29ZUMFysRjaw579Z3yx1980pIFWBb/b4` |
| `crm-ops` | `/Users/neil/practice/software/.secrets/ssh-keys/crm-ops-volcengine-sh-20260603` | `SHA256:PHl9ZXjKKPzI5oiWrll9Jj60X04+5S7/TMpV1q3AYQA` |

Evidence: `docs/release/evidence/operator-access-transcript-2026-06-03.txt`.
