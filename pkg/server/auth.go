package server

import (
	"cakes-database-app/pkg/models"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// Handler's high-level method
func (h *Handler) SignUp(c *context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.sign-up"
        log := log.With(
			slog.String("op", op),
			slog.String("time", time.Now().Format("2024-10-29 21:03:54")),
		)

		var req models.User
		defer r.Body.Close()
        err := render.DecodeJSON(r.Body, &req); 
        if err != nil {
            log.Info("error operation: %s: %s", op, err.Error())
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }

        id, err := h.services.CreateUser(req)
        if err != nil {
            log.Info("error operation: %s: %s", op, err.Error())
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }
        log.Info("new user registered: ", req.FullName, 1)
        w.WriteHeader(http.StatusOK)
        jsonResponse := map[string]interface{}{
            "id": id,
        }
        json.NewEncoder(w).Encode(jsonResponse)
	}
}

func (h *Handler) SignIn(c *context.Context, log *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.sign-in"
        log := log.With(
			slog.String("op", op),
			slog.String("time", time.Now().Format("2024-10-29 21:03:54")),
		)

        var input signInInput
        defer r.Body.Close()
        err := render.DecodeJSON(r.Body, &input); 
        if err != nil {
            log.Info("error operation: %s: %s", op, err.Error())
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }

        token, err := h.services.GenerateToken(input.Username, input.Password)
        if err != nil {
            log.Info("error operation: %s: %s", op, err.Error())
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return 
        }

        log.Info("generated token for user %s: %s", input.Username, token)
        w.WriteHeader(http.StatusOK)
        jsonResponse := map[string]interface{}{
            "token": token,
        }
        json.NewEncoder(w).Encode(jsonResponse)
    }
}
	
type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}