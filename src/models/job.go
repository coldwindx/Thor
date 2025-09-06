package models

import "time"

type Job struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	TaskId    int64     `json:"task_id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	Retry     int       `json:"retry"`
	AwakenAt  time.Time `json:"awaken_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   time.Time `json:"deleted"`
}
