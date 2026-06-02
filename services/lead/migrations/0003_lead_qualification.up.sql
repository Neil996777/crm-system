ALTER TABLE lead.leads
  ADD COLUMN IF NOT EXISTS invalid_reason text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS converted_account_id text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS converted_opportunity_id text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS conversion_idempotency_key text NOT NULL DEFAULT '';

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'leads_invalid_reason_required'
      AND conrelid = 'lead.leads'::regclass
  ) THEN
    ALTER TABLE lead.leads
      ADD CONSTRAINT leads_invalid_reason_required
      CHECK (status <> 'Invalid' OR invalid_reason <> '');
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'leads_converted_links_required'
      AND conrelid = 'lead.leads'::regclass
  ) THEN
    ALTER TABLE lead.leads
      ADD CONSTRAINT leads_converted_links_required
      CHECK (
        status <> 'Converted To Opportunity'
        OR (converted_account_id <> '' AND converted_opportunity_id <> '' AND conversion_idempotency_key <> '')
      );
  END IF;
END
$$;

CREATE UNIQUE INDEX IF NOT EXISTS leads_conversion_idempotency_key_idx
  ON lead.leads (conversion_idempotency_key)
  WHERE conversion_idempotency_key <> '';

GRANT SELECT, INSERT, UPDATE ON lead.leads TO crm_lead_user;
REVOKE DELETE ON lead.leads FROM crm_lead_user;
