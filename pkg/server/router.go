package server

import (
	"cakes-database-app/pkg/service"
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) NewRouter(ctx *context.Context) http.Handler {
	router := chi.NewRouter()
	// Auth group
	router.Route("/auth", func(r chi.Router) {
		// TODO: implement handleFuncs 
		r.Post("/sign-up", h.SignUp(ctx))
		// r.Post("/sign-in", func(w http.ResponseWriter, r *http.Request) {})
	})

	// TODO: add other groups

	return router
}