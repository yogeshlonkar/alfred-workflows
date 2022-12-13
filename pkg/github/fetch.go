package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/config"
	"github.com/yogeshlonkar/alfred-workflows/pkg/misc"
)

const perPageCount = 100

type Query struct {
	RepositoryOwner struct {
		Repositories struct {
			Nodes    []Node
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"repositories(first: $per_page, after: $endCursor, ownerAffiliations: OWNER)"`
	} `graphql:"repositoryOwner(login: $owner)"`
}

type Node struct {
	Name          githubv4.String
	NameWithOwner githubv4.String
	Description   githubv4.String
	UpdatedAt     githubv4.DateTime
	Issues        struct {
		TotalCount githubv4.Int
	} `graphql:"issues(states: OPEN)"`
	PullRequests struct {
		TotalCount githubv4.Int
	} `graphql:"pullRequests(states: OPEN)"`
}

func (n Node) RelativeUpdateDate() string {
	return misc.RelativeDate(n.UpdatedAt.Time)
}

func (c *client) Fetch(ctx context.Context) ([]alfred.Item, error) {
	log.Debug().Msg("fetching all github repositories")
	var githubCnf config.GitHub
	if err := config.Get(&githubCnf); err != nil {
		return nil, err
	}
	ghGQLClient := githubv4.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubCnf.Token},
	)))
	suggestions := make([]alfred.Item, 0)
	for _, owner := range githubCnf.Owners {
		items, err := c.FetchForOwner(ctx, owner, ghGQLClient)
		if err != nil {
			return nil, err
		}
		suggestions = append(suggestions, items...)
	}
	log.Debug().Int("count", len(suggestions)).Msg("fetched github repos")
	return suggestions, nil
}

func (c *client) FetchForOwner(ctx context.Context, owner string, ghGQLClient *githubv4.Client) ([]alfred.Item, error) {
	variables := map[string]interface{}{
		"owner":     githubv4.String(owner),
		"per_page":  githubv4.Int(perPageCount),
		"endCursor": (*githubv4.String)(nil),
	}
	suggestions := make([]alfred.Item, 0)
	var query Query
	for {
		if err := ghGQLClient.Query(ctx, &query, variables); err != nil {
			return nil, fmt.Errorf("error querying graphql on github: %w", err)
		}
		for _, repo := range query.RepositoryOwner.Repositories.Nodes {
			name := string(repo.Name)
			nameWithOwner := string(repo.NameWithOwner)
			url := fmt.Sprintf("https://github.com/%s", string(repo.NameWithOwner))
			suggestions = append(suggestions, alfred.Item{
				UID:          nameWithOwner,
				Title:        name,
				Subtitle:     fmt.Sprintf("\uEB0C %d | \uEA64 %d | \uF5F5 %s", repo.Issues.TotalCount, repo.PullRequests.TotalCount, repo.RelativeUpdateDate()),
				Arg:          url,
				Match:        name,
				Autocomplete: nameWithOwner,
				Mods: &alfred.Mods{
					Alt: &alfred.Mod{
						Valid:    true,
						Subtitle: "Open actions",
						Arg:      fmt.Sprintf("%s/actions", url),
						Variables: map[string]string{
							"action": "actions",
						},
					},
					Shift: &alfred.Mod{
						Valid:    true,
						Subtitle: "Open pull requests",
						Arg:      fmt.Sprintf("%s/pulls", url),
						Variables: map[string]string{
							"action": "pull_request",
						},
					},
					Ctrl: &alfred.Mod{
						Valid:    true,
						Subtitle: "Open issues",
						Arg:      fmt.Sprintf("%s/issues", url),
						Variables: map[string]string{
							"action": "issues",
						},
					},

					Cmd: &alfred.Mod{
						Valid:    true,
						Subtitle: "Open in github.dev Web IDE",
						Arg:      strings.ReplaceAll(url, ".com/", ".dev/"),
						Variables: map[string]string{
							"action": "github.dev",
						},
					},
				},
				Text: &alfred.Text{
					Copy:      url,
					LargeType: fmt.Sprintf("%s\n\uEB0C %d | \uEA64 %d | \uF5F5 %s", string(repo.NameWithOwner), repo.Issues.TotalCount, repo.PullRequests.TotalCount, repo.RelativeUpdateDate()),
				},
			})
		}
		if !query.RepositoryOwner.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["endCursor"] = githubv4.NewString(query.RepositoryOwner.Repositories.PageInfo.EndCursor)
	}
	return suggestions, nil
}
