# _Ship it!_ [![Build Status](https://travis-ci.com/Wattpad/ship-it.svg?branch=master)](https://travis-ci.com/Wattpad/ship-it) [![Go Report Card](https://goreportcard.com/badge/github.com/Wattpad/ship-it)](https://goreportcard.com/report/github.com/Wattpad/ship-it) [![codecov](https://codecov.io/gh/Wattpad/ship-it/branch/master/graph/badge.svg)](https://codecov.io/gh/Wattpad/ship-it)

<img src="https://media.giphy.com/media/143vPc6b08locw/giphy.gif" width="300">

_Ship it!_ is Wattpad's tool for continuous deployment to Kubernetes.

The technical background and architecture of Ship-it can be found in the overview documentation [here](./docs/OVERVIEW.md)

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

5. Connect to the service

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
