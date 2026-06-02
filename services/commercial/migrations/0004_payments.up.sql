CREATE TABLE IF NOT EXISTS commercial.payment_plans (
  id text PRIMARY KEY,
  contract_id text NOT NULL REFERENCES commercial.contracts(id),
  due_amount numeric(14,2) NOT NULL,
  due_date date NOT NULL,
  currency text NOT NULL DEFAULT 'CNY',
  status text NOT NULL,
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT payment_plans_amount_positive CHECK (due_amount > 0),
  CONSTRAINT payment_plans_currency_single CHECK (currency = 'CNY'),
  CONSTRAINT payment_plans_status_allowed CHECK (status IN ('Unpaid', 'PartiallyPaid', 'Paid', 'Overdue')),
  CONSTRAINT payment_plans_version_positive CHECK (version >= 1)
);

CREATE TABLE IF NOT EXISTS commercial.actual_payments (
  id text PRIMARY KEY,
  contract_id text NOT NULL REFERENCES commercial.contracts(id),
  idempotency_key text NOT NULL,
  amount numeric(14,2) NOT NULL,
  payment_date date NOT NULL,
  note text NOT NULL DEFAULT '',
  currency text NOT NULL DEFAULT 'CNY',
  payment_status text NOT NULL,
  remaining_amount numeric(14,2) NOT NULL,
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT actual_payments_amount_positive CHECK (amount > 0),
  CONSTRAINT actual_payments_remaining_nonnegative CHECK (remaining_amount >= 0),
  CONSTRAINT actual_payments_currency_single CHECK (currency = 'CNY'),
  CONSTRAINT actual_payments_status_allowed CHECK (payment_status IN ('PartiallyPaid', 'Paid')),
  CONSTRAINT actual_payments_version_positive CHECK (version >= 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS actual_payments_idempotency_unique
  ON commercial.actual_payments (contract_id, idempotency_key);
CREATE INDEX IF NOT EXISTS payment_plans_contract_idx ON commercial.payment_plans (contract_id);
CREATE INDEX IF NOT EXISTS actual_payments_contract_idx ON commercial.actual_payments (contract_id);

GRANT SELECT, INSERT, UPDATE ON commercial.payment_plans TO crm_commercial_user;
GRANT SELECT, INSERT, UPDATE ON commercial.actual_payments TO crm_commercial_user;
REVOKE DELETE ON commercial.payment_plans FROM crm_commercial_user;
REVOKE DELETE ON commercial.actual_payments FROM crm_commercial_user;
