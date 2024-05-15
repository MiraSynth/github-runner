package live

import (
	"github.com/gin-gonic/gin"
)

func RegisterController(routerGroup *gin.RouterGroup, rootpath string) {
	routerGroup.GET(rootpath, func(c *gin.Context) {
		c.String(200, "live")
	})
}
