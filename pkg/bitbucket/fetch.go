package bitbucket

import (
	"context"
	"fmt"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/config"
)

const perPageCount = 100

func (c *client) Fetch(ctx context.Context) ([]alfred.Item, error) {
	done := make(chan []alfred.Item, 1)
	var items []alfred.Item
	var err error
	go func() {
		items, err = c.fetch()
		close(done)
	}()
	select {
	case <-done:
		return items, err
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	}
}

func (c *client) fetch() ([]alfred.Item, error) {
	log.Debug().Msg("fetching bitbucket repositories")
	var bitbucketCnf config.Bitbucket
	if err := config.Get(&bitbucketCnf); err != nil {
		return nil, err
	}
	bClient := bitbucket.NewBasicAuth(bitbucketCnf.Username, bitbucketCnf.Password)
	bClient.Pagelen = perPageCount
	ro := &bitbucket.RepositoriesOptions{
		Owner: "xivart",
	}
	res, err := bClient.Repositories.ListForAccount(ro)
	if err != nil {
		return nil, fmt.Errorf("error listing repositories: %w", err)
	}
	suggestions := make([]alfred.Item, 0)
	for _, repo := range res.Items {
		href := links(repo.Links["html"])["href"].(string)
		suggestions = append(suggestions, alfred.Item{
			UID:          repo.Uuid,
			Title:        repo.Name,
			Subtitle:     subtitle(repo),
			Arg:          href,
			Match:        repo.Name,
			Autocomplete: repo.Slug,
			Mods: &alfred.Mods{
				Shift: &alfred.Mod{
					Valid:    true,
					Subtitle: "Open pull requests",
					Arg:      fmt.Sprintf("%s/pull-requests", href),
					Variables: map[string]string{
						"action": "pull_request",
					},
				},
			},
			Text: &alfred.Text{
				Copy:      href,
				LargeType: fmt.Sprintf("%s - %s\n%s", repo.Name, repo.Slug, href),
			},
		})
	}
	log.Debug().Int("count", len(suggestions)).Msg("fetched bitbucket repos")
	return suggestions, nil
}

func links(m interface{}) map[string]interface{} {
	return m.(map[string]interface{})
}

func subtitle(repo bitbucket.Repository) string {
	if repo.Description != "" {
		return repo.Description
	}
	return repo.Slug
}
