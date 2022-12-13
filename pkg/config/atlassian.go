package config

type Atlassian struct {
	noDefault
	Username string `env:"atlassian_user,required"`
	Token    string `env:"atlassian_token,required"`
}

func (a *Atlassian) name() string {
	return "atlassian"
}
