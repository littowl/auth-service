package utils

import "os"

type Config struct {
	HTTPAddr string
	GRPCAddr string
	PgHost   string
	PgPort   string
	PgDB     string
	PgUser   string
	PgPass   string
}

func NewConfig() *Config {
	return &Config{
		HTTPAddr: fromEnv("HTTP_ADDR", "127.0.0.1:8081"),
		GRPCAddr: fromEnv("GRPC_ADDR", "127.0.0.1:3201"),
		PgHost:   fromEnv("POSTGRES_HOST", "localhost"),
		PgPort:   fromEnv("POSTGRES_PORT", "5432"),
		PgDB:     fromEnv("POSTGRES_DB", "postgres"),
		PgUser:   fromEnv("POSTGRES_USER", "postgres"),
		PgPass:   fromEnv("POSTGRES_PASSWORD", "6355"),
	}
}

func fromEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
