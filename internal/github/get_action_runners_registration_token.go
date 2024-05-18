package github

import (
	"bytes"
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

func (c *ClientImplementation) GetActionRunnersRegistrationToken(options *GetActionRunnersRegistrationTokenOptions) (*GetActionRunnersRegistrationTokenResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runners/registration-token", options.Username, options.Repository)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		return nil, err
	}

	defaultHeadersToken(request, options.Token)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("github action runner registration token could not be fetched")
	}

	var result *GetActionRunnersRegistrationTokenResponse
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
