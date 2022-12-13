package jira

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/cache"
)

const Keyword = "jira"

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
	return c.cache.Update(ctx, Keyword, sync, c.Fetch)
}

func (c *client) Cached(ctx context.Context) ([]byte, error) {
	return c.cache.Cached(ctx, Keyword)
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
