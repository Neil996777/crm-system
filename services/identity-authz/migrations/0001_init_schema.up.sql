CREATE SCHEMA IF NOT EXISTS identity_authz;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crm_identity_authz_user') THEN
    CREATE ROLE crm_identity_authz_user LOGIN PASSWORD 'crm_identity_authz_dev_password';
  END IF;
END
$$;

GRANT USAGE, CREATE ON SCHEMA identity_authz TO crm_identity_authz_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA identity_authz GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crm_identity_authz_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA identity_authz GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO crm_identity_authz_user;
