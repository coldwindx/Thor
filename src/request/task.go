package request

import "Thor/tools"

type TaskCreateReq struct {
	Pipeline string `json:"pipeline" binding:"required"`
	Input    string `json:"input"`
	Output   string `json:"output"`
	Options  string `json:"options"`
}

func (TaskCreateReq) GetValidateMessage() tools.ValidatorMessages {
	return tools.ValidatorMessages{
		"pipeline.required": "pipeline不能为空",
	}
}
