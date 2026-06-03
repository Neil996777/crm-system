CREATE TABLE IF NOT EXISTS import_export.export_runs (
  run_id text PRIMARY KEY,
  object_type text NOT NULL,
  filename text NOT NULL,
  status text NOT NULL,
  actor_id text NOT NULL,
  actor_role text NOT NULL,
  team_id text NOT NULL DEFAULT 'single-team',
  include_archived boolean NOT NULL DEFAULT false,
  exported_count integer NOT NULL DEFAULT 0,
  operation_log_status text NOT NULL DEFAULT 'not_configured',
  cleanup_status text NOT NULL DEFAULT 'pending',
  retained_until timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  completed_at timestamptz
);

GRANT SELECT, INSERT, UPDATE ON import_export.export_runs TO crm_import_export_user;
