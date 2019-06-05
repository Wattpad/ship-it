package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func New() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", health)

	r.Mount("/dashboard", http.FileServer(http.Dir("")))
	r.Mount("/static", http.FileServer(http.Dir("dashboard")))

	return r
}
