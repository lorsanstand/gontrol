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
    VALUES (id=$1, agent_id=$2, status=$3, command=$4, args=$5, exec_ms=$6, result=$7)
        RETURNING id, agent_id, status, command, args, exec_ms, result, created_at;

-- name: EditStatusTaskByID :one
UPDATE Tasks
    SET status=$2 WHERE id=$1
        RETURNING id, status;

-- name: EditResultTaskByID :exec
UPDATE Tasks
    SET status=$2, result=$3, exec_ms=$4 WHERE id=$1;

-- name: DeleteTaskByID :exec
DELETE FROM Tasks
    WHERE id=$1;