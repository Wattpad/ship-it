# _Ship it!_ [![Build Status](https://travis-ci.com/Wattpad/ship-it.svg?branch=master)](https://travis-ci.com/Wattpad/ship-it) [![Go Report Card](https://goreportcard.com/badge/github.com/Wattpad/ship-it)](https://goreportcard.com/report/github.com/Wattpad/ship-it) [![codecov](https://codecov.io/gh/Wattpad/ship-it/branch/master/graph/badge.svg)](https://codecov.io/gh/Wattpad/ship-it)

<img src="https://media.giphy.com/media/143vPc6b08locw/giphy.gif" width="300">

_Ship it!_ is Wattpad's tool for continuously deploying code.

## Local Development

This project uses skaffold for local development and testing. To get started,
you'll need to install these tools locally:

* helm: https://github.com/helm/helm
* minikube: https://github.com/kubernetes/minikube
* skaffold: https://github.com/GoogleContainerTools/skaffold

1. Start minikube cluster

```bash
$ minikube start
```

2. Update kubectl's cluster context

```bash
$ minikube update-context
$ kubectl config current-context
minikube
```

3. Install tiller in minikube

```bash
$ helm init
```

4. Run skaffold in development mode

```bash
$ skaffold dev
```

5. Get the service address

```bash
minikube service ship-it-api
```
