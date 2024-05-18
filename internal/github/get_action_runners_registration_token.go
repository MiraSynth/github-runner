package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GetActionRunnersRegistrationTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetActionRunnersRegistrationTokenOptions struct {
	Username     string `json:"username"`
	Organization string `json:"organization"`
	Repository   string `json:"repository"`
	Token        string `json:"token"`
}

// GetActionRunnersRegistrationToken returns a registration token to be used when registering a self-hosted runner
// on GitHub.
// https://mirasynth.stream/ghapiredir#get-a-self-hosted-runner-for-a-repository
func (c *ClientImplementation) GetActionRunnersRegistrationToken(options *GetActionRunnersRegistrationTokenOptions) (*GetActionRunnersRegistrationTokenResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runners/registration-token", options.Username, options.Repository)

	response, err := c.startRequest(url, http.MethodPost, false, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("github action runner registration token could not be fetched")
	}

	var result *GetActionRunnersRegistrationTokenResponse
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
