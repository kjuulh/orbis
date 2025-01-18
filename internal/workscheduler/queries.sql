-- name: Ping :one
SELECT 1;

-- name: GetCurrentQueueSize :one
SELECT 
    COUNT(*) current_queue_size
FROM
    work_schedule
WHERE
        worker_id = $1
    AND state <> 'archived';

-- name: InsertQueueItem :exec
INSERT INTO work_schedule
    (
          schedule_id
        , worker_id
        , start_run
        , end_run
        , state
    )
VALUES
    (
          $1
        , $2
        , $3
        , $4
        , 'pending'
    );


-- name: GetNext :one
SELECT 
    *
FROM 
    work_schedule
WHERE
        worker_id = $1
    AND state = 'pending'
ORDER BY updated_at DESC
LIMIT 1;

-- name: StartProcessing :exec
UPDATE work_schedule
SET
    state = 'processing'
WHERE
    schedule_id = $1;

-- name: Archive :exec
UPDATE work_schedule
SET
    state = 'archived'
WHERE
    schedule_id = $1;
   
