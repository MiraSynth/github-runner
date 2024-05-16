package github

import (
	"github.com/gin-gonic/gin"
	"mirasynth.stream/github-runner/internal/server/github/webhook"
)

func RegisterController(routerGroup *gin.RouterGroup) {
	githubRouterGroup := routerGroup.Group("/github")

	webhook.RegisterController(githubRouterGroup, "/webhook")
}
