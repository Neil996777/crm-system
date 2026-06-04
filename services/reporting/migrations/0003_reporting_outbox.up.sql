CREATE TABLE IF NOT EXISTS reporting.outbox_events (
  id text PRIMARY KEY,
  event_type text NOT NULL,
  aggregate_type text NOT NULL,
  aggregate_id text NOT NULL,
  payload jsonb NOT NULL,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_reporting_outbox_unpublished
  ON reporting.outbox_events (occurred_at, id)
  WHERE published_at IS NULL;
