package jira

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/config"
)

const (
	storyIcon   = "\uF02E" // 
	taskIcon    = "\uF631" // 
	bugIcon     = "\uEAD8" // 
	epicIcon    = "\uF0E7" // 
	defaultIcon = "\uF059" // 
	subtaskIcon = "\uF637" // 
	statusIcon  = "\uF058" // 
	dateIcon    = "\uF5F5" // 
	timeIcon    = "\uF64F" // 
)

//go:embed preview.tmpl
var previewTemplate string

func (c *client) Fetch(ctx context.Context) ([]alfred.Item, error) {
	log.Debug().Msg("fetching jira issues")
	suggestions := make([]alfred.Item, 0)
	var jiraCnf config.Jira
	var atlassianCnf config.Atlassian
	if err := config.Get(&jiraCnf, &atlassianCnf); err != nil {
		return nil, err
	}
	tp := jira.BasicAuthTransport{
		Username: atlassianCnf.Username,
		Password: atlassianCnf.Token,
	}
	client, err := jira.NewClient(tp.Client(), jiraCnf.URL)
	if err != nil {
		return nil, fmt.Errorf("error creating jira client: %w", err)
	}
	options := &jira.SearchOptions{
		StartAt:       0,
		MaxResults:    0,
		Expand:        "",
		Fields:        nil,
		ValidateQuery: "",
	}
	mapIssues := func(ji jira.Issue) error {
		issue := &Issue{ji}
		item, err := issue.ToItem(jiraCnf.URL)
		if err != nil {
			return err
		}
		suggestions = append(suggestions, *item)
		return nil
	}
	err = client.Issue.SearchPagesWithContext(ctx, fmt.Sprintf(`project IN (%s) AND status != Closed AND status != Completed AND status != Cancelled ORDER BY created DESC`, jiraCnf.ProjectKey), options, mapIssues)
	if err != nil {
		return nil, fmt.Errorf("error search issue pages: %w", err)
	}
	log.Debug().Int("count", len(suggestions)).Msg("fetched issues")
	return suggestions, nil
}
