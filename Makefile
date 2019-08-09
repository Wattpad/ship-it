.PHONY: build docs jsonschema push

REGISTRY := 723255503624.dkr.ecr.us-east-1.amazonaws.com
VERSION := $(shell git rev-parse HEAD)

TARGET_IMAGE := $(TARGET):$(VERSION)
LATEST_IMAGE := $(TARGET):latest

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

# empty target
internal/api/models/*.go:

# api docs should be rebuilt when model code changes
api/*.json: internal/api/models/*.go
	go run tools/jsonschema/main.go

jsonschema: api/*.json

docs/operator-release-states.png:
	dot -Tpng docs/operator-release-states.dot -o docs/operator-release-states.png

docs: api/*.json docs/operator-release-states.png
