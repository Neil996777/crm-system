DROP INDEX IF EXISTS audit_history.events_retention_idx;
ALTER TABLE audit_history.events
  DROP COLUMN IF EXISTS retain_until,
  DROP COLUMN IF EXISTS retention_policy;
