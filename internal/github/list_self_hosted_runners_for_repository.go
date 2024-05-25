package github

import (
	"fmt"
	"net/http"
)

type ListSelfHostedRunnersForRepositoryResponse struct {
	TotalCount int      `json:"total_count"`
	Runners    []Runner `json:"runners"`
}

type LListSelfHostedRunnersForRepositoryOptions struct {
	Username   string `json:"username"`
	Repository string `json:"repository"`
}

// ListSelfHostedRunnersForRepository returns a list of all the GitHub self-hosted runners for a repository
// https://mirasynth.stream/ghapiredir#list-self-hosted-runners-for-a-repository
func (c *ClientImplementation) ListSelfHostedRunnersForRepository(options *LListSelfHostedRunnersForRepositoryOptions) (*ListSelfHostedRunnersForRepositoryResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runners", options.Username, options.Repository)

	return startRequest(c, &startRequestOptions[ListSelfHostedRunnersForRepositoryResponse]{
		URL:      url,
		Method:   http.MethodGet,
		UseToken: true,
		StatusCodes: map[int]statusCode{
			http.StatusOK: {},
			defaultStatusCode: {
				"github action runner registration token could not be fetched",
			},
		},
		Pagination: &pagination[ListSelfHostedRunnersForRepositoryResponse]{
			PerPage:   5,
			StartPage: 1,
			PageReducer: func(accumulator ListSelfHostedRunnersForRepositoryResponse, result ListSelfHostedRunnersForRepositoryResponse) *ListSelfHostedRunnersForRepositoryResponse {
				if len(result.Runners) <= 0 {
					return nil
				}
				accumulator.Runners = append(accumulator.Runners, result.Runners...)
				return &accumulator
			},
		},
	})
}
