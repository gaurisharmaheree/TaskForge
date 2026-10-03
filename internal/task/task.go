package task

type Task struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
	Retries int                    `json:"retries"`
}

type Metrics struct {
	TotalJobsInQueue int64 `json:"total_jobs_in_queue"`
	JobsDone         int   `json:"jobs_done"`
	JobsFailed       int   `json:"jobs_failed"`
}