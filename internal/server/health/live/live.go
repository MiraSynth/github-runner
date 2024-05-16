package live

import (
	"github.com/gin-gonic/gin"
)

func RegisterController(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/live", func(c *gin.Context) {
		c.String(200, "live")
	})
}
