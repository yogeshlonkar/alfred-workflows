package alfred

import (
	"encoding/json"
	"fmt"
)

type Text struct {
	Copy      string `json:"copy,omitempty"`
	LargeType string `json:"largetype,omitempty"`
}

type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

type Mod struct {
	Valid     bool              `json:"valid"`
	Arg       string            `json:"arg,omitempty"`
	Subtitle  string            `json:"subtitle,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

type Mods struct {
	Alt    *Mod `json:"alt,omitempty"`
	Cmd    *Mod `json:"cmd,omitempty"`
	Shift  *Mod `json:"shift,omitempty"`
	Ctrl   *Mod `json:"ctrl,omitempty"`
	CmdAlt *Mod `json:"cmd+alt,omitempty"`
}

type Item struct {
	UID          string `json:"uid,omitempty"`
	Title        string `json:"title,omitempty"`
	Subtitle     string `json:"subtitle,omitempty"`
	Arg          string `json:"arg,omitempty"`
	Icon         *Icon  `json:"icon,omitempty"`
	Match        string `json:"match,omitempty"`
	Autocomplete string `json:"autocomplete,omitempty"`
	Text         *Text  `json:"text,omitempty"`
	Mods         *Mods  `json:"mods"`
}

// {
//  "uit":      "update-$owner",
//  "title":    "Update keyword for $owner",
//  "subtitle": "\uF6D9 Fetch suggestions from $owner",
//  "arg":      "update-suggestions/github/$owner",
//  "match":    "update $owner"
// }

func (i Item) ToJSON() ([]byte, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item: %w", err)
	}
	return data, nil
}
