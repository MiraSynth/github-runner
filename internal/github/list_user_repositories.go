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

	result, err := startRequest[ListUserRepositoriesResponse](c, &startRequestOptions{
		URL:      url,
		Method:   http.MethodGet,
		UseToken: false,
		StatusCodes: map[int]statusCode{
			http.StatusOK: {},
			defaultStatusCode: {
				"github action runner registration token could not be fetched",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
