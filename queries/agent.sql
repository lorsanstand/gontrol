-- name: GetAgentByID :one
SELECT id, os, hostname, name, description, created_at FROM Agents
    WHERE id=$1;

-- name: GetAllAgents :many
SELECT id, os, hostname, name, description, created_at FROM Agents;

-- name: CreateAgent :one
INSERT INTO Agents (id, name, os, hostname, description)
    VALUES ($1, $2, $3, $4, $5)
        RETURNING id, os, hostname, name, description, created_at;

-- name: DeleteAgentByID :execresult
DELETE FROM Agents
    WHERE id=$1;

-- name: EditDescAgentByID :execresult
UPDATE Agents
    SET description = $2 WHERE id=$1;

-- name: EditNameAgentByID :execresult
    UPDATE Agents
        SET name = $2 WHERE id=$1;