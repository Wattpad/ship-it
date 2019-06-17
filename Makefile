.PHONY: build push run jsonschema

ECR_REGISTRY := 723255503624.dkr.ecr.us-east-1.amazonaws.com
PROJECT_NAME := ship-it
VERSION := $(shell git rev-parse HEAD)

IMAGE := $(PROJECT_NAME):$(VERSION)
LATEST_IMAGE := $(PROJECT_NAME):latest

build:
	docker build -t $(IMAGE) .

push: build
	docker tag $(IMAGE) $(ECR_REGISTRY)/$(IMAGE)
	docker tag $(IMAGE) $(ECR_REGISTRY)/$(LATEST_IMAGE)
	docker push $(ECR_REGISTRY)/$(IMAGE)
	docker push $(ECR_REGISTRY)/$(LATEST_IMAGE)

run: build
	docker run -p 8080:80 \
	    -e AWS_REGION="us-east-1" \
	    -e QUEUE_NAME="foo" \
	    -e DOGSTATSD_HOST="localhost" \
	    -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	    -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	    -e AWS_SECURITY_TOKEN=${AWS_SECURITY_TOKEN} \
	    -e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
	    $(shell docker images -q $(PROJECT_NAME) | head -n 1)

jsonschema:
	go run cmd/jsonschema/main.go
	# docker run -it --rm -v $(shell pwd):/workspace -w /workspace golang:1.12 go run cmd/jsonschema/main.go
