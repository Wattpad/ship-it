.PHONY: build push

ECR_REGISTRY := 723255503624.dkr.ecr.us-east-1.amazonaws.com
PROJECT_NAME := ship-it
VERSION := $(shell git rev-parse HEAD)

IMAGE := $(PROJECT_NAME):$(VERSION)

build:
	docker build -t $(IMAGE) .

push: build
	docker tag $(IMAGE) $(ECR_REGISTRY)/$(IMAGE)
	docker push $(ECR_REGISTRY)/$(IMAGE)
