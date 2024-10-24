package httpserver

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	// Auth group
	router.Route("/auth", func(r chi.Router) {
		// TODO: implement handleFuncs 
		r.Post("/sign-up", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/sign-in", func(w http.ResponseWriter, r *http.Request) {})
	})

	// TODO: add other groups

	return router
}