package github

import (
	"fmt"
	"net/http"
	"runtime"

	"mirasynth.stream/github-runner/internal/config"
	"mirasynth.stream/github-runner/internal/utils/generate_jwt"
)

type Client interface {
	GetInstallationAccessToken(*GetInstallationAccessTokenOptions) (*GetInstallationAccessTokenResponse, error)
	GetActionRunnersRegistrationToken(*GetActionRunnersRegistrationTokenOptions) (*GetActionRunnersRegistrationTokenResponse, error)
}

type ClientImplementation struct {
	token          string
	jwt            string
	installationId int
}

func CreateClient() (Client, error) {
	client := ClientImplementation{}

	jwt, err := generate_jwt.Generate(config.GetGitHubClientId(), config.GetGitHubSecret())
	if err != nil {
		return nil, err
	}
	client.jwt = jwt

	response, err := client.GetInstallationAccessToken(&GetInstallationAccessTokenOptions{
		RequestData: &GetActionRunnersRegistrationTokenRequest{
			Repositories: []string{},
			Permissions:  map[string]string{},
		},
	})

	if err != nil {
		return nil, err
	}

	client.token = response.Token

	return &client, nil
}

func defaultHeaders(request *http.Request) {
	request.Header.Set("Accept", "Accept")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("User-Agent", fmt.Sprintf("mirasynth/0.0.0-alpha GoLang/%s (%s; %s)", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	request.Header.Set("Content-Type", "application/json")
}

func defaultHeadersJWT(request *http.Request) error {
	defaultHeaders(request)

	jwt, err := generate_jwt.Generate(config.GetGitHubClientId(), config.GetGitHubSecret())
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))

	return nil
}

func defaultHeadersToken(request *http.Request, token string) {
	defaultHeaders(request)
	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
}
