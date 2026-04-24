package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config punya alasan berubah saat kebutuhan konfigurasi aplikasi berubah,
// bukan saat aturan bisnis item atau query database berubah.
type Config struct {
	AppPort           string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()

	dbMaxOpenConns, err := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	if err != nil {
		return Config{}, fmt.Errorf("nilai DB_MAX_OPEN_CONNS tidak valid: %w", err)
	}

	dbMaxIdleConns, err := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	if err != nil {
		return Config{}, fmt.Errorf("nilai DB_MAX_IDLE_CONNS tidak valid: %w", err)
	}

	dbConnMaxLifetime, err := getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	if err != nil {
		return Config{}, fmt.Errorf("nilai DB_CONN_MAX_LIFETIME tidak valid: %w", err)
	}

	cfg := Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "inventory_db"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns:    dbMaxOpenConns,
		DBMaxIdleConns:    dbMaxIdleConns,
		DBConnMaxLifetime: dbConnMaxLifetime,
	}

	if cfg.AppPort == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return Config{}, fmt.Errorf("konfigurasi aplikasi belum lengkap")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func getEnvAsDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}
