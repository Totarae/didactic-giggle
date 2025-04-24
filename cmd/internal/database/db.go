package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

// DB представляет подключение к БД
type DB struct {
	Pool   *pgxpool.Pool
	Logger *zap.Logger
}

// NewDB создает новое подключение к БД
func NewDB(logger *zap.Logger, dsn string) (*DB, error) {
	if dsn == "" {
		logger.Fatal("DATABASE_DSN is not set")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &DB{Pool: pool, Logger: logger}, nil
}

// Ping проверяет соединение с БД
func (db *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return db.Pool.Ping(ctx)
}

// Close закрывает соединение с БД
func (db *DB) Close() {
	db.Pool.Close()
}
