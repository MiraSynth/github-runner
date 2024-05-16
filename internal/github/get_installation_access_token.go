package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mirasynth.stream/github-runner/internal/config"
)

type InstallationAccessTokenOptions struct {
	Repositories []string          `json:"repositories"`
	Permissions  map[string]string `json:"permissions"`
}

func GetInstallationToken(options *InstallationAccessTokenOptions) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", config.GetGitHubInstallationId())

	optionsBytes, err := json.Marshal(options)
	if err != nil {
		return false, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(optionsBytes))
	if err != nil {
		return false, err
	}

	err = defaultHeadersJWT(request)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("installation token could not be fetched")
	}

	return true, err
}
