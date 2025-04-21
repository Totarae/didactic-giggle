package router

import (
	"didactic-giggle/cmd/internal/database"
	"didactic-giggle/cmd/internal/handlers"
	"didactic-giggle/cmd/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler *handlers.Handler, logger *zap.Logger, db *database.DB, withdrawHandler *handlers.WithdrawHandler,
	userHandler *handlers.UserHandler, orderHandler *handlers.OrderHandler, balanceHandler *handlers.BalanceHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/user",
		func(r chi.Router) {
			r.Post("/register", userHandler.RegisterHandler)
			r.Post("/login", userHandler.LoginHandler)

			// Protected группа — только для аутентифицированных пользователей
			r.Group(func(protected chi.Router) {
				protected.Use(middleware.AuthMiddleware(db))
				protected.Post("/orders", orderHandler.UploadOrderHandler)
				protected.Get("/orders", orderHandler.GetOrdersHandler)
				protected.Get("/balance", balanceHandler.GetBalanceHandler)
				protected.Post("/balance/withdraw", withdrawHandler.WithdrawHandler)
				protected.Get("/withdrawals", withdrawHandler.GetWithdrawalsHandler)
			})
		})

	return r
}
