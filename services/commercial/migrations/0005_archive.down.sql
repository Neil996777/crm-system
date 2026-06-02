ALTER TABLE commercial.payment_plans
  DROP COLUMN IF EXISTS archive_reason,
  DROP COLUMN IF EXISTS archived_by,
  DROP COLUMN IF EXISTS archived_at;

ALTER TABLE commercial.contracts
  DROP COLUMN IF EXISTS archive_reason,
  DROP COLUMN IF EXISTS archived_by,
  DROP COLUMN IF EXISTS archived_at;
