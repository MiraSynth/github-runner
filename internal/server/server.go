package server

import (
	"github.com/gin-gonic/gin"
	"mirasynth.stream/github-runner/internal/server/health"
)

func StartServer() {
	ginEngine := gin.Default()

	routerGroup := ginEngine.Group("/api/v1")

	health.RegisterController(routerGroup, "/health")

	ginEngine.Run(":3038")
}
