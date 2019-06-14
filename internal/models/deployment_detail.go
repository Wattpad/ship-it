package models

import (
	"time"
)

type DeploymentDetail struct {
	Name         string    `json:"name"`
	Created      string    `json:"created"`
	LastDeployed time.Time `json:"lastDeployed"`
	Owner        struct {
		Team  string `json:"team"`
		Slack string `json:"slack"`
	} `json:"owner"`
	AutoDeploy bool `json:"autoDeploy"`
	Code       struct {
		Github string `json:"github"`
		Ref    string `json:"ref"`
	} `json:"code"`
	Build struct {
		Travis string `json:"travis"`
	} `json:"build"`
	Monitoring struct {
		Datadog struct {
			Dashboard string `json:"dashboard"`
			Monitors  string `json:"monitors"`
		} `json:"datadog"`
		Sumologic string `json:"sumologic"`
	} `json:"monitoring"`
	Artifacts struct {
		Docker struct {
			Image string `json:"image"`
			Tag   string `json:"tag"`
		} `json:"docker"`
		Chart struct {
			Path    string `json:"path"`
			Version string `json:"version"`
		} `json:"chart"`
	} `json:"artifacts"`
	Deployment struct {
		Status string `json:"status"`
	} `json:"deployment"`
}
