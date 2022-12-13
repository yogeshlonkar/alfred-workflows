package app

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/config"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/bitbucket"
	"github.com/yogeshlonkar/alfred-workflows/pkg/confluence"
	"github.com/yogeshlonkar/alfred-workflows/pkg/github"
	"github.com/yogeshlonkar/alfred-workflows/pkg/jira"
)

var (
	start     = time.Now()
	workflows = map[string]func() (alfred.Suggester, error){
		github.Keyword:     github.NewClient,
		bitbucket.Keyword:  bitbucket.NewClient,
		jira.Keyword:       jira.NewClient,
		confluence.Keyword: confluence.NewClient,
	}
	Keywords = fmt.Sprintf("%s, %s, %s, %s", github.Keyword, bitbucket.Keyword, jira.Keyword, confluence.Keyword)
)

func App(keyword string, sync bool, timeout time.Duration) error {
	log.Debug().Bool("sync", sync).Str("keyword", keyword).Send()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	workflow, ok := workflows[keyword]
	if !ok {
		return fmt.Errorf(`invalid argument '%s' to --type flag. Expected one of github, bitbucket, jira, confluence, all`, keyword)
	}
	log.Debug().Str("type", keyword).Msg("generating suggestions")
	suggester, err := workflow()
	if err != nil {
		return err
	}
	defer suggester.Close()
	if sync {
		return synchronise(ctx, suggester)
	}
	return suggest(ctx, suggester)
}

func synchronise(ctx context.Context, suggester alfred.Suggester) error {
	result, _, err := suggester.Update(ctx, true)
	if err != nil {
		return err
	}
	var alfredCnf config.Alfred
	if err = config.Get(&alfredCnf); err != nil {
		return err
	}
	msg := fmt.Sprintf("%s workflow updated ✓\n$$", alfredCnf.WorkflowName)
	if result.Keyword != "" {
		msg += fmt.Sprintf("\n%s ⬅", result.Keyword)
	}
	msg += fmt.Sprintf("\n%d items, %s ⏱", len(result.Items), time.Since(start).Round(time.Millisecond))
	fmt.Print(msg)
	return nil
}

func suggest(ctx context.Context, suggester alfred.Suggester) error {
	data, err := suggester.Cached(ctx)
	if err != nil {
		return err
	}
	if data == nil {
		if _, data, err = suggester.Update(ctx, true); err != nil {
			return err
		}
	}
	fmt.Print(string(data))
	return nil
}
