package k8s

import (
	"fmt"
	shipitv1beta1 "ship-it-operator/api/v1beta1"
	"strconv"
)

func annotationFor(k string) string {
	return fmt.Sprintf("helmreleases.%s/%s", shipitv1beta1.GroupVersion.Group, k)
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
