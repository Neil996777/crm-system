CREATE TABLE IF NOT EXISTS commercial.quotes (
  id text PRIMARY KEY,
  opportunity_id text NOT NULL,
  customer_id text NOT NULL,
  amount numeric(14,2) NOT NULL,
  status text NOT NULL,
  validity_end date NOT NULL,
  owner_id text NOT NULL,
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT quotes_opportunity_required CHECK (opportunity_id <> ''),
  CONSTRAINT quotes_customer_required CHECK (customer_id <> ''),
  CONSTRAINT quotes_amount_positive CHECK (amount > 0),
  CONSTRAINT quotes_owner_required CHECK (owner_id <> ''),
  CONSTRAINT quotes_status_allowed CHECK (status IN ('Draft', 'Sent', 'Accepted', 'Rejected', 'Expired')),
  CONSTRAINT quotes_version_positive CHECK (version >= 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS quotes_opportunity_unique ON commercial.quotes (opportunity_id);
CREATE INDEX IF NOT EXISTS quotes_owner_idx ON commercial.quotes (owner_id);
CREATE INDEX IF NOT EXISTS quotes_status_idx ON commercial.quotes (status);

CREATE TABLE IF NOT EXISTS commercial.outbox_events (
  id text PRIMARY KEY,
  event_type text NOT NULL,
  aggregate_id text NOT NULL,
  payload jsonb NOT NULL,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS commercial_outbox_unpublished_idx
  ON commercial.outbox_events (occurred_at)
  WHERE published_at IS NULL;

GRANT SELECT, INSERT, UPDATE ON commercial.quotes TO crm_commercial_user;
REVOKE DELETE ON commercial.quotes FROM crm_commercial_user;
GRANT SELECT, INSERT, UPDATE ON commercial.outbox_events TO crm_commercial_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA commercial TO crm_commercial_user;
