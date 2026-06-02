CREATE SCHEMA IF NOT EXISTS opportunity;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crm_opportunity_user') THEN
    CREATE ROLE crm_opportunity_user LOGIN PASSWORD 'crm_opportunity_dev_password';
  END IF;
END
$$;

GRANT USAGE, CREATE ON SCHEMA opportunity TO crm_opportunity_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA opportunity GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crm_opportunity_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA opportunity GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO crm_opportunity_user;
