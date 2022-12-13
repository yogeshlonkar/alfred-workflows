package alfred

import (
	"context"
)

type Suggester interface {
	Keyword() string
	Fetch(ctx context.Context) ([]Item, error)
	Update(ctx context.Context, sync bool) (*Result, []byte, error)
	Cached(ctx context.Context) ([]byte, error)
	Close()
}
