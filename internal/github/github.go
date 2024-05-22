package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultStatusCode = 0
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
	37 * time.Second,
	101 * time.Second,
}

type Client interface {
	GetInstallationAccessToken(*GetInstallationAccessTokenOptions) (*GetInstallationAccessTokenResponse, error)
	GetActionRunnersRegistrationToken(*GetActionRunnersRegistrationTokenOptions) (*GetActionRunnersRegistrationTokenResponse, error)

	ListUserRepositories(*ListUserRepositoriesOptions) (*ListUserRepositoriesResponse, error)

	refreshToken() error
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
	context        context.Context
}

func CreateClient(ctx context.Context, options *ClientOptions) (Client, error) {
	client := ClientImplementation{
		options: options,
		context: ctx,
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
	requestDataBytes, err := json.Marshal(options.RequestData)
	if err != nil {
		return nil, err
	}

	var returnResult *T

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

			if len(statusCodeBehaviour.ErrorMessage) <= 0 {
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

		pageReducerResult := options.Pagination.PageReducer(*returnResult, result)
		if pageReducerResult != nil {
			returnResult = pageReducerResult
			response.Body.Close()
			return returnResult, nil
		}

		response.Body.Close()
		return returnResult, nil
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
		err = c.defaultHeadersJWT(request)
	} else {
		err = c.defaultHeadersToken(request)
	}

	if err != nil {
		return nil, err
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
