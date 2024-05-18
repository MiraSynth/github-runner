package github

import (
	"bytes"
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
	Repositories []string          `json:"repositories"`
	Permissions  map[string]string `json:"permissions"`
}

type GetInstallationAccessTokenOptions struct {
	RequestData *GetActionRunnersRegistrationTokenRequest
}

func (c *ClientImplementation) GetInstallationAccessToken(options *GetInstallationAccessTokenOptions) (*GetInstallationAccessTokenResponse, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", c.installationId)

	optionsBytes, err := json.Marshal(options.RequestData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(optionsBytes))
	if err != nil {
		return nil, err
	}

	err = defaultHeadersJWT(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("installation token could not be fetched")
	}

	var result *GetInstallationAccessTokenResponse
	resultBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resultBytes, result)
	if err != nil {
		return nil, err
	}

	return result, err
}
