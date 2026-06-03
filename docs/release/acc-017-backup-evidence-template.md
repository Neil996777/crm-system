# ACC-017 Backup And Restore Evidence

| Field | Evidence |
|---|---|
| Runtime host | `srv-volcengine-sh-01` / `118.196.44.193` |
| Off-server target | `srv-aliyun-bj-01` / `47.95.119.211` |
| Local backup path | `/opt/crm-system/backups/postgres` |
| Backup script | `deploy/backup/backup.sh` |
| Offsite copy script | `deploy/backup/offsite-copy.sh` |
| Schedule | `deploy/ops/crm-backup.timer` daily at `02:00`; enabled and active on `srv-volcengine-sh-01`, next run `2026-06-04 02:00:00 CST` |
| Encryption | OpenSSL AES-256-CBC with PBKDF2; passphrase file outside repository |
| Retention | Local encrypted backups older than 7 days pruned |
| Latest encrypted backup | `/opt/crm-system/backups/postgres/crm-postgres-20260603T104620Z.sql.gz.enc` (`9680` bytes, mode `600`) |
| Latest checksum | `f7381eaa9d246126cac93b304a147ea01721de404f52de76b340e3cfa9ba9d2a` |
| Off-server copy result | Passed. `srv-aliyun-bj-01` has the encrypted backup and `.sha256`; `sha256sum -c` returned `OK` on 2026-06-03. |
| Restore rehearsal | Passed. Run `20260603T104837Z` restored the off-server encrypted backup into a controlled PostgreSQL target and verified roles, database, service schemas, users table, audit-history tables, and service permission roles. |
| Operator | Codex acting as `infrastructure-ops` |
| Result | Passed for TASK-040 / ACC-017 backup and restore release evidence |

Evidence files:

- `docs/release/evidence/backup-restore-rehearsal-2026-06-03.json`
- `docs/release/evidence/backup-restore-transcript-2026-06-03.txt`

## Restore rehearsal

The restore rehearsal must verify:

- users/roles survive restore
- CRM records survive restore
- history/logs survive restore
- service DB permissions remain scoped
- checksum verification passes before decryption
- decrypted SQL is removed after the controlled restore

2026-06-03 rehearsal result:

- Source backup: off-server encrypted copy from `srv-aliyun-bj-01`.
- Restore database: `crm_system`.
- Checksum: passed before decryption.
- Roles: `10`; service permission roles: `9`.
- Service schemas: `9`.
- Identity users table: present.
- Audit-history tables: present.
- Decrypted SQL retention: removed after the controlled restore; the rehearsal
  directory retains only the encrypted backup, checksum, and `restore-result.json`.

G12 rework captured the real `sha256sum -c` output, OpenSSL decrypt + `gzip -t`
integrity check, restore-result file, `psql \du`, schema list, and
`audit_history` table listing in
`docs/release/evidence/backup-restore-transcript-2026-06-03.txt`.
