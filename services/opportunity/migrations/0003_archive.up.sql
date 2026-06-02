ALTER TABLE opportunity.opportunities
  ADD COLUMN IF NOT EXISTS archived_at timestamptz,
  ADD COLUMN IF NOT EXISTS archived_by text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS archive_reason text NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS opportunities_archived_idx ON opportunity.opportunities (archived_at);

GRANT SELECT, INSERT, UPDATE ON opportunity.opportunities TO crm_opportunity_user;
