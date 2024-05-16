package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"mirasynth.stream/github-runner/internal/config"
	"mirasynth.stream/github-runner/internal/utils/generate_jwt"
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

	request.Header.Set("Accept", "Accept")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("User-Agent", fmt.Sprintf("mirasynth/0.0.0-alpha GoLang/%s (%s; %s)", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	request.Header.Set("Content-Type", "application/json")

	jwt, err := generate_jwt.Generate(config.GetGitHubClientId(), config.GetGitHubSecret())
	if err != nil {
		return false, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))

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
