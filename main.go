package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	http.ListenAndServe(":80", http.HandlerFunc(hello))
}
