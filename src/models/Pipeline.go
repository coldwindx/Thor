package models

import "strings"

type PipelineNode struct {
	JobId   int64   `json:"job_id"`
	JobName string  `json:"job_name"`
	Next    []int64 `json:"next"`
}

type DAG struct {
	Dag []PipelineNode `json:"dag"`
}

type Pipeline struct {
	Name        string `json:"name"`
	WorkflowStr string `json:"workflow"`
	Workflow    map[string][]string
}

func ParseWorkflowToDag(workflow string) map[string][]string {
	var table map[string][]string
	links := strings.Split(workflow, ",")
	for _, link := range links {
		edge := strings.SplitN(link, "->", 0)
		table[edge[0]] = append(table[edge[0]], edge[1])
	}
	return table
}
