package ready

import (
	"github.com/gin-gonic/gin"
)

func RegisterController(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/ready", func(c *gin.Context) {
		c.String(200, "ready")
	})
}
