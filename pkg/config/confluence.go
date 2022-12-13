package config

type Confluence struct {
	noDefault
	URL    string `env:"confluence_url,required"`
	Spaces CSV    `env:"confluence_spaces,required"`
}

func (c *Confluence) name() string {
	return "confluence"
}
