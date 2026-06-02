CREATE TABLE IF NOT EXISTS commercial.contracts (
  id text PRIMARY KEY,
  quote_id text NOT NULL REFERENCES commercial.quotes(id),
  opportunity_id text NOT NULL,
  customer_id text NOT NULL,
  amount numeric(14,2) NOT NULL,
  status text NOT NULL,
  contract_note text NOT NULL,
  expected_signed_date date NOT NULL,
  signed_effective_date date,
  amount_difference_reason text,
  owner_id text NOT NULL,
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT contracts_quote_required CHECK (quote_id <> ''),
  CONSTRAINT contracts_opportunity_required CHECK (opportunity_id <> ''),
  CONSTRAINT contracts_customer_required CHECK (customer_id <> ''),
  CONSTRAINT contracts_amount_positive CHECK (amount > 0),
  CONSTRAINT contracts_status_allowed CHECK (status IN ('Pending Signature', 'Signed', 'Active', 'Completed', 'Terminated')),
  CONSTRAINT contracts_note_required CHECK (contract_note <> ''),
  CONSTRAINT contracts_signed_date_required CHECK (
    status IN ('Pending Signature', 'Terminated') OR signed_effective_date IS NOT NULL
  ),
  CONSTRAINT contracts_owner_required CHECK (owner_id <> ''),
  CONSTRAINT contracts_version_positive CHECK (version >= 1),
  CONSTRAINT contracts_amount_difference_reason_required CHECK (
    amount_difference_reason IS NULL OR amount_difference_reason <> ''
  )
);

CREATE UNIQUE INDEX IF NOT EXISTS contracts_quote_unique ON commercial.contracts (quote_id);
CREATE INDEX IF NOT EXISTS contracts_owner_idx ON commercial.contracts (owner_id);
CREATE INDEX IF NOT EXISTS contracts_status_idx ON commercial.contracts (status);

GRANT SELECT, INSERT, UPDATE ON commercial.contracts TO crm_commercial_user;
REVOKE DELETE ON commercial.contracts FROM crm_commercial_user;
