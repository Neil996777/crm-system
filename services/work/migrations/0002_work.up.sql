CREATE TABLE IF NOT EXISTS work.activities (
  id text PRIMARY KEY,
  related_type text NOT NULL,
  related_id text NOT NULL,
  activity_type text NOT NULL,
  content text NOT NULL,
  actor_id text NOT NULL,
  owner_id text NOT NULL,
  occurred_at timestamptz NOT NULL,
  version integer NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS work.notes (
  id text PRIMARY KEY,
  related_type text NOT NULL,
  related_id text NOT NULL,
  content text NOT NULL,
  actor_id text NOT NULL,
  owner_id text NOT NULL,
  occurred_at timestamptz NOT NULL,
  version integer NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS work.tasks (
  id text PRIMARY KEY,
  related_type text NOT NULL,
  related_id text NOT NULL,
  title text NOT NULL,
  due_date date NOT NULL,
  status text NOT NULL,
  actor_id text NOT NULL,
  owner_id text NOT NULL,
  completed_at timestamptz,
  cancelled_at timestamptz,
  cancellation_reason text,
  version integer NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS activities_related_idx ON work.activities (related_type, related_id);
CREATE INDEX IF NOT EXISTS notes_related_idx ON work.notes (related_type, related_id);
CREATE INDEX IF NOT EXISTS tasks_related_idx ON work.tasks (related_type, related_id);
CREATE INDEX IF NOT EXISTS tasks_owner_status_idx ON work.tasks (owner_id, status, due_date);

CREATE TABLE IF NOT EXISTS work.outbox_events (
  id text PRIMARY KEY,
  event_type text NOT NULL,
  aggregate_id text NOT NULL,
  payload jsonb NOT NULL,
  occurred_at timestamptz NOT NULL,
  published_at timestamptz
);
