package github

import (
	"bytes"
	"fmt"
	"net/http"
)

type ActionRunnersRegistrationTokenOptions struct {
	Username     string `json:"username"`
	Organization string `json:"organization"`
	Repository   string `json:"repository"`
	Token        string `json:"token"`
}

func GetActionRunnersRegistrationToken(options *ActionRunnersRegistrationTokenOptions) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runners/registration-token", options.Username, options.Repository)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		return false, err
	}

	defaultHeadersToken(request, options.Token)

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
