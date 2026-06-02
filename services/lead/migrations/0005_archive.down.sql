ALTER TABLE lead.leads
  DROP COLUMN IF EXISTS archive_reason,
  DROP COLUMN IF EXISTS archived_by,
  DROP COLUMN IF EXISTS archived_at;
