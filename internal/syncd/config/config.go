package config

import (
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/kelseyhightower/envconfig"
)

// Config provides the service's configuration options.
type Config struct {
	AWSRegion            string `envconfig:"AWS_REGION" required:"true"`
	DogstatsdHost        string `split_words:"true" required:"true"`
	DogstatsdPort        string `split_words:"true" default:"8125"`
	EcrQueue             string `split_words:"true" required:"true"`
	GithubAppID          int    `envconfig:"GITHUB_APP_ID" required:"true"`
	GithubAppSecret      string `split_words:"true" required:"true"`
	GithubInstallationID int    `envconfig:"GITHUB_INSTALLATION_ID" required:"true"`
	GithubOrg            string `split_words:"true" required:"true"`
	GithubQueue          string `split_words:"true" required:"true"`
	HelmTimeoutSeconds   int64  `split_words:"true" default:"10"`
	Namespace            string `split_words:"true" default:"default"`
	OperationsRepository string `split_words:"true" required:"true"`
	RegistryChartPath    string `split_words:"true" required:"true"`
	ReleaseBranch        string `split_words:"true" default:"master"`
	ReleaseName          string `split_words:"true" required:"true"`
	TillerHost           string `split_words:"true" required:"true"`
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

func (c *Config) HelmTimeout() time.Duration {
	return time.Duration(c.HelmTimeoutSeconds) * time.Second
}

// FromEnv returns a config using environment values.
func FromEnv() (*Config, error) {
	env := new(Config)
	if err := envconfig.Process("", env); err != nil {
		return nil, err
	}
	return env, nil
}
