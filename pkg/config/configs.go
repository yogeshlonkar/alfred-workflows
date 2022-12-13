package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Configuration interface {
	name() string
	defaults() error
}

type noDefault struct {
}

func (n noDefault) defaults() error {
	// nothing to set as default
	return nil
}

func Get(cnf ...Configuration) error {
	for _, cnf := range cnf {
		if cnf == nil {
			return fmt.Errorf("configuration should not be nil")
		}
		if err := env.Parse(cnf); err != nil {
			return fmt.Errorf("error parsing %s config: %w", cnf.name(), err)
		}
		if err := cnf.defaults(); err != nil {
			return fmt.Errorf("error setting %s config defaults: %w", cnf.name(), err)
		}
	}
	return nil
}
