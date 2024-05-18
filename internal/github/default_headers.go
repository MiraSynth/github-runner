package github

import (
	"fmt"
	"mirasynth.stream/github-runner/internal/config"
	"mirasynth.stream/github-runner/internal/utils/generate_jwt"
	"net/http"
	"runtime"
)

func defaultHeaders(request *http.Request) {
	request.Header.Set("Accept", "Accept")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("User-Agent", fmt.Sprintf("mirasynth/0.0.0-alpha GoLang/%s (%s; %s)", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	request.Header.Set("Content-Type", "application/json")
}

func (c *ClientImplementation) defaultHeadersJWT(request *http.Request) error {
	defaultHeaders(request)

	jwt, err := generate_jwt.Generate(config.GetGitHubClientId(), config.GetGitHubSecret())
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))

	return nil
}

func (c *ClientImplementation) defaultHeadersToken(request *http.Request) error {
	defaultHeaders(request)

	err := c.refreshToken()
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.auth.token))
	return nil
}
