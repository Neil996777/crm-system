CREATE TABLE IF NOT EXISTS opportunity.opportunities (
  id text PRIMARY KEY,
  customer_id text NOT NULL,
  owner_id text NOT NULL,
  stage text NOT NULL,
  expected_amount numeric(14,2) NOT NULL,
  expected_close_date date NOT NULL,
  title text NOT NULL DEFAULT '',
  close_date date,
  won_contract_id text,
  lost_reason_code text,
  lost_reason_detail text,
  closed_at timestamptz,
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT opportunities_customer_required CHECK (customer_id <> ''),
  CONSTRAINT opportunities_owner_required CHECK (owner_id <> ''),
  CONSTRAINT opportunities_amount_positive CHECK (expected_amount > 0),
  CONSTRAINT opportunities_version_positive CHECK (version >= 1),
  CONSTRAINT opportunities_stage_allowed CHECK (stage IN (
    'New Opportunity',
    'Needs Confirmed',
    'Quote',
    'Contract Negotiation',
    'Won',
    'Lost'
  )),
  CONSTRAINT opportunities_won_contract_required CHECK (stage <> 'Won' OR won_contract_id IS NOT NULL),
  CONSTRAINT opportunities_lost_reason_required CHECK (stage <> 'Lost' OR lost_reason_code IS NOT NULL),
  CONSTRAINT opportunities_close_date_required CHECK (stage NOT IN ('Won', 'Lost') OR close_date IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS opportunities_owner_idx ON opportunity.opportunities (owner_id);
CREATE INDEX IF NOT EXISTS opportunities_customer_idx ON opportunity.opportunities (customer_id);
CREATE INDEX IF NOT EXISTS opportunities_stage_idx ON opportunity.opportunities (stage);

CREATE TABLE IF NOT EXISTS opportunity.outbox_events (
  id text PRIMARY KEY,
  event_type text NOT NULL,
  aggregate_id text NOT NULL,
  payload jsonb NOT NULL,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS opportunity_outbox_unpublished_idx
  ON opportunity.outbox_events (occurred_at)
  WHERE published_at IS NULL;

GRANT SELECT, INSERT, UPDATE ON opportunity.opportunities TO crm_opportunity_user;
REVOKE DELETE ON opportunity.opportunities FROM crm_opportunity_user;
GRANT SELECT, INSERT, UPDATE ON opportunity.outbox_events TO crm_opportunity_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA opportunity TO crm_opportunity_user;
