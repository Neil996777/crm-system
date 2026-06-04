DROP INDEX IF EXISTS account.accounts_lead_conversion_idempotency_key_idx;

ALTER TABLE account.accounts
  DROP COLUMN IF EXISTS lead_conversion_idempotency_key;
