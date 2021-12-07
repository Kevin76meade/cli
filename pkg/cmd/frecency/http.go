package frecency

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/ghrepo"
)

// get most recent PRs
// returns api.Issue for easier DB insertion
func getPullRequests(c *http.Client, repo ghrepo.Interface) ([]api.Issue, error) {
	apiClient := api.NewClientFromHTTP(c)
	query := `query GetPRs($owner: String!, $repo: String!) {
  repository(owner: $owner, name: $repo) {
    pullRequests(first: 100, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes {
        number
        title
      }
    }
  }
}`
	variables := map[string]interface{}{
		"owner": repo.RepoOwner(),
		"repo":  repo.RepoName(),
	}
	type responseData struct {
		Repository struct {
			PullRequests struct {
				Nodes []api.Issue
			}
		}
	}

	var resp responseData
	err := apiClient.GraphQL(repo.RepoHost(), query, variables, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Repository.PullRequests.Nodes, nil
}

// get PRs created after specified time
func getPullRequestsSince(c *http.Client, repo ghrepo.Interface, since time.Time) ([]api.Issue, error) {
	apiClient := api.NewClientFromHTTP(c)
	repoName := ghrepo.FullName(repo)

	// using search since gql `pullRequests` can't be filtered by date
	timeFmt := since.UTC().Format("2006-01-02T15:04:05-0700")
	searchQuery := fmt.Sprintf("repo:%s is:pr is:open created:>%s", repoName, timeFmt)
	query := `query GetPRsSince($query: String!) {
  search(query: $query, type: ISSUE, first: 100) {
    nodes {
      ... on PullRequest {
        title
		number
      }
    }
  }
}`
	variables := map[string]interface{}{"query": searchQuery}
	type responseData struct {
		Search struct {
			Nodes []api.Issue
		}
	}

	var resp responseData
	err := apiClient.GraphQL(repo.RepoHost(), query, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Search.Nodes, nil
}

// get issues created after specified time
func getIssuesSince(c *http.Client, repo ghrepo.Interface, since time.Time) ([]api.Issue, error) {
	apiClient := api.NewClientFromHTTP(c)
	query := `query GetIssuesSince($owner: String!, $repo: String!, $since: DateTime!, $limit: Int!) {
  repository(owner: $owner, name: $repo) {
    issues(first: $limit, orderBy: {field: CREATED_AT, direction: DESC}, filterBy: {since: $since, states: [OPEN]}) {
      nodes {
        number
        title
      }
    }
  }
}`
	variables := map[string]interface{}{
		"owner": repo.RepoOwner(),
		"repo":  repo.RepoName(),
		"since": since.UTC().Format("2006-01-02T15:04:05-0700"),
		"limit": 100,
	}
	type responseData struct {
		Repository struct {
			Issues struct {
				Nodes []api.Issue
			}
		}
	}
	var resp responseData
	err := apiClient.GraphQL(repo.RepoHost(), query, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Repository.Issues.Nodes, nil
}

func getIssues(c *http.Client, repo ghrepo.Interface) ([]api.Issue, error) {
	apiClient := api.NewClientFromHTTP(c)
	query := `query GetIssues($owner: String!, $repo: String!){
  repository(owner: $owner, name: $repo) {
    issues(first: 100, orderBy: {field: CREATED_AT, direction: DESC}, filterBy: {states: [OPEN]}) {
      nodes {
        number
        title
      }
    }
  }
}`
	variables := map[string]interface{}{
		"owner": repo.RepoOwner(),
		"repo":  repo.RepoName(),
	}
	type responseData struct {
		Repository struct {
			Issues struct {
				Nodes []api.Issue
			}
		}
	}
	var resp responseData
	err := apiClient.GraphQL(repo.RepoHost(), query, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Repository.Issues.Nodes, nil
}
