package k8s

import (
	"strconv"
)

func annotationFor(k string) string {
	// FIXME this should be dynamic
	return "helmreleases.shipit.wattpad.com/" + k
}

type helmReleaseAnnotations map[string]string

func (a helmReleaseAnnotations) get(k string) string {
	return a[annotationFor(k)]
}

func (a helmReleaseAnnotations) AutoDeploy() bool {
	autoDeploy, err := strconv.ParseBool(a.get("autodeploy"))
	if err != nil {
		return false
	}

	return autoDeploy
}

func (a helmReleaseAnnotations) Code() string {
	return a.get("code")
}

func (a helmReleaseAnnotations) Datadog() string {
	return a.get("datadog")
}

func (a helmReleaseAnnotations) Squad() string {
	return a.get("squad")
}

func (a helmReleaseAnnotations) Slack() string {
	return a.get("slack")
}

func (a helmReleaseAnnotations) Sumologic() string {
	return a.get("sumologic")
}
