CREATE TABLE IF NOT EXISTS account.duplicate_warning_tokens (
  token text PRIMARY KEY,
  target_type text NOT NULL,
  normalized_signature text NOT NULL,
  actor_user_id text NOT NULL,
  used_at timestamptz,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT duplicate_warning_target_required CHECK (target_type <> ''),
  CONSTRAINT duplicate_warning_signature_required CHECK (normalized_signature <> ''),
  CONSTRAINT duplicate_warning_actor_required CHECK (actor_user_id <> '')
);

CREATE INDEX IF NOT EXISTS account_duplicate_warning_actor_idx
  ON account.duplicate_warning_tokens (actor_user_id, target_type, expires_at);

GRANT SELECT, INSERT, UPDATE ON account.duplicate_warning_tokens TO crm_account_user;
REVOKE DELETE ON account.duplicate_warning_tokens FROM crm_account_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA account TO crm_account_user;
