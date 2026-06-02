CREATE SCHEMA IF NOT EXISTS reporting;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crm_reporting_user') THEN
    CREATE ROLE crm_reporting_user LOGIN PASSWORD 'crm_reporting_dev_password';
  END IF;
END
$$;

GRANT USAGE, CREATE ON SCHEMA reporting TO crm_reporting_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA reporting GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crm_reporting_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA reporting GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO crm_reporting_user;
