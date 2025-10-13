package models

import (
	"time"
)

type JobStatus int

const (
	UNKNOWN        JobStatus = iota // 禁止使用，仅用于标注数据库查询
	INITING        JobStatus = iota // 正在初始化
	WAITING                         // 等待执行中
	START_FAILURE                   // 启动失败
	RUNNING                         // 正在执行中
	RUNNED_FAILURE                  // 执行失败
	SUCCESS                         // 执行成功
)

type Job struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	TaskId    int64     `json:"task_id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	Status    int       `json:"status"`
	Locked    int       `json:"locked"`
	Retry     int       `json:"retry"`
	AwakenAt  time.Time `json:"awaken_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   int       `json:"deleted"`
}

type JobQuery struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	TaskId       int64     `json:"task_id"`
	Input        string    `json:"input"`
	Output       string    `json:"output"`
	Options      string    `json:"options"`
	Status       int       `json:"status"`
	Retry        int       `json:"retry"`
	Locked       int       `json:"locked"`
	AwakenAt     time.Time `json:"awaken_at"`
	CreatedAfter time.Time `json:"created_after"`
	UpdatedAt    time.Time `json:"updated_at"`
	Deleted      int       `json:"deleted"`
	PageSize     int       `json:"page_size"`
}
