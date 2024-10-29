package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

const (
    authorizationHeader = "Authorization"
    userCtx             = "userId"
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
            if len(headerParts) != 2 || headerParts[0] != "Bearer" {
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

func NewLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}