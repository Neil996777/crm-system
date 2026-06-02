CREATE TABLE IF NOT EXISTS audit_history.events (
  sequence_id bigserial PRIMARY KEY,
  event_uid text NOT NULL UNIQUE,
  event_id text NOT NULL,
  event_version integer NOT NULL,
  producer_service text NOT NULL,
  surfaces text[] NOT NULL,
  actor_user_id text NOT NULL,
  actor_role text NOT NULL,
  actor_display text NOT NULL,
  action text NOT NULL,
  resource_type text NOT NULL,
  resource_id text NOT NULL,
  parent_resource_type text NOT NULL DEFAULT '',
  parent_resource_id text NOT NULL DEFAULT '',
  result text NOT NULL,
  reason_code text NOT NULL DEFAULT '',
  before_summary jsonb NOT NULL DEFAULT '{}'::jsonb,
  after_summary jsonb NOT NULL DEFAULT '{}'::jsonb,
  diff_classification text NOT NULL DEFAULT '',
  scope_summary text NOT NULL DEFAULT '',
  safe_summary text NOT NULL,
  correlation_id text NOT NULL DEFAULT '',
  causation_id text NOT NULL DEFAULT '',
  acceptance_ids text[] NOT NULL,
  occurred_at timestamptz NOT NULL,
  prev_hash text NOT NULL DEFAULT '',
  event_hash text NOT NULL
);

CREATE INDEX IF NOT EXISTS events_record_history_idx
  ON audit_history.events (resource_type, resource_id, occurred_at)
  WHERE 'record_history' = ANY(surfaces);

CREATE INDEX IF NOT EXISTS events_operation_log_idx
  ON audit_history.events (occurred_at)
  WHERE 'operation_log' = ANY(surfaces);

GRANT SELECT, INSERT ON audit_history.events TO crm_audit_history_user;
GRANT USAGE, SELECT ON SEQUENCE audit_history.events_sequence_id_seq TO crm_audit_history_user;
REVOKE UPDATE, DELETE ON audit_history.events FROM crm_audit_history_user;
