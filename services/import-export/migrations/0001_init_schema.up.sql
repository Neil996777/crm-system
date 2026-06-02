CREATE SCHEMA IF NOT EXISTS import_export;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crm_import_export_user') THEN
    CREATE ROLE crm_import_export_user LOGIN PASSWORD 'crm_import_export_dev_password';
  END IF;
END
$$;

GRANT USAGE, CREATE ON SCHEMA import_export TO crm_import_export_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA import_export GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crm_import_export_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA import_export GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO crm_import_export_user;
