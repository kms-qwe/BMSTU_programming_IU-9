package config

import (
	"log"
	"os"
	"strconv"
)

const (
	defaultHTTPHost = "0.0.0.0"
	defaultHTTPPort = 8080
	defaultAppEnv   = "local"
)

type Config struct {
	AppEnv   string
	HTTP     HTTPConfig
	Database DatabaseConfig
}

type HTTPConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func MustLoad() Config {
	return Config{
		AppEnv: getEnv("APP_ENV", defaultAppEnv),
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", defaultHTTPHost),
			Port: getEnvAsInt("HTTP_PORT", defaultHTTPPort),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "postgres"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "coffee"),
			Password:        getEnv("DB_PASSWORD", "coffee"),
			Name:            getEnv("DB_NAME", "coffee"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 30),
		},
	}
}

func (c HTTPConfig) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c DatabaseConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.Name + "?sslmode=" + c.SSLMode
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid integer value for %s: %v", key, err)
	}

	return number
}
