package initializer

import (
	"Thor/bootstrap"
	"github.com/gin-gonic/gin"
)

func init() {
	v := &RouterInitializer{name: "router", order: 500}
	bootstrap.Manager[v.name] = v
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
