CREATE TABLE worker_register (
      worker_id UUID PRIMARY KEY NOT NULL
    , capacity INTEGER NOT NULL
    , heart_beat TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE model_schedules (
      model_name TEXT PRIMARY KEY NOT NULL
    , last_run TIMESTAMPTZ
);

CREATE TABLE work_schedule (
      schedule_id UUID PRIMARY KEY NOT NULL
    , worker_id UUID NOT NULL
    , start_run TIMESTAMPTZ NOT NULL
    , end_run TIMESTAMPTZ NOT NULL
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , state TEXT NOT NULL
);

CREATE INDEX idx_work_schedule_worker ON work_schedule (worker_id);
CREATE INDEX idx_work_schedule_worker_updated ON work_schedule (worker_id, updated_at DESC);
