package models

import "time"

type Task struct {
	Id        int64     `json:"id"`
	Pipeline  string    `json:"pipeline"`
	DagStr    string    `json:"dag_str"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   int       `json:"deleted"`
}
