package jira

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/andygrunwald/go-jira"

	"github.com/yogeshlonkar/alfred-workflows/pkg/alfred"
	"github.com/yogeshlonkar/alfred-workflows/pkg/misc"
)

type Issue struct {
	jira.Issue
}

func (i *Issue) Description() string {
	cmd := exec.Command("pandoc", "--from=jira", "--to=markdown_strict")
	cmd.Stdin = strings.NewReader(i.Fields.Description)
	data, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(bytes.TrimSpace(data))
}

func (i *Issue) RelativeUpdateDate() string {
	return misc.RelativeDate(time.Time(i.Fields.Updated))
}

func (i *Issue) AbsoluteUpdateTime() string {
	return time.Time(i.Fields.Updated).Format("15:04:05")
}

func (i *Issue) TypeIcon() string {
	switch i.Fields.Type.Name {
	case "Story":
		return storyIcon
	case "Epic":
		return epicIcon
	case "Task":
		return taskIcon
	case "Sub-task":
		return subtaskIcon
	case "Bug":
		return bugIcon
	default:
		return defaultIcon
	}
}

func (i *Issue) ToItem(url string) (*alfred.Item, error) {
	typeIcon := i.TypeIcon()
	issueTime := i.AbsoluteUpdateTime()
	issueDate := i.RelativeUpdateDate()
	tmpl, err := template.New("issue").Parse(fmt.Sprintf(previewTemplate, typeIcon, statusIcon, dateIcon, issueDate+" "+issueTime))
	if err != nil {
		return nil, fmt.Errorf("error parsing issue preview template: %w", err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, i)
	if err != nil {
		return nil, fmt.Errorf("error executing issue preview template: %w", err)
	}
	return &alfred.Item{
		UID:          i.Key,
		Title:        i.Fields.Summary,
		Subtitle:     fmt.Sprintf("%s %s | %s %s | %s %s | %s %s", typeIcon, i.Key, statusIcon, i.Fields.Status.Name, dateIcon, issueDate, timeIcon, issueTime),
		Arg:          i.Key,
		Match:        fmt.Sprintf("%s %s", i.Fields.Summary, i.Key),
		Autocomplete: i.Key,
		Mods: &alfred.Mods{
			Alt: &alfred.Mod{
				Valid:    true,
				Subtitle: "Preview",
				Variables: map[string]string{
					"preview": "true",
				},
				Arg: tpl.String(),
			},
		},
		Text: &alfred.Text{
			Copy:      fmt.Sprintf("%s/browse/%s", url, i.Key),
			LargeType: tpl.String(),
		},
	}, nil
}
