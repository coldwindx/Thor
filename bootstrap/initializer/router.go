package initializer

import (
	"Thor/bootstrap"
	"Thor/utils/inject"
	"github.com/gin-gonic/gin"
)

func init() {
	v := &RouterInitializer{name: "RouterInitializer", order: 500}
	bootstrap.Beans.Provide(&inject.Object{Name: v.GetName(), Value: v})
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
	bootstrap.Router = gin.Default()
}

func (*RouterInitializer) Close() {

}
