package config

import (
	"os"
	"strconv"
)

// Config contém todas as configurações da aplicação
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig contém configurações do banco de dados
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig contém configurações do servidor
type ServerConfig struct {
	Port    string
	GinMode string
}

// Load carrega as configurações das variáveis de ambiente
func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "alderaan"),
			Password: getEnv("DB_PASSWORD", "alderaan123"),
			DBName:   getEnv("DB_NAME", "alderaan_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
	}
}

// getEnv retorna o valor da variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retorna o valor da variável de ambiente como int ou um valor padrão
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
