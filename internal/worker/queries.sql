-- name: Ping :one
SELECT 1;

-- name: RegisterWorker :exec
INSERT INTO worker_register (worker_id, capacity)
VALUES (
      $1
    , $2
);

-- name: GetWorkers :many
SELECT 
      worker_id
    , capacity
FROM
    worker_register;

-- name: UpdateWorkerHeartbeat :exec
UPDATE worker_register
SET
    heart_beat = now()
WHERE
    worker_id = $1;

