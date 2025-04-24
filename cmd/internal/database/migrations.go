package database

import (
	"database/sql"
	"didactic-giggle/cmd/internal/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"path/filepath"
)

// runPgMigrations runs Postgres migrations
func RunPgMigrations(cfg *config.Config) error {
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
