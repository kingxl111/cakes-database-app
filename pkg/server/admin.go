package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
)

func (h *Handler) AdminSignIn(ctx *context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.sign-in"
        log := log.With(
			slog.String("op", op),
			slog.String("time", time.Now().Format("2024-10-29 21:03:54")),
		)
        var input signInInput
        defer r.Body.Close()
        err := render.DecodeJSON(r.Body, &input); 
        if err != nil {
            h.services.Logger.WriteLog("ERROR", op + err.Error())
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }

        token, err := h.services.GenerateAdminToken(input.Username, input.Password)
        if err != nil {
            log.Info("error operation: %s: %s", op, err.Error())
            h.services.Logger.WriteLog("ERROR", op + err.Error())
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return 
        }

        log.Info("generated token for user %s: %s", input.Username, token)
        h.services.Logger.WriteLog("INFO", op + ": admin:" + input.Username)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        jsonResponse := map[string]interface{}{
            "token": token,
        }
        json.NewEncoder(w).Encode(jsonResponse)
	}
}

func (h *Handler) ShowUsers(ctx *context.Context, log *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.admin.show-users: "

        users, err := h.services.AdminService.GetUsers()
        if err != nil {
            log.Info("error operation: %s: %s", op, err.Error())
            h.services.Logger.WriteLog("ERROR", op + err.Error())
            newErrorResponse(w, http.StatusInternalServerError, err.Error())
            return 
        }
        adminID := r.Context().Value(userCtx).(int)
        log.Info(op, "adminID: %v", adminID)        
        h.services.Logger.WriteLog("INFO", op + "adminID: " + strconv.Itoa(adminID))

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(users)
    }
}