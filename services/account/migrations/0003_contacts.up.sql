CREATE TABLE IF NOT EXISTS account.contacts (
  id text PRIMARY KEY,
  account_id text NOT NULL REFERENCES account.accounts(id),
  contact_name text NOT NULL,
  email text NOT NULL DEFAULT '',
  phone text NOT NULL DEFAULT '',
  role_note text NOT NULL DEFAULT '',
  version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT contacts_account_required CHECK (account_id <> ''),
  CONSTRAINT contacts_name_required CHECK (contact_name <> ''),
  CONSTRAINT contacts_method_or_note_required CHECK (email <> '' OR phone <> '' OR role_note <> ''),
  CONSTRAINT contacts_version_positive CHECK (version >= 1)
);

CREATE INDEX IF NOT EXISTS contacts_account_idx ON account.contacts (account_id);

GRANT SELECT, INSERT, UPDATE ON account.contacts TO crm_account_user;
REVOKE DELETE ON account.contacts FROM crm_account_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA account TO crm_account_user;
