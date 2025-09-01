package bootstrap

import (
	"Thor/ctx"
	"github.com/gin-gonic/gin"
)

func init() {
	v := &RouterInitializer{name: "router", order: 500}
	Manager[v.name] = v
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
	ctx.Router = gin.Default()
}

func (*RouterInitializer) Close() {

}
