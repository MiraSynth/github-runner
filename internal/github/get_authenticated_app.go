package github

import (
	"fmt"
	"net/http"
	"time"
)

type GetAuthenticatedAppResponse struct {
	Id          int                     `json:"id"`
	Slug        string                  `json:"slug"`
	NodeId      string                  `json:"node_id"`
	Owner       Owner                   `json:"owner"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ExternalUrl string                  `json:"external_url"`
	HtmlUrl     string                  `json:"html_url"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	Permissions InstallationPermissions `json:"permissions"`
	Events      []string                `json:"events"`
}

type GetAuthenticatedAppOptions struct {
}

// GetAuthenticatedApp return the app that belongs to the client data
// https://mirasynth.stream/ghapiredir#get-the-authenticated-app
func (c *ClientImplementation) GetAuthenticatedApp(_ *GetAuthenticatedAppOptions) (*GetAuthenticatedAppResponse, error) {
	url := fmt.Sprintf("https://api.github.com/app")

	return startRequest(c, &startRequestOptions[GetAuthenticatedAppResponse]{
		URL:      url,
		Method:   http.MethodGet,
		UseToken: false,
		StatusCodes: map[int]statusCode{
			http.StatusOK: {},
			defaultStatusCode: {
				ErrorMessage: "installation token could not be fetched",
			},
		},
	})
}
