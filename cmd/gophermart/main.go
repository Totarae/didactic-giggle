package main

import (
	"context"
	"didactic-giggle/cmd/internal/config"
	"didactic-giggle/cmd/internal/database"
	"didactic-giggle/cmd/internal/handlers"
	"didactic-giggle/cmd/internal/router"
	"didactic-giggle/cmd/internal/services"
	"go.uber.org/zap"
	"net/http"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}
	defer logger.Sync()

	// Инициализация конфигурации
	cfg := config.NewConfig()

	logger.Info("Конфигурация сервера",
		zap.String("ServerAddress", cfg.ServerAddress),
		zap.String("DatabaseDSN", cfg.DatabaseDSN),
		zap.String("AccrualSystemAddress", cfg.AccrualSystemAddress),
	)

	var db *database.DB
	db, err = database.NewDB(logger, cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := database.RunPgMigrations(cfg); err != nil {
		logger.Fatal("runPgMigrations failed: ", zap.Error(err))
	}

	// стартуем демона
	orderSvc := services.NewOrderService(db, cfg.AccrualSystemAddress, logger)
	if cfg.Mode == "prod" {
		orderSvc.Start(context.Background())
		logger.Info("Воркер accrual-системы запущен")
	} else {
		logger.Info("Воркер отключён (MODE != prod)", zap.String("mode", cfg.Mode))
	}

	// Передача базового URL в обработчики
	handler := handlers.NewHandler("/", logger, "mode")
	withdrawHandler := handlers.NewWithdrawHandler(db, logger)
	userHandler := handlers.NewUserHandler(db, logger)
	orderHandler := handlers.NewOrderHandler(db, logger)
	balanceHandler := handlers.NewBalanceHandler(db, logger)

	// Соберем все обработчики в одну структуру для лаконичности
	appHandlers := &handlers.AppHandlers{
		User:     userHandler,
		Order:    orderHandler,
		Balance:  balanceHandler,
		Withdraw: withdrawHandler,
	}

	r := router.NewRouter(handler, logger, db, appHandlers)

	logger.Info("Сервер запущен на ", zap.String("address", cfg.ServerAddress))
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		logger.Fatal("Ошибка при запуске сервера: ", zap.Error(err))
	}
}
