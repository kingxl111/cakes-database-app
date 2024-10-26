package server

import (
	"cakes-database-app/pkg/models"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)


// Handler's high-level method
func (h *Handler) SignUp(c *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.User
		ct := r.Context()
		_ = ct

		defer r.Body.Close()
        if err := render.DecodeJSON(r.Body, &req); err != nil {
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

// func (h *Handler) signIn(c *context.Context) {

// }
	
// type signInInput struct {
// 	Username string `json:"username" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }