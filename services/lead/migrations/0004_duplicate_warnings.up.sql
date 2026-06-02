ALTER TABLE lead.leads
  ADD COLUMN IF NOT EXISTS email text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS phone text NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS lead.duplicate_warning_tokens (
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

CREATE INDEX IF NOT EXISTS lead_duplicate_warning_actor_idx
  ON lead.duplicate_warning_tokens (actor_user_id, target_type, expires_at);

GRANT SELECT, INSERT, UPDATE ON lead.duplicate_warning_tokens TO crm_lead_user;
GRANT SELECT, INSERT, UPDATE ON lead.leads TO crm_lead_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA lead TO crm_lead_user;
