package server

import (
	"cakes-database-app/pkg/models"
	"context"
	"encoding/json"
	"net/http"
	"log"

	"github.com/go-chi/render"
)

func (h *Handler) MakeOrder(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.MakeOrderRequest
		
		userID := r.Context().Value(userCtx)  // getting user id from middleware's context
		log.Println("userID from contextL ", userID)
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
		}
		defer r.Body.Close()
		
		log.Printf("userID: %v\npayment_method: %v\n", req.UserID, req.PaymentMethod)
		orderID, err := h.services.OrderManager.CreateOrder(
			userID.(int),
			req.Delivery,
			req.Cakes,
			req.PaymentMethod)
		if err != nil {
			log.Printf("Failed to create order: %v", err)
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return 
		}

		w.WriteHeader(http.StatusOK)
        jsonResponse := map[string]interface{}{
            "orderID": orderID,
        }
        json.NewEncoder(w).Encode(jsonResponse)
	}
}