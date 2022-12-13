package config

type Jira struct {
	noDefault
	ProjectKey CSV    `env:"jira_projects,required"`
	URL        string `env:"jira_url,required"`
}

func (j *Jira) name() string {
	return "jira"
}
