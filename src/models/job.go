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
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   int       `json:"deleted"`
}
type JobQuery struct {
	Id        int64     `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Options   string    `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   int       `json:"deleted"`
	PageSize  int       `json:"page_size"`
}
