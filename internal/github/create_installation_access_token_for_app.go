package github

import (
	"fmt"
	"net/http"
	"time"
)

type CreateInstallationAccessTokenForAppResponse struct {
	Token               string       `json:"token"`
	ExpiresAt           time.Time    `json:"expires_at"`
	Permissions         Permissions  `json:"permissions"`
	RepositorySelection string       `json:"repository_selection"`
	Repositories        []Repository `json:"repositories"`
}

type CreateInstallationAccessTokenForAppRequest struct {
	Repositories []ClientRepository `json:"repositories"`
	Permissions  ClientPermissions  `json:"permissions"`
}

type CreateInstallationAccessTokenForAppOptions struct {
	RequestData *CreateInstallationAccessTokenForAppRequest
}

// CreateInstallationAccessTokenForApp returns an access token and the time it expires
// https://mirasynth.stream/ghapiredir#create-an-installation-access-token-for-an-app
func (c *ClientImplementation) CreateInstallationAccessTokenForApp(options *CreateInstallationAccessTokenForAppOptions) (*CreateInstallationAccessTokenForAppResponse, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", c.installation.Id)

	return startRequest(c, &startRequestOptions[CreateInstallationAccessTokenForAppResponse]{
		URL:         url,
		Method:      http.MethodPost,
		UseToken:    false,
		RequestData: options.RequestData,
		StatusCodes: map[int]statusCode{
			http.StatusCreated: {},
			defaultStatusCode: {
				ErrorMessage: "installation token could not be fetched",
			},
		},
	})
}
