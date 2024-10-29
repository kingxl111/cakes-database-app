package server

import (
	"cakes-database-app/pkg/service"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) NewRouter(ctx *context.Context, log *slog.Logger, env string) http.Handler {
	log.Info(
		"starting TheSweetsOfLifeApp",
		slog.String("env", env),
		slog.String("version", "1.0"),
	)
	log.Debug("debug messages are enabled")

	router := chi.NewRouter()
	router.Use()
	// Auth group
	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.SignUp(ctx, log))
		r.Post("/sign-in", h.SignIn(ctx, log))
	})
	
	// our new secure router will require jwt
	apiRouter := chi.NewRouter()
	apiRouter.Use(h.UserIdentityMiddleware()) 	// validate users
	apiRouter.Use(NewLogger(log))
	apiRouter.Route("/api", func(r chi.Router) {
		r.Post("/make-order", h.MakeOrder(ctx))
		r.Get("/view-orders", h.ViewOrders(ctx))
		r.Post("/change-order", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/delete-order", func(w http.ResponseWriter, r *http.Request) {})
	})

	router.Mount("/", apiRouter)

	return router
}

