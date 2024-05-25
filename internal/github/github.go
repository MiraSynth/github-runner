package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mirasynth.stream/github-runner/internal/config"
	"net/http"
	"time"
)

const (
	defaultStatusCode = 0
)

var instance *ClientImplementation

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
	37 * time.Second,
	101 * time.Second,
}

type Client interface {
	GetAuthenticatedApp(*GetAuthenticatedAppOptions) (*GetAuthenticatedAppResponse, error)
	CreateInstallationAccessTokenForApp(*CreateInstallationAccessTokenForAppOptions) (*CreateInstallationAccessTokenForAppResponse, error)
	GetActionRunnersRegistrationToken(*GetActionRunnersRegistrationTokenOptions) (*GetActionRunnersRegistrationTokenResponse, error)

	GetInstallationForAuthenticatedApp(*GetInstallationForAuthenticatedAppOptions) (*GetInstallationForAuthenticatedAppResponse, error)
	ListInstallationsForAuthenticatedApp(*ListInstallationsForAuthenticatedAppOptions) (*ListInstallationsForAuthenticatedAppResponse, error)

	ListUserRepositories(*ListUserRepositoriesOptions) (*ListUserRepositoriesResponse, error)

	refreshToken() error
	defaultHeadersJWT(request *http.Request) error
	defaultHeadersToken(request *http.Request) error
}

type ClientOptions struct {
	Repositories []ClientRepository `json:"repositories"`
}

type ClientInstallation struct {
	Id          int
	Permissions ClientPermissions
}

type ClientImplementation struct {
	options      *ClientOptions
	auth         *ClientToken
	installation *ClientInstallation
	context      context.Context
}

func GetClient(ctx context.Context, options *ClientOptions) (Client, error) {
	if instance != nil && options == nil {
		return instance, nil
	}

	if options == nil {
		options = &ClientOptions{}
	}

	instance = &ClientImplementation{
		options:      options,
		context:      ctx,
		installation: &ClientInstallation{},
	}

	err := validateClientOptions(options)
	if err != nil {
		return nil, err
	}

	authenticatedApps, err := instance.ListInstallationsForAuthenticatedApp(&ListInstallationsForAuthenticatedAppOptions{})
	if err != nil {
		return nil, err
	}

	for _, authenticatedApp := range *authenticatedApps {
		if authenticatedApp.AppId == config.GetGitHubAppId() {
			instance.installation.Id = authenticatedApp.Id
			instance.installation.Permissions = authenticatedApp.Permissions
			break
		}
	}

	if instance.installation.Id == 0 {
		return nil, fmt.Errorf("app no installed")
	}

	err = instance.refreshToken()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func validateClientOptions(options *ClientOptions) error {
	if options == nil {
		return fmt.Errorf("options argument must be provided to the client")
	}

	return nil
}

func (c *ClientImplementation) refreshToken() error {
	if c.auth != nil && c.auth.TokenExpiresAt.After(time.Now()) {
		return nil
	}

	response, err := c.CreateInstallationAccessTokenForApp(&CreateInstallationAccessTokenForAppOptions{
		RequestData: &CreateInstallationAccessTokenForAppRequest{
			Repositories: c.options.Repositories,
			Permissions:  c.installation.Permissions,
		},
	})

	if err != nil {
		return err
	}

	if c.auth == nil {
		c.auth = &ClientToken{}
	}

	c.auth.Token = response.Token
	c.auth.TokenExpiresAt = response.ExpiresAt

	return nil
}

type statusCode struct {
	ErrorMessage string
}

type pagination[T any] struct {
	PerPage     int
	StartPage   int
	PageReducer func(accumulator T, result T) *T
}

type startRequestOptions[T any] struct {
	URL         string
	Method      string
	UseToken    bool
	RequestData any
	StatusCodes map[int]statusCode
	Pagination  *pagination[T]
}

func startRequest[T any](c *ClientImplementation, options *startRequestOptions[T]) (*T, error) {
	var requestDataBytes []byte
	var err error
	if options.RequestData != nil {
		requestDataBytes, err = json.Marshal(options.RequestData)
		if err != nil {
			return nil, err
		}
	}

	var returnResult T

	var response *http.Response
	for {
		for _, backoff := range backoffSchedule {
			response, err = singleRequest(c, options, &requestDataBytes)

			if response == nil {
				return nil, err
			}

			statusCodeBehaviour, ok := options.StatusCodes[response.StatusCode]
			if !ok {
				statusCodeBehaviour, ok = options.StatusCodes[defaultStatusCode]
			}
			if !ok {
				statusCodeBehaviour = statusCode{
					ErrorMessage: "an error has occurred",
				}
			}

			strLen := len(statusCodeBehaviour.ErrorMessage)
			if strLen == 0 {
				break
			}

			err = handleError(response, &statusCodeBehaviour)
			time.Sleep(backoff)
		}

		if response == nil || err != nil {
			return nil, err
		}

		var result T
		resultBytes, resultBytesErr := io.ReadAll(response.Body)
		err = resultBytesErr
		if err != nil {
			response.Body.Close()
			return nil, err
		}

		err = json.Unmarshal(resultBytes, &result)
		if err != nil {
			response.Body.Close()
			return nil, err
		}

		if options.Pagination == nil {
			response.Body.Close()
			return &result, err
		}

		pageReducerResult := options.Pagination.PageReducer(returnResult, result)
		if pageReducerResult != nil {
			returnResult = *pageReducerResult
			response.Body.Close()
			options.Pagination.StartPage++
			continue
		}

		response.Body.Close()
		return &returnResult, nil
	}
}

func singleRequest[T any](c *ClientImplementation, options *startRequestOptions[T], requestDataBytes *[]byte) (*http.Response, error) {
	request, err := http.NewRequest(options.Method, options.URL, bytes.NewBuffer(*requestDataBytes))
	if err != nil {
		return nil, err
	}

	deadlineContext, cancel := context.WithTimeout(c.context, time.Second*5)
	defer cancel()
	request.WithContext(deadlineContext)

	if options.UseToken {
		err = c.defaultHeadersToken(request)
	} else {
		err = c.defaultHeadersJWT(request)
	}

	if err != nil {
		return nil, err
	}

	if options.Pagination != nil {
		query := request.URL.Query()
		query.Add("page", fmt.Sprintf("%d", options.Pagination.StartPage))
		query.Add("per_page", fmt.Sprintf("%d", options.Pagination.PerPage))
		request.URL.RawQuery = query.Encode()
	}

	client := &http.Client{}
	return client.Do(request)
}

func handleError(response *http.Response, options *statusCode) error {
	defer response.Body.Close()

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
