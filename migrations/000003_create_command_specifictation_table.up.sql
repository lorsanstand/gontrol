CREATE TABLE IF NOT EXISTS Commands (
    agent_id UUID PRIMARY KEY,
    commands jsonb NOT NULL,

    FOREIGN KEY (agent_id) REFERENCES Agents (id)
    ON DELETE CASCADE ON UPDATE CASCADE
)