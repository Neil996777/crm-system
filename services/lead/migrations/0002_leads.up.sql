CREATE TABLE IF NOT EXISTS lead.leads (
  id text PRIMARY KEY,
  lead_name text NOT NULL DEFAULT '',
  company_name text NOT NULL DEFAULT '',
  source text NOT NULL,
  status text NOT NULL,
  owner_id text NOT NULL DEFAULT '',
  need_summary text NOT NULL DEFAULT '',
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT leads_name_or_company_required CHECK (lead_name <> '' OR company_name <> ''),
  CONSTRAINT leads_source_required CHECK (source <> ''),
  CONSTRAINT leads_status_allowed CHECK (status IN ('Unassigned', 'Pending Qualification', 'Valid', 'Invalid', 'Converted To Opportunity')),
  CONSTRAINT leads_owner_required_after_unassigned CHECK (status = 'Unassigned' OR owner_id <> ''),
  CONSTRAINT leads_version_positive CHECK (version >= 1)
);

CREATE INDEX IF NOT EXISTS leads_owner_idx ON lead.leads (owner_id);
CREATE INDEX IF NOT EXISTS leads_search_idx ON lead.leads (lead_name, company_name);

CREATE TABLE IF NOT EXISTS lead.outbox_events (
  id text PRIMARY KEY,
  event_type text NOT NULL,
  aggregate_id text NOT NULL,
  payload jsonb NOT NULL,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS lead_outbox_unpublished_idx
  ON lead.outbox_events (occurred_at)
  WHERE published_at IS NULL;

GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA lead TO crm_lead_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA lead TO crm_lead_user;
