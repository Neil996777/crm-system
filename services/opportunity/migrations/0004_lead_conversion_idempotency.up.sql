ALTER TABLE opportunity.opportunities
  ADD COLUMN IF NOT EXISTS lead_conversion_idempotency_key text;

CREATE UNIQUE INDEX IF NOT EXISTS opportunities_lead_conversion_idempotency_key_idx
  ON opportunity.opportunities (lead_conversion_idempotency_key)
  WHERE lead_conversion_idempotency_key IS NOT NULL;
