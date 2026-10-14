-- name: GetAgentByID :one
SELECT id, os, hostname, name, description, created_at FROM Agents
    WHERE id=$1;

-- name: GetAllAgents :many
SELECT id, os, hostname, name, description, created_at FROM Agents;

-- name: CreateAgent :execresult
INSERT INTO Agents (id, name, os, hostname, description)
    VALUES ($1, $2, $3, $4, $5)
        RETURNING id, os, hostname, name, description, created_at;

-- name: DeleteAgentByID :exec
DELETE FROM Agents
    WHERE id=$1;

-- name: EditDescAgentByID :execresult
UPDATE agents
    SET description = $2 WHERE id=$1
        RETURNING id, os, hostname, name, description, created_at;