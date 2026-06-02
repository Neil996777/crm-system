DROP INDEX IF EXISTS lead.leads_conversion_idempotency_key_idx;
ALTER TABLE lead.leads
  DROP CONSTRAINT IF EXISTS leads_converted_links_required,
  DROP CONSTRAINT IF EXISTS leads_invalid_reason_required,
  DROP COLUMN IF EXISTS conversion_idempotency_key,
  DROP COLUMN IF EXISTS converted_opportunity_id,
  DROP COLUMN IF EXISTS converted_account_id,
  DROP COLUMN IF EXISTS invalid_reason;
