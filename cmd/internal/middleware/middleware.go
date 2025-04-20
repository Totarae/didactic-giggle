package middleware

import (
	"context"
	"didactic-giggle/cmd/internal/database"
	"net/http"
)

type contextKey string

const loginKey contextKey = "login"
const userIDKey contextKey = "user_id"

// LoginFromContext возвращает логин из контекста
func LoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(loginKey).(string)
	return login, ok
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

// AuthMiddleware проверяет наличие куки auth и добавляет логин в контекст
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

			ctx := context.WithValue(r.Context(), loginKey, login)
			ctx = context.WithValue(ctx, userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
