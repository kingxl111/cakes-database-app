package server

import (
	"cakes-database-app/pkg/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

// Handler's high-level method
func (h *Handler) SignUp(c *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.User

		defer r.Body.Close()
        err := render.DecodeJSON(r.Body, &req); 
        if err != nil {
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }

        id, err := h.services.CreateUser(req)
        if err != nil {
            newErrorResponse(w, http.StatusInternalServerError, err.Error())
            return
        }

        w.WriteHeader(http.StatusOK)
        jsonResponse := map[string]interface{}{
            "id": id,
        }
        json.NewEncoder(w).Encode(jsonResponse)
	}
}

func (h *Handler) SignIn(c *context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        var input signInInput
        defer r.Body.Close()
        err := render.DecodeJSON(r.Body, &input); 
        if err != nil {
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return
        }

        token, err := h.services.GenerateToken(input.Username, input.Password)
        if err != nil {
            newErrorResponse(w, http.StatusBadRequest, err.Error())
            return 
        }

        log.Printf("generated token for user %s: %s", input.Username, token)
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