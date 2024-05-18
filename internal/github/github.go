package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mirasynth.stream/github-runner/internal/config"
	"mirasynth.stream/github-runner/internal/utils/generate_jwt"
	"net/http"
	"runtime"
	"time"
)

type Client interface {
	GetInstallationAccessToken(*GetInstallationAccessTokenOptions) (*GetInstallationAccessTokenResponse, error)
	GetActionRunnersRegistrationToken(*GetActionRunnersRegistrationTokenOptions) (*GetActionRunnersRegistrationTokenResponse, error)

	refreshToken() error
	startRequest(url string, method string, useToken bool, requestData interface{}) (*http.Response, error)
	defaultHeadersJWT(request *http.Request) error
	defaultHeadersToken(request *http.Request) error
}

type ClientToken struct {
	token          string
	tokenExpiresAt time.Time
}

type ClientImplementation struct {
	options        *ClientOptions
	auth           *ClientToken
	installationId int
}

const PermissionRead ClientPermissionType = "read"
const PermissionWrite ClientPermissionType = "write"

const PermissionAdministration ClientPermissionScope = "administration"

type ClientRepository string
type ClientPermissionScope string
type ClientPermissionType string
type ClientPermissions map[ClientPermissionScope]ClientPermissionType

type ClientOptions struct {
	Repositories []ClientRepository `json:"repositories"`
	Permissions  ClientPermissions  `json:"permissions"`
}

func CreateClient(options *ClientOptions) (Client, error) {
	client := ClientImplementation{
		options: options,
	}

	err := validateClientOptions(options)
	if err != nil {
		return nil, err
	}

	err = client.refreshToken()
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func validateClientOptions(options *ClientOptions) error {
	if options == nil {
		return fmt.Errorf("options argument must be provided to the client")
	}

	if options.Repositories == nil || len(options.Repositories) < 1 {
		return fmt.Errorf("you must provide atleast one repository in the options argument")
	}

	if options.Permissions == nil || len(options.Permissions) < 1 {
		return fmt.Errorf("you must provide atleast one permission in the options argument")
	}

	return nil
}

func (c *ClientImplementation) refreshToken() error {
	if c.auth.tokenExpiresAt.After(time.Now()) {
		return nil
	}

	response, err := c.GetInstallationAccessToken(&GetInstallationAccessTokenOptions{
		RequestData: &GetActionRunnersRegistrationTokenRequest{
			Repositories: c.options.Repositories,
			Permissions:  c.options.Permissions,
		},
	})

	if err != nil {
		return err
	}

	c.auth.token = response.Token
	c.auth.tokenExpiresAt = response.ExpiresAt

	return nil
}

func (c *ClientImplementation) startRequest(url string, method string, useToken bool, requestData interface{}) (*http.Response, error) {
	optionsBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(optionsBytes))
	if err != nil {
		return nil, err
	}

	if useToken {
		err = c.defaultHeadersJWT(request)
	} else {
		err = c.defaultHeadersToken(request)
	}

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

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
