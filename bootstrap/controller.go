package bootstrap

import "github.com/gin-gonic/gin"

// Controller is the interface that wraps the Routes method.
// Routes adds the routes for the controller to the given engine.
type Controller interface {
	Routes(engine *gin.Engine)
}
