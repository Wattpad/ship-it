package main

import (
	"net"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AWSRegion     string `envconfig:"AWS_REGION" required:"true"`
	DogstatsdHost string `split_words:"true" required:"true"`
	DogstatsdPort string `split_words:"true" default:"8125"`
	QueueName     string `split_words:"true" required:"true"`
	ServicePort   string `split_words:"true" default:"80"`
}

func (e *Config) DataDogAddress() string {
	return net.JoinHostPort(e.DogstatsdHost, e.DogstatsdPort)
}

func FromEnv() (*Config, error) {
	var env *Config
	if err := envconfig.Process("", env); err != nil {
		return nil, err
	}
	return env, nil
}
