-- name: GetTaskByID :one
SELECT id, agent_id, status, command, args, exec_ms, result, created_at FROM Tasks
    WHERE id=$1;

-- name: GetTasksByAgentID :many
SELECT id, agent_id, status, command, args, exec_ms, result, created_at FROM Tasks
    WHERE agent_id=$1;

-- name: NextTaskByAgentID :one
SELECT id, agent_id, status, command, args, exec_ms, result, created_at FROM Tasks
    WHERE agent_id=$1 AND status='pending'
        ORDER BY created_at ASC LIMIT 1;

-- name: CreateTask :one
INSERT INTO Tasks (id, agent_id, status, command, args, exec_ms, result)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, agent_id, status, command, args, exec_ms, result, created_at;

-- name: EditStatusTaskByID :execresult
UPDATE Tasks
    SET status=$2 WHERE id=$1;

-- name: EditResultTaskByID :execresult
UPDATE Tasks
    SET status=$2, result=$3, exec_ms=$4 WHERE id=$1;

-- name: DeleteTaskByID :execresult
DELETE FROM Tasks
    WHERE id=$1;