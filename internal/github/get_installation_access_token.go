package github

import (
	"fmt"
	"net/http"
	"time"
)

type GetInstallationAccessTokenResponse struct {
	Token               string       `json:"token"`
	ExpiresAt           time.Time    `json:"expires_at"`
	Permissions         Permissions  `json:"permissions"`
	RepositorySelection string       `json:"repository_selection"`
	Repositories        []Repository `json:"repositories"`
}

type GetActionRunnersRegistrationTokenRequest struct {
	Repositories []ClientRepository `json:"repositories"`
	Permissions  ClientPermissions  `json:"permissions"`
}

type GetInstallationAccessTokenOptions struct {
	RequestData *GetActionRunnersRegistrationTokenRequest
}

// GetInstallationAccessToken returns an access token and the time it expires
// https://mirasynth.stream/ghapiredir#create-an-installation-access-token-for-an-app
func (c *ClientImplementation) GetInstallationAccessToken(options *GetInstallationAccessTokenOptions) (*GetInstallationAccessTokenResponse, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", c.installationId)

	return startRequest(c, &startRequestOptions[GetInstallationAccessTokenResponse]{
		URL:      url,
		Method:   http.MethodPost,
		UseToken: true,
		StatusCodes: map[int]statusCode{
			http.StatusCreated: {},
			defaultStatusCode: {
				ErrorMessage: "installation token could not be fetched",
			},
		},
	})
}
