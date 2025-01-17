CREATE TABLE worker_register (
    worker_id UUID PRIMARY KEY NOT NULL
    , heart_beat TIMESTAMPTZ NOT NULL DEFAULT now()
);
