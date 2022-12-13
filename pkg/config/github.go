package config

import (
	"flag"
	"fmt"
	"strings"
)

var (
	ghOwner = flag.String("gh-owner", "", "GitHub owner to fetch/ suggest repositories for")
)

type GitHub struct {
	Token         string `env:"GH_TOKEN,required"`
	Owners        CSV    `env:"github_owners,required"`
	SelectedOwner string
}

func (g *GitHub) name() string {
	return "github"
}

func (g *GitHub) defaults() error {
	if ghOwner != nil {
		g.SelectedOwner = strings.ReplaceAll(*ghOwner, "github/", "")
		return nil
	}
	return fmt.Errorf("expected owner flag to be not nil")
}
