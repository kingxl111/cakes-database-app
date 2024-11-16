package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

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
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error(err.Error())
			}
		}()

		orderID, err := h.services.OrderManager.CreateOrder(
			userID.(int),
			req.Delivery,
			req.Cakes,
			req.PaymentMethod)
		if err != nil {
			log.Error("Failed to create order:", err)
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Info(op, "user: ", userID)
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
		log.Info(op, "user: ", userID)
	}
}

func (h *Handler) CancelOrder(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.cancel-order: "
		var req models.CancelOrderRequest

		userID := r.Context().Value(userCtx)
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			log.Error(op, "failed to decode request", err)
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error(err.Error())
			}
		}()

		log.Info(op, "user_id: ", userID, ", order_id: ", req.OrderID)

		err = h.services.DeleteOrder(userID.(int), req.OrderID)
		if err != nil {
			log.Error(op, "failed to delete order", err)
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) DeliveryPoints(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.delivery-points: "
		userID := r.Context().Value(userCtx)

		delPoints, err := h.services.GetDeliveryPoints()
		if err != nil {
			log.Error(op, "failed to get delivery points: ", err)
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		jsonData, err := json.Marshal(delPoints)
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
		log.Info(op, "user: ", userID)
	}
}

func (h *Handler) UpdateOrder(ctx *context.Context, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.update-order: "
		var req models.UpdateOrderRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			log.Error(op, "failed to decode request", err)
		}

		userID := r.Context().Value(userCtx).(int)
		err = h.services.UpdateOrder(userID, req.OrderID, req.PaymentMethod)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			log.Error(op, "failed to update order", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
