package main

import (
	"net/http"
	"ship-it/internal/api"
)

func main() {
	http.ListenAndServe(":80", api.New())
}
