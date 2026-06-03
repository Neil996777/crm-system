CREATE TABLE IF NOT EXISTS import_export.import_runs (
  run_id text PRIMARY KEY,
  object_type text NOT NULL,
  filename text NOT NULL,
  status text NOT NULL,
  actor_id text NOT NULL,
  actor_role text NOT NULL,
  team_id text NOT NULL DEFAULT 'single-team',
  total_rows integer NOT NULL DEFAULT 0,
  success_count integer NOT NULL DEFAULT 0,
  failure_count integer NOT NULL DEFAULT 0,
  operation_log_status text NOT NULL DEFAULT 'not_configured',
  cleanup_status text NOT NULL DEFAULT 'pending',
  retained_until timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  completed_at timestamptz
);

CREATE TABLE IF NOT EXISTS import_export.import_row_results (
  id bigserial PRIMARY KEY,
  run_id text NOT NULL REFERENCES import_export.import_runs(run_id) ON DELETE CASCADE,
  row_number integer NOT NULL,
  success boolean NOT NULL,
  field text,
  code text,
  safe_message text,
  target_record_id text,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS import_row_results_run_idx
  ON import_export.import_row_results (run_id, row_number);

GRANT SELECT, INSERT, UPDATE ON import_export.import_runs TO crm_import_export_user;
GRANT SELECT, INSERT, UPDATE ON import_export.import_row_results TO crm_import_export_user;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE import_export.import_row_results_id_seq TO crm_import_export_user;
