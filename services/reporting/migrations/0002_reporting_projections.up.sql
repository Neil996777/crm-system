CREATE TABLE IF NOT EXISTS reporting.record_projections (
  id bigserial PRIMARY KEY,
  source_service text NOT NULL,
  record_type text NOT NULL,
  record_id text NOT NULL,
  owner_id text NOT NULL,
  team_id text NOT NULL DEFAULT 'single-team',
  status text,
  stage text,
  amount numeric(14,2) NOT NULL DEFAULT 0,
  archived_at timestamptz,
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (record_type, record_id)
);

CREATE INDEX IF NOT EXISTS record_projections_team_type_idx
  ON reporting.record_projections (team_id, record_type)
  WHERE archived_at IS NULL;

GRANT SELECT, INSERT, UPDATE ON reporting.record_projections TO crm_reporting_user;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE reporting.record_projections_id_seq TO crm_reporting_user;
