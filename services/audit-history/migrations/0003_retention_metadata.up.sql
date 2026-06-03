ALTER TABLE audit_history.events
  ADD COLUMN IF NOT EXISTS retention_policy text NOT NULL DEFAULT 'business_record_or_operation_log_min_7y',
  ADD COLUMN IF NOT EXISTS retain_until timestamptz NOT NULL DEFAULT (now() + interval '7 years');

CREATE INDEX IF NOT EXISTS events_retention_idx
  ON audit_history.events (retain_until);

GRANT SELECT, INSERT ON audit_history.events TO crm_audit_history_user;
REVOKE UPDATE, DELETE ON audit_history.events FROM crm_audit_history_user;
