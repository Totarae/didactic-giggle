package router

import (
	"didactic-giggle/cmd/internal/database"
	"didactic-giggle/cmd/internal/handlers"
	"didactic-giggle/cmd/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler *handlers.Handler, logger *zap.Logger, db *database.DB, appHandlers *handlers.AppHandlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware(logger))

	r.Route("/api/user",
		func(r chi.Router) {
			r.Post("/register", appHandlers.User.RegisterHandler)
			r.Post("/login", appHandlers.User.LoginHandler)

			// Protected группа — только для аутентифицированных пользователей
			r.Group(func(protected chi.Router) {
				protected.Use(middleware.AuthMiddleware(db))
				// Собрал API методы по группам (матрешка)
				protected.Route("/orders", func(order chi.Router) {
					order.Post("/", appHandlers.OrderWriter.UploadOrderHandler)
					order.Get("/", appHandlers.OrderReader.GetOrdersHandler)
				})

				protected.Route("/balance", func(balance chi.Router) {
					balance.Get("/", appHandlers.Balance.GetBalanceHandler)
					balance.Post("/withdraw", appHandlers.Withdraw.WithdrawHandler)
				})

				protected.Get("/withdrawals", appHandlers.Withdraw.GetWithdrawalsHandler)
			})
		})

	return r
}
