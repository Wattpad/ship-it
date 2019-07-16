package main

import (
	"fmt"
	shipitv1beta1 "ship-it-operator/api/v1beta1"
)

func main() {
	helmRelease := shipitv1beta1.HelmRelease{}

	fmt.Println("helm release test", helmRelease)
}
