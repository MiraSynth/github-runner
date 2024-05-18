package github

import (
	"encoding/json"
	"fmt"
	"io"
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

	response, err := c.startRequest(url, http.MethodPost, true, options.RequestData)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("installation token could not be fetched")
	}

	var result *GetInstallationAccessTokenResponse
	resultBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resultBytes, result)
	if err != nil {
		return nil, err
	}

	return result, err
}
