package main

import (
	"context"
	"database/sql"
	"didactic-giggle/cmd/internal/config"
	"didactic-giggle/cmd/internal/database"
	"didactic-giggle/cmd/internal/handlers"
	"didactic-giggle/cmd/internal/router"
	"didactic-giggle/cmd/internal/services"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
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

	if err := runPgMigrations(cfg); err != nil {
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
	handler := handlers.NewHandler("/", logger, "mode", db)

	r := router.NewRouter(handler, logger, db)

	logger.Info("Сервер запущен на ", zap.String("address", cfg.ServerAddress))
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		logger.Fatal("Ошибка при запуске сервера: ", zap.Error(err))
	}
}

// runPgMigrations runs Postgres migrations
func runPgMigrations(cfg *config.Config) error {
	if cfg.PgMigrationsPath == "" {
		return nil
	}
	if cfg.DatabaseDSN == "" {
		return errors.New("no cfg.DatabaseDSN provided")
	}

	absPath, err := filepath.Abs(cfg.PgMigrationsPath)
	if err != nil {
		return fmt.Errorf("не удалось получить абсолютный путь: %w", err)
	}

	// Source: file://
	sourceDriver, err := (&file.File{}).Open("file://" + filepath.ToSlash(absPath))
	if err != nil {
		return fmt.Errorf("ошибка открытия миграций: %w", err)
	}

	// Подключение через sql.DB
	sqlDB, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	defer sqlDB.Close()

	// БД-драйвер
	dbDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("ошибка создания postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("file", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("ошибка создания мигратора: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("ошибка применения миграции: %w", err)
	}

	return nil
}
