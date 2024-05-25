package github

import (
	"fmt"
	"net/http"
	"time"
)

type ListInstallationsForAuthenticatedAppResponse []struct {
	Id                     int               `json:"id"`
	Account                Account           `json:"account"`
	AccessTokensUrl        string            `json:"access_tokens_url"`
	RepositoriesUrl        string            `json:"repositories_url"`
	HtmlUrl                string            `json:"html_url"`
	AppId                  int               `json:"app_id"`
	TargetId               int               `json:"target_id"`
	TargetType             string            `json:"target_type"`
	Permissions            ClientPermissions `json:"permissions"`
	Events                 []string          `json:"events"`
	SingleFileName         string            `json:"single_file_name"`
	HasMultipleSingleFiles bool              `json:"has_multiple_single_files"`
	SingleFilePaths        []string          `json:"single_file_paths"`
	RepositorySelection    string            `json:"repository_selection"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
	AppSlug                string            `json:"app_slug"`
	SuspendedAt            interface{}       `json:"suspended_at"`
	SuspendedBy            interface{}       `json:"suspended_by"`
}

type ListInstallationsForAuthenticatedAppOptions struct {
}

// ListInstallationsForAuthenticatedApp return the installed apps that belongs to the client data
// https://mirasynth.stream/ghapiredir#get-the-authenticated-app
func (c *ClientImplementation) ListInstallationsForAuthenticatedApp(_ *ListInstallationsForAuthenticatedAppOptions) (*ListInstallationsForAuthenticatedAppResponse, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations")

	return startRequest(c, &startRequestOptions[ListInstallationsForAuthenticatedAppResponse]{
		URL:      url,
		Method:   http.MethodGet,
		UseToken: false,
		Pagination: &pagination[ListInstallationsForAuthenticatedAppResponse]{
			PerPage:   10,
			StartPage: 1,
			PageReducer: func(accumulator ListInstallationsForAuthenticatedAppResponse, result ListInstallationsForAuthenticatedAppResponse) *ListInstallationsForAuthenticatedAppResponse {
				if len(result) <= 0 {
					return nil
				}
				accumulator = append(accumulator, result...)
				return &accumulator
			},
		},
		StatusCodes: map[int]statusCode{
			http.StatusOK: {},
			defaultStatusCode: {
				ErrorMessage: "app installation token could not be fetched",
			},
		},
	})
}
