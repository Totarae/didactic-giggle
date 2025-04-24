package middleware

import (
	"context"
	"didactic-giggle/cmd/internal/database"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type contextKey string

const (
	loginKey  contextKey = "login"
	userIDKey contextKey = "user_id"
)

// WithLogin добавляет логин в контекст
func WithLogin(ctx context.Context, login string) context.Context {
	return context.WithValue(ctx, loginKey, login)
}

// WithUserID добавляет ID пользователя в контекст
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// LoginFromContext возвращает логин из контекста
func LoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(loginKey).(string)
	return login, ok
}

// UserIDFromContext возвращает ID пользователя из контекста
func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

// AuthMiddleware проверяет наличие куки auth и добавляет логин и userID в контекст
func AuthMiddleware(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth")
			if err != nil || cookie.Value == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			login := cookie.Value

			userID, err := db.GetUserIDByLogin(r.Context(), login)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := WithLogin(r.Context(), login)
			ctx = WithUserID(ctx, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LoggingMiddleware логирует запросы с помощью zap.Logger
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("duration", duration),
			)
		})
	}
}
