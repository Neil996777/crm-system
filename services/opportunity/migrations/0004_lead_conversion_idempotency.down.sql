DROP INDEX IF EXISTS opportunity.opportunities_lead_conversion_idempotency_key_idx;

ALTER TABLE opportunity.opportunities
  DROP COLUMN IF EXISTS lead_conversion_idempotency_key;
