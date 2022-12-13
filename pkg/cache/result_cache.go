package cache

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
)

type resultCache struct {
	Cache[alfred.Result]
}

func (rc *resultCache) Update(ctx context.Context, keyword string, sync bool, fetch Fetch) (*alfred.Result, []byte, error) {
	var result alfred.Result
	var data []byte
	if sync || rc.Stale() {
		suggestions, err := fetch(ctx)
		if err != nil {
			return nil, nil, err
		}
		result = alfred.Result{Items: suggestions}
		data, err = rc.Save(keyword, result)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating %s suggestions: %w", keyword, err)
		}
		log.Debug().Msgf("updated %s suggestions", keyword)
	}
	return &result, data, nil
}

func (rc *resultCache) Cached(ctx context.Context, keyword string) ([]byte, error) {
	return rc.Cache.GetJSON(keyword)
}
