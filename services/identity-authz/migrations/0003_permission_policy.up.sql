ALTER TABLE identity_authz.users
  ADD COLUMN IF NOT EXISTS team_id text NOT NULL DEFAULT 'single-team';

CREATE TABLE IF NOT EXISTS identity_authz.permission_policy (
  action text PRIMARY KEY,
  description text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO identity_authz.permission_policy (action, description)
VALUES
  ('lead.create', 'Sales may create leads; managers and administrators operate in broader scope.'),
  ('lead.update', 'Record update authorization follows administrator all, manager team, sales owned scope.'),
  ('opportunity.update', 'Opportunity update authorization follows administrator all, manager team, sales owned scope.'),
  ('operation_log.read', 'Administrator-only global operation log read.'),
  ('user.create', 'Administrator-only user creation.'),
  ('user.change_role', 'Administrator-only role change with last-admin guard.'),
  ('user.change_status', 'Administrator-only status change with last-admin guard.'),
  ('*.hard_delete', 'Hard delete is forbidden for core CRM records.')
ON CONFLICT (action) DO NOTHING;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA identity_authz TO crm_identity_authz_user;
