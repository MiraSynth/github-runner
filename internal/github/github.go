package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultStatusCode = 0
)

type Client interface {
	GetInstallationAccessToken(*GetInstallationAccessTokenOptions) (*GetInstallationAccessTokenResponse, error)
	GetActionRunnersRegistrationToken(*GetActionRunnersRegistrationTokenOptions) (*GetActionRunnersRegistrationTokenResponse, error)

	refreshToken() error
	startRequest(options *startRequestOptions) (*http.Response, error)
	defaultHeadersJWT(request *http.Request) error
	defaultHeadersToken(request *http.Request) error
}

type ClientOptions struct {
	Repositories []ClientRepository `json:"repositories"`
	Permissions  ClientPermissions  `json:"permissions"`
}

type ClientImplementation struct {
	options        *ClientOptions
	auth           *ClientToken
	installationId int
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

type statusCode struct {
	ErrorMessage string
}

type startRequestOptions struct {
	URL         string
	Method      string
	UseToken    bool
	RequestData any
	StatusCodes map[int]statusCode
}

func (c *ClientImplementation) startRequest(options *startRequestOptions) (*http.Response, error) {
	optionsBytes, err := json.Marshal(options.RequestData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(options.Method, options.URL, bytes.NewBuffer(optionsBytes))
	if err != nil {
		return nil, err
	}

	if options.UseToken {
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

func startRequest[T any](c *ClientImplementation, options *startRequestOptions) (*T, error) {
	response, err := c.startRequest(options)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	statusCodeBehaviour, ok := options.StatusCodes[response.StatusCode]
	if !ok {
		statusCodeBehaviour, ok = options.StatusCodes[defaultStatusCode]
	}
	if !ok {
		statusCodeBehaviour = statusCode{
			ErrorMessage: "an error has occurred",
		}
	}

	if len(statusCodeBehaviour.ErrorMessage) > 0 {
		err = handleError(response, &statusCodeBehaviour)
		return nil, err
	}

	var result T
	resultBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func handleError(response *http.Response, options *statusCode) error {
	var result Error
	resultBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read body from error response, %s", err)
	}

	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return fmt.Errorf("could not parse json body from error response, %s", err)
	}

	return fmt.Errorf("%s. github said; %s %s", options.ErrorMessage, result.Message, result.DocumentationUrl)
}
