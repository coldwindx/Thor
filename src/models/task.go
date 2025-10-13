package models

import "time"

type Task struct {
	Id        int64     `json:"id"`
	Pipeline  string    `json:"pipeline"`
	Dag       string    `json:"dag"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   int       `json:"deleted"`
}

type TaskQuery struct {
	Id       int64 `json:"id"`
	PageSize int   `json:"page_size"`
}
