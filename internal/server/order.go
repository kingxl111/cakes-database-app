package server

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/kingxl111/cakes-database-app/internal/models"

	"github.com/go-chi/render"
)

func (h *Handler) MakeOrder(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.make-order: "
		var req models.MakeOrderRequest

		userID := r.Context().Value(userCtx) // getting user id from middleware's context
		log.Info("userID from contextL ", userID)
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error(op, "failed to decode request", err)
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		orderID, err := h.services.OrderManager.CreateOrder(
			userID.(int),
			req.Delivery,
			req.Cakes,
			req.PaymentMethod)
		if err != nil {
			log.Error("Failed to create order: %v", err)
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Info(op + "user: " + strconv.Itoa(userID.(int)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse := models.MakeOrderResponse{
			OrderID:      orderID,
			DeliveryTime: "few days",
		}
		if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
			log.Error(err.Error())
		}
	}
}

func (h *Handler) ViewOrders(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.view-orders: "

		userID := r.Context().Value(userCtx)
		var resp models.GetOrdersResponse
		resp, err := h.services.OrderManager.GetOrders(userID.(int))
		if err != nil {
			log.Error("error from operation: " + op + err.Error())
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonData, err := json.Marshal(resp)
		if err != nil {
			log.Error("error encoding into JSON: " + err.Error())
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(jsonData); err != nil {
			log.Error(err.Error())
		}
		log.Info(op + "user: " + strconv.Itoa(userID.(int)))
	}
}
