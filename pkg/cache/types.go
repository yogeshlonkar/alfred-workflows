package cache

import (
	"context"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/config"
)

type Cache[K any] interface {
	Close() error
	Delete() error
	Exists() bool
	Get(key string) (K, error)
	GetJSON(key string) (json []byte, err error)
	Open() error
	Path() string
	Save(key string, value K) ([]byte, error)
	Stale() bool
}

type Fetch = func(ctx context.Context) ([]alfred.Item, error)

type ResultCache interface {
	Cache[alfred.Result]
	Update(ctx context.Context, keyword string, sync bool, fetch Fetch) (*alfred.Result, []byte, error)
	Cached(ctx context.Context, keyword string) ([]byte, error)
}

func NewCache[K any](bucket, path string) (Cache[K], error) {
	s := &cache[K]{bucket: bucket, path: path}
	err := s.Open()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func NewResultCache(keyword string) (ResultCache, error) {
	var alfredCnf config.Alfred
	if err := config.Get(&alfredCnf); err != nil {
		return nil, err
	}
	cache, err := NewCache[alfred.Result](keyword, alfredCnf.WorkflowCache)
	if err != nil {
		return nil, err
	}
	return &resultCache{cache}, nil
}
