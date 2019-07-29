# Ship-it! Overview

This document provides a summary of Velocity's Kubernetes deployment tool. Documentation pretaining to registering and deploying microservices with Ship-it can be found [here](./REGISTER.md). 

## Background

Before Ship-it, there was [Kube-deploy](https://github.com/Wattpad/kube-deploy). Kube-deploy
was built to serve the simple purpose of
deploying docker images to Wattpad's Kubernetes cluster. As a result, the final product did 
not have the extensibility that Velocity required to make continuous deployment (CD)
improvements later on. The technical debt introduced by kube-deploy came to a head when
Velocity attempted to introduce an automated rollaback feature for failed releases into
kube-deploy. We ultimately discovered that to achieve this, a large swoth of kube-deploy
would need to be rewritten.
Instead of rewriting kube-deploy, Velocity decided to either build or use a pre-existing tool.
The final decision was to build Ship-it, which would have the features required for
developers to more confidently push changes to production while providing full visibility at
each stage of the deployment process. Specifically, Ship-it would implement the following 
features and practices:  

- Automated rollbacks based on release state and Data Dog monitors
- A web UI to monitor deployments and their state (deploying, deployed, rolled back etc.)
- Standardized and extensible Helm charts for deployments instead of a k8s.yml file
- Use of Custom Resource Definitions and Kubernetes Operators to reconcile the state of miranda with the state of the cluster
- Event driven architecture for detecting developer actions as opposed to polling various external dependencies done in kube-deploy.
- Agile development delivering functionality in vertical slices
- Unit testing and test driven development
- Open source development
- Cloud agnostic project where business logic is decoupled from the API calls it depends upon

## Diagram
![Architecture](./arch.png)

## Deployment Process Overview
This section follows the architecture diagram above describing each step Ship-it takes behind the scenes to take newly merged code and deploy it to kubernetes. It is 

1. Microservices containing code changes are filtered out by our existing CI pipeline and built into new docker images. The newly built images are then pushed to a docker registry in ECR.
2. The image push event is picked up by a CloudWatch rule monitoring ECR. The rule loads the push event into an SQS Queue. 
3. Ship-it has an SQS consumer implmentation, which takes the image push event, parses out the  image tag and pushes the updated tag to the Kubernetes Custom Resource YAML in miranda corresponding to that release. This step ensures miranda accurately represents the version of the service running in production
4. To deploy the service, the Helm chart for the service is downloaded from the repository in S3 and used to deploy the image to the cluster using the Helm API.
5. The Custom Resource Definition (CRD) is updated using the Kubernetes operator within ship-it to change the state of the release to deploying status.
6. The operator will update the CRD each time the release state changes.
7. The API server reads off the CRD to get the state of each Helm Release and exposes a REST API to read information about the releases.
8. The Web UI consumes this API to present engineers with the state and configuration of their deployment.
9. If all goes according to plan the state of the release should advance from "deploying" to "deployed", if there is an issue, a rollback is triggered and the "rolled back" state is reflected on the UI.

