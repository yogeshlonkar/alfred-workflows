package config

import (
	"fmt"
	"os/user"
	"path/filepath"
)

type Alfred struct {
	WorkflowBundleID         string `env:"alfred_workflow_bundleid"`
	Theme                    string `env:"alfred_theme"`
	ThemeSubtext             string `env:"alfred_theme_subtext"`
	VersionBuild             string `env:"alfred_version_build"`
	WorkflowVersion          string `env:"alfred_workflow_version"`
	WorkflowName             string `env:"alfred_workflow_name"`
	Version                  string `env:"alfred_version"`
	WorkflowUID              string `env:"alfred_workflow_uid"`
	WorkflowCache            string `env:"alfred_workflow_cache"`
	WorkflowData             string `env:"alfred_workflow_data"`
	ThemeSelectionBackground string `env:"alfred_theme_selection_background"`
	Debug                    int    `env:"alfred_debug"`
	Preferences              string `env:"alfred_preferences"`
	ThemeBackground          string `env:"alfred_theme_background"`
	PreferencesLocalhash     string `env:"alfred_preferences_localhash"`
}

func (a *Alfred) name() string {
	return "alfred"
}

func (a *Alfred) defaults() error {
	if a.WorkflowCache == "" {
		curUser, err := user.Current()
		if err != nil {
			return fmt.Errorf("error getting current user: %w", err)
		}
		a.WorkflowCache = filepath.Join(curUser.HomeDir, ".suggestions")
	}
	return nil
}
