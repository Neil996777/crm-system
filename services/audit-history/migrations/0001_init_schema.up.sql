CREATE SCHEMA IF NOT EXISTS audit_history;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crm_audit_history_user') THEN
    CREATE ROLE crm_audit_history_user LOGIN PASSWORD 'crm_audit_history_dev_password';
  END IF;
END
$$;

GRANT USAGE, CREATE ON SCHEMA audit_history TO crm_audit_history_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA audit_history GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crm_audit_history_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA audit_history GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO crm_audit_history_user;
