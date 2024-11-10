package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"time"

	"github.com/go-chi/render"
)

func (h *Handler) AdminSignIn(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.sign-in"
		log := log.WithFields(logrus.Fields{
			"op":   op,
			"time": time.Now().Format(time.RFC3339),
		})
		var input signInInput

		err := render.DecodeJSON(r.Body, &input)
		if err != nil {
			log.Error(op, "can't decode json", err)
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		token, err := h.services.GenerateAdminToken(input.Username, input.Password)
		if err != nil {
			log.Error(op, "can't generate admin token", err.Error())
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("generated token for admin:", input.Username, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse := map[string]interface{}{
			"token": token,
		}
		if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
			log.Error(err.Error())
		}
	}
}

func (h *Handler) ShowUsers(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.show-users: "

		users, err := h.services.AdminService.GetUsers()
		if err != nil {
			log.Error("error operation:", op, err.Error())
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		adminID := r.Context().Value(userCtx).(int)
		log.Info(op, "adminID:", adminID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Error(err.Error())
		}
	}
}
