package server

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/kingxl111/cakes-database-app/internal/service"

	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) NewRouter(ctx *context.Context, log *logrus.Logger, env string) http.Handler {
	log.Info(
		"STARTING TheSweetsOfLifeApp",
		log.WithFields(logrus.Fields{
			"env":     env,
			"version": "1.0",
		}),
	)
	log.Debug("debug messages are enabled")

	router := chi.NewRouter()
	router.Use(NewLogger(log))
	// Auth group
	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.SignUp(ctx, log))
		r.Post("/sign-in", h.SignIn(ctx, log))
	})

	// our new secure router will require jwt
	apiRouter := chi.NewRouter()
	apiRouter.Use(h.UserIdentityMiddleware()) // validate users
	apiRouter.Use(NewLogger(log))
	apiRouter.Route("/api", func(r chi.Router) {
		r.Post("/make-order", h.MakeOrder(ctx, log))
		r.Get("/view-orders", h.ViewOrders(ctx, log))
		r.Post("/change-order", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/delete-order", func(w http.ResponseWriter, r *http.Request) {})

		r.Get("/cakes", h.Cakes(ctx, log))
	})

	adminRouter := chi.NewRouter()
	adminRouter.Use(NewLogger(log))
	adminRouter.Route("/", func(r chi.Router) {

		r.Post("/sign-in", h.AdminSignIn(ctx, log))

		admManagerRouter := chi.NewRouter()
		admManagerRouter.Use(h.AdminIdentityMiddleware())
		admManagerRouter.Use(NewLogger(log))

		admManagerRouter.Route("/manage-users", func(r chi.Router) {
			r.Get("/users", h.ShowUsers(ctx, log))
			r.Post("/delete-user/{id}", func(w http.ResponseWriter, r *http.Request) {})
		})

		admManagerRouter.Route("/manage-cakes", func(r chi.Router) {
			r.Get("/cakes", h.Cakes(ctx, log))
			r.Post("/add-cakes", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/remove-cakes", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/update-cake/{id}", func(w http.ResponseWriter, r *http.Request) {})
		})

		admManagerRouter.Route("/database", func(r chi.Router) {
			// backup - database dump
			r.Post("/backup", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/recovery", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Mount("/", admManagerRouter)
	})

	router.Mount("/", apiRouter)
	router.Mount("/adm", adminRouter)

	return router
}
