package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config хранит конфигурацию сервера
type Config struct {
	ServerAddress    string
	BaseURL          string
	FileStoragePath  string
	DatabaseDSN      string
	PgMigrationsPath string
	Mode             string
}

// NewConfig инициализирует конфигурацию на основе аргументов командной строки
func NewConfig() *Config {

	viper.SetDefault("SERVER_ADDRESS", "localhost:8080") // Значения по умолчанию
	viper.SetDefault("BASE_URL", "http://localhost:8080")
	viper.SetDefault("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("PG_MIGRATIONS_PATH", "migrations")

	viper.AutomaticEnv()

	// Если переменные окружения заданы — они имеют высший приоритет
	cfg := &Config{
		ServerAddress:    viper.GetString("SERVER_ADDRESS"),
		BaseURL:          viper.GetString("BASE_URL"),
		FileStoragePath:  viper.GetString("FILE_STORAGE_PATH"),
		DatabaseDSN:      viper.GetString("DATABASE_DSN"),
		PgMigrationsPath: viper.GetString("PG_MIGRATIONS_PATH"),
	}

	log.Printf("Инициализация конфигурации: ServerAddress=%s", cfg.ServerAddress)
	log.Printf("Инициализация конфигурации: BaseURL=%s", cfg.BaseURL)
	log.Printf("Инициализация конфигурации: FileStoragePath=%s", cfg.FileStoragePath)
	log.Printf("Инициализация конфигурации: DatabaseDSN=%s", cfg.DatabaseDSN)
	log.Printf("Инициализация конфигурации: PgMigrationsPath=%s", cfg.PgMigrationsPath)
	log.Printf("Инициализация конфигурации: Mode=%s", cfg.Mode)

	return cfg

}
