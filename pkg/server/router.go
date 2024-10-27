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
		r.Post("/sign-up", h.SignUp(ctx))
		r.Post("/sign-in", h.SignIn(ctx))
	})
	
	// our new secure router will require jwt
	apiRouter := chi.NewRouter()
	apiRouter.Use(h.UserIdentityMiddleware()) 	// validate users
	apiRouter.Route("/api", func(r chi.Router) {
		r.Post("/make-order", h.MakeOrder(ctx))
		r.Get("/view-order", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/change-order", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/delete-order", func(w http.ResponseWriter, r *http.Request) {})
	})

	router.Mount("/", apiRouter)

	return router
}