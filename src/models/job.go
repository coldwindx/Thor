package models

import (
	"time"
)

type Job struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	TaskId    int64     `json:"task_id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	Retry     int       `json:"retry"`
	Status    int       `json:"status"`
	AwakenAt  time.Time `json:"awaken_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   int       `json:"deleted"`
}

type JobQuery struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	TaskId    int64     `json:"task_id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	Retry     int       `json:"retry"`
	Status    int       `json:"status"`
	AwakenAt  time.Time `json:"awaken_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   int       `json:"deleted"`
}
