package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
