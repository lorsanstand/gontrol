-- name: GetCommandsByAgentID :one
SELECT commands FROM Commands
    WHERE agent_id=$1;

-- name: SetCommandsByAgentID :exec
INSERT INTO Commands (agent_id, commands)
    VALUES ($1, $2)
    ON CONFLICT (agent_id)
    DO UPDATE SET commands = EXCLUDED.commands;