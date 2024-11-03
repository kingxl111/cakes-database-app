package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) Cakes(ctx *context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.get-cakes: "

		userID := r.Context().Value(userCtx)
		cakes, err := h.services.CakeManager.GetCakes()
		if err != nil {
			// h.services.Logger.WriteLog("ERROR", op + err.Error())
			newErrorResponse(w, http.StatusInternalServerError, "can't get cakes")
			return
		}

		jsonData, err := json.Marshal(cakes)
		if err != nil {
			// h.services.Logger.WriteLog("ERROR", op + err.Error())
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(jsonData); err != nil {
			log.Error(err.Error())
		}
		log.Info(op + "user: " + strconv.Itoa(userID.(int)))
		//h.services.Logger.WriteLog("INFO", op+"user: "+strconv.Itoa(userID.(int)))
	}
}
