package initializer

import (
	"Thor/bootstrap"
	"Thor/bootstrap/inject"
	"github.com/gin-gonic/gin"
)

func init() {
	v := &RouterInitializer{name: "RouterInitializer", order: 500}
	bootstrap.Beans.Provide(&inject.Object{Name: v.GetName(), Value: v, Completed: true})
}

type RouterInitializer struct {
	name  string
	order int
}

func (t *RouterInitializer) GetName() string {
	return t.name
}
func (t *RouterInitializer) GetOrder() int {
	return t.order
}
func (*RouterInitializer) Initialize() {
	// step 初始化路由
	bootstrap.Beans.Provide(&inject.Object{Name: "Router", Value: gin.Default(), Completed: true})
}

func (*RouterInitializer) Close() {

}
