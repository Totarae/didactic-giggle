package router

import (
	"didactic-giggle/cmd/internal/database"
	"didactic-giggle/cmd/internal/handlers"
	"didactic-giggle/cmd/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler *handlers.Handler, logger *zap.Logger, db *database.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/user",
		func(r chi.Router) {
			r.Post("/register", handler.RegisterHandler)
			r.Post("/login", handler.LoginHandler)

			// Protected группа — только для аутентифицированных пользователей
			r.Group(func(protected chi.Router) {
				protected.Use(middleware.AuthMiddleware(db))
				protected.Post("/orders", handler.UploadOrderHandler)
				protected.Get("/orders", handler.GetOrdersHandler)
				protected.Get("/balance", handler.GetBalanceHandler)
				protected.Post("/balance/withdraw", handler.WithdrawHandler)
				protected.Get("/withdrawals", handler.GetWithdrawalsHandler)
			})

			//r.Get("/balance", getBalanceHandler)
			//r.Post("/balance/withdraw", withdrawHandler)
			//r.Get("/withdrawals", getWithdrawalsHandler)
		})

	return r
}
