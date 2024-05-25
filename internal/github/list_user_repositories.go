package github

import (
	"fmt"
	"net/http"
)

type ListUserRepositoriesResponse []Repository

type ListUserRepositoriesOptions struct {
	Username string `json:"username"`
}

// ListUserRepositories returns a list of all the repositories that belong to the specified user on GitHub.
// https://mirasynth.stream/ghapiredir#list-repositories-for-a-user
func (c *ClientImplementation) ListUserRepositories(options *ListUserRepositoriesOptions) (*ListUserRepositoriesResponse, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", options.Username)

	return startRequest(c, &startRequestOptions[ListUserRepositoriesResponse]{
		URL:      url,
		Method:   http.MethodGet,
		UseToken: true,
		StatusCodes: map[int]statusCode{
			http.StatusOK: {},
			defaultStatusCode: {
				"github action runner registration token could not be fetched",
			},
		},
		Pagination: &pagination[ListUserRepositoriesResponse]{
			PerPage:   5,
			StartPage: 1,
			PageReducer: func(accumulator ListUserRepositoriesResponse, result ListUserRepositoriesResponse) *ListUserRepositoriesResponse {
				if len(result) <= 0 {
					return nil
				}
				accumulator = append(accumulator, result...)
				return &accumulator
			},
		},
	})
}
