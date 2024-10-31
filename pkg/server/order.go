package server

import (
	"cakes-database-app/pkg/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

func (h *Handler) MakeOrder(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.make-order: "
		var req models.MakeOrderRequest
		
		userID := r.Context().Value(userCtx)  // getting user id from middleware's context
		// log.Println("userID from contextL ", userID)
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.services.Logger.WriteLog( "ERROR", op + err.Error())
			newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
		}
		defer r.Body.Close()
		
		// log.Printf("userID: %v\npayment_method: %v\n", req.UserID, req.PaymentMethod)
		orderID, err := h.services.OrderManager.CreateOrder(
			userID.(int),
			req.Delivery,
			req.Cakes,
			req.PaymentMethod)
		if err != nil {
			// log.Printf("Failed to create order: %v", err)
			h.services.Logger.WriteLog("ERROR", op + err.Error())
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return 
		}

		h.services.Logger.WriteLog("INFO", op + "user: " + strconv.Itoa(userID.(int)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
        jsonResponse := models.MakeOrderResponse {
            OrderID: orderID,
			DeliveryTime: "few days",
        }
        json.NewEncoder(w).Encode(jsonResponse)
	}
}

func (h *Handler) ViewOrders(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.view-orders: "

		var req models.ViewOrdersRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Printf("error from operation: %s: %s", op, err.Error())
			h.services.Logger.WriteLog("ERROR", op + err.Error())
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		userID  := r.Context().Value(userCtx)
		var resp models.GetOrdersResponse
		resp, err = h.services.OrderManager.GetOrders(userID.(int))
		if err != nil {
			log.Printf("error from operation: %s: %s", op, err.Error())
			h.services.Logger.WriteLog( "ERROR", op + err.Error())
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return 
		}

        jsonData, err := json.Marshal(resp)
		if err != nil {
			log.Printf("error encoding into JSON: %v", err)
			h.services.Logger.WriteLog("ERROR", op + err.Error())
			newErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		h.services.Logger.WriteLog("INFO", op + "user: " + strconv.Itoa(userID.(int)))
	}
}