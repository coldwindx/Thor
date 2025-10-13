package request

import (
	"Thor/tools"
	"time"
)

type JobInsertReq struct {
	Name      string    `json:"name" form:"name" binding:"required"`
	Input     string    `json:"input" form:"input"`
	Output    string    `json:"output" form:"output"`
	Options   string    `json:"options" form:"options"`
	AwakenAt  time.Time `json:"awaken_at" form:"awaken_at"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

// GetMessage 自定义错误信息
func (JobInsertReq) GetValidateMessage() tools.ValidatorMessages {
	return tools.ValidatorMessages{
		"name.required": "Job名称不能为空",
	}
}

type JobQueryReq struct {
	Name         string    `json:"name" form:"name" binding:"required"`
	Status       int       `json:"status" form:"status"`
	CreatedAfter time.Time `json:"created_after" form:"created_after"`
}

func (JobQueryReq) GetValidateMessage() tools.ValidatorMessages {
	return tools.ValidatorMessages{
		"name.required": "Job名称不能为空",
	}
}

type JobDeleteReq struct {
	Id int64 `json:"id" form:"id" binding:"required"`
}

func (JobDeleteReq) GetValidateMessage() tools.ValidatorMessages {
	return tools.ValidatorMessages{
		"id.required": "ID不能为空",
	}
}
