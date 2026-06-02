CREATE SCHEMA IF NOT EXISTS commercial;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crm_commercial_user') THEN
    CREATE ROLE crm_commercial_user LOGIN PASSWORD 'crm_commercial_dev_password';
  END IF;
END
$$;

GRANT USAGE, CREATE ON SCHEMA commercial TO crm_commercial_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA commercial GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crm_commercial_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA commercial GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO crm_commercial_user;
