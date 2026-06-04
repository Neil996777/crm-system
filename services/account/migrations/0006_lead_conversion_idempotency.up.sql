ALTER TABLE account.accounts
  ADD COLUMN IF NOT EXISTS lead_conversion_idempotency_key text;

CREATE UNIQUE INDEX IF NOT EXISTS accounts_lead_conversion_idempotency_key_idx
  ON account.accounts (lead_conversion_idempotency_key)
  WHERE lead_conversion_idempotency_key IS NOT NULL;
