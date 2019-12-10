# _Ship it!_ [![Build Status](https://travis-ci.com/Wattpad/ship-it.svg?branch=master)](https://travis-ci.com/Wattpad/ship-it) [![Go Report Card](https://goreportcard.com/badge/github.com/Wattpad/ship-it)](https://goreportcard.com/report/github.com/Wattpad/ship-it) [![codecov](https://codecov.io/gh/Wattpad/ship-it/branch/master/graph/badge.svg)](https://codecov.io/gh/Wattpad/ship-it)

<img src="https://media.giphy.com/media/143vPc6b08locw/giphy.gif" width="300">

_Ship it!_ is Wattpad's tool for continuous deployment to Kubernetes.

The technical background and architecture of Ship-it can be found in the overview documentation [here](./docs/OVERVIEW.md)

## Operations

### Continuous Deployment

Ship-it watches `HelmRelease` custom resources in the namespace in which it is
deployed. Whenever a `HelmRelease` resource is created/updated/destroyed,
Ship-it performs the appropriate install/upgrade/delete Helm operation on the
associated Helm release.

An example of a minimal `HelmRelease` definition:

```
apiVersion: shipit.wattpad.com/v1beta1
kind: HelmRelease
metadata:
  name: my-service-name
  annotations:
    helmreleases.shipit.wattpad.com/autodeploy: "true"
spec:
  releaseName: my-release-name

  chart:
    repository: my-chart-repository
    name: my-chart-name
    version: my-chart-version

  values: {}
```

As a quick overview, there are 4 important sections to take note of.

`metadata.annotations` enrich the `HelmRelease` with additional information.
Annotations are always optional, however the `HelmRelease` will not autodeploy
as the default is `false`.

`spec.releaseName` defines the desired name of the Helm release

`spec.chart` provides the desired chart name and version

`spec.values` provides the desired chart values. In this example the service
does not provide any overriding values for the chart, however the field is still
explicitly required

A much more thorough documentation of the `HelmRelease` custom resource can be
found in the resource
[definition](./operator/config/crd/bases/shipit.wattpad.com_helmreleases.yaml),
or by `kubectl describe crd/helmreleases.shipit.wattpadhq.com` in a cluster
namespace where the CRD exists.

### Automatic Rollback

Ship-it initially has limited support for automatic rollbacks. When a Helm
release enters the `FAILED` state, Ship-it will perform a Helm rollback
operation to the most recent successful release version. Most often, this will
be the immediately previous release version.

Developers should take note that a `FAILED` Helm release is caused by
_deployment_ errors, not application runtime errors. Some non-exhaustive
examples are: incomplete chart values or corrupted manifests produced by invalid
templates, failure to create resources in the desired namespace, or tiller
internal error. A Helm release will not fail because a Pod crashes after its
deployment is successfully rolled out.

In the future, we plan to add support for more advanced rollback strategies
such as rollbacks triggered by failing liveness/readiness health probes or
developer defined conditional expressions on service metrics.

## Local Development

This project uses kind for local development and testing. To get started,
you'll need to install these tools to your development machine:

* helm: https://github.com/helm/helm
* kind: https://github.com/kubernetes-sigs/kind

1. Create a new kind cluster

```bash
$ make kind-up
```

2. Update kubectl's cluster context

```bash
$ export KUBECONFIG=$(kind get kubeconfig-path --name="ship-it-dev")
$ kubectl config current-context
kubernetes-admin@ship-it-dev
```

3. Build local images

```bash
$ TARGET=ship-it-api make build
$ TARGET=ship-it-syncd make build
$ cd operator && make docker-build
```

4. Deploy local ship-it images

```bash
$ make kind-deploy
```

5. Connect to the service's node port (check the chart values)

```bash
curl -i http://localhost:31901/api/releases
```

or

```bash
open http://localhost:31901
```

6. Destroy cluster when finished

```bash
make kind-down
```
