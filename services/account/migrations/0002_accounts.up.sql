CREATE TABLE IF NOT EXISTS account.accounts (
  id text PRIMARY KEY,
  company_name text NOT NULL,
  customer_status text NOT NULL,
  owner_id text NOT NULL,
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT accounts_company_name_required CHECK (company_name <> ''),
  CONSTRAINT accounts_customer_status_required CHECK (customer_status <> ''),
  CONSTRAINT accounts_owner_required CHECK (owner_id <> ''),
  CONSTRAINT accounts_version_positive CHECK (version >= 1)
);

CREATE INDEX IF NOT EXISTS accounts_owner_idx ON account.accounts (owner_id);
CREATE INDEX IF NOT EXISTS accounts_company_search_idx ON account.accounts (company_name);
CREATE INDEX IF NOT EXISTS accounts_customer_status_idx ON account.accounts (customer_status);

CREATE TABLE IF NOT EXISTS account.outbox_events (
  id text PRIMARY KEY,
  event_type text NOT NULL,
  aggregate_id text NOT NULL,
  payload jsonb NOT NULL,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS account_outbox_unpublished_idx
  ON account.outbox_events (occurred_at)
  WHERE published_at IS NULL;

GRANT SELECT, INSERT, UPDATE ON account.accounts TO crm_account_user;
REVOKE DELETE ON account.accounts FROM crm_account_user;
GRANT SELECT, INSERT, UPDATE ON account.outbox_events TO crm_account_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA account TO crm_account_user;
