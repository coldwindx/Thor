package do

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Job struct {
	ID        int64                 `gorm:"primaryKey"`
	Type      string                `gorm:"column:type"`
	Name      string                `gorm:"column:name"`
	Input     string                `gorm:"column:input"`
	Output    string                `gorm:"column:output"`
	Options   string                `gorm:"column:options"`
	CreatedAt time.Time             `gorm:"column:created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at"`
	Deleted   soft_delete.DeletedAt `gorm:"column:deleted;softDelete:flag"`
}

// TableName returns the table name of the job.
func (j Job) TableName() string {
	return "t_job"
}

type JobQuery struct {
	Job Job `gorm:"embedded" deepcopier:"ignore"`
}
