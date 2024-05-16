package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"mirasynth.stream/github-runner/internal/config"
)

func RegisterController(routerGroup *gin.RouterGroup) {
	webhookRouterGroup := routerGroup.Group("/webhook", verifySignatureMiddleware())

	webhookRouterGroup.POST("/webhook", func(c *gin.Context) {
		// do stuff
	})
}

func verifySignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := verifySignature(c.Request, config.GetGitHubWebhookSecret())

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Next()
	}
}

// original code snippet from https://stackoverflow.com/questions/53242837/validating-github-webhook-hmac-signature-in-go
// modified to return errors instead of logging them
func verifySignature(request *http.Request, key string) error {
	// Assuming a non-empty header
	gotHash := strings.SplitN(request.Header.Get("X-Hub-Signature"), "=", 2)
	if gotHash[0] != "sha1" {
		return fmt.Errorf("the X-Hub-Signature header contains invalid data")
	}
	defer request.Body.Close()

	b, err := io.ReadAll(request.Body)
	if err != nil {
		return fmt.Errorf("cannot read the request body: %s", err)
	}

	hash := hmac.New(sha1.New, []byte(key))
	if _, err := hash.Write(b); err != nil {
		return fmt.Errorf("cannot compute the HMAC for request: %s", err)
	}

	expectedHash := hex.EncodeToString(hash.Sum(nil))
	if gotHash[1] != expectedHash {
		return fmt.Errorf("signatures do not match")
	}

	return nil
}
