package config

import (
	"net"

	"github.com/kelseyhightower/envconfig"
)

// Config provides the service's configuration options.
type Config struct {
	DogstatsdHost string `split_words:"true" required:"true"`
	DogstatsdPort string `split_words:"true" default:"8125"`
	ServicePort   string `split_words:"true" default:"80"`
}

// DataDogAddress returns the local address of the datadog agent.
func (c *Config) DataDogAddress() string {
	return net.JoinHostPort(c.DogstatsdHost, c.DogstatsdPort)
}

// FromEnv returns a config using environment values.
func FromEnv() (*Config, error) {
	env := new(Config)
	if err := envconfig.Process("", env); err != nil {
		return nil, err
	}
	return env, nil
}
