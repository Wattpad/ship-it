#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

REPO_ROOT="${REPO_ROOT:-$(git rev-parse --show-toplevel)}"
cd "${REPO_ROOT}"

go mod vendor
export GO111MODULE="off"

# fake being in a gopath
FAKE_GOPATH="$(mktemp -d)"
trap 'rm -rf ${FAKE_GOPATH}' EXIT

FAKE_REPOPATH="${FAKE_GOPATH}/src/ship-it"
mkdir -p "$(dirname "${FAKE_REPOPATH}")" && ln -s "${REPO_ROOT}" "${FAKE_REPOPATH}"

export GOPATH="${FAKE_GOPATH}"
cd "${FAKE_REPOPATH}"

chmod +x vendor/k8s.io/code-generator/generate-groups.sh

./vendor/k8s.io/code-generator/generate-groups.sh \
    "deepcopy,client,lister,informer" \
    ship-it/pkg/generated ship-it/pkg/apis \
    "k8s.wattpad.com:v1alpha1" \
    --go-header-file tools/k8s-code-gen/boilerplate.go.txt

export GO111MODULE="on"
cd $REPO_ROOT
