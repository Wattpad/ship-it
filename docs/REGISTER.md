# Service Registration with Ship-it!

This document provides instructions on registering a GoLang microservice for deployment with Ship-it. It is assumed that there the service we are deploying is ready to be built into a docker image and sent to a docker repository. At Wattpad, this means the service (or changes to an existing one) are in PR and ready to be merged into highlander.  

To register the service with Ship-it, it needs to be added to the `ship-it-registry` chart, which is found in the repository providing the source of truth for the state of Kubernetes. At Wattpad, this repository is miranda. To configure the service for deployment, an entry must be added to the `templates` folder of this chart. the name of the file should match that of the service itself. 
