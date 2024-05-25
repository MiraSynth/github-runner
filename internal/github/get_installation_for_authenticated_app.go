package github

import (
	"fmt"
	"net/http"
	"time"
)

type GetInstallationForAuthenticatedAppResponse struct {
	Id                     int                     `json:"id"`
	Account                Account                 `json:"account"`
	AccessTokensUrl        string                  `json:"access_tokens_url"`
	RepositoriesUrl        string                  `json:"repositories_url"`
	HtmlUrl                string                  `json:"html_url"`
	AppId                  int                     `json:"app_id"`
	TargetId               int                     `json:"target_id"`
	TargetType             string                  `json:"target_type"`
	Permissions            InstallationPermissions `json:"permissions"`
	Events                 []string                `json:"events"`
	SingleFileName         string                  `json:"single_file_name"`
	HasMultipleSingleFiles bool                    `json:"has_multiple_single_files"`
	SingleFilePaths        []string                `json:"single_file_paths"`
	RepositorySelection    string                  `json:"repository_selection"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
	AppSlug                string                  `json:"app_slug"`
	SuspendedAt            interface{}             `json:"suspended_at"`
	SuspendedBy            interface{}             `json:"suspended_by"`
}

type GetInstallationForAuthenticatedAppOptions struct {
	InstallationId int `json:"installationId"`
}

// GetInstallationForAuthenticatedApp return the installed app that belongs to the client data
// https://mirasynth.stream/ghapiredir#get-the-authenticated-app
func (c *ClientImplementation) GetInstallationForAuthenticatedApp(options *GetInstallationForAuthenticatedAppOptions) (*GetInstallationForAuthenticatedAppResponse, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d", options.InstallationId)

	return startRequest(c, &startRequestOptions[GetInstallationForAuthenticatedAppResponse]{
		URL:      url,
		Method:   http.MethodGet,
		UseToken: false,
		StatusCodes: map[int]statusCode{
			http.StatusOK: {},
			defaultStatusCode: {
				ErrorMessage: "app installation token could not be fetched",
			},
		},
	})
}
