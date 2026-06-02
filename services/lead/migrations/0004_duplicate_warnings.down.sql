DROP TABLE IF EXISTS lead.duplicate_warning_tokens;

ALTER TABLE lead.leads
  DROP COLUMN IF EXISTS email,
  DROP COLUMN IF EXISTS phone;
