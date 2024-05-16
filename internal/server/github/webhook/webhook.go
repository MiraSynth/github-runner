package webhook

import "github.com/gin-gonic/gin"

func RegisterController(routerGroup *gin.RouterGroup) {
	webhookRouterGroup := routerGroup.Group("/webhook")

	webhookRouterGroup.POST("/webhook", func(c *gin.Context) {
		c.String(200, "live")
	})
}
