package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"github.com/xujiajun/nutsdb"
	"golang.org/x/oauth2"

	"github.com/yogeshlonkar/alfred-workflows/pkg/config"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/cache"
)

const Keyword = "github"

type client struct {
	cache cache.ResultCache
}

func NewClient() (alfred.Suggester, error) {
	cache, err := cache.NewResultCache(Keyword)
	if err != nil {
		return nil, err
	}
	return &client{cache}, nil
}

func (c *client) Update(ctx context.Context, sync bool) (*alfred.Result, []byte, error) {
	var githubCnf config.GitHub
	if err := config.Get(&githubCnf); err != nil {
		return nil, nil, err
	}
	ghGQLClient := githubv4.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubCnf.Token},
	)))
	log.Debug().Str("SelectedOwner", githubCnf.SelectedOwner).Msg("fetching github repositories")
	results := &alfred.Result{}
	for _, owner := range githubCnf.Owners {
		keyword := keywordFor(owner)
		if githubCnf.SelectedOwner == "" || githubCnf.SelectedOwner == owner {
			_, data, err := c.cache.Update(ctx, keyword, sync, func(ctx context.Context) ([]alfred.Item, error) {
				items, err := c.FetchForOwner(ctx, owner, ghGQLClient)
				if err != nil {
					return nil, err
				}
				results.Items = append(results.Items, items...)
				return items, nil
			})
			if err != nil {
				return nil, nil, err
			}
			if githubCnf.SelectedOwner != "" {
				results.Keyword = keyword
				return results, data, nil
			}
		}
	}
	results.Keyword = Keyword
	data, err := json.Marshal(results)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling results: %w", err)
	}
	return results, data, nil
}

func keywordFor(owner string) string {
	return Keyword + "/" + owner
}

func (c *client) Cached(context.Context) ([]byte, error) {
	var githubCnf config.GitHub
	if err := config.Get(&githubCnf); err != nil {
		return nil, err
	}
	data, err := c.cache.GetJSON(keywordFor(githubCnf.SelectedOwner))
	if err != nil && !errors.Is(err, nutsdb.ErrKeyNotFound) {
		return nil, fmt.Errorf("error getting cached data: %w", err)
	}
	return data, nil
}

func (c *client) Close() {
	err := c.cache.Close()
	if err != nil {
		log.Error().Err(err).Msgf("error closing %s cache", Keyword)
	}
}

func (c *client) Keyword() string {
	return Keyword
}
