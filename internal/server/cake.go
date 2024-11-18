package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

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

func (h *Handler) Cake(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.getCake"
		id := chi.URLParam(r, "id")
		if id == "" {
			log.Error("ID not found in URL")
			newErrorResponse(w, http.StatusBadRequest, "id not found in URL")
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Error("wrong id format: ", id)
			newErrorResponse(w, http.StatusBadRequest, "wrong id format")
			return
		}

		cake, err := h.services.GetCake(idInt)
		if err != nil {
			log.Error(op, "failed to get cake", err)
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonData, err := json.Marshal(cake)
		if err != nil {
			log.Error(op, "failed to marshal cake", err)
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(jsonData); err != nil {
			log.Error("failed to write response", err)
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
