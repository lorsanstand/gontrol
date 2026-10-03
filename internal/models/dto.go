package models

type Task struct {
	ID      string            `json:"id"`
	Command string            `json:"command"`
	Args    map[string]string `json:"args"`
}

type TaskResult struct {
	TimeExecMS int    `json:"time_exec_ms"`
	Result     string `json:"result"`
}

type HeartbeatRequest struct {
	AgentID  string `json:"agent_id"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

type HeartbeatResponse struct {
	AgentID string `json:"agent_id"`
	Task    *Task  `json:"task"`
}
