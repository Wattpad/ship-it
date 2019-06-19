package models

import (
	"time"
)

type Release struct {
	Name         string     `json:"name" jsonschema:"description=The name of the release"`
	Created      time.Time  `json:"created" jsonschema:"description=The time when the release was created"`
	LastDeployed time.Time  `json:"lastDeployed" jsonschema:"description=The time when the release was last deployed"`
	Owner        Owner      `json:"owner" jsonschema:"description=Ownership and contact information"`
	AutoDeploy   bool       `json:"autoDeploy" jsonschema:"description=The state of the release's auto-deployment option"`
	Code         SourceCode `json:"code" jsonschema:"description=The repository and branch ref of the release's source code"`
	Build        Build      `json:"build" jsonschema:"description=The CI build page of current release,required=true"`
	Monitoring   Monitoring `json:"monitoring" jsonschema:"description=The monitoring resources for the release"`
	Artifacts    Artifacts  `json:"artifacts" jsonschema:"description=The build artifacts of the release"`
	Status       string     `json:"status" jsonschema:"description=The status of the release,example=deployed,example=failed,example=pending_rollback,example=pending_install,example=pending_upgrade"`
}

type Owner struct {
	Squad string `json:"squad"`
	Slack string `json:"slack"`
}

type SourceCode struct {
	Github string `json:"github" jsonschema:"format=uri"`
	Ref    string `json:"ref" jsonschema:"format=uri"`
}

type Build struct {
	Travis string `json:"travis" jsonschema:"format=uri"`
}

type Monitoring struct {
	Datadog   Datadog `json:"datadog"`
	Sumologic string  `json:"sumologic" jsonschema:"format=uri"`
}

type Datadog struct {
	Dashboard string `json:"dashboard" jsonschema:"format=uri"`
	Monitors  string `json:"monitors" jsonschema:"format=uri"`
}

type Artifacts struct {
	Docker DockerArtifact `json:"docker"`
	Chart  HelmArtifact   `json:"chart"`
}

type DockerArtifact struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
}

type HelmArtifact struct {
	Path    string `json:"path" jsonschema:"format=uri"`
	Version string `json:"version" jsonschema:"example=1.2.3"`
}
