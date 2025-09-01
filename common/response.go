package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 响应结构体
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{0, data, "OK"})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{code, nil, message})
}
