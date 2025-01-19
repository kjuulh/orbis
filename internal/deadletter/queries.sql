-- name: Ping :one
SELECT 1;

-- name: InsertDeadLetter :exec
INSERT INTO dead_letter
    (
        schedule_id
    )
VALUES
    (
        $1
    );
