package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Salvar variáveis de ambiente originais
	originalEnv := make(map[string]string)
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD",
		"DB_NAME", "DB_SSLMODE", "SERVER_PORT",
	}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Limpar após o teste
	defer func() {
		for key, val := range originalEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	t.Run("load with custom environment variables", func(t *testing.T) {
		// Configurar variáveis de ambiente
		os.Setenv("DB_HOST", "testhost")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_SSLMODE", "require")
		os.Setenv("SERVER_PORT", "9090")

		cfg := Load()

		if cfg.Database.Host != "testhost" {
			t.Errorf("DB_HOST = %v, want testhost", cfg.Database.Host)
		}
		if cfg.Database.Port != 5433 {
			t.Errorf("DB_PORT = %v, want 5433", cfg.Database.Port)
		}
		if cfg.Database.User != "testuser" {
			t.Errorf("DB_USER = %v, want testuser", cfg.Database.User)
		}
		if cfg.Database.Password != "testpass" {
			t.Errorf("DB_PASSWORD = %v, want testpass", cfg.Database.Password)
		}
		if cfg.Database.DBName != "testdb" {
			t.Errorf("DB_NAME = %v, want testdb", cfg.Database.DBName)
		}
		if cfg.Database.SSLMode != "require" {
			t.Errorf("DB_SSLMODE = %v, want require", cfg.Database.SSLMode)
		}
		if cfg.Server.Port != "9090" {
			t.Errorf("SERVER_PORT = %v, want 9090", cfg.Server.Port)
		}
	})

	t.Run("load with default values", func(t *testing.T) {
		// Limpar todas as variáveis
		for _, key := range envVars {
			os.Unsetenv(key)
		}

		cfg := Load()

		// Verificar valores padrão
		if cfg.Database.Host != "localhost" {
			t.Errorf("default DB_HOST = %v, want localhost", cfg.Database.Host)
		}
		if cfg.Database.Port != 5432 {
			t.Errorf("default DB_PORT = %v, want 5432", cfg.Database.Port)
		}
		if cfg.Database.User != "alderaan" {
			t.Errorf("default DB_USER = %v, want alderaan", cfg.Database.User)
		}
		if cfg.Database.Password != "alderaan123" {
			t.Errorf("default DB_PASSWORD = %v, want alderaan123", cfg.Database.Password)
		}
		if cfg.Database.DBName != "alderaan_db" {
			t.Errorf("default DB_NAME = %v, want alderaan_db", cfg.Database.DBName)
		}
		if cfg.Database.SSLMode != "disable" {
			t.Errorf("default DB_SSLMODE = %v, want disable", cfg.Database.SSLMode)
		}
		if cfg.Server.Port != "8080" {
			t.Errorf("default SERVER_PORT = %v, want 8080", cfg.Server.Port)
		}
	})

	t.Run("load with partial environment variables", func(t *testing.T) {
		// Limpar todas as variáveis
		for _, key := range envVars {
			os.Unsetenv(key)
		}

		// Configurar apenas algumas
		os.Setenv("DB_HOST", "customhost")
		os.Setenv("DB_PORT", "6543")

		cfg := Load()

		// Verificar valores customizados
		if cfg.Database.Host != "customhost" {
			t.Errorf("DB_HOST = %v, want customhost", cfg.Database.Host)
		}
		if cfg.Database.Port != 6543 {
			t.Errorf("DB_PORT = %v, want 6543", cfg.Database.Port)
		}

		// Verificar valores padrão
		if cfg.Database.User != "alderaan" {
			t.Errorf("DB_USER = %v, want alderaan (default)", cfg.Database.User)
		}
		if cfg.Server.Port != "8080" {
			t.Errorf("SERVER_PORT = %v, want 8080 (default)", cfg.Server.Port)
		}
	})

	t.Run("load with invalid port number", func(t *testing.T) {
		os.Unsetenv("DB_PORT")
		os.Setenv("DB_PORT", "invalid")

		cfg := Load()

		// Deve usar o valor padrão quando a conversão falhar
		if cfg.Database.Port != 5432 {
			t.Errorf("DB_PORT with invalid value = %v, want 5432 (default)", cfg.Database.Port)
		}
	})
}

func TestConfig_Struct(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "testhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
		Server: ServerConfig{
			Port: "3000",
		},
	}

	t.Run("database config", func(t *testing.T) {
		if cfg.Database.Host != "testhost" {
			t.Errorf("Host = %v, want testhost", cfg.Database.Host)
		}
		if cfg.Database.Port != 5432 {
			t.Errorf("Port = %v, want 5432", cfg.Database.Port)
		}
		if cfg.Database.User != "testuser" {
			t.Errorf("User = %v, want testuser", cfg.Database.User)
		}
		if cfg.Database.Password != "testpass" {
			t.Errorf("Password = %v, want testpass", cfg.Database.Password)
		}
		if cfg.Database.DBName != "testdb" {
			t.Errorf("DBName = %v, want testdb", cfg.Database.DBName)
		}
		if cfg.Database.SSLMode != "disable" {
			t.Errorf("SSLMode = %v, want disable", cfg.Database.SSLMode)
		}
	})

	t.Run("server config", func(t *testing.T) {
		if cfg.Server.Port != "3000" {
			t.Errorf("Port = %v, want 3000", cfg.Server.Port)
		}
	})
}

// Benchmark
func BenchmarkLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Load()
	}
}
