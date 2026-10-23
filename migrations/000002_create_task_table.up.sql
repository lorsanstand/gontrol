CREATE TABLE IF NOT EXISTS Tasks (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,
    status VARCHAR(15) NOT NULL,
    command VARCHAR(10) NOT NULL,
    args jsonb,
    exec_ms INTEGER,
    result TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (agent_id) REFERENCES Agents (id)
    ON DELETE CASCADE ON UPDATE CASCADE
)