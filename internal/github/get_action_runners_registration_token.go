package github

import (
	"fmt"
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

	return startRequest(c, &startRequestOptions[GetActionRunnersRegistrationTokenResponse]{
		URL:      url,
		Method:   http.MethodPost,
		UseToken: true,
		StatusCodes: map[int]statusCode{
			http.StatusOK: {},
			http.StatusUnauthorized: {
				"the authorization details provided where invalid",
			},
			http.StatusForbidden: {
				"the request was forbidden",
			},
			http.StatusNotFound: {
				"the resource being requested was not found",
			},
			http.StatusUnprocessableEntity: {
				"the entity could not be processed, see additional error for information",
			},
			defaultStatusCode: {
				"github action runner registration token could not be fetched",
			},
		},
	})
}
