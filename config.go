package main

import (
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/kelseyhightower/envconfig"
)

// Config provides the service's configuration options.
type Config struct {
	AWSRegion     string `envconfig:"AWS_REGION" required:"true"`
	DogstatsdHost string `split_words:"true" required:"true"`
	DogstatsdPort string `split_words:"true" default:"8125"`
	QueueName     string `split_words:"true" required:"true"`
	ServicePort   string `split_words:"true" default:"80"`
	GithubToken   string `split_words:"true" required:"true"`
	GithubOrg     string `split_words:"true" required:"true"`
}

// DataDogAddress returns the local address of the datadog agent.
func (c *Config) DataDogAddress() string {
	return net.JoinHostPort(c.DogstatsdHost, c.DogstatsdPort)
}

// AWS returns an AWS config using the service's config values.
func (c *Config) AWS() *aws.Config {
	return &aws.Config{
		Region: aws.String(c.AWSRegion),
	}
}

// FromEnv returns a config using environment values.
func FromEnv() (*Config, error) {
	env := new(Config)
	if err := envconfig.Process("", env); err != nil {
		return nil, err
	}
	return env, nil
}
