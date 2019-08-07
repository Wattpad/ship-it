package models

type ReleaseResources struct {
	Name      string `json:"name" jsonschema:"description=The name of the release"`
	Resources string `json:"status" jsonschema:"description=The kubernetes resources belonging to the release"`
}
