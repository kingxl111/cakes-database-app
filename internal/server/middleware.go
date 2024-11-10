package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	bear                = "Bearer"
)

func (h *Handler) UserIdentityMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get(authorizationHeader)
			if header == "" {
				http.Error(w, "empty auth header", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != bear {
				http.Error(w, "invalid auth header", http.StatusUnauthorized)
				return
			}

			if len(headerParts[1]) == 0 {
				http.Error(w, "token is empty", http.StatusUnauthorized)
				return
			}

			// main part: validation token
			userId, err := h.services.Authorization.ParseToken(headerParts[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userCtx, userId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func (h *Handler) AdminIdentityMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get(authorizationHeader)
			if header == "" {
				http.Error(w, "empty auth header", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(w, "invalid auth header", http.StatusUnauthorized)
				return
			}

			if len(headerParts[1]) == 0 {
				http.Error(w, "token is empty", http.StatusUnauthorized)
				return
			}

			// main part: validation token
			adminID, err := h.services.AdminAuthorization.ParseAdminToken(headerParts[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userCtx, adminID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func NewLogger(log *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		log.Info("logger middleware enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Создаем новый лог-запись с дополнительными полями
			entry := log.WithFields(logrus.Fields{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
				"request_id":  middleware.GetReqID(r.Context()),
			})

			// Оборачиваем ResponseWriter для захвата статуса и байтов
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Запоминаем время начала запроса
			t1 := time.Now()
			defer func() {
				// Логируем завершение запроса
				entry.WithFields(logrus.Fields{
					"status":   ww.Status(),
					"bytes":    ww.BytesWritten(),
					"duration": time.Since(t1).String(),
				}).Info("request completed")
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
