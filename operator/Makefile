
# Image URL to use all building/pushing image targets
IMG ?= 723255503624.dkr.ecr.us-east-1.amazonaws.com/ship-it-operator:$(shell git rev-parse HEAD)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

CHART_PATH ?= "../deploy/ship-it"

all: manager

# Run tests
test: generate fmt vet manifests
	kill -9 $(shell ps aux | grep 'port-forward pod/tiller-deploy' | grep -v 'grep' | awk '{print $$2}') 2>/dev/null || true
	kubectl -n kube-system port-forward $(shell kubectl get po -n kube-system -o name -lapp=helm) 44134 &
	go test -v ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	kill -9 $(shell ps aux | grep 'port-forward pod/tiller-deploy' | grep -v 'grep' | awk '{print $$2}') 2>/dev/null || true
	kubectl -n kube-system port-forward $(shell kubectl get po -n kube-system -o name -lapp=helm) 44134 &
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	cp config/crd/bases/*.yaml $(CHART_PATH)/templates

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./api/...

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push: docker-build
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0-beta.2
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
