package github

import (
	"fmt"
	"mirasynth.stream/github-runner/internal/config"
	"mirasynth.stream/github-runner/internal/version"
	"net/http"
	"runtime"
)

func defaultHeaders(request *http.Request) {
	request.Header.Set("Accept", "Accept")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("User-Agent", fmt.Sprintf("mirasynth/%s GoLang/%s (%s; %s)", version.GetVersion(), runtime.Version(), runtime.GOOS, runtime.GOARCH))
	request.Header.Set("Content-Type", "application/json")
}

func (c *ClientImplementation) defaultHeadersJWT(request *http.Request) error {
	defaultHeaders(request)

	jwt, err := generateJwt(config.GetGitHubClientId(), config.GetGitHubKey())
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

	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.auth.Token))
	return nil
}
