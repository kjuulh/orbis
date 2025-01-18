-- name: Ping :one
SELECT 1;

-- name: GetLast :one
SELECT last_run
FROM
    model_schedules
WHERE
    model_name = $1
LIMIT 1;

-- name: UpsertModel :exec
INSERT INTO model_schedules (model_name, last_run)
VALUES ($1, $2)
ON CONFLICT (model_name)
DO UPDATE SET
last_run = excluded.last_run;

