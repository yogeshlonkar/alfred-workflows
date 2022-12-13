package confluence

import (
	"context"
	"fmt"

	"github.com/essentialkaos/go-confluence/v5"
	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/config"
	"github.com/yogeshlonkar/alfred-workflows/pkg/misc"
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
	log.Debug().Msg("fetching confluence pages")
	var atlassianCnf config.Atlassian
	var confluenceCnf config.Confluence
	if err := config.Get(&confluenceCnf, &atlassianCnf); err != nil {
		return nil, err
	}
	cClient, err := confluence.NewAPI(confluenceCnf.URL, atlassianCnf.Username, atlassianCnf.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating confluence client: %w", err)
	}
	suggestions := make([]alfred.Item, 0)
	start := 0
	for {
		result, err := cClient.Search(confluence.SearchParameters{
			CQL:   fmt.Sprintf("space IN (%s) and type=page", confluenceCnf.Spaces),
			Start: start,
			Limit: perPageCount,
		})
		if err != nil {
			return nil, fmt.Errorf("error searching in space %s: %w", "SE", err)
		}
		for _, page := range result.Results {
			url := confluenceCnf.URL + page.URL
			suggestions = append(suggestions, alfred.Item{
				UID:          page.Content.ID,
				Title:        page.Title,
				Subtitle:     fmt.Sprintf("ﮮ %s \uEB26 %s", misc.RelativeDate(page.LastModified.Time), page.Excerpt),
				Arg:          url,
				Match:        page.Title,
				Autocomplete: url,
				Text: &alfred.Text{
					Copy:      url,
					LargeType: content(page),
				},
			})
		}
		if len(suggestions) >= result.TotalSize {
			break
		}
		start += result.Size
	}
	log.Debug().Int("count", len(suggestions)).Msg("fetched confluence pages")
	return suggestions, nil
}

func content(page *confluence.SearchEntity) string {
	if page.Content != nil {
		if page.Content.Body != nil {
			if page.Content.Body.View != nil {
				return page.Content.Body.View.Value
			}
		}
	}
	return fmt.Sprintf("\uF0A9 %s\n\uEB26 %s", page.Title, page.Excerpt)
}
