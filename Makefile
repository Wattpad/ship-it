.PHONY: build docs jsonschema push

CHART_VERSION = $(shell helm inspect deploy/ship-it | awk '/version/{ print $$2; }')
CHART_REPOSITORY = $(shell helm repo list | awk '/wattpad/{ print $$1; }')

REGISTRY := 723255503624.dkr.ecr.us-east-1.amazonaws.com
VERSION := $(shell git rev-parse HEAD)

TARGET_IMAGE := $(TARGET):$(VERSION)
LATEST_IMAGE := $(TARGET):latest

KIND_CLUSTER_NAME := ship-it-dev

env-target:
ifndef TARGET
	    $(error TARGET is undefined)
endif

build: env-target
	docker build -t $(TARGET_IMAGE) -f cmd/$(TARGET)/Dockerfile .

push: build
	docker tag $(TARGET_IMAGE) $(REGISTRY)/$(TARGET_IMAGE)
	docker tag $(TARGET_IMAGE) $(REGISTRY)/$(LATEST_IMAGE)
	docker push $(REGISTRY)/$(TARGET_IMAGE)
	docker push $(REGISTRY)/$(LATEST_IMAGE)

chart:
	helm package deploy/ship-it
	helm s3 push ship-it-$(CHART_VERSION).tgz $(CHART_REPOSITORY)

# empty target
internal/api/models/*.go:

# api docs should be rebuilt when model code changes
api/*.json: internal/api/models/*.go
	go run tools/jsonschema/main.go

jsonschema: api/*.json

docs/operator-release-states.png: docs/operator-release-states.dot
	dot -Tpng docs/operator-release-states.dot -o docs/operator-release-states.png

docs: api/*.json docs/operator-release-states.png

kind-up:
	@echo Creating the $(KIND_CLUSTER_NAME) cluster...
	kind create cluster --config hack/$(KIND_CLUSTER_NAME).yaml --name $(KIND_CLUSTER_NAME)
	$(eval KUBECONFIG := $(shell kind get kubeconfig-path --name $(KIND_CLUSTER_NAME)))
	KUBECONFIG=$(KUBECONFIG) kubectl apply -f hack/tiller/rbac.yaml
	KUBECONFIG=$(KUBECONFIG) kubectl apply -f hack/github/secret.yaml
	KUBECONFIG=$(KUBECONFIG) helm init --service-account tiller
	KUBECONFIG=$(KUBECONFIG) kubectl rollout status deployment -n kube-system tiller-deploy
	@echo Done! Set your kubectl context:
	@echo
	@echo export KUBECONFIG=$(KUBECONFIG)

kind-down:
	@echo Destroying the $(KIND_CLUSTER_NAME) cluster...
	kind delete cluster --name $(KIND_CLUSTER_NAME)

kind-deploy:
	kind load docker-image --name $(KIND_CLUSTER_NAME) ship-it-api:$(VERSION)
	kind load docker-image --name $(KIND_CLUSTER_NAME) ship-it-syncd:$(VERSION)
	kind load docker-image --name $(KIND_CLUSTER_NAME) ship-it-operator:$(VERSION)
	helm upgrade --install localstack deploy/localstack
	helm upgrade --install ship-it deploy/ship-it --set api.image.tag=$(VERSION),api.image.repository=ship-it-api,syncd.image.tag=$(VERSION),syncd.image.repository=ship-it-syncd,operator.image.tag=$(VERSION),operator.image.repository=ship-it-operator,devEnv.DOGSTATSD_HOST="localhost"
