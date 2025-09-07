package models

type Pipeline struct {
	Name     string `json:"name"`
	Workflow string `json:"workflow"`
}

type WorkNode struct {
	JobId   int64   `json:"job_id"`
	JobName string  `json:"job_name"`
	Next    []int64 `json:"next"`
}
