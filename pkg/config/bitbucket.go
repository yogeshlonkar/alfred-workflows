package config

type Bitbucket struct {
	noDefault
	Username   string `env:"bitbucket_username,required"`
	Password   string `env:"bitbucket_app_password,required"`
	URL        string `env:"bitbucket_url,required"`
	Workspaces CSV    `env:"bitbucket_workspaces,required"`
}

func (b *Bitbucket) name() string {
	return "bitbucket"
}
