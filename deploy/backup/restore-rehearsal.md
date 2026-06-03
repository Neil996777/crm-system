# Restore Rehearsal Procedure

TASK-040 restore rehearsal must use an encrypted backup copied off-server to
`srv-aliyun-bj-01`, restore it into a controlled PostgreSQL target, and record
the result in `docs/release/acc-017-backup-evidence-template.md`.

Required record:

- operator
- timestamp
- source encrypted backup file
- checksum file and `sha256sum -c` result
- decryption command and access subject
- restore target container/host
- verification commands and result

Controlled restore outline:

1. Copy the selected encrypted backup and `.sha256` file from the off-server
   target into a restricted local restore directory.
2. Verify checksum with `sha256sum -c`.
3. Decrypt with `openssl enc -d -aes-256-cbc -pbkdf2` using the backup
   passphrase file stored outside the repository.
4. Start a temporary PostgreSQL container that is not exposed publicly.
5. Restore the decrypted `pg_dumpall` SQL into the temporary container.
6. Verify users/roles with `\du` or catalog queries.
7. Verify CRM records exist in service schemas, including at least one core
   record row where present.
8. Verify history/logs survive restore by checking `audit_history` tables.
9. Verify service DB permissions remain scoped to service roles, not a shared
   cross-service owner.
10. Destroy the temporary restore target and decrypted SQL after evidence is
    recorded.

Restricted data from the rehearsal must not be copied into product UI, public
logs, committed repository files, or chat output. Evidence records command
names, timestamps, checksums, and pass/fail results only.
