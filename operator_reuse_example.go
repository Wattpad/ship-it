package main

import (
	"context"
	"fmt"
	"os"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	runtime "k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	config "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	scheme := runtime.NewScheme()
	shipitv1beta1.AddToScheme(scheme)

	cl, err := client.New(config.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})

	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}

	rlsList := &shipitv1beta1.HelmReleaseList{}

	err = cl.List(context.Background(), rlsList, client.InNamespace("default"))

	if err != nil {
		fmt.Printf("failed to list helm releases in namespace default: %v\n", err)
		os.Exit(1)
	}

	for _, rls := range rlsList.Items {
		fmt.Println("Found HelmRelease:", rls.ObjectMeta.Name)
	}
}
