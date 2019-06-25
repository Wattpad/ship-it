.PHONY: build push run

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
	docker tag $(IMAGE) $(REGISTRY)/$(TARGET_IMAGE)
	docker tag $(IMAGE) $(REGISTRY)/$(LATEST_IMAGE)
	docker push $(REGISTRY)/$(TARGET_IMAGE)
	docker push $(REGISTRY)/$(LATEST_IMAGE)

run: build
	docker run -p 8080:80 \
	    -e AWS_REGION="us-east-1" \
	    -e QUEUE_NAME="foo" \
	    -e DOGSTATSD_HOST="localhost" \
	    -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	    -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	    -e AWS_SECURITY_TOKEN=${AWS_SECURITY_TOKEN} \
	    -e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
	    -e GITHUB_TOKEN=fake \
	    -e GITHUB_ORG="wattpad" \
	    $(shell docker images -q $(TARGET) | head -n 1)
