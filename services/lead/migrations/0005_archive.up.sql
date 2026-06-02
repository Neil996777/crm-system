ALTER TABLE lead.leads
  ADD COLUMN IF NOT EXISTS archived_at timestamptz,
  ADD COLUMN IF NOT EXISTS archived_by text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS archive_reason text NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS leads_archived_idx ON lead.leads (archived_at);

GRANT SELECT, INSERT, UPDATE ON lead.leads TO crm_lead_user;
