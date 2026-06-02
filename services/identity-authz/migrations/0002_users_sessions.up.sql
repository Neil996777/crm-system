CREATE TABLE IF NOT EXISTS identity_authz.roles (
  name text PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT roles_name_allowed CHECK (name IN ('Administrator', 'Sales Manager', 'Sales'))
);

INSERT INTO identity_authz.roles (name)
VALUES ('Administrator'), ('Sales Manager'), ('Sales')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS identity_authz.users (
  id text PRIMARY KEY,
  email text NOT NULL UNIQUE,
  display_name text NOT NULL,
  password_hash text NOT NULL,
  role_name text NOT NULL REFERENCES identity_authz.roles(name),
  status text NOT NULL DEFAULT 'Active',
  authz_version integer NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT users_status_allowed CHECK (status IN ('Active', 'Disabled')),
  CONSTRAINT users_authz_version_positive CHECK (authz_version >= 1)
);

CREATE TABLE IF NOT EXISTS identity_authz.sessions (
  id text PRIMARY KEY,
  user_id text NOT NULL REFERENCES identity_authz.users(id),
  authz_version_at_issue integer NOT NULL,
  expires_at timestamptz NOT NULL,
  idle_expires_at timestamptz NOT NULL,
  revoked_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  last_seen_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS sessions_user_active_idx
  ON identity_authz.sessions (user_id)
  WHERE revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS identity_authz.outbox_events (
  id text PRIMARY KEY,
  event_type text NOT NULL,
  aggregate_type text NOT NULL,
  aggregate_id text NOT NULL,
  payload jsonb NOT NULL,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS outbox_events_unpublished_idx
  ON identity_authz.outbox_events (occurred_at)
  WHERE published_at IS NULL;

INSERT INTO identity_authz.users (
  id, email, display_name, password_hash, role_name, status, authz_version
)
VALUES (
  'usr_seed_admin',
  'admin@example.com',
  'Seed Administrator',
  '$2a$10$oVcOJvmFKSwki4bGtHKA7eBn8rXFOK6YUPvxObiKFzcRGwRLd.RsW',
  'Administrator',
  'Active',
  1
)
ON CONFLICT (id) DO NOTHING;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA identity_authz TO crm_identity_authz_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA identity_authz TO crm_identity_authz_user;
