ALTER TABLE account.accounts
  ADD COLUMN IF NOT EXISTS archived_at timestamptz,
  ADD COLUMN IF NOT EXISTS archived_by text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS archive_reason text NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS accounts_archived_idx ON account.accounts (archived_at);

GRANT SELECT, INSERT, UPDATE ON account.accounts TO crm_account_user;
