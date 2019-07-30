# The Ship-it Kubernetes Operator

The ship-it operator is a [custom controller](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#custom-controllers) that reconciles changes to a [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) we have defined of type `HelmRelease`.


## Requirements

This component uses the [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) SDK to generate the scaffolding for a custom controller. Follow their [Installation](https://book.kubebuilder.io/quick-start.html#installation) instructions to install the `kubebuilder` binary.


## Navigating this component

The main parts of the code you're likely to be interested in here are:

- The type definition for the [HelmRelease](./api/v1beta1/helmrelease_types.go)  type, from which the CRD is generated.
- The [Reconcile](./controllers/helmrelease_controller.go) func which contains the logic to run whenever a change is made to a `HelmRelease` resource. 
 

## Developing

Run `make install` to install the CRD into the cluster (you will need a local cluster running).

Run `make run` to run the controller. 
