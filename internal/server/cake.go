package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (h *Handler) Cakes(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.get-cakes: "

		userID := r.Context().Value(userCtx)
		cakes, err := h.services.CakeManager.GetCakes()
		if err != nil {
			log.Error(op, "failed to get cakes", err)
			newErrorResponse(w, http.StatusInternalServerError, "can't get cakes")
			return
		}

		jsonData, err := json.Marshal(cakes)
		if err != nil {
			log.Error(op, "failed to marshal cakes", err)
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(jsonData); err != nil {
			log.Error("failed to write response", err)
		}
		log.Info(op + "user: " + strconv.Itoa(userID.(int)))
	}
}
