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
		"STARTING TheSweetsOfLifeApp",
		slog.String("env", env),
		slog.String("version", "1.0"),
	)
	log.Debug("debug messages are enabled")
	if err := h.services.Logger.WriteLog("INFO", "new router started!"); err != nil {
		log.Info(err.Error())
	}

	router := chi.NewRouter()
	router.Use(NewLogger(log))
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

		r.Get("/cakes", h.Cakes(ctx))
	})


	adminRouter := chi.NewRouter()
	adminRouter.Use(NewLogger(log))
	adminRouter.Route("/", func(r chi.Router) {
		// TODO: handler funcs
		r.Post("/sign-in", h.AdminSignIn(ctx, log))
		
		admManagerRouter := chi.NewRouter()
		admManagerRouter.Use(h.AdminIdentityMiddleware())
		admManagerRouter.Use(NewLogger(log))

		admManagerRouter.Route("/manage-users", func(r chi.Router) {
			r.Get("/users", h.ShowUsers(ctx, log))
			r.Post("/delete-user/{id}", func(w http.ResponseWriter, r *http.Request) {})
		})

		admManagerRouter.Route("/manage-cakes", func(r chi.Router) {
			r.Get("/cakes", func(w http.ResponseWriter, r *http.Request) {})
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

