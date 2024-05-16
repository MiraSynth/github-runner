package health

import (
	"github.com/gin-gonic/gin"
	"mirasynth.stream/github-runner/internal/server/health/live"
	"mirasynth.stream/github-runner/internal/server/health/ready"
)

func RegisterController(routerGroup *gin.RouterGroup) {
	healthRouterGroup := routerGroup.Group("/health")

	live.RegisterController(healthRouterGroup)
	ready.RegisterController(healthRouterGroup)
}
