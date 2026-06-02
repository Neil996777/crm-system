ALTER TABLE commercial.contracts
  ADD COLUMN IF NOT EXISTS archived_at timestamptz,
  ADD COLUMN IF NOT EXISTS archived_by text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS archive_reason text NOT NULL DEFAULT '';

ALTER TABLE commercial.payment_plans
  ADD COLUMN IF NOT EXISTS archived_at timestamptz,
  ADD COLUMN IF NOT EXISTS archived_by text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS archive_reason text NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS contracts_archived_idx ON commercial.contracts (archived_at);
CREATE INDEX IF NOT EXISTS payment_plans_archived_idx ON commercial.payment_plans (archived_at);

GRANT SELECT, INSERT, UPDATE ON commercial.contracts TO crm_commercial_user;
GRANT SELECT, INSERT, UPDATE ON commercial.payment_plans TO crm_commercial_user;
