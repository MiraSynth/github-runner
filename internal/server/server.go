package server

import (
	"github.com/gin-gonic/gin"
	"mirasynth.stream/github-runner/internal/server/github"
	"mirasynth.stream/github-runner/internal/server/health"
)

func StartServer() {
	ginEngine := gin.Default()

	routerGroup := ginEngine.Group("/api/v1")

	health.RegisterController(routerGroup, "/health")
	github.RegisterController(routerGroup, "/github")

	ginEngine.Run(":3038")
}
